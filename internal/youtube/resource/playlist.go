package resource

import (
	"context"
	"fmt"
	"sync"

	"github.com/Arvind215271/askito/internal/youtube"
	"github.com/Arvind215271/askito/internal/youtube/metadata"
	"github.com/Arvind215271/askito/internal/youtube/pipeline"
	"github.com/Arvind215271/askito/internal/youtube/stats"
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
	playlist, _ := ProcessPlaylistWithStats(ctx, metadataService, pipelineService, playlistID, req, concurrency)
	return playlist
}

// ProcessPlaylistWithStats fetches playlist items and processes them, returning playlist and resource stats.
func ProcessPlaylistWithStats(
	ctx context.Context,
	metadataService *metadata.Service,
	pipelineService *pipeline.Service,
	playlistID string,
	req *pipeline.Request,
	concurrency int,
) (*youtube.Playlist, stats.ResourceStats) {
	playlist := &youtube.Playlist{
		ID: playlistID,
	}
	var resourceStats stats.ResourceStats

	if metadataService == nil {
		playlist.Errors = append(playlist.Errors, youtube.Error{Message: "metadata service is nil"})
		return playlist, resourceStats
	}

	// 1. Fetch playlist items with stats
	items, metaStats, err := metadataService.GetPlaylistItemsWithStats(ctx, playlistID, metadata.ProviderYTDLP)
	resourceStats.Metadata.Add(metaStats)
	if err != nil {
		items, metaStats, err = metadataService.GetPlaylistItemsWithStats(ctx, playlistID, metadata.ProviderAPI)
		resourceStats.Metadata.Add(metaStats)
		if err != nil {
			playlist.Errors = append(playlist.Errors, youtube.Error{Message: fmt.Sprintf("failed to fetch playlist items: %v", err)})
			return playlist, resourceStats
		}
	}

	playlist.Items = items
	playlist.ItemCount = len(items)

	if len(items) == 0 {
		return playlist, resourceStats
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
		resourceStats.VideosProcessed += len(items)
		return playlist, resourceStats
	}

	if concurrency <= 0 {
		concurrency = 1
	}

	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	var mu sync.Mutex

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
				mu.Lock()
				resourceStats.VideosProcessed++
				resourceStats.VideoFailures++
				mu.Unlock()
				return
			}
			defer func() { <-sem }()

			baseVideo := ItemToVideo(pi)
			var processedVideo *youtube.Video
			var videoResStats stats.ResourceStats
			if pipelineService != nil {
				rc := &pipeline.ResourceContext{
					Video:   &baseVideo,
					Request: req,
					Plan:    req.ExecutionPlan,
				}
				pVideo, pStats := pipelineService.ProcessResourceWithStats(ctx, rc, true)
				processedVideo = pVideo
				videoResStats.VideosProcessed++
				videoResStats.Metadata.Add(pStats.Metadata)
				videoResStats.Subtitle.Add(pStats.Subtitle)
				if pVideo != nil && len(pVideo.Errors) > 0 {
					videoResStats.VideoFailures++
				}
			} else {
				processedVideo = &baseVideo
				videoResStats.VideosProcessed++
			}

			if processedVideo == nil {
				processedVideo = &baseVideo
			}

			playlistVideos[index] = youtube.PlaylistVideo{
				Video:    *processedVideo,
				Position: pi.Position,
				AddedAt:  pi.AddedAt,
			}

			mu.Lock()
			resourceStats.Add(videoResStats)
			mu.Unlock()
		}(i, item)
	}

	wg.Wait()
	playlist.Videos = playlistVideos
	return playlist, resourceStats
}
