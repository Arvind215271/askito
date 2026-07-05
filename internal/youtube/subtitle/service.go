package subtitle

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

type SubtitleService struct {
}



func NewSubtitleService() *SubtitleService {
	return &SubtitleService{}
}

func (s *SubtitleService) DownloadSubtitle(ctx context.Context, req DownloadRequest, meta SubtitleMetadata) (*SubtitleResult, error) {
	if req.Format == "" {
		req.Format = "json3"
	}

	// Validate type and check if language exists
	var found bool
	if req.Type == "manual" {
		for _, t := range meta.Manual {
			if t.LanguageCode == req.Language {
				found = true
				break
			}
		}
	} else if req.Type == "automatic" {
		for _, t := range meta.Automatic {
			if t.LanguageCode == req.Language {
				found = true
				break
			}
		}
	} else {
		return nil, fmt.Errorf("invalid subtitle type: %s", req.Type)
	}

	if !found {
		return nil, fmt.Errorf("subtitle track not found for type %s and language %s", req.Type, req.Language)
	}

	// Download the track using yt-dlp
	cmd := exec.CommandContext(ctx, "yt-dlp", "--skip-download", "--write-sub", "--sub-lang", req.Language, "--sub-format", req.Format, "-o", "-", req.VideoID)
	
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to download subtitle: %w", err)
	}

	return &SubtitleResult{
		Content: stdout.Bytes(),
		Format:  req.Format,
		Language: req.Language,
	}, nil
}
