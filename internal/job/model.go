package job

import (
	"errors"
	"time"
)

// JobType represents the type of a job.
type JobType string

const (
	JobTypeExport     JobType = "export"
	JobTypeSubtitle   JobType = "subtitle"
	JobTypeTranscript JobType = "transcript"
)

// JobStatus represents the current state of a job.
type JobStatus string

const (
	StatusQueued    JobStatus = "queued"
	StatusRunning   JobStatus = "running"
	StatusCompleted JobStatus = "completed"
	StatusFailed    JobStatus = "failed"
	StatusCancelled JobStatus = "cancelled"
)

// IsTerminal returns true if the status is a terminal state.
func (s JobStatus) IsTerminal() bool {
	return s == StatusCompleted || s == StatusFailed || s == StatusCancelled
}

// Job represents an asynchronous job model.
type Job struct {
	ID         string     `json:"id"`
	Type       JobType    `json:"type"`
	Status     JobStatus  `json:"status"`
	OwnerID    string     `json:"owner_id"`
	CreatedAt  time.Time  `json:"created_at"`
	StartedAt  *time.Time `json:"started_at,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
	Error      string     `json:"error,omitempty"`
}

// Sentinel errors for job management.
var (
	ErrJobNotFound       = errors.New("job not found")
	ErrInvalidTransition = errors.New("invalid job status transition")
	ErrActiveJobExists   = errors.New("user already has an active job")
)
