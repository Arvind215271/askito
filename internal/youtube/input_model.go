// ./internal/youtube/input/input_model.go
package youtube

type PlaylistInput struct {
	ID            string
	OriginalURL   string
	NormalizedURL string
}

type VideoInput struct {
	ID            string
	OriginalURL   string
	NormalizedURL string
}