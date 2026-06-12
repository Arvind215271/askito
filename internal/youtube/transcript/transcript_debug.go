package transcript

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"strings"

// 	"github.com/Arvind215271/askito/internal/logger"
// 	transcript_parser "github.com/Arvind215271/askito/internal/youtube/transcript/parser"
// 	ytdlp "github.com/Arvind215271/askito/internal/youtube/transcript/provider/ytdlp"
// )

// func DebugTranscript(
// 	ctx context.Context,
// 	logger *logger.Logger,
// ) {

// 	const videoID = "dQw4w9WgXcQ" // Ah, the classic test video.

// 	logger.Info(
// 		"starting transcript debug",
// 		"video_id",
// 		videoID,
// 	)

// 	// ------------------------------------
// 	// validate yt-dlp
// 	// ------------------------------------
// 	ytdlpClient := ytdlp.Client{}

// 	if err := ytdlpClient.ValidateYTDLP(); err != nil {
// 		logger.Error(
// 			"yt-dlp not found",
// 			"error",
// 			err,
// 		)
// 		return
// 	}

// 	fmt.Printf("\n================ YT-DLP =================\n")
// 	fmt.Println("yt-dlp found")

// 	// ------------------------------------
// 	// Method A: Standard VTT Extraction
// 	// ------------------------------------
// 	logger.Info("fetching standard vtt transcript...")
// 	vttTranscript, err := ytdlpClient.GetTranscript(ctx, videoID)
// 	if err != nil {
// 		logger.Error(
// 			"failed to fetch VTT transcript",
// 			"error",
// 			err,
// 		)
// 		// Don't return early; try to run the JSON extraction anyway
// 	}

// 	// ------------------------------------
// 	// Method B: Flat JSON3 Extraction (The New Optimal Method)
// 	// ------------------------------------
// 	logger.Info("fetching flat json3 transcript...")
// 	jsonTranscript, err :=ytdlpClient.GetTranscript(ctx, videoID)
// 	if err != nil {
// 		logger.Error(
// 			"failed to fetch JSON3 transcript",
// 			"error",
// 			err,
// 		)
// 	}

// 	// Helper function to print comparisons cleanly
// 	printComparison := func(name string, text string) {
// 		fmt.Printf("\n================ %s =================\n", strings.ToUpper(name))
// 		fmt.Printf("Text Size: %d bytes\n\n", len(text))

// 		preview := text
// 		if len(preview) > 1000 {
// 			preview = preview[:1000] + "...\n[TRUNCATED]"
// 		}
// 		fmt.Println(preview)
// 	}

// 	// ------------------------------------
// 	// Print VTT Process Previews & Heuristics
// 	// ------------------------------------
// 	if vttTranscript != nil {
// 		fmt.Printf("\n================ METADATA (VTT) =================\n")
// 		fmt.Println("Language :", vttTranscript.Language)
// 		fmt.Println("Source   :", vttTranscript.Source)
// 		fmt.Println("Raw Size :", len(vttTranscript.Raw), "bytes")

// 		fmt.Printf("\n================ RAW VTT PREVIEW =================\n")
// 		rawPreview := vttTranscript.Raw
// 		if len(rawPreview) > 1000 {
// 			rawPreview = rawPreview[:1000] + "...\n[TRUNCATED]"
// 		}
// 		fmt.Println(rawPreview)

// 		// VTT Algorithmic Parsers
// 		originalText := transcript_parser.ExtractCleanTranscriptText(vttTranscript.Raw)
// 		screenBufferText := transcript_parser.ExtractScreenBufferTranscript(vttTranscript.Raw)
// 		tokenStreamText := transcript_parser.ExtractTokenStreamTranscript(vttTranscript.Raw)

// 		printComparison("Method 0: Original Heuristic", originalText)
// 		printComparison("Method 1: Screen Buffer State Machine", screenBufferText)
// 		printComparison("Method 2: Token Stream parser", tokenStreamText)

// 		// Export JSON structure representation using the legacy pipeline format
// 		vttTranscript.Text = screenBufferText
// 		exportJSON, _ := json.MarshalIndent(vttTranscript, "", "  ")
// 		fmt.Printf("\n================ EXPORT JSON (Using Method 1 State Machine) =================\n")
// 		jsonPreview := string(exportJSON)
// 		if len(jsonPreview) > 1000 {
// 			jsonPreview = jsonPreview[:1000] + "...\n[TRUNCATED]"
// 		}
// 		fmt.Println(jsonPreview)
// 		fmt.Printf("\nFinal serialized VTT structure size : %.2f KB\n", float64(len(exportJSON))/1024)
// 	}

// 	// ------------------------------------
// 	// Print Native JSON3 Pipeline Previews
// 	// ------------------------------------
// 	if jsonTranscript != nil {
// 		fmt.Printf("\n================ METADATA (NATIVE JSON3) =================\n")
// 		fmt.Println("Language :", jsonTranscript.Language)
// 		fmt.Println("Source   :", jsonTranscript.Source)
// 		fmt.Println("Raw Size :", len(jsonTranscript.Raw), "bytes")

// 		fmt.Printf("\n================ RAW YT JSON3 BACKEND PREVIEW =================\n")
// 		rawJson3Preview := jsonTranscript.Raw
// 		if len(rawJson3Preview) > 1000 {
// 			rawJson3Preview = rawJson3Preview[:1000] + "...\n[TRUNCATED]"
// 		}
// 		fmt.Println(rawJson3Preview)

// 		// Print the final string produced by the native struct parser block
// 		printComparison("Method 3: Native JSON3 Struct Unmarshal", jsonTranscript.Text)

// 		// Export internal structure representation using the JSON3 pipeline format
// 		exportJSON, _ := json.MarshalIndent(jsonTranscript, "", "  ")
// 		fmt.Printf("\n================ EXPORT JSON (Using Native JSON3 Parser) =================\n")
// 		jsonPreview := string(exportJSON)
// 		if len(jsonPreview) > 1000 {
// 			jsonPreview = jsonPreview[:1000] + "...\n[TRUNCATED]"
// 		}
// 		fmt.Println(jsonPreview)
// 		fmt.Printf("\nFinal serialized JSON3 structure size : %.2f KB\n", float64(len(exportJSON))/1024)
// 	}
// 	fmt.Println()
// }