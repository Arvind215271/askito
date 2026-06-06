// ./internal/youtube/input/input_model.go
package youtube

type InputType string

const (
	InputTypeVideo    InputType = "video"
	InputTypePlaylist InputType = "playlist"
)

type YouTubeInput struct {
	InputType    InputType
	ID           string
	OriginalURL  string
	NormalizedURL string
}