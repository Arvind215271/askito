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

// ProcessingRequest defines how an already obtained transcript should be processed.
type ProcessingRequest struct {
	WindowSize float64 `json:"windowSize"` // optional: for grouping
	Output     string  `json:"output"`     // required: timeline-text | plain-text | segments
}
