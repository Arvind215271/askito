package subtitle

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Arvind215271/askito/internal/cache"
	"github.com/Arvind215271/askito/internal/logger"
)

type SubtitleService struct {
	cache  *cache.Manager
	logger *logger.Logger
}

func NewSubtitleService(cache *cache.Manager, logger *logger.Logger) *SubtitleService {
	return &SubtitleService{
		cache:  cache,
		logger: logger,
	}
}

func (s *SubtitleService) DownloadSubtitle(
	ctx context.Context,
	req DownloadRequest,
	meta SubtitleMetadata,
) (*SubtitleResult, error) {

	req, err := validateRequest(req)
	if err != nil {
		return nil, err
	}

	if err := validateTrack(req, meta); err != nil {
		return nil, err
	}

	cacheKey := subtitleCacheKey(req)

	if cached, err := s.cache.Get(req.VideoID, cacheKey); err == nil {
		s.logger.Debug(
			"subtitle cache hit",
			"videoID", req.VideoID,
			"language", req.Language,
		)

		return &SubtitleResult{
			Content:  cached,
			Format:   req.Format,
			Language: req.Language,
		}, nil
	}

	data, err := s.download(ctx, req)
	if err != nil {
		return nil, err
	}

	if err := s.cache.Save(req.VideoID, cacheKey, data); err != nil {
		s.logger.Warn(
			"failed to cache subtitle",
			"videoID", req.VideoID,
			"error", err,
		)
	}

	return &SubtitleResult{
		Content:  data,
		Format:   req.Format,
		Language: req.Language,
	}, nil
}

func validateRequest(req DownloadRequest) (DownloadRequest, error) {

	if req.VideoID == "" {
		return req, fmt.Errorf("video ID is required")
	}

	if req.Language == "" {
		return req, fmt.Errorf("subtitle language is required")
	}

	switch req.Type {
	case "manual", "automatic":
	default:
		return req, fmt.Errorf(
			"invalid subtitle type: %s",
			req.Type,
		)
	}

	if req.Format == "" {
		req.Format = "json3"
	}

	switch req.Format {
	case "json3", "vtt":
	default:
		return req, fmt.Errorf(
			"unsupported subtitle format: %s",
			req.Format,
		)
	}

	return req, nil
}

func validateTrack(
	req DownloadRequest,
	meta SubtitleMetadata,
) error {

	var tracks []SubtitleTrack

	switch req.Type {
	case "manual":
		tracks = meta.Manual

	case "automatic":
		tracks = meta.Automatic
	}

	for _, track := range tracks {
		if track.LanguageCode != req.Language {
			continue
		}

		if !formatSupported(track, req.Format) {
			return fmt.Errorf(
				"subtitle format %s is not available for language %s",
				req.Format,
				req.Language,
			)
		}

		return nil
	}

	return fmt.Errorf(
		"subtitle track not found for type %s and language %s",
		req.Type,
		req.Language,
	)
}

func formatSupported(
	track SubtitleTrack,
	format string,
) bool {

	if len(track.Formats) == 0 {
		return true
	}

	for _, supported := range track.Formats {
		if supported == format {
			return true
		}
	}

	return false
}

func subtitleCacheKey(req DownloadRequest) string {
	return fmt.Sprintf(
		"subtitle.%s.%s",
		req.Language,
		req.Format,
	)
}

func (s *SubtitleService) download(
	ctx context.Context,
	req DownloadRequest,
) ([]byte, error) {

	tempDir, err := os.MkdirTemp(
		"",
		"askito-subtitles-*",
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create subtitle temp directory: %w",
			err,
		)
	}

	defer os.RemoveAll(tempDir)

	outputTemplate := filepath.Join(
		tempDir,
		"%(id)s",
	)

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

	cmd := exec.CommandContext(
		ctx,
		"yt-dlp",
		args...,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to download subtitle: %w\n%s",
			err,
			string(output),
		)
	}

	return readDownloadedSubtitle(
		tempDir,
		req,
	)
}

func readDownloadedSubtitle(
	tempDir string,
	req DownloadRequest,
) ([]byte, error) {

	pattern := filepath.Join(
		tempDir,
		fmt.Sprintf(
			"*.%s.%s",
			req.Language,
			req.Format,
		),
	)

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to locate subtitle file: %w",
			err,
		)
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf(
			"subtitle file was not created",
		)
	}

	if len(matches) > 1 {
		return nil, fmt.Errorf(
			"multiple subtitle files found",
		)
	}

	data, err := os.ReadFile(matches[0])
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read subtitle file: %w",
			err,
		)
	}

	return data, nil
}
