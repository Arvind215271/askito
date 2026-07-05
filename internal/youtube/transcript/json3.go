package transcript

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// YouTubeJSON3 maps YouTube's JSON3 transcript structure.
type YouTubeJSON3 struct {
	Events []JSON3Event `json:"events"`
}

type JSON3Event struct {
	// YouTube stores timestamps in milliseconds.
	StartMs    int64          `json:"tStartMs"`
	DurationMs int64          `json:"dDurationMs"`
	Segments   []JSON3Segment `json:"segs"`
}

type JSON3Segment struct {
	UTF8 string `json:"utf8"`
}

// ExtractTextFromJSON3 converts JSON3 into a readable timestamped transcript.
func ExtractTextFromJSON3(jsonData []byte) (string, error) {
	var yt YouTubeJSON3

	if err := json.Unmarshal(jsonData, &yt); err != nil {
		return "", err
	}

	var transcriptLines []string

	for _, event := range yt.Events {
		var lineBuilder strings.Builder

		for _, seg := range event.Segments {
			lineBuilder.WriteString(seg.UTF8)
		}

		cleanLine := strings.TrimSpace(lineBuilder.String())

		if cleanLine == "" {
			continue
		}

		d := time.Duration(event.StartMs) * time.Millisecond

		h := d / time.Hour
		d -= h * time.Hour

		m := d / time.Minute
		d -= m * time.Minute

		s := d / time.Second

		timestamp := fmt.Sprintf("%02d:%02d:%02d", h, m, s)

		transcriptLines = append(
			transcriptLines,
			fmt.Sprintf("[%s] %s", timestamp, cleanLine),
		)
	}

	return strings.Join(transcriptLines, "\n"), nil
}

// ParseJSON3ToSegments converts JSON3 into transcript segments.
func ParseJSON3ToSegments(jsonData []byte) ([]TranscriptSegment, error) {
	var yt YouTubeJSON3

	if err := json.Unmarshal(jsonData, &yt); err != nil {
		return nil, err
	}

	segments := make([]TranscriptSegment, 0, len(yt.Events))

	for _, event := range yt.Events {

		var b strings.Builder

		for _, seg := range event.Segments {
			b.WriteString(seg.UTF8)
		}

		text := strings.TrimSpace(b.String())

		if text == "" {
			continue
		}

		start := float64(event.StartMs) / 1000.0

		end := start
		if event.DurationMs > 0 {
			end += float64(event.DurationMs) / 1000.0
		}

		segments = append(segments, TranscriptSegment{
			Start: start,
			End:   end,
			Text:  text,
		})
	}

	return segments, nil
}