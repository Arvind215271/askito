package export

import (
	"sync"
	"github.com/Arvind215271/askito/internal/youtube"
)

// Service is the application layer that orchestrates
// building export data + choosing exporter.
type Service struct {
	exporters map[Format]Exporter
	mu        sync.RWMutex
}

// NewService creates the export service with registered exporters.
func NewService() *Service {
	return &Service{
		exporters: make(map[Format]Exporter),
	}
}

// RegisterExporter allows plugging in new formats (JSON, CSV, etc.)
func (s *Service) RegisterExporter(format Format, exporter Exporter) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.exporters[format] = exporter
}

// ExportPlaylist is the main orchestration function.
func (s *Service) ExportPlaylist(
	playlist youtube.Playlist,
	req PlaylistExportRequest,
) ([]byte, error) {

	data, err := BuildPlaylistExport(playlist, req.VideoFields)
	if err != nil {
		return nil, err
	}

	return s.exportData(req.Format, data)
}

// ExportVideo is the main orchestration function.
func (s *Service) ExportVideo(
	video youtube.Video,
	req VideoExportRequest,
) ([]byte, error) {

	data, err := BuildVideoExport(video, req.Fields)
	if err != nil {
		return nil, err
	}

	return s.exportData(req.Format, data)
}

// internal shared logic
func (s *Service) exportData(format Format, data ExportData) ([]byte, error) {
	s.mu.RLock()
	exporter, ok := s.exporters[format]
	s.mu.RUnlock()

	if !ok {
		return nil, youtube.Err.Export.InvalidFormat()
	}

	return exporter.Export(data)
}