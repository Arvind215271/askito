package job

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"

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
	RegisterRoutes(e.Group("/jobs"), h)

	// Create jobs in various states
	// 1. Queued job
	queuedJob, err := manager.Create(domainJob.JobTypeExport)
	if err != nil {
		t.Fatalf("failed to create queued job: %v", err)
	}

	// 2. Running job
	runningJob, err := manager.Create(domainJob.JobTypeSubtitle)
	if err != nil {
		t.Fatalf("failed to create running job: %v", err)
	}
	_, err = manager.Transition(runningJob.ID, domainJob.StatusRunning, nil)
	if err != nil {
		t.Fatalf("failed to transition to running: %v", err)
	}
	runningJob, _ = manager.Get(runningJob.ID)

	// 3. Completed job
	completedJob, err := manager.Create(domainJob.JobTypeTranscript)
	if err != nil {
		t.Fatalf("failed to create completed job: %v", err)
	}
	_, err = manager.Transition(completedJob.ID, domainJob.StatusRunning, nil)
	if err != nil {
		t.Fatalf("failed to transition to running: %v", err)
	}
	_, err = manager.Transition(completedJob.ID, domainJob.StatusCompleted, nil)
	if err != nil {
		t.Fatalf("failed to transition to completed: %v", err)
	}
	completedJob, _ = manager.Get(completedJob.ID)

	// 4. Failed job
	failedJob, err := manager.Create(domainJob.JobTypeExport)
	if err != nil {
		t.Fatalf("failed to create failed job: %v", err)
	}
	_, err = manager.Transition(failedJob.ID, domainJob.StatusRunning, nil)
	if err != nil {
		t.Fatalf("failed to transition to running: %v", err)
	}
	_, err = manager.Transition(failedJob.ID, domainJob.StatusFailed, errors.New("processing failed"))
	if err != nil {
		t.Fatalf("failed to transition to failed: %v", err)
	}
	failedJob, _ = manager.Get(failedJob.ID)

	// 5. Cancelled job
	cancelledJob, err := manager.Create(domainJob.JobTypeExport)
	if err != nil {
		t.Fatalf("failed to create cancelled job: %v", err)
	}
	_, err = manager.Transition(cancelledJob.ID, domainJob.StatusCancelled, nil)
	if err != nil {
		t.Fatalf("failed to transition to cancelled: %v", err)
	}
	cancelledJob, _ = manager.Get(cancelledJob.ID)

	tests := []struct {
		name           string
		jobID          string
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name:           "Queued Job",
			jobID:          queuedJob.ID,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var j domainJob.Job
				if err := json.Unmarshal(body, &j); err != nil {
					t.Fatalf("failed to unmarshal: %v", err)
				}
				if j.Status != domainJob.StatusQueued {
					t.Errorf("expected status queued, got %s", j.Status)
				}
			},
		},
		{
			name:           "Running Job with started_at",
			jobID:          runningJob.ID,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var j domainJob.Job
				if err := json.Unmarshal(body, &j); err != nil {
					t.Fatalf("failed to unmarshal: %v", err)
				}
				if j.Status != domainJob.StatusRunning {
					t.Errorf("expected status running, got %s", j.Status)
				}
				if j.StartedAt == nil {
					t.Error("expected started_at to be set")
				}
			},
		},
		{
			name:           "Completed Job with finished_at",
			jobID:          completedJob.ID,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var j domainJob.Job
				if err := json.Unmarshal(body, &j); err != nil {
					t.Fatalf("failed to unmarshal: %v", err)
				}
				if j.Status != domainJob.StatusCompleted {
					t.Errorf("expected status completed, got %s", j.Status)
				}
				if j.FinishedAt == nil {
					t.Error("expected finished_at to be set")
				}
			},
		},
		{
			name:           "Failed Job with error and finished_at",
			jobID:          failedJob.ID,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var j domainJob.Job
				if err := json.Unmarshal(body, &j); err != nil {
					t.Fatalf("failed to unmarshal: %v", err)
				}
				if j.Status != domainJob.StatusFailed {
					t.Errorf("expected status failed, got %s", j.Status)
				}
				if j.FinishedAt == nil {
					t.Error("expected finished_at to be set")
				}
				if j.Error == "" {
					t.Error("expected error message to be set")
				}
			},
		},
		{
			name:           "Cancelled Job with finished_at",
			jobID:          cancelledJob.ID,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var j domainJob.Job
				if err := json.Unmarshal(body, &j); err != nil {
					t.Fatalf("failed to unmarshal: %v", err)
				}
				if j.Status != domainJob.StatusCancelled {
					t.Errorf("expected status cancelled, got %s", j.Status)
				}
				if j.FinishedAt == nil {
					t.Error("expected finished_at to be set")
				}
			},
		},
		{
			name:           "Unknown Job ID",
			jobID:          "non-existent-id",
			expectedStatus: http.StatusNotFound,
			checkResponse:  func(t *testing.T, body []byte) {},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/jobs/"+tc.jobID, nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			if rec.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d (body: %s)", tc.expectedStatus, rec.Code, rec.Body.String())
			}

			if rec.Code == http.StatusOK {
				tc.checkResponse(t, rec.Body.Bytes())
			}
		})
	}
}
