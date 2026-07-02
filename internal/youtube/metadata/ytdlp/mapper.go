package ytdlp

import (
	"fmt"
	"time"

	"github.com/Arvind215271/askito/internal/youtube"
)

func MapVideo(meta YTOutput) youtube.Video {
	durationSeconds := int64(meta.Duration)
	duration := time.Duration(durationSeconds) * time.Second
	
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	
	durationTimestamp := fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	if hours == 0 {
		durationTimestamp = fmt.Sprintf("%02d:%02d", minutes, seconds)
	}

	return youtube.Video{
		ID:                meta.ID,
		Title:             meta.Title,
		Description:       meta.Description,
		DurationSeconds:   durationSeconds,
		DurationMinutes:   duration.Minutes(),
		DurationTimestamp: durationTimestamp,
		Duration:          durationTimestamp,
		ViewCount:         meta.ViewCount,
		LikeCount:         meta.LikeCount,
		CommentCount:      meta.CommentCount,
		ChannelTitle:      meta.Channel,
		ChannelID:         meta.ChannelID,
		ThumbnailURL:      meta.Thumbnail,
		Tags:              meta.Tags,
		CategoryID:        "", // Need to map categories if available
		PrivacyStatus:     meta.Availability,
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
