package subtitle

import "fmt"


func (r *DownloadRequest) Validate() error {
	if r.VideoID == "" {
		return fmt.Errorf("video ID is required")
	}

	if r.Language == "" {
		return fmt.Errorf("subtitle language is required")
	}

	switch r.Type {
	case "manual", "automatic":
	default:
		return fmt.Errorf("invalid subtitle type: %s", r.Type)
	}

	if r.Format == "" {
		r.Format = "json3"
	}

	switch r.Format {
	case "json3", "vtt":
	default:
		return fmt.Errorf("unsupported subtitle format: %s", r.Format)
	}

	return nil
}