package resource

import (
	"context"

	"github.com/Arvind215271/askito/internal/youtube"
	"github.com/Arvind215271/askito/internal/youtube/pipeline"
)

// ProcessVideo handles processing a single video URL/ID using the pipeline.
func ProcessVideo(ctx context.Context, pipelineService *pipeline.Service, videoID string, req *pipeline.Request) *youtube.Video {
	if pipelineService == nil {
		return &youtube.Video{
			ID:     videoID,
			Errors: []youtube.Error{{Message: "pipeline service is nil"}},
		}
	}
	return pipelineService.ProcessFaultTolerant(ctx, videoID, req)
}
