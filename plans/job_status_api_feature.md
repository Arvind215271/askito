# Architectural Design & Implementation Plan: Feature 3 (Job Status API)

## 1. Current System Understanding

- **`internal/job/model.go`**: Defines `Job` struct (`id`, `type`, `status`, `created_at`, `started_at`, `finished_at`, `error`), `JobStatus` values (`queued`, `running`, `completed`, `failed`, `cancelled`), and sentinel error `ErrJobNotFound`.
- **`internal/job/manager.go`**: Provides thread-safe `JobManager` with methods `Create`, `Get(id)`, and `Transition`.
- **`internal/api/`**: Echo v5 framework usage with centralized `api.AppError`, `api.NewError`, and `api.ErrorHandler` converting `AppError` to JSON error responses with appropriate HTTP status codes (e.g. 404 for not found).
- **Existing API packages**: Follow a consistent pattern of `internal/api/<domain>/` containing `handler.go`, `routes.go`, `errors.go` (or inline errors), and `handler_test.go`.

---

## 2. Architectural Design of Feature 3 (Job Status API)

### A. Package Structure (`internal/api/job/`)
1. **`internal/api/job/errors.go`**: Defines API-level error constructors wrapping domain errors (e.g. mapping `job.ErrJobNotFound` to `JOB_NOT_FOUND` with HTTP 404 status).
2. **`internal/api/job/handler.go`**: Implements HTTP handler taking `job.JobManager`, exposing `Get(c *echo.Context) error`.
3. **`internal/api/job/routes.go`**: Registers route `GET /:id` under the `/jobs` route group.
4. **`internal/api/job/handler_test.go`**: Unit/integration tests using Echo test context and real `job.JobManager` instances.

### B. Request & Response Flow
- **Endpoint**: `GET /jobs/:id`
- **Path Parameter**: `id` (string UUID of the job).
- **Execution flow**:
  1. Extract `:id` from route parameters.
  2. Call `jobManager.Get(id)`.
  3. If not found (`errors.Is(err, job.ErrJobNotFound)`), return `Err.JobNotFound()` (HTTP 404).
  4. On success, return `c.JSON(http.StatusOK, job)` with JSON serialization of `job.Job`.

### C. Wiring in `main.go`
- Instantiate `jobManager := job.NewManager()` (or share instance if used by runner/services).
- Instantiate `jobHandler := job.NewHandler(jobManager)`.
- Register routes: `job.RegisterRoutes(e.Group("/jobs"), jobHandler)`.

---

## 3. Implementation Tasks & Acceptance Criteria

### Task 1: Create `internal/api/job/errors.go`
- **Objective**: Define package-level error mapper for job domain errors.
- **Affected Files**: `internal/api/job/errors.go`
- **Acceptance Criteria**:
  - `JobNotFound` error returns status 404 and code `JOB_NOT_FOUND`.

### Task 2: Create `internal/api/job/handler.go`
- **Objective**: Implement `Handler` and `Get` method.
- **Affected Files**: `internal/api/job/handler.go`
- **Acceptance Criteria**:
  - Successfully retrieves jobs by ID and returns HTTP 200 with JSON payload.
  - Returns HTTP 404 when job does not exist.

### Task 3: Create `internal/api/job/routes.go`
- **Objective**: Register `GET /:id` on Echo group.
- **Affected Files**: `internal/api/job/routes.go`
- **Acceptance Criteria**:
  - Group route successfully maps `GET /jobs/:id` to `Get`.

### Task 4: Create `internal/api/job/handler_test.go`
- **Objective**: Write comprehensive HTTP tests covering all job states.
- **Affected Files**: `internal/api/job/handler_test.go`
- **Acceptance Criteria**:
  - Tests verify responses for:
    - `queued` job
    - `running` job
    - `completed` job
    - `failed` job
    - `cancelled` job
    - unknown / non-existent job (HTTP 404)

### Task 5: Wire Job API in `main.go`
- **Objective**: Instantiate `JobManager`, `Handler`, and register routes in `main.go`.
- **Affected Files**: `main.go`
- **Acceptance Criteria**:
  - Server starts successfully with `/jobs/:id` endpoint available.

---

## 4. Verification Plan
- Run `go test ./internal/api/job/... -v -race` to verify unit tests pass under race detection.
