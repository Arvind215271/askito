package job

import (
	"context"
	"fmt"

	"github.com/Arvind215271/askito/internal/logger"
)

// JobRunner manages asynchronous background job execution using JobManager.
type JobRunner struct {
	manager *JobManager
	logger  *logger.Logger
}

// NewRunner creates a new JobRunner instance with an optional logger.
func NewRunner(manager *JobManager, loggers ...*logger.Logger) *JobRunner {
	var log *logger.Logger
	if len(loggers) > 0 {
		log = loggers[0]
	}
	return &JobRunner{
		manager: manager,
		logger:  log,
	}
}

// Submit creates a new job in queued state and spawns a background goroutine to execute the work.
// It returns the created Job snapshot immediately without waiting for execution to complete.
func (r *JobRunner) Submit(jobType JobType, work func(context.Context) error) (*Job, error) {
	job, err := r.manager.Create(jobType)
	if err != nil {
		return nil, err
	}

	go r.execute(job.ID, work)

	return job, nil
}

// execute runs the job lifecycle in a background goroutine.
func (r *JobRunner) execute(jobID string, work func(context.Context) error) {
	if _, err := r.manager.Transition(jobID, StatusRunning, nil); err != nil {
		if r.logger != nil {
			r.logger.Error("failed to transition job to running", "job_id", jobID, "error", err)
		}
		return
	}

	ctx := context.Background()

	// Panic recovery
	defer func() {
		if p := recover(); p != nil {
			err := fmt.Errorf("job panicked: %v", p)
			if r.logger != nil {
				r.logger.Error("job panicked during execution", "job_id", jobID, "panic", p, "error", err)
			}
			_, _ = r.manager.Transition(jobID, StatusFailed, err)
		}
	}()

	if err := work(ctx); err != nil {
		_, _ = r.manager.Transition(jobID, StatusFailed, err)
	} else {
		_, _ = r.manager.Transition(jobID, StatusCompleted, nil)
	}
}
