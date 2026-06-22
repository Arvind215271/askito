// ./internal/youtube/debug/video.go
package video

import (
	"fmt"
	
	"github.com/Arvind215271/askito/internal/youtube"
)

type VideoTestResult struct {
	Name string `json:"name"`

	OriginalSize int `json:"original_size"`
	ResultSize   int `json:"result_size"`

	ReductionPercent float64 `json:"reduction_percent"`

	Output string `json:"output"`
}

type VideoTestSuite struct {
	VideoID string `json:"video_id"`
	Title   string `json:"title"`

	Results []VideoTestResult `json:"results"`
}

// DebugVideoTranscript executes all transcript-related experiments
// against a single video and exports the resulting artifacts.
func DebugVideoTranscript(video youtube.Video) error {
	if video.Transcript == nil {
		return fmt.Errorf("video transcript not available")
	}

	if err := LoadHeavyStopWords(
		"debug/video/heavy_stopwords.txt",
	); err != nil {
		return err
	}

	results := []VideoTestResult{
		// RunConceptDistillation(video ),

		// RunLightStopWordReductionTest(video),
		// RunHeavyStopWordReductionTest(video),

		// RunUniqueWordReductionTest(video),
		// RunHeavyStopWordUniqueWordTest(video),

		// RunWordStatsTest(video),
		// RunHeavyWordStatsTest(video),
		// RunHeavyWordStatsBucketFormatTest(video),

		RunWindowWordStatsTest(video),
		// RunWindowMergeWordStatsTest(video),
		RunWindowCompressedConceptStatsTest(video),
		
		// RunWindowWordStatsEncoded(video),
		// RunWindowFingerprint(video),
		// RunWindowCompressedConceptCodeMapTest(video),

		// NEW correct layer
		// RunTranscriptWindowTest(*video.Transcript, video.ID),
		// RunTranscriptTopicSignalTestV2(*video.Transcript, video.ID),
		// RunConceptDistillationTest(video),
	}

	suite := VideoTestSuite{
		VideoID: video.ID,
		Title:   video.Title,
		Results: results,
	}

	return ExportVideoTestResults(suite)
}



