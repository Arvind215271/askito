package yt_bench

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func BenchmarkPlaylistDump(ctx context.Context, playlistID string, outputDir string) (BenchmarkResult, error) {
	start := time.Now()
	
	// Ensure log directory exists
	logDir := "debug/yt_bench/output/logs"
	os.MkdirAll(logDir, 0755)
	logFile, err := os.Create(filepath.Join(logDir, "playlist_dump.log"))
	if err != nil {
		return BenchmarkResult{}, err
	}
	defer logFile.Close()
	logger := log.New(logFile, "", log.LstdFlags)

	fetchStart := time.Now()
	fullOutput, err := FetchPlaylistFullMetadata(ctx, playlistID)
	fetchDuration := time.Since(fetchStart)
	if err != nil {
		return BenchmarkResult{}, err
	}
	
	writeStart := time.Now()
	os.WriteFile(filepath.Join(outputDir, "B_full_dump.json"), fullOutput, 0644)
	writeDuration := time.Since(writeStart)
	
	msg := fmt.Sprintf("PlaylistDump: FetchDuration=%v WriteDuration=%v\n", fetchDuration, writeDuration)
	fmt.Print(msg)
	logger.Print(msg)
	
	res := BenchmarkResult{
		PlaylistID:            playlistID,
		Strategy:              "PlaylistDump",
		PlaylistSize:          0, // Unknown/Irrelevant here
		TotalExecutionTime:    time.Since(start),
		StartedAt:             start,
		ProcessesSpawned:      1,
		FlatPlaylistFetchTime: fetchDuration,
		DiskWriteTime:         writeDuration,
	}
	SaveResult(res)
	return res, nil
}
