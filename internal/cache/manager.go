package cache

import (
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/Arvind215271/askito/internal/logger"
)

type Config struct {
	CacheDir string
	TTLDays  int
	MaxFiles int
}

type Manager struct {
	config Config
	logger *logger.Logger
}

func NewManager(cfg Config, logger *logger.Logger) *Manager {
	return &Manager{config: cfg, logger: logger}
}

func (m *Manager) GetPath(id, filename string) string {
	return filepath.Join(m.config.CacheDir, id, filename)
}

func (m *Manager) Get(id, filename string) ([]byte, error) {
	path := m.GetPath(id, filename)
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, err
	}
	if time.Since(info.ModTime()) > time.Duration(m.config.TTLDays)*24*time.Hour {
		return nil, os.ErrNotExist
	}
	return os.ReadFile(path)
}

func (m *Manager) Save(id, filename string, data []byte) error {
	dir := filepath.Dir(m.GetPath(id, filename))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(m.GetPath(id, filename), data, 0644)
}

func (m *Manager) VideoKey() string {
	return "metadata." +  ".json"
}

func (m *Manager) PlaylistKey() string {
	return "playlist." + ".json"
}

func (m *Manager) SubtitleKey(subType, language, format string) string {
	return "subtitles." + subType + "." + language + "." + format
}

func (m *Manager) SubtitlePath(subType string) string {
	return "subtitles." + subType
}

// Cleanup generalizes the cleanup logic to remove expired directories
// based on the age of a specific "marker" file (e.g., metadata.json).
func (m *Manager) Cleanup(markerFilename string) error {
	entries, err := os.ReadDir(m.config.CacheDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	type fileInfo struct {
		path    string
		modTime time.Time
	}
	var directories []fileInfo

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		markerPath := filepath.Join(m.config.CacheDir, entry.Name(), markerFilename)
		info, err := os.Stat(markerPath)
		if err != nil {
			continue
		}

		// TTL-based eviction
		if time.Since(info.ModTime()) > time.Duration(m.config.TTLDays)*24*time.Hour {
			os.RemoveAll(filepath.Join(m.config.CacheDir, entry.Name()))
			continue
		}

		directories = append(directories, fileInfo{
			path:    filepath.Join(m.config.CacheDir, entry.Name()),
			modTime: info.ModTime(),
		})
	}

	// Capacity-based eviction
	if len(directories) > m.config.MaxFiles {
		sort.Slice(directories, func(i, j int) bool {
			return directories[i].modTime.Before(directories[j].modTime)
		})

		numToDelete := len(directories) - (m.config.MaxFiles * 9 / 10) // Keep 90%
		for i := 0; i < numToDelete; i++ {
			os.RemoveAll(directories[i].path)
		}
	}

	return nil
}
