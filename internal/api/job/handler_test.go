package job

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Arvind215271/askito/internal/api"
	domainJob "github.com/Arvind215271/askito/internal/job"
	"github.com/Arvind215271/askito/internal/logger"
)

func TestJobHandler_Get(t *testing.T) {
	manager := domainJob.NewManager()
	h := NewHandler(manager)

	log := logger.New("test")
	errorHandler := api.NewErrorHandler(log)

	e := echo.New()
	e.HTTPErrorHandler = errorHandler.Handle
	RegisterRoutes(e.Group("/jobs", UserIdentityMiddleware()), h)

	// Create jobs in various states
	// 1. Queued job
	queuedJob, err := manager.Create(domainJob.JobTypeExport)
	require.NoError(t, err)

	// 2. Running job
	runningJob, err := manager.Create(domainJob.JobTypeSubtitle)
	require.NoError(t, err)
	_, err = manager.Transition(runningJob.ID, domainJob.StatusRunning, nil)
	require.NoError(t, err)
	runningJob, _ = manager.Get(runningJob.ID)

	// 3. Completed job
	completedJob, err := manager.Create(domainJob.JobTypeTranscript)
	require.NoError(t, err)
	_, err = manager.Transition(completedJob.ID, domainJob.StatusRunning, nil)
	require.NoError(t, err)
	_, err = manager.Transition(completedJob.ID, domainJob.StatusCompleted, nil)
	require.NoError(t, err)
	completedJob, _ = manager.Get(completedJob.ID)

	// 4. Failed job
	failedJob, err := manager.Create(domainJob.JobTypeExport)
	require.NoError(t, err)
	_, err = manager.Transition(failedJob.ID, domainJob.StatusRunning, nil)
	require.NoError(t, err)
	_, err = manager.Transition(failedJob.ID, domainJob.StatusFailed, errors.New("processing failed"))
	require.NoError(t, err)
	failedJob, _ = manager.Get(failedJob.ID)

	// 5. Cancelled job
	cancelledJob, err := manager.Create(domainJob.JobTypeExport)
	require.NoError(t, err)
	_, err = manager.Transition(cancelledJob.ID, domainJob.StatusCancelled, nil)
	require.NoError(t, err)
	cancelledJob, _ = manager.Get(cancelledJob.ID)

	validUserID := uuid.New().String()

	tests := []struct {
		name           string
		jobID          string
		userIDHeader   string
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name:           "Queued Job with valid user ID",
			jobID:          queuedJob.ID,
			userIDHeader:   validUserID,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var j domainJob.Job
				require.NoError(t, json.Unmarshal(body, &j))
				assert.Equal(t, domainJob.StatusQueued, j.Status)
			},
		},
		{
			name:           "Running Job with started_at",
			jobID:          runningJob.ID,
			userIDHeader:   validUserID,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var j domainJob.Job
				require.NoError(t, json.Unmarshal(body, &j))
				assert.Equal(t, domainJob.StatusRunning, j.Status)
				assert.NotNil(t, j.StartedAt)
			},
		},
		{
			name:           "Completed Job with finished_at",
			jobID:          completedJob.ID,
			userIDHeader:   validUserID,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var j domainJob.Job
				require.NoError(t, json.Unmarshal(body, &j))
				assert.Equal(t, domainJob.StatusCompleted, j.Status)
				assert.NotNil(t, j.FinishedAt)
			},
		},
		{
			name:           "Failed Job with error and finished_at",
			jobID:          failedJob.ID,
			userIDHeader:   validUserID,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var j domainJob.Job
				require.NoError(t, json.Unmarshal(body, &j))
				assert.Equal(t, domainJob.StatusFailed, j.Status)
				assert.NotNil(t, j.FinishedAt)
				assert.NotEmpty(t, j.Error)
			},
		},
		{
			name:           "Cancelled Job with finished_at",
			jobID:          cancelledJob.ID,
			userIDHeader:   validUserID,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var j domainJob.Job
				require.NoError(t, json.Unmarshal(body, &j))
				assert.Equal(t, domainJob.StatusCancelled, j.Status)
				assert.NotNil(t, j.FinishedAt)
			},
		},
		{
			name:           "Unknown Job ID",
			jobID:          "non-existent-id",
			userIDHeader:   validUserID,
			expectedStatus: http.StatusNotFound,
			checkResponse:  func(t *testing.T, body []byte) {},
		},
		{
			name:           "Missing X-User-ID header",
			jobID:          queuedJob.ID,
			userIDHeader:   "",
			expectedStatus: http.StatusBadRequest,
			checkResponse:  func(t *testing.T, body []byte) {},
		},
		{
			name:           "Invalid X-User-ID header (not UUID)",
			jobID:          queuedJob.ID,
			userIDHeader:   "not-a-valid-uuid",
			expectedStatus: http.StatusBadRequest,
			checkResponse:  func(t *testing.T, body []byte) {},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/jobs/"+tc.jobID, nil)
			if tc.userIDHeader != "" {
				req.Header.Set("X-User-ID", tc.userIDHeader)
			}
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code, "body: %s", rec.Body.String())

			if rec.Code == http.StatusOK {
				tc.checkResponse(t, rec.Body.Bytes())
			}
		})
	}
}

func TestUserIDFromContext(t *testing.T) {
	ctx := context.Background()
	_, ok := UserIDFromContext(ctx)
	assert.False(t, ok)

	uid := uuid.New().String()
	ctx = context.WithValue(ctx, userIDKey, uid)
	val, ok := UserIDFromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, uid, val)
}
