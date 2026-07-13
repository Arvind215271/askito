package yt_bench

import (
	"sync/atomic"
	"time"
)

type WorkerExecution struct {
	VideoID    string        `json:"video_id"`
	StartedAt  time.Time     `json:"started_at"`
	FinishedAt time.Time     `json:"finished_at"`
	Duration   time.Duration `json:"duration"`
}

type BenchmarkResult struct {
	PlaylistID            string        `json:"playlist_id"`
	Strategy              string        `json:"strategy"`
	PlaylistSize          int           `json:"playlist_size"`
	ProcessesSpawned      int           `json:"processes_spawned"`
	CPUUsage              float64       `json:"cpu_usage"`
	PeakMemoryUsage       uint64        `json:"peak_memory_usage"`
	MetadataReturned      int           `json:"metadata_returned"`

	// Original fields (kept for compatibility, should be migrated)
	TotalExecutionTime    time.Duration `json:"total_execution_time"`
	StartedAt             time.Time     `json:"started_at"`
	FlatPlaylistFetchTime time.Duration `json:"flat_playlist_fetch_time"`
	VideoMetadataFetchTime time.Duration `json:"video_metadata_fetch_time"` // Deprecated: use ActualMetadataFetchDuration
	DiskWriteTime         time.Duration `json:"disk_write_time"`
	WorkerSpawnTime       time.Duration `json:"worker_spawn_time"`
	WorkerWarmupTime      time.Duration `json:"worker_warmup_time"`

	// New detailed timing fields
	BenchmarkStartedAt      time.Time     `json:"benchmark_started_at"`
	BenchmarkFinishedAt     time.Time     `json:"benchmark_finished_at"`
	FlatFetchStartedAt      time.Time     `json:"flat_fetch_started_at"`
	FlatFetchFinishedAt     time.Time     `json:"flat_fetch_finished_at"`
	WorkerSpawnStartedAt    time.Time     `json:"worker_spawn_started_at"`
	WorkerSpawnFinishedAt   time.Time     `json:"worker_spawn_finished_at"`
	MetadataDispatchStartedAt time.Time   `json:"metadata_dispatch_started_at"`
	MetadataDispatchFinishedAt time.Time  `json:"metadata_dispatch_finished_at"`
	FirstWorkerStartedAt    time.Time     `json:"first_worker_started_at"`
	LastWorkerFinishedAt    time.Time     `json:"last_worker_finished_at"`
	DiskWriteStartedAt      time.Time     `json:"disk_write_started_at"`
	DiskWriteFinishedAt     time.Time     `json:"disk_write_finished_at"`

	// Derived durations
	FlatFetchDuration           time.Duration `json:"flat_fetch_duration"`
	WorkerSpawnDuration         time.Duration `json:"worker_spawn_duration"`
	DispatchDuration            time.Duration `json:"dispatch_duration"`
	ActualMetadataFetchDuration time.Duration `json:"actual_metadata_fetch_duration"` // LastWorkerFinishedAt - FirstWorkerStartedAt
	DiskWriteDuration           time.Duration `json:"disk_write_duration"`
	TotalBenchmarkDuration      time.Duration `json:"total_benchmark_duration"`

	// Aggregate worker time (sum of all worker durations)
	AggregateWorkerTime time.Duration `json:"aggregate_worker_time"`

	// Warmup/Sleep timing (for warm/sleep benchmarks)
	WarmupStartedAt  time.Time     `json:"warmup_started_at,omitempty"`
	WarmupFinishedAt time.Time     `json:"warmup_finished_at,omitempty"`
	SleepStartedAt   time.Time     `json:"sleep_started_at,omitempty"`
	SleepFinishedAt  time.Time     `json:"sleep_finished_at,omitempty"`
	WarmupDuration   time.Duration `json:"warmup_duration,omitempty"`
	SleepDuration    time.Duration `json:"sleep_duration,omitempty"`

	// Atomic counters for concurrent updates
	firstWorkerStartedAt atomic.Value // time.Time
	lastWorkerFinishedAt atomic.Value // time.Time
	firstDiskWriteAt     atomic.Value // time.Time
	lastDiskWriteAt      atomic.Value // time.Time
}
