package transcript

type TranscriptSource string

const (
	TranscriptSourceYTDLP TranscriptSource = "yt-dlp"
	TranscriptSourceJSON3 TranscriptSource = "json3"
)

type TranscriptSegment struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Text string `json:"text"`
}

type Transcript struct {
	Language string           `json:"language"`
	Source   TranscriptSource `json:"source"`
	Segments []TranscriptSegment `json:"segments,omitempty"`
}