package youtubeapi

import (
	yt "google.golang.org/api/youtube/v3"
	youtube "github.com/Arvind215271/askito/internal/youtube"
)

func MapVideo(
	video *yt.Video,
) youtube.Video {

	return youtube.Video{
		ID:          video.Id,
		Title:       video.Snippet.Title,
		Description: video.Snippet.Description,

		ChannelID:    video.Snippet.ChannelId,
		ChannelTitle: video.Snippet.ChannelTitle,

		ThumbnailURL: getVideoThumbnail(video),

		PublishedAt: parseTime(video.Snippet.PublishedAt),

		Duration: video.ContentDetails.Duration,

		ViewCount:    video.Statistics.ViewCount,
		LikeCount:    video.Statistics.LikeCount,
		CommentCount: video.Statistics.CommentCount,

		Tags: video.Snippet.Tags,

		CategoryID: video.Snippet.CategoryId,

		CaptionAvailable: video.ContentDetails.Caption == "true",

		PrivacyStatus:       video.Status.PrivacyStatus,
		LiveBroadcastStatus: video.Snippet.LiveBroadcastContent,
	}
}

// maps the fetched Youtube API video type to domain logic video type
func MapPlaylistVideos(
	items []*yt.PlaylistItem,
	videos []*yt.Video,
) []youtube.PlaylistVideo {

	// create a map of yt.video and total video we have/
	videoMap := make(map[string]*yt.Video, len(videos))

	// for each video. map that in map with its video ID.
	// why? Because each video in Playlist Video have a position associated to it. Thus, we would need to map the video to it. Thus, it would be far easier to do so.  
	for _, video := range videos {
		videoMap[video.Id] = video
	}

	result := make([]youtube.PlaylistVideo, 0, len(items))

	for _, item := range items {

		if item == nil ||
			item.Snippet == nil ||
			item.ContentDetails == nil {
			continue
		}

		video := videoMap[item.ContentDetails.VideoId]

		if video == nil ||
			video.Snippet == nil ||
			video.ContentDetails == nil ||
			video.Statistics == nil ||
			video.Status == nil {
			continue
		}

		result = append(result, youtube.PlaylistVideo{
			Video: MapVideo(video),

			Position: int(item.Snippet.Position),
			AddedAt:  parseTime(item.Snippet.PublishedAt),
		})
	}

	return result
}