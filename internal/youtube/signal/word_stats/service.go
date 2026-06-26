package wordstats

import transcript "github.com/Arvind215271/askito/internal/youtube/transcript"

type AnalysisConfig struct {
	UseHeavyStopWords bool

	MinFreq int

	Depth float64 // 0.0 - 1.0

	MaxChars int
}

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) AnalyzeVideo(
	t *transcript.Transcript,
	cfg AnalysisConfig,
) Result {

	// 1. raw extraction
	raw := BuildWordStats(t)

	// 2. scoring
	scored := ScoreWords(raw)

	// 3. optional heavy stopword filtering
	if cfg.UseHeavyStopWords {
		scored = FilterHeavyStopWords(scored)
	}

	// 4. main filtering rules
	scored = FilterWords(scored, FilterConfig{
		MinFreq:  cfg.MinFreq,
		Depth:    cfg.Depth,
		MaxChars: cfg.MaxChars,
	})

	// 5. bucketization
	buckets := CreateBuckets(scored)

	// 6. result
	return BuildResult(t, scored, buckets)
}

