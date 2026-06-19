package parser

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	transcript "github.com/Arvind215271/askito/internal/youtube/transcript"
)

// YouTubeJSON3 Structure mapping YouTube's dynamic text layouts
type YouTubeJSON3 struct {
	Events []JSON3Event `json:"events"`
}

type JSON3Event struct {
	// YoutubeJSON3 stores start and duration in miliseconds
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



func ParseJSON3ToSegments(jsonData []byte) ([]transcript.TranscriptSegment, error) {
	var yt YouTubeJSON3
	if err := json.Unmarshal(jsonData, &yt); err != nil {
		return nil, err
	}

	segments := make([]transcript.TranscriptSegment, 0, len(yt.Events))

	for _, event := range yt.Events {
		// build text from segs
		var b strings.Builder

		for _, seg := range event.Segments {
			b.WriteString(seg.UTF8)
		}

		text := strings.TrimSpace(b.String())
		if text == "" {
			continue
		}

		start := float64(event.StartMs) / 1000.0
		end := start + float64(event.DurationMs)/1000.0

		if event.DurationMs > 0 {
			end = start + float64(event.DurationMs)/1000.0
		} else {
			end = start
		}

		segments = append(segments, transcript.TranscriptSegment{
			Start: start,
			End:   end,
			Text:  text,
		})
	}

	return segments, nil
}