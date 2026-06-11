// ./internal/youtube/transcript.go
package youtube

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"regexp"
	"strings"

	"encoding/json"
	"time"
)

// check if yt-dlp exist in the system or not
func ValidateYTDLP() error {
	_, err := exec.LookPath("yt-dlp")
	return err
}

// get the transcript using YTDLP as it is the only otpion to get those at cheap cost than youtube API quota
func GetTranscript(
	ctx context.Context,
	videoID string,
) (*Transcript, error) {

	videoURL := fmt.Sprintf(
		"https://www.youtube.com/watch?v=%s",
		videoID,
	)

	// create a temporray folder
	// ytdlp gives us a .vtt file. Thus, we need to download it somewhere and then store it.
	// the * means any valid name that is possible ig? 
	tmpDir, err := os.MkdirTemp("", "askito-transcript-*")
	if err != nil {
		return nil, err
	}

	// remove this folder when the funciton exist 
	// as it is no longer needed to be used afterwards
	defer os.RemoveAll(tmpDir)

	// execute the command 	
	// auto-subs and subs both can be downlaoded.. 
	// here we have en.* means we only download english ones that exist here.
	// skip-download here skips the video from downloading and only download the subtitle here.
	// -o means output to be sotred in tmpDir and whatever id and extension we are using for it.
	// 1. The optimal yt-dlp command configuration
	cmd := exec.CommandContext(
		ctx,
		"yt-dlp",
		"--write-auto-subs",
		"--write-subs",
		"--sub-langs", "en,en-orig",
		"--sub-format", "vtt", // Pull the flat data structures first
		"--convert-subs", "vtt",           // Let yt-dlp convert it to clean VTT locally
		"--skip-download",
		"--no-playlist",
		"--extractor-args", "youtube:player_client=android",
		"-o", filepath.Join(tmpDir, "%(id)s.%(ext)s"),
		videoURL,
	)
	// here... this is for storing the error message
	// as error might occur. Why buffer? Because we are downloading it continuously... that is why
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	

	// this comamnd acutally runs the command. and chek for error here
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf(
			"yt-dlp failed: %w: %s",
			err,
			stderr.String(),
		)
	}

	fmt.Println("TMP DIR:", tmpDir)

	filesX, _ := os.ReadDir(tmpDir)
	for _, f := range filesX {
		fmt.Println("FOUND:", f.Name())
	}

	// find the files in the path... if they exist. and if yes.. get the ones with .vtt at the  end.
	files, err := filepath.Glob(
		filepath.Join(tmpDir, "*.vtt"),
	)
	if err != nil {
		return nil, err
	}

	entries, _ := os.ReadDir(tmpDir)

	fmt.Println("TMP DIR:", tmpDir)

	for _, e := range entries {
		fmt.Println("FILE:", e.Name())
	}

	if len(files) == 0 {
		return nil, fmt.Errorf(
			"transcript not available",
		)
	}

	raw, err := os.ReadFile(files[0])
	if err != nil {
		return nil, err
	}

	

	// create a text version for this file.
	text := ExtractCleanTranscriptTextAlreadyClean(string(raw))

	return &Transcript{
		Language: "en",
		Source:   TranscriptSourceYTDLP,
		Raw:      string(raw),
		Text:     text,
	}, nil
}

// GetTranscriptJSON fetches the pre-flattened structural JSON format from YouTube.
// This bypasses rolling text duplications and prevents post-processing conversion failures.
func GetTranscriptJSON(
	ctx context.Context,
	videoID string,
) (*Transcript, error) {

	videoURL := fmt.Sprintf(
		"https://www.youtube.com/watch?v=%s",
		videoID,
	)

	// Create a temporary folder to capture the download artifact
	tmpDir, err := os.MkdirTemp("", "askito-transcript-json-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	// Execute the command targeting json3 specifically
	cmd := exec.CommandContext(
		ctx,
		"yt-dlp",
		"--write-auto-subs",
		"--write-subs",
		"--sub-langs", "en,en-orig",
		"--sub-format", "json3", // Target the raw layout tree directly
		"--skip-download",
		"--no-playlist",
		"--extractor-args", "youtube:player_client=android",
		"-o", filepath.Join(tmpDir, "%(id)s.%(ext)s"),
		videoURL,
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf(
			"yt-dlp failed: %w: %s",
			err,
			stderr.String(),
		)
	}

	// Glob lookups search for the downloaded json3 asset
	files, err := filepath.Glob(
		filepath.Join(tmpDir, "*.json3"),
	)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("transcript JSON not available")
	}

	raw, err := os.ReadFile(files[0])
	if err != nil {
		return nil, err
	}

	// Parse and flatten the raw json3 tree down into a linear transcript string
	text, err := ExtractTextFromJSON3(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse json3 transcript: %w", err)
	}

	return &Transcript{
		Language: "en",
		Source:   TranscriptSourceYTDLP,
		Raw:      string(raw),
		Text:     text,
	}, nil
}


// YouTubeJSON3 Structure mapping YouTube's dynamic text layouts
type YouTubeJSON3 struct {
	Events []JSON3Event `json:"events"`
}

