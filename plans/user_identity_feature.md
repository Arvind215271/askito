# Architectural Design & Implementation Plan: Feature 4 (User Identity for Job System)

## 1. Current System Understanding
- **`internal/api/job/`**: Contains job API routes (`routes.go`), handlers (`handler.go`), and error definitions (`errors.go`).
- **`internal/api/app_error.go`**: Defines centralized `api.AppError` structure with codes, messages, HTTP status codes, and optional fields.
- **Echo v5**: Used for routing, groups, middleware, and request context handling.
- **`github.com/google/uuid`**: Used across the domain (e.g. `job.Manager`) for UUID generation and validation.

---

## 2. Architectural Design of Feature 4 (User Identity)

### A. Components & Package Structure (`internal/api/job/`)
1. **`internal/api/job/errors.go`**: Add `MissingOrInvalidUserID()` error constructor returning `api.AppError` with HTTP status 400 (`BAD_REQUEST` or `INVALID_USER_ID`).
2. **`internal/api/job/middleware.go`**: Implement:
   - Request-level Echo middleware `UserIdentityMiddleware()` that:
     - Reads `X-User-ID` header from request.
     - Validates that the header value is a valid UUID using `github.com/google/uuid`.
     - If missing or invalid, returns `api.AppError` (HTTP 400).
     - If valid, injects the validated user ID into the Go request context (`c.Request().WithContext(...)`) using an unexported context key type to prevent collisions.
   - Public helper function `UserIDFromContext(ctx context.Context) (string, bool)` (or returning `string`) to retrieve the user identifier from context.
3. **`main.go`**: Register the user identity middleware on the job route group (`e.Group("/jobs", apiJob.UserIdentityMiddleware())`).
4. **`internal/api/job/handler_test.go`**: Update existing job route tests to supply a valid `X-User-ID` header, and add new unit tests covering:
   - Valid `X-User-ID` header (passes middleware and propagates to context).
   - Missing `X-User-ID` header (returns HTTP 400).
   - Invalid UUID `X-User-ID` header (returns HTTP 400).
   - Context helper `UserIDFromContext` behavior.

---

## 3. Implementation Tasks & Acceptance Criteria

### Task 1: Update `internal/api/job/errors.go`
- **Objective**: Add `MissingOrInvalidUserID` error constructor.
- **Affected Files**: `internal/api/job/errors.go`
- **Acceptance Criteria**:
  - Returns `api.AppError` with HTTP status 400 and code `INVALID_USER_ID`.

### Task 2: Create `internal/api/job/middleware.go`
- **Objective**: Implement user identity middleware and context helper.
- **Affected Files**: `internal/api/job/middleware.go`
- **Acceptance Criteria**:
  - Validates `X-User-ID` header using `uuid.Parse`.
  - Rejects missing or malformed headers with 400 Bad Request.
  - Stores validated ID in request context using package-private key.
  - Provides public `UserIDFromContext(ctx)` helper function.

### Task 3: Wire Middleware in `main.go`
- **Objective**: Apply `UserIdentityMiddleware()` to job routes.
- **Affected Files**: `main.go`
- **Acceptance Criteria**:
  - Job routes group `/jobs` enforces user identity middleware.

### Task 4: Add Unit Tests in `internal/api/job/handler_test.go` (or `middleware_test.go`)
- **Objective**: Test valid identity, missing identity, invalid identity, and context helper.
- **Affected Files**: `internal/api/job/handler_test.go` (and/or `middleware_test.go`)
- **Acceptance Criteria**:
  - All test cases pass successfully.

---

## 4. Verification Plan
- Run `go test ./internal/api/job/... -v -race` to verify unit tests pass.
