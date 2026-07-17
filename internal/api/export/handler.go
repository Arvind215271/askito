package export

import (
	"net/http"
	"encoding/json"

	"github.com/labstack/echo/v5"

	"github.com/Arvind215271/askito/internal/api"
	"github.com/Arvind215271/askito/internal/youtube/export"
	"github.com/Arvind215271/askito/internal/youtube/fields"
	youtubeurl "github.com/Arvind215271/askito/internal/youtube/input"
	"github.com/Arvind215271/askito/internal/youtube/pipeline"
	"github.com/Arvind215271/askito/internal/youtube/metadata"
	"github.com/Arvind215271/askito/internal/youtube"

)

type Handler struct {
	metadataService *metadata.Service
	pipelineService *pipeline.Service
	exportService   *export.Service
}

func NewHandler(
	metadataService *metadata.Service,
	pipelineService *pipeline.Service,
	exportService *export.Service,
) *Handler {
	return &Handler{
		metadataService: metadataService,
		pipelineService: pipelineService,
		exportService:   exportService,
	}
}

func parseFormat(s string) (export.Format, error) {
        format := export.FormatJSON

        if s != "" {
                format = export.Format(s)
        }

        switch format {
        case export.FormatJSON:
                return format, nil

        default:
                return "", ErrInvalidFormat
        }
}


func (h *Handler) ExportVideo(c *echo.Context) error {
	var req VideoExportRequest

	if err := c.Bind(&req); err != nil {
		return api.Err.Common.
			BadRequest("invalid request body").
			Wrap(err)
	}

	if req.Input == "" {
		return api.Err.Common.BadRequest("input is required")
	}

	format, err := parseFormat(req.Format)
	if err != nil {
		return err
	}

	parsedInput, err := youtubeurl.Parse(req.Input)
	if err != nil {
		return api.Err.Common.
			BadRequest("invalid youtube input").
			Wrap(err)
	}

	if parsedInput.InputType != youtubeurl.InputTypeVideo {
		return ErrInvalidInputType
	}

	fieldPlanner, err := fields.NewPlanner(req.Fields)
	if err != nil {
		return err
	}

	pipelinePlanner := pipeline.NewPlanner(fieldPlanner)

	pipelineReq := &pipeline.Request{
		Fields: req.Fields,

		Subtitle:   req.Subtitle,
		Transcript: req.Transcript,
		Signal:     req.Signal,
	}

	video, err := h.pipelineService.Process(
		c.Request().Context(),
		parsedInput.ID,
		pipelineReq,
		pipelinePlanner,
	)
	if err != nil {
		return err
	}

	exportReq := export.VideoExportRequest{
		VideoID: video.ID,
		Fields:  fieldPlanner,
		Format:  export.Format(format),
	}

	data, err := h.exportService.ExportVideo(
		*video,
		exportReq,
	)
	if err != nil {
		return err
	}

	return c.Blob(
		http.StatusOK,
		"application/json",
		data,
	)
}

func (h *Handler) ExportPlaylist(c *echo.Context) error {
	var req PlaylistExportRequest

	if err := c.Bind(&req); err != nil {
		return api.Err.Common.
			BadRequest("invalid request body").
			Wrap(err)
	}

	if req.Input == "" {
		return api.Err.Common.BadRequest("input is required")
	}

	// format, err := parseFormat(req.Format)
	// if err != nil {
	// 	return err
	// }

	parsedInput, err := youtubeurl.Parse(req.Input)
	if err != nil {
		return api.Err.Common.
			BadRequest("invalid youtube input").
			Wrap(err)
	}

	if parsedInput.InputType != youtubeurl.InputTypePlaylist {
		return ErrInvalidInputType
	}

	fieldPlanner, err := fields.NewPlanner(req.VideoFields)
	if err != nil {
		return err
	}

	pipelineReq := &pipeline.Request{
		Fields:     req.VideoFields,
		Subtitle:   req.Subtitle,
		Transcript: req.Transcript,
		Signal:     req.Signal,
	}

	pipelinePlanner := pipeline.NewPlanner(fieldPlanner)

	ctx := c.Request().Context()

	// 1. Fetch playlist metadata.
	playlist, err := h.metadataService.GetPlaylistMetadata(
		ctx,
		parsedInput.ID,
		metadata.ProviderYTDLP,
	)
	if err != nil {
		return err
	}

	// 2. Fetch playlist items.
	items, err := h.metadataService.GetPlaylistItems(
		ctx,
		parsedInput.ID,
		metadata.ProviderYTDLP,
	)
	if err != nil {
		return err
	}

	// 3. Extract video IDs in playlist order.
	videoIDs := make([]string, len(items))

	for i, item := range items {
		videoIDs[i] = item.VideoID
	}

	// 4. Process all videos concurrently.
	videos := h.pipelineService.ProcessVideos(
		ctx,
		videoIDs,
		pipelineReq,
		pipelinePlanner,
	)

	// 5. Wrap processed videos back into PlaylistVideo.
	playlist.Videos = make([]youtube.PlaylistVideo, len(videos))

	for i, video := range videos {
		playlist.Videos[i] = youtube.PlaylistVideo{
			Video:    *video,
			Position: items[i].Position,
			AddedAt:  items[i].AddedAt,
		}
	}

	// 6. Return the processed playlist as JSON.
	data, err := json.Marshal(playlist)
	if err != nil {
		return err
	}

	return c.Blob(
		http.StatusOK,
		"application/json",
		data,
	)
}