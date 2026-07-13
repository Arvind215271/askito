package yt_bench

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func RunBenchmarks(playlistIDs []string) {
	ctx := context.Background()
	outputDir := "debug/yt_bench/output"
	// os.RemoveAll(outputDir)
	os.MkdirAll(outputDir, 0755)

	for _, playlistID := range playlistIDs {
		fmt.Printf("Starting Benchmarks for: %s\n", playlistID)

		// //Strategy A
		// resA, err := BenchmarkSequential(ctx, playlistID, outputDir)
		// if err != nil {
		// 	log.Fatalf("Strategy A failed: %v", err)
		// }
		// fmt.Printf("Strategy A (Sequential) Time: %v\n", resA.TotalExecutionTime)

		// // Strategy B
		// resB, err := BenchmarkPlaylistDump(ctx, playlistID, outputDir)
		// if err != nil {
		// 	log.Fatalf("Strategy B failed: %v", err)
		// }
		// fmt.Printf("Strategy B (Full Dump) Time: %v\n", resB.TotalExecutionTime)

		// Strategy C
		// for _, w := range []int{16, 32, 64, 128} {
		// 	resC, err := BenchmarkParallel(ctx, playlistID, outputDir, w)
		// 	if err != nil {
		// 		log.Fatalf("Strategy C (%d) failed: %v", w, err)
		// 	}
		// 	fmt.Printf("Strategy C (Parallel %d) Time: %v\n", w, resC.TotalExecutionTime)
		// }
		workers := []int{16,32,64,128,256}

		// Strategy D
		// for _, w := range workers {
		// 	resD, err := BenchmarkPythonPool(ctx, playlistID, outputDir, w)
		// 	if err != nil {
		// 		log.Fatalf("Strategy D (%d) failed: %v", w, err)
		// 	}
		// 	fmt.Printf("Strategy D (Python Pool %d) Time: %v\n", w, resD.TotalBenchmarkDuration)
		// }

		for _, w := range workers {
			resD, err := BenchmarkPythonPoolWarm(ctx, playlistID, outputDir, w)
			if err != nil {
				log.Fatalf("Strategy D.2 (%d) failed: %v", w, err)
			}
			fmt.Printf("Strategy D.2 (Python Pool Warm %d) Time: %v\n", w, resD.TotalBenchmarkDuration)
		}

		// for _, w := range workers {
		// 	resD, err := BenchmarkPythonPoolWarmSleep(ctx, playlistID, outputDir, w)
		// 	if err != nil {
		// 		log.Fatalf("Strategy D.3 (%d) failed: %v", w, err)
		// 	}
		// 	fmt.Printf("Strategy D.3 (Python Pool Warm Sleep %d) Time: %v\n", w, resD.TotalBenchmarkDuration)
		// }

		// // Strategy E
		// for _, w := range []int{1,4,8,16} {
		// 	for _, b := range []int{4,8} {
		// 		resE, err := BenchmarkPythonBatch(ctx, playlistID, outputDir, w, b)
		// 		if err != nil {
		// 			log.Fatalf("Strategy E (W=%d B=%d) failed: %v", w, b, err)
		// 		}
		// 		fmt.Printf("Strategy E (Python Batch W=%d B=%d) Time: %v\n", w, b, resE.TotalExecutionTime)
		// 	}
		// }
	}
}

func SaveResult(res BenchmarkResult) {
	dir := "docs"
	os.MkdirAll(dir, 0755)
	historyFile := filepath.Join(dir, "benchmark_history.json")

	var history []BenchmarkResult
	if data, err := os.ReadFile(historyFile); err == nil {
		json.Unmarshal(data, &history)
	}
	history = append(history, res)
	data, _ := json.MarshalIndent(history, "", "  ")
	err := os.WriteFile(historyFile, data, 0644)
	if err != nil {
		fmt.Printf("ERROR: Failed to save result to %s: %v\n", historyFile, err)
		return
	}
	fmt.Printf("Saved result for %s strategy: %s\n", res.Strategy, res.PlaylistID)
}
