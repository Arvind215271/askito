package wordstats



type WordStats struct {
	Word string `json:"word"`

	// raw signals
	Count    int     `json:"count"`
	Duration float64 `json:"duration"`

	// computed importance
	Score float64 `json:"score"`

	// derived structure
	Bucket int `json:"bucket"`

	// optional: debugging / filtering support
	Pruned bool `json:"pruned,omitempty"`
}


type Bucket struct {
	ID    int      `json:"id"`
	Words []string `json:"words"`

	CountMin int `json:"count_min"`
	CountMax int `json:"count_max"`

	DurationMin float64 `json:"duration_min"`
	DurationMax float64 `json:"duration_max"`
}



type Result struct {
	// core output structure
	Buckets []Bucket `json:"buckets"`

	// compression metadata
	OriginalChars int     `json:"original_chars"`
	OutputChars   int     `json:"output_chars"`
	CompressionRatio float64 `json:"compression_ratio"`

	// word-level stats
	OriginalWords int `json:"original_words"`
	KeptWords     int `json:"kept_words"`
}

