package pipeline

import (
	"context"
	"fmt"

	// "github.com/Arvind215271/askito/internal/youtube"
	"github.com/Arvind215271/askito/internal/youtube/description"
	"github.com/Arvind215271/askito/internal/youtube/metadata"
	"github.com/Arvind215271/askito/internal/youtube/signal"
	wordstats "github.com/Arvind215271/askito/internal/youtube/signal/word_stats"
	"github.com/Arvind215271/askito/internal/youtube/subtitle"
	"github.com/Arvind215271/askito/internal/youtube/transcript"
)

func ProcessMetadata(ctx context.Context, rc *ResourceContext, metaSvc *metadata.Service) error {
	if rc.Video == nil || rc.Video.ID == "" {
		return fmt.Errorf("video ID is required for metadata processing")
	}

	meta, err := metaSvc.GetVideo(ctx, rc.Video.ID, metadata.ProviderYTDLP)
	if err != nil {
		return fmt.Errorf("failed to fetch base metadata for video %s: %w", rc.Video.ID, err)
	}

	*rc.Video = meta
	return nil
}

func ProcessDescription(ctx context.Context, rc *ResourceContext, descSvc *description.Service) error {
	if rc.Video == nil {
		return fmt.Errorf("video is nil in context")
	}

	descMeta, err := descSvc.GetDescription(ctx, rc.Video.Description)
	if err != nil {
		return fmt.Errorf("failed to fetch description for video %s: %w", rc.Video.ID, err)
	}

	MapDescription(rc.Video, &descMeta)
	return nil
}

func ProcessSubtitleAndTranscript(ctx context.Context, rc *ResourceContext, subSvc *subtitle.SubtitleService, transSvc *transcript.Service) (*transcript.Transcript, error) {
	if rc.Video == nil {
		return nil, fmt.Errorf("video is nil in context")
	}

	req := rc.Request
	var subReq subtitle.DownloadRequest

	if req.Subtitle != nil {
		subReq = *req.Subtitle
		subReq.VideoID = rc.Video.ID
	} else {
		prefs := req.Preferences
		if len(prefs) == 0 {
			prefs = subtitle.DefaultPreferences()
		}
		resolvedReq, err := subSvc.ResolveDownloadRequest(rc.Video.SubtitleMetadata, prefs, req.Format)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve subtitle download request for video %s: %w", rc.Video.ID, err)
		}
		resolvedReq.VideoID = rc.Video.ID
		subReq = resolvedReq
	}

	if err := subReq.Validate(); err != nil {
		return nil, fmt.Errorf("invalid subtitle request for video %s: %w", rc.Video.ID, err)
	}

	sub, err := subSvc.DownloadSubtitle(ctx, subReq, rc.Video.SubtitleMetadata)
	if err != nil {
		return nil, fmt.Errorf("subtitle error for video %s: %w", rc.Video.ID, err)
	}

	trans, err := transSvc.Parse(sub)
	if err != nil {
		return nil, fmt.Errorf("transcript error for video %s: %w", rc.Video.ID, err)
	}

	if req.Transcript != nil {
		processed, err := transSvc.Process(trans, req.Transcript)
		if err != nil {
			return nil, fmt.Errorf("transcript processing error for video %s: %w", rc.Video.ID, err)
		}

		rc.Video.Transcript = trans
		rc.Video.TranscriptText = processed
	} else {
		MapTranscript(rc.Video, trans)
	}

	return trans, nil
}

func ProcessSignal(ctx context.Context, rc *ResourceContext, sigSvc *signal.SignalService, trans *transcript.Transcript) error {
	if rc.Video == nil || trans == nil {
		return fmt.Errorf("video or transcript is nil in context")
	}

	req := rc.Request
	sigReq := req.Signal

	if sigReq == nil {
		defaultReq := signal.DefaultSignalRequest(rc.Video.ID)
		sigReq = &defaultReq
	}

	if err := sigReq.Validate(); err != nil {
		return fmt.Errorf("invalid signal request for video %s: %w", rc.Video.ID, err)
	}

	sig := sigSvc.AnalyzeWordStats(trans, wordstats.DefaultWordStatsConfig())
	MapSignal(rc.Video, &sig)
	return nil
}
