package yt_bench

// import (
// 	"context"
// 	"os"
// 	"testing"
// )

// func TestSequential(t *testing.T) {
// 	playlistID := "PLWKjhJtqVAbnqBxcdjVGgT3uVR10bzTEB"
// 	outputDir := "debug/yt_bench/output"
// 	os.RemoveAll(outputDir)
// 	os.MkdirAll(outputDir, 0755)
// 	res, err := BenchmarkSequential(context.Background(), playlistID, outputDir)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Logf("Result: %v", res.TotalExecutionTime)
// }
