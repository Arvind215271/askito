package yt_bench

import (
	"context"
	"os/exec"
)

// FetchFullVideoMetadata: Function 1 & 3 - Get full JSON data for a video
func FetchFullVideoMetadata(ctx context.Context, videoID string) ([]byte, error) {
	url := "https://www.youtube.com/watch?v=" + videoID
	return exec.CommandContext(ctx, "yt-dlp", "--skip-download", "--dump-single-json", "--no-cache-dir", url).CombinedOutput()
}

// FetchPlaylistFlatMetadata: Function 2 - Get only basic info (ID, title) for all videos in playlist
func FetchPlaylistFlatMetadata(ctx context.Context, playlistID string) ([]byte, error) {
	url := "https://www.youtube.com/playlist?list=" + playlistID
	return exec.CommandContext(ctx, "yt-dlp", "-j", "--flat-playlist", "--no-cache-dir", url).CombinedOutput()
}

// FetchPlaylistFullMetadata: Alternative for Function 3 - Get full JSON data for all videos in playlist
func FetchPlaylistFullMetadata(ctx context.Context, playlistID string) ([]byte, error) {
	url := "https://www.youtube.com/playlist?list=" + playlistID
	return exec.CommandContext(ctx, "yt-dlp", "--skip-download", "--dump-single-json", "--no-cache-dir", url).CombinedOutput()
}
