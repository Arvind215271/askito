// ./internal/youtube/debug_summary.go
package debug

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"path/filepath"
	"strings"

	"github.com/Arvind215271/askito/internal/logger"
	"github.com/Arvind215271/askito/internal/youtube"
	"github.com/Arvind215271/askito/internal/youtube/chapter"
	"github.com/Arvind215271/askito/internal/youtube/transcript"
)

// DebugLLMContext fetches your playlist, runs your context reduction strategy,
// and exports both a dense JSON and an ultra-compressed plain text file.
func DebugLLMContext(
	ctx context.Context,
	log *logger.Logger,
	youtubeSvc *youtube.Service,
	transcriptSvc *transcript.Service,
	playlistID string,
) {
	fmt.Printf("\n================ COMPRESSION TEST RUN ================\n")

	// 1. Fetch the raw playlist data using your existing service architecture
	playlist, err := youtubeSvc.GetPlaylist(ctx, playlistID)
	if err != nil {
		log.Error("failed to fetch playlist for summary optimization test", "playlist_id", playlistID, "error", err)
		return
	}

	// Calculate baseline size before enrichment/compression metrics
	printSize("Raw Playlist Structural Metadata", playlist)

	// 2. Enrich the structural entities loop exactly like your layout
	for i := range playlist.Videos {
		// Populate Chapters
		playlist.Videos[i].Video.Chapters = chapter.ExtractChapters(playlist.Videos[i].Video.Description)

		// Fallback to Transcripts if Chapters do not exist
		if len(playlist.Videos[i].Video.Chapters.Text) == 0 && playlist.Videos[i].Video.CaptionAvailable {
			transcriptData, err := transcriptSvc.Get(ctx, playlist.Videos[i].Video.ID)
			if err == nil {
				playlist.Videos[i].Video.Transcript = transcriptData
			} else {
				log.Warn("transcript download skipped during fallback pipeline extraction", "video_id", playlist.Videos[i].Video.ID, "error", err)
			}
		}
	}

	// 3. Process data using your newly designed algorithm
	summaryPayload := GenerateLLMContext(playlist)

	// 4. Output Results and Size Comparisons
	fmt.Printf("\n-------------- VERIFICATION METRICS --------------\n")
	fmt.Println("Processed Title  :", summaryPayload.Title)
	fmt.Println("Channel Identity :", summaryPayload.ChannelTitle)
	fmt.Println("Total Curriculum :", summaryPayload.TotalVideos, "Videos Map Items")

	// --- A. Dense Compact JSON Generation ---
	compressedJSON, err := json.Marshal(summaryPayload)
	if err != nil {
		log.Error("failed to serialize dense optimized output payload matching target design schemas", "error", err)
		return
	}
	fmt.Printf("Dense JSON Payload Structural Size : %.2f KB (Length: %d chars)\n", float64(len(compressedJSON))/1024, len(compressedJSON))

	// --- B. Plain Text Curriculum Generation ---
	rawTextPayload := BuildRawCurriculumText(summaryPayload)
	fmt.Printf("Raw Text Payload Structural Size   : %.2f KB (Length: %d chars)\n", float64(len(rawTextPayload))/1024, len(rawTextPayload))

	// 5. Save Output Artifacts
	baseDir := filepath.Join("testdata", "youtube", "optimized_summaries")

	// Save the Dense JSON version
	jsonPath := filepath.Join(baseDir, playlistID+"_llm_curriculum.json")
	if err := saveFile(jsonPath, compressedJSON); err != nil {
		log.Error("failed storing generated dense optimization summary context json artifact", "path", jsonPath, "error", err)
	} else {
		fmt.Println("Successfully Exported JSON Profile:", jsonPath)
	}

	// Save the Ultra-Compressed Raw Text version
	textPath := filepath.Join(baseDir, playlistID+"_llm_curriculum.txt")
	if err := saveFile(textPath, []byte(rawTextPayload)); err != nil {
		log.Error("failed storing generated dense optimization raw text artifact", "path", textPath, "error", err)
	} else {
		fmt.Println("Successfully Exported Text Profile:", textPath)
	}
}

// BuildRawCurriculumText evaluates data purely based on text relationships, formatting, 
// and structural similarity—making it universally valid for any topic or playlist domain.
func BuildRawCurriculumText(summary PlaylistSummary) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Playlist: %s\nChannel: %s\n\n", summary.Title, summary.ChannelTitle))

	// Global universal timeline structural noise to filter out
	universalFluff := map[string]bool{
		"introduction": true, "intro": true, "outro": true, "conclusion": true,
		"summary": true, "recap": true, "agenda": true, "lesson": true, "part": true,
	}

	for _, v := range summary.Curriculum {
		var activeTopics []string
		seen := make(map[string]bool)

		// Clean up video branding/suffixes before running comparisons
		cleanVideoTitle := v.Title
		if idx := strings.IndexAny(cleanVideoTitle, "|-"); idx != -1 {
			cleanVideoTitle = strings.TrimSpace(cleanVideoTitle[:idx])
		}
		lowVideoTitle := strings.ToLower(cleanVideoTitle)

		for _, rawTopic := range v.Topics {
			// Isolate words to filter layout punctuation safely
			words := strings.Fields(strings.ToLower(rawTopic))
			var cleanWords []string

			for _, w := range words {
				// Clean off structural brackets/punctuation from the word boundaries
				w = strings.Trim(w, " -+:,[]()️/.\\•*")
				
				// Drop purely numeric markers, single characters, or generic timeline words
				if w == "" || len(w) <= 1 || universalFluff[w] {
					continue
				}
				cleanWords = append(cleanWords, w)
			}

			// Reconstruct the sanitized topic phrase
			if len(cleanWords) == 0 {
				continue
			}
			lowTopic := strings.Join(cleanWords, " ")

			// 1. Two-Way Containment & Similarity Test
			// Drops the topic if it heavily duplicates the main video title (or vice-versa)
			if strings.Contains(lowVideoTitle, lowTopic) || strings.Contains(lowTopic, lowVideoTitle) {
				continue
			}
			if computeIntersectionScore(lowTopic, lowVideoTitle) > 0.65 {
				continue
			}

			if !seen[lowTopic] {
				seen[lowTopic] = true
				// Capitalize first letter of words for clean presentation
				activeTopics = append(activeTopics, strings.Title(lowTopic))
			}
		}

		// Write output to string builder inline
		if len(activeTopics) > 0 {
			b.WriteString(fmt.Sprintf("#%d %s (%s)\n", v.Position, cleanVideoTitle, strings.Join(activeTopics, ", ")))
		} else {
			b.WriteString(fmt.Sprintf("#%d %s\n", v.Position, cleanVideoTitle))
		}
	}
	return b.String()
}

// computeIntersectionScore checks how heavily two topic tokens overlap structurally
func computeIntersectionScore(s1, s2 string) float64 {
	w1 := strings.Fields(s1)
	w2 := strings.Fields(s2)
	if len(w1) == 0 || len(w2) == 0 {
		return 0.0
	}

	matches := 0
	targetMap := make(map[string]bool)
	for _, word := range w2 {
		targetMap[word] = true
	}

	for _, word := range w1 {
		if targetMap[word] {
			matches++
		}
	}

	minSize := math.Min(float64(len(w1)), float64(len(w2)))
	return float64(matches) / minSize
}