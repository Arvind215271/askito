package ytdlp

import (
	"fmt"
	"time"

	"github.com/Arvind215271/askito/internal/youtube"
	"github.com/Arvind215271/askito/internal/youtube/subtitle"
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
		SubtitleMetadata:  mapSubtitleMetadata(meta),
	}
}

func mapSubtitleMetadata(meta YTOutput) subtitle.SubtitleMetadata {
	return subtitle.SubtitleMetadata{
		Manual:    mapTracks(meta.Subtitles),
		Automatic: mapTracks(meta.AutomaticCaptions),
	}
}

func mapTracks(raw map[string][]SubtitleFormat) []subtitle.SubtitleTrack {
	var tracks []subtitle.SubtitleTrack
	for lang, formats := range raw {
		track := subtitle.SubtitleTrack{
			LanguageCode: lang,
			LanguageName: lang, // Could be enhanced with a lookup
			Formats:      []string{},
		}
		seen := make(map[string]bool)
		for _, f := range formats {
			if !seen[f.Ext] {
				track.Formats = append(track.Formats, f.Ext)
				seen[f.Ext] = true
			}
		}
		tracks = append(tracks, track)
	}
	return tracks
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
