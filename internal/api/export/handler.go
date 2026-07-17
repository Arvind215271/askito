package export

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/Arvind215271/askito/internal/api"
	"github.com/Arvind215271/askito/internal/youtube/export"
	"github.com/Arvind215271/askito/internal/youtube/fields"
	youtubeurl "github.com/Arvind215271/askito/internal/youtube/input"
	"github.com/Arvind215271/askito/internal/youtube/pipeline"
)

type Handler struct {
	pipelineService *pipeline.Service
	exportService   *export.Service
}

func NewHandler(
	pipelineService *pipeline.Service,
	exportService *export.Service,
) *Handler {
	return &Handler{
		pipelineService: pipelineService,
		exportService:   exportService,
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

	if parsedInput.InputType != youtubeurl.InputTypePlaylist {
		return ErrInvalidInputType
	}

	fieldPlanner, err := fields.NewPlanner(req.VideoFields)
	if err != nil {
		return err
	}

	// Playlist processing is not implemented yet.
	//
	// When playlist processing is added, each video should receive
	// the same feature-specific request configuration:
	//
	//     req.Subtitle
	//     req.Transcript
	//     req.Signal
	//
	// through a pipeline.Request.

	_ = pipeline.Request{
		Fields:     req.VideoFields,
		Subtitle:   req.Subtitle,
		Transcript: req.Transcript,
		Signal:     req.Signal,
	}

	exportReq := export.PlaylistExportRequest{
		PlaylistID:  parsedInput.ID,
		VideoFields: fieldPlanner,
		Format:      export.Format(format),
	}

	// TODO: Process playlist and export results.
	_ = exportReq

	return nil
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