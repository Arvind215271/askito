# JobRunner Architecture & Implementation Design

## 1. Overview & Objectives
The `JobRunner` component in package `internal/job` provides asynchronous background execution for jobs managed by `JobManager`. It bridges job creation/lifecycle management and asynchronous goroutine execution without introducing generic worker pools, job queues, task scheduling, or HTTP request context dependencies.

## 2. Public API & Struct Design (`runner.go`)

```go
package job

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/arvind-saini/askito/internal/logger"
)

// JobRunner coordinates asynchronous execution of jobs using JobManager.
type JobRunner struct {
	manager *JobManager
	logger  *logger.Logger
}

// NewRunner creates a new JobRunner instance.
func NewRunner(manager *JobManager, log *logger.Logger) *JobRunner {
	return &JobRunner{
		manager: manager,
		logger:  log,
	}
}

// Submit creates a new job in the JobManager and launches it asynchronously in a background goroutine.
func (r *JobRunner) Submit(jobType JobType, work func(context.Context) error) (*Job, error) {
	job, err := r.manager.Create(jobType)
	if err != nil {
		return nil, fmt.Errorf("failed to create job for submission: %w", err)
	}

	// Launch background execution
	go r.execute(job.ID, work)

	return job, nil
}
```

## 3. Execution Flow & Lifecycle Integration (`execute`)

The execution goroutine follows a strict lifecycle state machine:
1. **Queued -> Running**:
   - Calls `r.manager.Transition(jobID, StatusRunning, nil)`.
   - If this transition fails (e.g. if the job was cancelled while queued), execution is aborted immediately without running `work`.
2. **Work Execution**:
   - Uses an independent background context (`context.Background()`), decoupled from any incoming HTTP request lifecycle.
   - Protected by `defer func()` for panic recovery.
3. **Completion / Failure**:
   - If `work` returns `nil` (and no panic occurred), transitions to `StatusCompleted`: `r.manager.Transition(jobID, StatusCompleted, nil)`.
   - If `work` returns an error `err`, transitions to `StatusFailed` with the error: `r.manager.Transition(jobID, StatusFailed, err)`.
   - If a panic occurs:
     - Recovered via `recover()`.
     - Captured as an error (e.g., formatted panic message with stack trace via `debug.Stack()`).
     - Transitions to `StatusFailed` with the panic error.
4. **Lifecycle Transition Error Handling**:
   - If transitioning to `StatusCompleted` or `StatusFailed` fails (e.g., due to concurrent cancellation), log the error via `r.logger` without panicking or entering infinite retry loops.

### Detailed `execute` implementation:

```go
func (r *JobRunner) execute(jobID string, work func(context.Context) error) {
	// 1. Transition to Running
	_, err := r.manager.Transition(jobID, StatusRunning, nil)
	if err != nil {
		r.logger.Error("failed to transition job to running", "job_id", jobID, "error", err)
		return
	}

	// 2. Setup panic recovery and execution context
	ctx := context.Background()

	var workErr error
	func() {
		defer func() {
			if p := recover(); p != nil {
				stack := string(debug.Stack())
				r.logger.Error("job execution panicked", "job_id", jobID, "panic", p, "stack", stack)
				workErr = fmt.Errorf("panic recovered: %v", p)
			}
		}()
		workErr = work(ctx)
	}()

	// 3. Transition to terminal state (Completed or Failed)
	if workErr != nil {
		_, transErr := r.manager.Transition(jobID, StatusFailed, workErr)
		if transErr != nil {
			r.logger.Error("failed to transition job to failed state", "job_id", jobID, "work_error", workErr, "transition_error", transErr)
		}
	} else {
		_, transErr := r.manager.Transition(jobID, StatusCompleted, nil)
		if transErr != nil {
			r.logger.Error("failed to transition job to completed state", "job_id", jobID, "transition_error", transErr)
		}
	}
}
```

## 4. Testing Strategy (`runner_test.go`)

Deterministic unit tests without arbitrary sleeps:
1. **Happy Path (`TestJobRunner_Success`)**:
   - Submit a job with a work function that completes successfully.
   - Wait for completion deterministically using channels or polling with timeout / synchronization primitives.
   - Verify final status is `completed`, `StartedAt` and `FinishedAt` are set.
2. **Failure Path (`TestJobRunner_Failure`)**:
   - Submit a job whose work function returns an error.
   - Verify final status is `failed`, error message is recorded correctly.
3. **Panic Recovery (`TestJobRunner_Panic`)**:
   - Submit a job whose work function panics (`panic("boom")`).
   - Verify panic is caught, job status transitions to `failed`, and error mentions the panic.
4. **Transition Failure Abort (`TestJobRunner_AbortedWhenQueuedCancelled`)**:
   - Create a job, transition it to `cancelled` manually before running.
   - Attempt execution or ensure work function is never executed if transition to running fails.

## 5. Feature 1 API Assessment
- Feature 1 (`JobManager`) provides `Create`, `Get`, and `Transition`, along with robust unit tests (`job_test.go`).
- No changes or corrections to Feature 1 are needed; its API cleanly supports atomic transitions and snapshot isolation.
