package transcript

type TranscriptSegment struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Text  string  `json:"text"`
}

type Transcript struct {
	Language string              `json:"language"`
	Segments []TranscriptSegment `json:"segments,omitempty"`
}