package signal

import "fmt"

func (r *SignalRequest) Validate() error {
	if r.VideoID == "" {
		return fmt.Errorf("video ID is required")
	}

	switch r.Analysis {
	case "word-stats", "window-stats":
	default:
		return fmt.Errorf("invalid analysis type: %s", r.Analysis)
	}

	if r.MinFreq < 0 {
		return fmt.Errorf("min_freq must be non-negative")
	}

	if r.Depth < 0 || r.Depth > 1 {
		return fmt.Errorf("depth must be between 0 and 1")
	}

	if r.WindowSize <= 0 {
		r.WindowSize = 300 // default
	}

	return nil
}
