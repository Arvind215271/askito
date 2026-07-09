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

func NewHandler(p *pipeline.Service, e *export.Service) *Handler {
	return &Handler{
		pipelineService: p,
		exportService:   e,
	}
}

func (h *Handler) ExportVideo(c *echo.Context) error {
	var req VideoExportRequest
	if err := c.Bind(&req); err != nil {
		return api.Err.Common.BadRequest("invalid request body").Wrap(err)
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
		return api.Err.Common.BadRequest("invalid youtube input").Wrap(err)
	}
	if parsedInput.InputType != youtubeurl.InputTypeVideo {
        return ErrInvalidInputType
    }

	fieldPlanner, err := fields.NewPlanner(req.Fields)
	if err != nil {
		return err
	}
    
    // Process using service
    pipelinePlanner := pipeline.NewPlanner(fieldPlanner)
	
    v, err := h.pipelineService.Process(c.Request().Context(), parsedInput.ID, pipelinePlanner)
    if err != nil {
        return err
    }

    // Convert request model to export request model
    exportReq := export.VideoExportRequest{
        VideoID: v.ID,
        Fields: fieldPlanner,
        Format: export.Format(format),
    }

    bytes, err := h.exportService.ExportVideo(*v, exportReq)
    if err != nil {
        return err
    }

	return c.Blob(http.StatusOK, "application/json", bytes)
}

func (h *Handler) ExportPlaylist(c *echo.Context) error {
	var req PlaylistExportRequest
	if err := c.Bind(&req); err != nil {
		return api.Err.Common.BadRequest("invalid request body").Wrap(err)
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
		return api.Err.Common.BadRequest("invalid youtube input").Wrap(err)
	}
	if parsedInput.InputType != youtubeurl.InputTypePlaylist {
        return ErrInvalidInputType
    }
    
    // Planner...
	fieldPlanner, err := fields.NewPlanner(req.VideoFields)
	if err != nil {
		return err
	}

    // Pipeline... Process Playlist here...
    // ...

    exportReq := export.PlaylistExportRequest{
        PlaylistID: parsedInput.ID,
        VideoFields: fieldPlanner,
        Format: export.Format(format),
    }
    
    // Export service ...
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
