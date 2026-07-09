package wordstats

import (
	// "fmt"
)

type AnalysisConfig struct {
	UseHeavyStopWords bool

	MinFreq int
	// minFreq of 0 means again include everything

	Depth float64 // 0.0 - 1.0
	// Depth 0 by default means inlcude everything.  

	MaxChars int
	// 0 means include everything

	WindowSize float64   // how many seconds each window is using
	// must be specified...
	// cannot be zero for window word stats
	BucketCount int      // how many bucekts to create per window
	// must be specified.
	// cannot be
}

const (
	DefaultBucketCount = 8
	DefaultWindowSize  = 60.0
	DefaultMinFreq     = 5
)

func DefaultWordStatsConfig() AnalysisConfig {
	return AnalysisConfig{
		UseHeavyStopWords: true,
		MinFreq:           DefaultMinFreq,
		BucketCount:       DefaultBucketCount,
	}
}

func DefaultWindowStatsConfig() AnalysisConfig {
	cfg := DefaultWordStatsConfig()
	cfg.WindowSize = DefaultWindowSize
	return cfg
}