package playlist

type OutputType string

const (
	OutputTypeID   OutputType = "id"
	OutputTypeURL  OutputType = "url"
	OutputTypeBoth OutputType = "both"
)

type PlaylistVideosRequest struct {
	URL      string     `json:"url" validate:"required"`
	Provider string     `json:"provider" validate:"omitempty"`
	Output   OutputType `json:"output" validate:"omitempty"`
}

type VideoEntry struct {
	ID       string `json:"id,omitempty"`
	URL      string `json:"url,omitempty"`
	Position int    `json:"position"`
	AddedAt  string `json:"added_at,omitempty"`
}

type PlaylistVideosResponse struct {
	PlaylistID string       `json:"playlist_id"`
	Total      int          `json:"total"`
	Videos     []VideoEntry `json:"videos"`
}
