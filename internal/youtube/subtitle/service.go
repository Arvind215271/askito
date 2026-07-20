package subtitle

import (
	"context"
	"fmt"

	"github.com/Arvind215271/askito/internal/cache"
	"github.com/Arvind215271/askito/internal/logger"
	"github.com/Arvind215271/askito/internal/youtube/metadata/ytdlp/python"
)

type SubtitleService struct {
	cache  *cache.Manager
	logger *logger.Logger
	pool   *python.SinglePool
}

func NewSubtitleService(
	cache *cache.Manager,
	logger *logger.Logger,
	pool *python.SinglePool,
) *SubtitleService {
	return &SubtitleService{
		cache:  cache,
		logger: logger,
		pool:   pool,
	}
}

func (s *SubtitleService) DownloadSubtitle(
	ctx context.Context,
	req DownloadRequest,
	meta SubtitleMetadata,
) (*SubtitleResult, error) {

	s.logger.Debug(
		"subtitle download requested",
		"videoID", req.VideoID,
		"language", req.Language,
		"type", req.Type,
		"format", req.Format,
	)

	req, err := validateRequest(req)
	if err != nil {
		s.logger.Warn(
			"subtitle request validation failed",
			"videoID", req.VideoID,
			"language", req.Language,
			"type", req.Type,
			"format", req.Format,
			"error", err,
		)

		return nil, err
	}

	s.logger.Debug(
		"subtitle request validated",
		"videoID", req.VideoID,
		"language", req.Language,
		"type", req.Type,
		"format", req.Format,
	)

	if err := validateTrack(req, meta); err != nil {
		s.logger.Warn(
			"subtitle track validation failed",
			"videoID", req.VideoID,
			"language", req.Language,
			"type", req.Type,
			"format", req.Format,
			"error", err,
		)

		return nil, err
	}

	s.logger.Debug(
		"subtitle track validated",
		"videoID", req.VideoID,
		"language", req.Language,
		"type", req.Type,
		"format", req.Format,
	)

	cacheKey := subtitleCacheKey(req)

	s.logger.Debug(
		"checking subtitle cache",
		"videoID", req.VideoID,
		"cacheKey", cacheKey,
	)

	if cached, err := s.cache.Get(req.VideoID, cacheKey); err == nil {

		s.logger.Debug(
			"subtitle cache hit",
			"videoID", req.VideoID,
			"language", req.Language,
			"type", req.Type,
			"format", req.Format,
			"cacheKey", cacheKey,
		)

		return &SubtitleResult{
			Content:  cached,
			Format:   req.Format,
			Language: req.Language,
		}, nil

	} else {

		s.logger.Debug(
			"subtitle cache miss",
			"videoID", req.VideoID,
			"cacheKey", cacheKey,
			"error", err,
		)
	}

	s.logger.Debug(
		"downloading subtitle",
		"videoID", req.VideoID,
		"language", req.Language,
		"type", req.Type,
		"format", req.Format,
	)

	data, err := s.pool.GetSubtitle(
		ctx,
		req.VideoID,
		req.Language,
		req.Type,
		req.Format,
	)
	if err != nil {

		s.logger.Warn(
			"subtitle download failed",
			"videoID", req.VideoID,
			"language", req.Language,
			"type", req.Type,
			"format", req.Format,
			"error", err,
		)

		return nil, err
	}

	s.logger.Debug(
		"subtitle downloaded",
		"videoID", req.VideoID,
		"language", req.Language,
		"type", req.Type,
		"format", req.Format,
		"bytes", len(data),
	)

	if err := s.cache.Save(req.VideoID, cacheKey, data); err != nil {

		s.logger.Warn(
			"failed to cache subtitle",
			"videoID", req.VideoID,
			"cacheKey", cacheKey,
			"error", err,
		)

	} else {

		s.logger.Debug(
			"subtitle cached",
			"videoID", req.VideoID,
			"cacheKey", cacheKey,
		)
	}

	s.logger.Debug(
		"subtitle download completed",
		"videoID", req.VideoID,
		"language", req.Language,
		"type", req.Type,
		"format", req.Format,
		"bytes", len(data),
	)

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