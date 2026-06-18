package ytdlp

import (
	"github.com/Arvind215271/askito/internal/youtube/transcript"

	"context"
)

type Provider struct {
	client *Client
}

func NewProvider(
	client *Client,
) *Provider {
	return &Provider{
		client: client,
	}
}

func (p *Provider) GetTranscript(
	ctx context.Context,
	videoID string,
) (*transcript.Transcript, error) {
	return p.client.GetTranscript(ctx, videoID)

}