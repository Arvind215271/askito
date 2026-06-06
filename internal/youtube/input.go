// ./internal/youtube/input.go
package youtube

import (
	"net/url"
	"strings"
)

type Input struct{}

var Inputs = Input{}

// ParsePlaylistURL validates a playlist URL,
// extracts the playlist ID,
// and returns a normalized playlist input.
func (Input) ParsePlaylistURL(raw string) (*PlaylistInput, error) {
	// get the raw URL and  trim it... because there might be some extra sspace around that would create problem if not in proper form 
	raw = strings.TrimSpace(raw)

	if raw == "" {
		return nil, Err.Playlist.InvalidURL()
	}

	u, err := url.Parse(raw)
	// here it will parse the url into  different object like way. Like it will have stuff like u.query, u.params,etc, that any url contain. So then we can easily extract infromatoin from it.
	if err != nil {
		return nil, Err.Playlist.InvalidURL().Wrap(err)
	}

	if !isYouTubeHost(u.Host) {
		return nil, Err.Playlist.InvalidDomain()
	}

	// youtube playlist have the query params of ?list=plylist_id that we need to to find that 
	playlistID := u.Query().Get("list")
	if playlistID == "" {
		return nil, Err.Playlist.MissingID()
	}

	// here, we are using a normalissed url because youtube have a lot of domain based on coutnry, type like youtube kids, music, etc. 
	return &PlaylistInput{
		ID:            playlistID,
		OriginalURL:   raw,
		NormalizedURL: normalizePlaylistURL(playlistID),
	}, nil
}

// isYouTubeHost validates supported YouTube hosts.
func isYouTubeHost(host string) bool {
	host = strings.ToLower(host)

	host = strings.TrimPrefix(host, "www.")
	host = strings.TrimPrefix(host, "m.")

	return host == "youtube.com" ||
		strings.HasPrefix(host, "youtube.")
}

// normalizePlaylistURL converts a playlist ID into a canonical URL.
func normalizePlaylistURL(id string) string {
	return "https://www.youtube.com/playlist?list=" + id
}