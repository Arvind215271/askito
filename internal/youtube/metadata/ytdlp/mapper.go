package ytdlp

import "github.com/Arvind215271/askito/internal/youtube"

func MapVideo(meta YTOutput) youtube.Video {
	return youtube.Video{
		ID:              meta.ID,
		Title:           meta.Title,
		Description:     meta.Description,
		DurationSeconds: int64(meta.Duration),
		ViewCount:       meta.ViewCount,
		LikeCount:       meta.LikeCount,
		CommentCount:    meta.CommentCount,
		ChannelTitle:    meta.Channel,
		ChannelID:       meta.ChannelID,
		ThumbnailURL:    meta.Thumbnail,
		Tags:            meta.Tags,
		CategoryID:      "", // Need to map categories if available
		PrivacyStatus:   meta.Availability,
	}
}

func MapPlaylist(meta YTPlaylistOutput) youtube.Playlist {
	return youtube.Playlist{
		ID:           meta.ID,
		Title:        meta.Title,
		Description:  meta.Description,
		ChannelTitle: meta.Channel,
		ChannelID:    meta.ChannelID,
		ThumbnailURL: meta.Thumbnail,
	}
}
