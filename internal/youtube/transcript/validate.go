package transcript

import (
	"fmt"
)

func ValidateProcessingRequest(req *ProcessingRequest) error {
	if req == nil {
		return nil
	}

	if req.Output == "" {
		return fmt.Errorf("output format is required")
	}

	switch req.Output {
	case "timeline-text", "plain-text", "segments":
		// valid
	default:
		return fmt.Errorf("unsupported output format: %s", req.Output)
	}

	if req.WindowSize < 0 {
		return fmt.Errorf("window size cannot be negative")
	}

	return nil
}
