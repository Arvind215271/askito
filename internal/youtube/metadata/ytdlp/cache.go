// Package ytdlp provides caching for youtube-dl metadata.
// The cache stores metadata files in: ./cache/ytdlp/[videoID]/metadata.json
// This structure isolates metadata from other assets (like subtitles) stored in the same folder.
package ytdlp

import (
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/Arvind215271/askito/internal/logger"
)

type Cache struct {
	config CacheConfig
	logger *logger.Logger
}

// NewCache creates a new cache instance.
func NewCache(cfg CacheConfig, logger *logger.Logger) *Cache {
	return &Cache{
		config: cfg,
		logger: logger,
	}
}

// Get retrieves metadata from cache if it exists and is within the TTL.
func (c *Cache) Get(videoID string) ([]byte, error) {
	path := c.getPath(videoID)
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		c.logger.Debug("cache miss", "videoID", videoID)
		return nil, err
	}

	// Check TTL
	if time.Since(info.ModTime()) > time.Duration(c.config.TTLDays)*24*time.Hour {
		c.logger.Debug("cache expired", "videoID", videoID)
		return nil, os.ErrNotExist
	}

	c.logger.Debug("cache hit", "videoID", videoID)
	return os.ReadFile(path)
}

// Save stores metadata JSON to disk.
// It ensures the directory structure exists using MkdirAll, which is safe for other files in the same directory.
func (c *Cache) Save(videoID string, data []byte) error {
	dir := filepath.Dir(c.getPath(videoID))
	if err := os.MkdirAll(dir, 0755); err != nil {
		c.logger.Error("failed to create cache directory", "error", err, "videoID", videoID)
		return err
	}
	err := os.WriteFile(c.getPath(videoID), data, 0644)
	if err != nil {
		c.logger.Error("failed to write cache file", "error", err, "videoID", videoID)
	}
	return err
}

// Cleanup performs maintenance by removing expired and excess files.
func (c *Cache) Cleanup() error {
	entries, err := os.ReadDir(c.config.CacheDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		c.logger.Error("failed to read cache directory for cleanup", "error", err)
		return err
	}

	type fileInfo struct {
		path    string
		modTime time.Time
	}
	var files []fileInfo
	deletedCount := 0

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		metadataPath := filepath.Join(c.config.CacheDir, entry.Name(), "metadata.json")
		info, err := os.Stat(metadataPath)
		if err != nil {
			continue
		}

		// TTL-based eviction
		if time.Since(info.ModTime()) > time.Duration(c.config.TTLDays)*24*time.Hour {
			os.RemoveAll(filepath.Join(c.config.CacheDir, entry.Name()))
			deletedCount++
			continue
		}

		files = append(files, fileInfo{path: filepath.Join(c.config.CacheDir, entry.Name()), modTime: info.ModTime()})
	}

	// Capacity-based eviction
	if len(files) > c.config.MaxFiles {
		sort.Slice(files, func(i, j int) bool {
			return files[i].modTime.Before(files[j].modTime)
		})

		numToDelete := len(files) - (c.config.MaxFiles * 9 / 10) // Keep 90%
		for i := 0; i < numToDelete; i++ {
			os.RemoveAll(files[i].path)
			deletedCount++
		}
	}
	
	if deletedCount > 0 {
		c.logger.Info("cache cleanup performed", "deletedCount", deletedCount)
	}

	return nil
}

func (c *Cache) getPath(videoID string) string {
	return filepath.Join(c.config.CacheDir, videoID, "metadata.json")
}
