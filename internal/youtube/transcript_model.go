// ./internal/youtube/transcript_model.go
package youtube

type TranscriptSource string

const (
	TranscriptSourceYTDLP TranscriptSource = "yt-dlp"
)

type TranscriptSegment struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`

	Text string `json:"text"`
}

type Transcript struct {
	Language string           `json:"language"`
	Source   TranscriptSource `json:"source"`

	// Original VTT content
	Raw string `json:"raw"`

	// Cleaned transcript text
	Text string `json:"text"`

	// Optional timestamped segments
	Segments []TranscriptSegment `json:"segments,omitempty"`
}