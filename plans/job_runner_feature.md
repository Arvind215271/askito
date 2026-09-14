# Architectural Design & Implementation Plan: Feature 2 (Asynchronous Job Processing & JobRunner)

## 1. Current System Understanding

- **`internal/job/model.go`**: Defines `JobType` (export, subtitle, transcript), `JobStatus` (`queued`, `running`, `completed`, `failed`, `cancelled`), `Job` struct (ID, Type, Status, CreatedAt, StartedAt, FinishedAt, Error), and sentinel errors (`ErrJobNotFound`, `ErrInvalidTransition`).
- **`internal/job/manager.go`**: Provides thread-safe `JobManager` using `sync.RWMutex` and `map[string]Job`. Supports `Create`, `Get`, and `Transition` with strict state transition validation (`validateTransition`).
- **`internal/job/runner.go`**: Implements `JobRunner`, bound to `JobManager` and optional `logger.Logger`.
  - `Submit(jobType JobType, work func(context.Context) error) (*Job, error)` creates a job in `queued` state and spawns a background goroutine (`go r.execute(job.ID, work)`).
  - Background execution flow:
    1. Transitions job from `queued` to `running` (recording `StartedAt`).
    2. Sets up an independent background context (`context.Background()`).
    3. Implements **defer panic recovery** to catch any goroutine panics, convert them to `fmt.Errorf("job panicked: %v", p)`, log them, and transition the job to `failed`.
    4. Executes the user work function. If it returns an error, transitions to `failed` with that error; otherwise, transitions to `completed`.
- **Tests**: `job_test.go` and `manager_test.go` thoroughly test job creation, valid/invalid state transitions, concurrency safety, and snapshot isolation.

---

## 2. Architectural Design of JobRunner

### A. Core Components & Responsibilities
- **`JobManager`**: In-memory store and state machine guarantor. Ensures jobs transition through valid paths without data races.
- **`JobRunner`**: Orchestrator of asynchronous execution. Decouples job initiation from HTTP/API response lifecycle by spawning background goroutines.

### B. API Specification
```go
type JobRunner struct {
	manager *JobManager
	logger  *logger.Logger
}

func NewRunner(manager *JobManager, logger ...*logger.Logger) *JobRunner
func (r *JobRunner) Submit(jobType JobType, work func(context.Context) error) (*Job, error)
```

### C. Lifecycle & Execution Flow
1. **Queueing**: Caller calls `Submit(jobType, work)`. `JobManager.Create` assigns a UUID, sets status to `queued`, and stores it.
2. **Dispatch**: A background goroutine (`go r.execute(...)`) is spawned immediately. `Submit` returns the initial queued `Job` snapshot to the caller right away.
3. **Execution & State Transition**:
   - Background worker calls `manager.Transition(jobID, StatusRunning, nil)`.
   - Executes `work(ctx)` with panic recovery wrapped around it.
   - If work succeeds, calls `manager.Transition(jobID, StatusCompleted, nil)`.
   - If work errors, calls `manager.Transition(jobID, StatusFailed, err)`.
   - If work panics, recover catches panic, formats error, and calls `manager.Transition(jobID, StatusFailed, panicErr)`.

---

## 3. Proposed Testing Strategy (`internal/job/runner_test.go`)

To ensure robust asynchronous execution, panic safety, and concurrency, we will add `internal/job/runner_test.go` covering:
1. **Successful Job Execution**: Submitting a job, waiting for completion via polling or synchronization, and verifying status becomes `completed` with `StartedAt` and `FinishedAt` populated.
2. **Failing Job Execution**: Submitting a job where `work` returns an error, verifying status becomes `failed` and `Error` field is populated.
3. **Panic Recovery**: Submitting a job where `work` panics (`panic("something went wrong")`), verifying panic is caught, status becomes `failed`, and error message indicates panic.
4. **Concurrent Job Runs**: Submitting multiple concurrent jobs via `JobRunner` under `-race` detector.

---

## 4. Implementation & Verification Tasks

1. Create `internal/job/runner_test.go` to thoroughly test `JobRunner` success, failure, panic recovery, and concurrency.
2. Run `go test -v -race ./internal/job/...` to verify all job package unit tests pass successfully.
