package pipeline

type Request struct {
	Fields []string

	Subtitle   *SubtitleRequest
	Transcript *TranscriptRequest
	Signal     *SignalRequest
}

type SubtitleRequest struct {
	Type     string
	Language string
	Format   string
}

type TranscriptRequest struct {
	Type     string
	Language string
	Format   string
}

type SignalRequest struct {
	Analysis          string
	UseHeavyStopWords bool
	MinFreq           int
	Depth             float64
	WindowSize        float64
	BucketCount       int
}
