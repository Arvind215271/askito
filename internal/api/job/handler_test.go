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

func TestJobHandler_GetUserJob(t *testing.T) {
	manager := domainJob.NewManager()
	h := NewHandler(manager)

	log := logger.New("test")
	errorHandler := api.NewErrorHandler(log)

	e := echo.New()
	e.HTTPErrorHandler = errorHandler.Handle
	RegisterRoutes(e.Group("/jobs", UserIdentityMiddleware()), h)

	userA := uuid.New().String()
	userB := uuid.New().String()

	// Create jobs in various states for userA
	queuedJob, err := manager.Create(domainJob.JobTypeExport, userA)
	require.NoError(t, err)

	runningJob, err := manager.Create(domainJob.JobTypeSubtitle, userA)
	require.NoError(t, err)
	_, err = manager.Transition(runningJob.ID, domainJob.StatusRunning, nil)
	require.NoError(t, err)
	runningJob, _ = manager.Get(runningJob.ID)

	completedJob, err := manager.Create(domainJob.JobTypeTranscript, userA)
	require.NoError(t, err)
	_, err = manager.Transition(completedJob.ID, domainJob.StatusRunning, nil)
	require.NoError(t, err)
	_, err = manager.Transition(completedJob.ID, domainJob.StatusCompleted, nil)
	require.NoError(t, err)
	completedJob, _ = manager.Get(completedJob.ID)

	failedJob, err := manager.Create(domainJob.JobTypeExport, userA)
	require.NoError(t, err)
	_, err = manager.Transition(failedJob.ID, domainJob.StatusRunning, nil)
	require.NoError(t, err)
	_, err = manager.Transition(failedJob.ID, domainJob.StatusFailed, errors.New("processing failed"))
	require.NoError(t, err)
	failedJob, _ = manager.Get(failedJob.ID)

	cancelledJob, err := manager.Create(domainJob.JobTypeExport, userA)
	require.NoError(t, err)
	_, err = manager.Transition(cancelledJob.ID, domainJob.StatusCancelled, nil)
	require.NoError(t, err)
	cancelledJob, _ = manager.Get(cancelledJob.ID)

	tests := []struct {
		name           string
		jobID          string
		userIDHeader   string
		expectedStatus int
		expectedCode   string
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name:           "Case 1: Owner can access their job successfully (HTTP 200)",
			jobID:          queuedJob.ID,
			userIDHeader:   userA,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var j domainJob.Job
				require.NoError(t, json.Unmarshal(body, &j))
				assert.Equal(t, domainJob.StatusQueued, j.Status)
				assert.Equal(t, userA, j.OwnerID)
			},
		},
		{
			name:           "Running Job with started_at",
			jobID:          runningJob.ID,
			userIDHeader:   userA,
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
			userIDHeader:   userA,
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
			userIDHeader:   userA,
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
			userIDHeader:   userA,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var j domainJob.Job
				require.NoError(t, json.Unmarshal(body, &j))
				assert.Equal(t, domainJob.StatusCancelled, j.Status)
				assert.NotNil(t, j.FinishedAt)
			},
		},
		{
			name:           "Case 2: Different user cannot access another user's job (HTTP 404, JOB_NOT_FOUND)",
			jobID:          queuedJob.ID,
			userIDHeader:   userB,
			expectedStatus: http.StatusNotFound,
			expectedCode:   "JOB_NOT_FOUND",
			checkResponse:  func(t *testing.T, body []byte) {},
		},
		{
			name:           "Case 3: Nonexistent job requested by user (HTTP 404, JOB_NOT_FOUND)",
			jobID:          uuid.New().String(),
			userIDHeader:   userA,
			expectedStatus: http.StatusNotFound,
			expectedCode:   "JOB_NOT_FOUND",
			checkResponse:  func(t *testing.T, body []byte) {},
		},
		{
			name:           "Case 4: Missing user ID header (HTTP 400, INVALID_USER_ID)",
			jobID:          queuedJob.ID,
			userIDHeader:   "",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "INVALID_USER_ID",
			checkResponse:  func(t *testing.T, body []byte) {},
		},
		{
			name:           "Case 5: Malformed user ID header (HTTP 400, INVALID_USER_ID)",
			jobID:          queuedJob.ID,
			userIDHeader:   "not-a-valid-uuid",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "INVALID_USER_ID",
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

			if tc.expectedCode != "" {
				var resp api.Response
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)
				if metaMap, ok := resp.Meta.(map[string]any); ok {
					assert.Equal(t, tc.expectedCode, metaMap["code"])
				} else {
					assert.Fail(t, "expected meta code in response")
				}
			}

			if rec.Code == http.StatusOK && tc.checkResponse != nil {
				tc.checkResponse(t, rec.Body.Bytes())
			}
		})
	}
}

func TestGetUserIDFromContext(t *testing.T) {
	ctx := context.Background()
	_, ok := GetUserIDFromContext(ctx)
	assert.False(t, ok)

	uid := uuid.New().String()
	ctx = context.WithValue(ctx, userIDKey, uid)
	val, ok := GetUserIDFromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, uid, val)
}
