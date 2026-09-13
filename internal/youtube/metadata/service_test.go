package metadata_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Arvind215271/askito/internal/youtube"
	"github.com/Arvind215271/askito/internal/youtube/metadata"
	"github.com/Arvind215271/askito/internal/youtube/stats"
)

type mockProvider struct {
	getVideoFunc         func(ctx context.Context, videoID string, st *stats.MetadataStats) (youtube.Video, error)
	getPlaylistItemsFunc func(ctx context.Context, playlistID string, st *stats.MetadataStats) ([]youtube.PlaylistItem, error)
	getPlaylistMetaFunc  func(ctx context.Context, playlistID string, st *stats.MetadataStats) (youtube.Playlist, error)
}

func (m *mockProvider) GetVideo(ctx context.Context, videoID string, st *stats.MetadataStats) (youtube.Video, error) {
	if m.getVideoFunc != nil {
		return m.getVideoFunc(ctx, videoID, st)
	}
	return youtube.Video{}, nil
}

func (m *mockProvider) GetPlaylistItems(ctx context.Context, playlistID string, st *stats.MetadataStats) ([]youtube.PlaylistItem, error) {
	if m.getPlaylistItemsFunc != nil {
		return m.getPlaylistItemsFunc(ctx, playlistID, st)
	}
	return nil, nil
}

func (m *mockProvider) GetPlaylistMetadata(ctx context.Context, playlistID string, st *stats.MetadataStats) (youtube.Playlist, error) {
	if m.getPlaylistMetaFunc != nil {
		return m.getPlaylistMetaFunc(ctx, playlistID, st)
	}
	return youtube.Playlist{}, nil
}

func TestMetadataService_GetVideoWithStats_Success(t *testing.T) {
	ytdlp := &mockProvider{
		getVideoFunc: func(ctx context.Context, videoID string, st *stats.MetadataStats) (youtube.Video, error) {
			if st != nil {
				st.YTDLPRequests++
				st.CacheHits++
			}
			return youtube.Video{ID: videoID, Title: "Test Video"}, nil
		},
	}

	svc := metadata.NewService(nil, ytdlp)
	video, st, err := svc.GetVideoWithStats(context.Background(), "vid1", metadata.ProviderYTDLP)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if video.ID != "vid1" {
		t.Errorf("expected video ID vid1, got %s", video.ID)
	}
	if st.Attempts != 1 || st.Successes != 1 || st.CacheHits != 1 || st.YTDLPRequests != 1 {
		t.Errorf("unexpected stats: %+v", st)
	}
}

func TestMetadataService_GetVideoWithStats_Fallback(t *testing.T) {
	ytdlp := &mockProvider{
		getVideoFunc: func(ctx context.Context, videoID string, st *stats.MetadataStats) (youtube.Video, error) {
			if st != nil {
				st.YTDLPRequests++
				st.CacheMisses++
				st.UpstreamFetches++
			}
			return youtube.Video{}, errors.New("ytdlp failed")
		},
	}
	api := &mockProvider{
		getVideoFunc: func(ctx context.Context, videoID string, st *stats.MetadataStats) (youtube.Video, error) {
			if st != nil {
				st.APIRequests++
				st.CacheMisses++
				st.UpstreamFetches++
			}
			return youtube.Video{ID: videoID, Title: "API Video"}, nil
		},
	}

	svc := metadata.NewService(api, ytdlp)
	video, st, err := svc.GetVideoWithStats(context.Background(), "vid1", metadata.ProviderYTDLP)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if video.Title != "API Video" {
		t.Errorf("expected API Video title, got %s", video.Title)
	}
	if st.Attempts != 1 || st.Successes != 1 || st.Fallbacks != 1 || st.YTDLPRequests != 1 || st.APIRequests != 1 {
		t.Errorf("unexpected stats: %+v", st)
	}
}

func TestMetadataService_GetVideoWithStats_Failure(t *testing.T) {
	ytdlp := &mockProvider{
		getVideoFunc: func(ctx context.Context, videoID string, st *stats.MetadataStats) (youtube.Video, error) {
			if st != nil {
				st.YTDLPRequests++
			}
			return youtube.Video{}, errors.New("total failure")
		},
	}

	svc := metadata.NewService(nil, ytdlp)
	_, st, err := svc.GetVideoWithStats(context.Background(), "vid1", metadata.ProviderYTDLP)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if st.Attempts != 1 || st.Failures != 1 || st.Successes != 0 {
		t.Errorf("unexpected stats: %+v", st)
	}
}
