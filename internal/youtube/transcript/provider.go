package transcript

import (
	"context"
)

type Provider interface {
	GetTranscript(
        ctx context.Context,
        videoID string,
    ) (*Transcript, error)
}