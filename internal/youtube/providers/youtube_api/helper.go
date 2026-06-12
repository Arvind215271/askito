package youtubeapi

import (
	"time"
	"strconv"
	yt "google.golang.org/api/youtube/v3"
)

func parseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

func parseUint(s string) uint64 {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}
	return v
}

func getPlaylistThumbnail(
	playlist *yt.Playlist,
) string {

	if playlist.Snippet == nil ||
		playlist.Snippet.Thumbnails == nil {
		return ""
	}

	if playlist.Snippet.Thumbnails.High != nil {
		return playlist.Snippet.Thumbnails.High.Url
	}

	if playlist.Snippet.Thumbnails.Medium != nil {
		return playlist.Snippet.Thumbnails.Medium.Url
	}

	if playlist.Snippet.Thumbnails.Default != nil {
		return playlist.Snippet.Thumbnails.Default.Url
	}

	return ""
}

func getVideoThumbnail(
	video *yt.Video,
) string {

	// there are different resolutoin type per video. Simply return whatever is present... Not much needed to be true. We can simply keep those field as empty.
	if video.Snippet == nil ||
		video.Snippet.Thumbnails == nil {
		return ""
	}

	if video.Snippet.Thumbnails.High != nil {
		return video.Snippet.Thumbnails.High.Url
	}

	if video.Snippet.Thumbnails.Medium != nil {
		return video.Snippet.Thumbnails.Medium.Url
	}

	if video.Snippet.Thumbnails.Default != nil {
		return video.Snippet.Thumbnails.Default.Url
	}

	return ""
}