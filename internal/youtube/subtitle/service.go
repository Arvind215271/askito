package subtitle

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type SubtitleService struct{}

func NewSubtitleService() *SubtitleService {
	return &SubtitleService{}
}

func (s *SubtitleService) DownloadSubtitle(
	ctx context.Context,
	req DownloadRequest,
	meta SubtitleMetadata,
) (*SubtitleResult, error) {

	if req.Format == "" {
		req.Format = "json3"
	}

	// Validate subtitle type and language.
	var found bool

	switch req.Type {
	case "manual":
		for _, t := range meta.Manual {
			if t.LanguageCode == req.Language {
				found = true
				break
			}
		}

	case "automatic":
		for _, t := range meta.Automatic {
			if t.LanguageCode == req.Language {
				found = true
				break
			}
		}

	default:
		return nil, fmt.Errorf("invalid subtitle type: %s", req.Type)
	}

	if !found {
		return nil, fmt.Errorf(
			"subtitle track not found for type %s and language %s",
			req.Type,
			req.Language,
		)
	}

	// Create a temporary directory for yt-dlp output.
	tempDir, err := os.MkdirTemp("", "askito-subtitles-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	outputTemplate := filepath.Join(tempDir, "%(id)s")

	args := []string{
		"--skip-download",
		"--sub-lang", req.Language,
		"--sub-format", req.Format,
		"-o", outputTemplate,
	}

	switch req.Type {
	case "manual":
		args = append(args, "--write-sub")
	case "automatic":
		args = append(args, "--write-auto-sub")
	}

	args = append(args, req.VideoID)

	cmd := exec.CommandContext(ctx, "yt-dlp", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to download subtitle: %w\n%s",
			err,
			string(output),
		)
	}

	// Find the downloaded subtitle file.
	pattern := filepath.Join(
		tempDir,
		fmt.Sprintf("*.%s.%s", req.Language, req.Format),
	)

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to locate subtitle file: %w", err)
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("subtitle file was not created")
	}

	if len(matches) > 1 {
		return nil, fmt.Errorf("multiple subtitle files found")
	}

	data, err := os.ReadFile(matches[0])
	if err != nil {
		return nil, fmt.Errorf("failed to read subtitle file: %w", err)
	}

	return &SubtitleResult{
		Content:  data,
		Format:   req.Format,
		Language: req.Language,
	}, nil
}

