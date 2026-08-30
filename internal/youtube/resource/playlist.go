package resource

import (
	"context"
	"fmt"
	"sync"

	"github.com/Arvind215271/askito/internal/youtube"
	"github.com/Arvind215271/askito/internal/youtube/metadata"
	"github.com/Arvind215271/askito/internal/youtube/pipeline"
)

// ProcessPlaylist fetches playlist items, decides whether the pipeline is needed based on execution plan,
// maps playlist items to videos, processes them concurrently, and attaches them back to the Playlist model.
func ProcessPlaylist(
	ctx context.Context,
	metadataService *metadata.Service,
	pipelineService *pipeline.Service,
	playlistID string,
	req *pipeline.Request,
	concurrency int,
) *youtube.Playlist {
	playlist := &youtube.Playlist{
		ID: playlistID,
	}

	if metadataService == nil {
		playlist.Errors = append(playlist.Errors, youtube.Error{Message: "metadata service is nil"})
		return playlist
	}

	// 1. Fetch playlist items
	items, err := metadataService.GetPlaylistItems(ctx, playlistID, metadata.ProviderYTDLP)
	if err != nil {
		items, err = metadataService.GetPlaylistItems(ctx, playlistID, metadata.ProviderAPI)
		if err != nil {
			playlist.Errors = append(playlist.Errors, youtube.Error{Message: fmt.Sprintf("failed to fetch playlist items: %v", err)})
			return playlist
		}
	}

	playlist.Items = items
	playlist.ItemCount = len(items)

	if len(items) == 0 {
		return playlist
	}

	// Determine if pipeline processing is needed
	needsPipeline := true
	if req != nil && req.ExecutionPlan != nil {
		plan := req.ExecutionPlan
		if !plan.NeedsMetadata() && !plan.NeedsDescription() && !plan.NeedsSubtitle() && !plan.NeedsTranscript() && !plan.NeedsSignal() {
			needsPipeline = false
		}
	}

	playlistVideos := make([]youtube.PlaylistVideo, len(items))

	if !needsPipeline {
		for i, item := range items {
			video := ItemToVideo(item)
			playlistVideos[i] = youtube.PlaylistVideo{
				Video:    video,
				Position: item.Position,
				AddedAt:  item.AddedAt,
			}
		}
		playlist.Videos = playlistVideos
		return playlist
	}

	if concurrency <= 0 {
		concurrency = 1
	}

	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for i, item := range items {
		wg.Add(1)
		go func(index int, pi youtube.PlaylistItem) {
			defer wg.Done()

			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				baseVideo := ItemToVideo(pi)
				baseVideo.Errors = append(baseVideo.Errors, youtube.Error{Message: "context cancelled"})
				playlistVideos[index] = youtube.PlaylistVideo{
					Video:    baseVideo,
					Position: pi.Position,
					AddedAt:  pi.AddedAt,
				}
				return
			}
			defer func() { <-sem }()

			baseVideo := ItemToVideo(pi)
			var processedVideo *youtube.Video
			if pipelineService != nil {
				rc := &pipeline.ResourceContext{
					Video:   &baseVideo,
					Request: req,
					Plan:    req.ExecutionPlan,
				}
				processedVideo = pipelineService.ProcessResource(ctx, rc, true)
			} else {
				processedVideo = &baseVideo
			}

			if processedVideo == nil {
				processedVideo = &baseVideo
			}

			playlistVideos[index] = youtube.PlaylistVideo{
				Video:    *processedVideo,
				Position: pi.Position,
				AddedAt:  pi.AddedAt,
			}
		}(i, item)
	}

	wg.Wait()
	playlist.Videos = playlistVideos
	return playlist
}
