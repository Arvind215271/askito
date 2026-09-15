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
- **Tests**: `job_test.go`, `manager_test.go`, and `runner_test.go` thoroughly test job creation, valid/invalid state transitions, concurrency safety, panic recovery, async submission, independent context handling, running state transitions before work execution, and concurrent execution under race detection.

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

func NewRunner(manager *JobManager, loggers ...*logger.Logger) *JobRunner
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

## 3. Test Coverage & Synchronization (`internal/job/runner_test.go`)

Robust asynchronous test suite (`internal/job/runner_test.go`) covers:
1. **`TestJobRunner_AsyncSubmission`**: Verifies non-blocking submission and background execution starting and completing.
2. **`TestJobRunner_SuccessfulWork`**: Verifies successful execution, completion status, and populated timestamps (`StartedAt`, `FinishedAt`).
3. **`TestJobRunner_FailedWork`**: Verifies error handling and `failed` status transition when work returns an error.
4. **`TestJobRunner_PanicRecovery`**: Verifies panic recovery catching panics, converting them to descriptive failure messages, and updating status to `failed`.
5. **`TestJobRunner_IndependentContext`**: Uses `ctxCh` channel synchronization to verify the background context passed to work is independent and not canceled.
6. **`TestJobRunner_MultipleConcurrentJobs`**: Uses `startedCh`, `releaseCh`, and `doneCh` coordination channels with zero polling loops or `time.Sleep` delays. Verifies unique IDs, terminal states, correct status outcomes, and timestamps strictly from the main test goroutine without calling assertion helpers inside spawned goroutines.
7. **`TestJobRunner_RunningTransitionBeforeWork`**: Uses deterministic channel synchronization (`jobIDCh`, `statusCh`) so `jobID` is sent after `Submit()` returns, avoiding any race condition where work runs before `jobID` variable is assigned.

---

## 4. Implementation & Verification

- Formatted code using `gofmt -w internal/job/*.go`.
- Verified all tests in `internal/job/` pass successfully under race detector.
