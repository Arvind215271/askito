package job

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJobCreation(t *testing.T) {
	manager := NewManager()

	job, err := manager.Create(JobTypeExport)
	require.NoError(t, err)
	require.NotNil(t, job)

	assert.NotEmpty(t, job.ID)
	assert.Equal(t, JobTypeExport, job.Type)
	assert.Equal(t, StatusQueued, job.Status)
	assert.False(t, job.CreatedAt.IsZero())
	assert.Equal(t, time.UTC, job.CreatedAt.Location())
	assert.Nil(t, job.StartedAt)
	assert.Nil(t, job.FinishedAt)
	assert.Empty(t, job.Error)
}

func TestJobStateTransitions_HappyPaths(t *testing.T) {
	tests := []struct {
		name  string
		steps []struct {
			target JobStatus
			hasErr string
		}
	}{
		{
			name: "Queued -> Running -> Completed",
			steps: []struct {
				target JobStatus
				hasErr string
			}{
				{target: StatusRunning},
				{target: StatusCompleted},
			},
		},
		{
			name: "Queued -> Running -> Failed",
			steps: []struct {
				target JobStatus
				hasErr string
			}{
				{target: StatusRunning},
				{target: StatusFailed, hasErr: "task failed"},
			},
		},
		{
			name: "Queued -> Running -> Cancelled",
			steps: []struct {
				target JobStatus
				hasErr string
			}{
				{target: StatusRunning},
				{target: StatusCancelled},
			},
		},
		{
			name: "Queued -> Cancelled",
			steps: []struct {
				target JobStatus
				hasErr string
			}{
				{target: StatusCancelled},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewManager()
			job, err := manager.Create(JobTypeTranscript)
			require.NoError(t, err)

			prevTime := job.CreatedAt

			for _, step := range tt.steps {
				var transErr error
				if step.hasErr != "" {
					transErr = errors.New(step.hasErr)
				}

				updated, err := manager.Transition(job.ID, step.target, transErr)
				require.NoError(t, err)
				require.NotNil(t, updated)

				assert.Equal(t, step.target, updated.Status)

				if step.target == StatusRunning {
					require.NotNil(t, updated.StartedAt)
					assert.Equal(t, time.UTC, updated.StartedAt.Location())
					assert.True(t, updated.StartedAt.Equal(prevTime) || updated.StartedAt.After(prevTime))
					prevTime = *updated.StartedAt
				}

				if step.target.IsTerminal() {
					require.NotNil(t, updated.FinishedAt)
					assert.Equal(t, time.UTC, updated.FinishedAt.Location())
					assert.True(t, updated.FinishedAt.Equal(prevTime) || updated.FinishedAt.After(prevTime))
					if step.hasErr != "" {
						assert.Equal(t, step.hasErr, updated.Error)
					}
				}

				job = updated
			}
		})
	}
}

func TestJobStateTransitions_InvalidTransitions(t *testing.T) {
	manager := NewManager()

	job, err := manager.Create(JobTypeSubtitle)
	require.NoError(t, err)

	// Invalid jump from Queued -> Completed
	updated, err := manager.Transition(job.ID, StatusCompleted, nil)
	assert.ErrorIs(t, err, ErrInvalidTransition)
	assert.Nil(t, updated)

	// Transition Queued -> Running (valid)
	runningJob, err := manager.Transition(job.ID, StatusRunning, nil)
	require.NoError(t, err)
	assert.Equal(t, StatusRunning, runningJob.Status)

	// Invalid jump from Running -> Queued
	updated, err = manager.Transition(job.ID, StatusQueued, nil)
	assert.ErrorIs(t, err, ErrInvalidTransition)
	assert.Nil(t, updated)

	// Transition Running -> Completed (valid terminal)
	completedJob, err := manager.Transition(job.ID, StatusCompleted, nil)
	require.NoError(t, err)
	assert.Equal(t, StatusCompleted, completedJob.Status)

	// Terminal state locking: try transitioning Completed -> Running
	updated, err = manager.Transition(job.ID, StatusRunning, nil)
	assert.ErrorIs(t, err, ErrInvalidTransition)
	assert.Nil(t, updated)

	// Terminal state locking: try transitioning Completed -> Cancelled
	updated, err = manager.Transition(job.ID, StatusCancelled, nil)
	assert.ErrorIs(t, err, ErrInvalidTransition)
	assert.Nil(t, updated)
}

func TestJobManager_NotFound(t *testing.T) {
	manager := NewManager()

	job, err := manager.Get("non-existent-id")
	assert.ErrorIs(t, err, ErrJobNotFound)
	assert.Nil(t, job)

	updated, err := manager.Transition("non-existent-id", StatusRunning, nil)
	assert.ErrorIs(t, err, ErrJobNotFound)
	assert.Nil(t, updated)
}

func TestJobManager_SnapshotIsolation(t *testing.T) {
	manager := NewManager()
	job, err := manager.Create(JobTypeExport)
	require.NoError(t, err)

	// Modify the returned snapshot directly
	job.Status = StatusCompleted
	job.Error = "tampered"

	// Retrieve again from manager and verify it remains unchanged in store
	retrieved, err := manager.Get(job.ID)
	require.NoError(t, err)
	assert.Equal(t, StatusQueued, retrieved.Status)
	assert.Empty(t, retrieved.Error)
}