type JSON3Event struct {
	StartMs    int64         `json:"tStartMs"`
	DurationMs int64         `json:"dDurationMs"`
	Segments   []JSON3Segment `json:"segs"`
}

type JSON3Segment struct {
	UTF8 string `json:"utf8"`
}

// ExtractTextFromJSON3 loops through the structural tree and formats entries cleanly
func ExtractTextFromJSON3(jsonData []byte) (string, error) {
	var ytTranscript YouTubeJSON3
	if err := json.Unmarshal(jsonData, &ytTranscript); err != nil {
		return "", err
	}

	var transcriptLines []string

	for _, event := range ytTranscript.Events {
		var lineBuilder strings.Builder

		// Collect string snippets safely inside the segment loop
		for _, seg := range event.Segments {
			lineBuilder.WriteString(seg.UTF8)
		}

		cleanLine := strings.TrimSpace(lineBuilder.String())

				
		// Filter structural noise elements like [Music] lines or empty events
		if cleanLine == "" || cleanLine == "\n" {
			continue
		}

		// Convert historical milliseconds to human-readable format: HH:MM:SS
		d := time.Duration(event.StartMs) * time.Millisecond
		h := d / time.Hour
		d -= h * time.Hour
		m := d / time.Minute
		d -= m * time.Minute
		s := d / time.Second

		timestamp := fmt.Sprintf("%02d:%02d:%02d", h, m, s)
		transcriptLines = append(transcriptLines, fmt.Sprintf("[%s] %s", timestamp, cleanLine))
	}

	return strings.Join(transcriptLines, "\n"), nil
}

var (
	timestampRegex = regexp.MustCompile(
		`^\d{2}:\d{2}:\d{2}\.\d+\s-->.*$`,
	)

	tagRegex = regexp.MustCompile(
		`<[^>]+>`,
	)
)


func isRollingDuplicate(prev, curr string) bool {
	if prev == "" {
		return false
	}

	// exact match
	if prev == curr {
		return true
	}

	// curr is suffix extension of prev (rolling behavior)
	if strings.HasSuffix(curr, prev) {
		return true
	}

	// prev is prefix of curr (stream continuation)
	if strings.HasPrefix(curr, prev) {
		return true
	}

	// small edit distance shortcut (cheap heuristic)
	if len(prev) > 0 && len(curr) > 0 {
		diff := abs(len(curr) - len(prev))
		if diff <= 3 && strings.Contains(curr, prev) {
			return true
		}
	}

	return false
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}



func ExtractCleanTranscriptText(vtt string) string {
	lines := strings.Split(vtt, "\n")

	var result []string

	var prevText string
	var prevStable string
	var buffer []string

	flush := func(time string, text string) {
		text = strings.TrimSpace(text)
		if text == "" {
			return
		}

		if text != prevStable {
			result = append(result, fmt.Sprintf("[%s] %s", time, text))
			prevStable = text
		}
	}

	var currentTime string

	for _, raw := range lines {
		line := strings.TrimSpace(raw)

		if line == "" {
			continue
		}

		// skip headers
		if line == "WEBVTT" ||
			strings.HasPrefix(line, "Kind:") ||
			strings.HasPrefix(line, "Language:") {
			continue
		}

		// timestamp
		if timestampRegex.MatchString(line) {
			// flush previous cue
			if len(buffer) > 0 {
				text := strings.Join(buffer, " ")
				text = tagRegex.ReplaceAllString(text, "")
				text = strings.TrimSpace(text)

				// 🔥 CORE LOGIC: rolling subtitle filter
				if !isRollingDuplicate(prevText, text) {
					flush(currentTime, text)
					prevText = text
				}
			}

			parts := strings.Split(line, " --> ")
			currentTime = strings.TrimSpace(parts[0])
			buffer = buffer[:0]
			continue
		}

		line = tagRegex.ReplaceAllString(line, "")
		line = strings.TrimSpace(line)

		if line != "" {
			buffer = append(buffer, line)
		}
	}

	// final flush
	if len(buffer) > 0 {
		text := strings.Join(buffer, " ")
		text = tagRegex.ReplaceAllString(text, "")
		text = strings.TrimSpace(text)

		if !isRollingDuplicate(prevText, text) {
			flush(currentTime, text)
		}
	}

	return strings.Join(result, "\n")
}


// ExtractCleanTranscriptText parses a pre-flattened VTT file (No duplication logic required!)
func ExtractCleanTranscriptTextAlreadyClean(vtt string) string {
	var result []string
	lines := strings.Split(vtt, "\n")
	var currentTime string

	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || line == "WEBVTT" || strings.HasPrefix(line, "Kind:") || strings.HasPrefix(line, "Language:") {
			continue
		}

		// Grab the timestamp
		if timestampRegex.MatchString(line) {
			parts := strings.Split(line, " --> ")
			currentTime = strings.TrimSpace(parts[0])
			continue
		}

		// This line is guaranteed to be clean, static text without rolling duplicates
		line = tagRegex.ReplaceAllString(line, "")
		line = strings.TrimSpace(line)

		if line != "" {
			result = append(result, fmt.Sprintf("[%s] %s", currentTime, line))
		}
	}

	return strings.Join(result, "\n")
}