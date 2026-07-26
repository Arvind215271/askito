package resource

import (
	"context"
	"sync"

	"github.com/Arvind215271/askito/internal/logger"
	"github.com/Arvind215271/askito/internal/youtube"
	"github.com/Arvind215271/askito/internal/youtube/input"
	"github.com/Arvind215271/askito/internal/youtube/metadata"
	"github.com/Arvind215271/askito/internal/youtube/pipeline"
)

type Service struct {
	metadataService *metadata.Service
	pipelineService *pipeline.Service
	logger          *logger.Logger
	concurrency     int
}

func NewService(
	metadataService *metadata.Service,
	pipelineService *pipeline.Service,
	logger *logger.Logger,
	concurrency int,
) *Service {
	if concurrency <= 0 {
		concurrency = 1
	}
	return &Service{
		metadataService: metadataService,
		pipelineService: pipelineService,
		logger:          logger,
		concurrency:     concurrency,
	}
}

// ProcessResources processes a slice of YouTube inputs (videos and playlists) while strictly
// preserving the input request order, dispatching to appropriate video or playlist handlers,
// and returning a uniform slice of youtube.Resource containers.
func (s *Service) ProcessResources(
	ctx context.Context,
	inputs []youtubeurl.YouTubeInput,
	req *pipeline.Request,
) []youtube.Resource {
	if len(inputs) == 0 {
		return nil
	}

	resources := make([]youtube.Resource, len(inputs))
	concurrency := s.concurrency
	if concurrency <= 0 {
		concurrency = 1
	}

	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for i, inputItem := range inputs {
		wg.Add(1)
		go func(index int, inp youtubeurl.YouTubeInput) {
			defer wg.Done()

			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				resources[index] = youtube.Resource{
					ID:   inp.ID,
					Type: youtube.ResourceTypeVideo,
					Video: &youtube.Video{
						ID:     inp.ID,
						Errors: []youtube.Error{{Message: "context cancelled"}},
					},
				}
				return
			}
			defer func() { <-sem }()

			switch inp.InputType {
			case youtubeurl.InputTypePlaylist:
				playlist := ProcessPlaylist(
					ctx,
					s.metadataService,
					s.pipelineService,
					inp.ID,
					req,
					concurrency,
				)
				resources[index] = youtube.Resource{
					ID:       inp.ID,
					Type:     youtube.ResourceTypePlaylist,
					Playlist: playlist,
				}
			case youtubeurl.InputTypeVideo:
				fallthrough
			default:
				video := ProcessVideo(
					ctx,
					s.pipelineService,
					inp.ID,
					req,
				)
				if video == nil {
					video = &youtube.Video{ID: inp.ID}
				}
				resources[index] = youtube.Resource{
					ID:    inp.ID,
					Type:  youtube.ResourceTypeVideo,
					Video: video,
				}
			}
		}(i, inputItem)
	}

	wg.Wait()
	return resources
}
