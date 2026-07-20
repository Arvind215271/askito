package video

type VideoRequest struct {
	URL      string `json:"url" validate:"required,url"`
	Provider string `json:"provider" validate:"omitempty"`
}

type VideoByIDRequest struct {
	ID       string `json:"id" validate:"required"`
	Provider string `json:"provider" validate:"omitempty"`
}

type VideoResponse struct {
	Video interface{} `json:"video"`
}
