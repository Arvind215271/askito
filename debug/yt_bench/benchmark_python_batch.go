package yt_bench

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/Arvind215271/askito/internal/youtube/metadata/ytdlp"
	"github.com/Arvind215271/askito/internal/youtube/metadata/ytdlp/python"
)

// BatchPlaylist helper function
func BatchPlaylist(videoIDs []string, batchSize int) [][]string {
	var batches [][]string
	for i := 0; i < len(videoIDs); i += batchSize {
		end := i + batchSize
		if end > len(videoIDs) {
			end = len(videoIDs)
		}
		batches = append(batches, videoIDs[i:end])
	}
	return batches
}

func BenchmarkPythonBatch(ctx context.Context, playlistID string, outputDir string, workers int, batchSize int) (BenchmarkResult, error) {
	start := time.Now()

	// Memory tracking
	var peakMemory uint64
	ticker := time.NewTicker(100 * time.Millisecond)
	go func() {
		for range ticker.C {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			if m.Alloc > peakMemory {
				atomic.StoreUint64(&peakMemory, m.Alloc)
			}
		}
	}()
	defer ticker.Stop()

	// Ensure log directory exists
	logDir := "debug/yt_bench/output/logs"
	os.MkdirAll(logDir, 0755)
	logFile, err := os.Create(filepath.Join(logDir, fmt.Sprintf("python_batch_w%d_b%d.log", workers, batchSize)))
	if err != nil {
		return BenchmarkResult{}, err
	}
	defer logFile.Close()
	logger := log.New(logFile, "", log.LstdFlags)

	flatStart := time.Now()
	flatOutput, err := FetchPlaylistFlatMetadata(ctx, playlistID)
	flatFetchTime := time.Since(flatStart)
	if err != nil {
		return BenchmarkResult{}, err
	}
	var flat ytdlp.YTPlaylistOutput
	scanner := bufio.NewScanner(bytes.NewReader(flatOutput))
	var videoIDs []string
	for scanner.Scan() {
		var entry ytdlp.YTPlaylistEntry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err == nil {
			flat.Entries = append(flat.Entries, entry)
			videoIDs = append(videoIDs, entry.ID)
		}
	}

	batches := BatchPlaylist(videoIDs, batchSize)

	pool, err := python.NewBatchPool(workers)
	if err != nil {
		return BenchmarkResult{}, err
	}
	defer pool.Close()

	var metadataFetchTime int64 // nanoseconds
	var diskWriteTime int64     // nanoseconds

	pool.Start(ctx, func(batch []string, results []map[string]any, duration time.Duration, err error) {
		if err != nil {
			logger.Printf("Batch error: %v", err)
			return
		}
		atomic.AddInt64(&metadataFetchTime, int64(duration))

		writeStart := time.Now()
		
		for i, videoID := range batch {
			if i < len(results) {
				res := results[i]
				if ok, _ := res["ok"].(bool); !ok {
					errMsg, _ := res["error"].(string)
					fmt.Fprintf(os.Stderr, "Failed to process video %s: %s\n", videoID, errMsg)
					logger.Printf("Failed to process video %s: %s", videoID, errMsg)
					continue
				}
				fileName := filepath.Join(outputDir, fmt.Sprintf("E_W%d_B%d_%s.json", workers, batchSize, videoID))
				resultData, _ := json.Marshal(res)
				os.WriteFile(fileName, resultData, 0644)
			}
		}
		atomic.AddInt64(&diskWriteTime, int64(time.Since(writeStart)))
		logger.Printf("Processed batch size=%d", len(batch))
	})

	for _, batch := range batches {
		pool.Submit(batch)
	}

	res := BenchmarkResult{
		PlaylistID:             playlistID,
		Strategy:               fmt.Sprintf("Python Batch (W=%d B=%d)", workers, batchSize),
		PlaylistSize:           len(flat.Entries),
		TotalExecutionTime:     time.Since(start),
		StartedAt:              start,
		ProcessesSpawned:       1 + workers + len(flat.Entries),
		PeakMemoryUsage:        atomic.LoadUint64(&peakMemory),
		FlatPlaylistFetchTime:  flatFetchTime,
		VideoMetadataFetchTime: time.Duration(atomic.LoadInt64(&metadataFetchTime)),
		DiskWriteTime:          time.Duration(atomic.LoadInt64(&diskWriteTime)),
	}
	SaveResult(res)
	return res, nil
}
