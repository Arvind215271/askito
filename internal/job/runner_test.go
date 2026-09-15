package job

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJobRunner_AsyncSubmission(t *testing.T) {
	manager := NewManager()
	runner := NewRunner(manager)

	blockChan := make(chan struct{})
	startedChan := make(chan struct{})
	ownerID := "123e4567-e89b-12d3-a456-426614174000"

	job, err := runner.Submit(JobTypeExport, ownerID, func(ctx context.Context) error {
		close(startedChan)
		<-blockChan
		return nil
	})

	require.NoError(t, err)
	require.NotNil(t, job)
	assert.Equal(t, StatusQueued, job.Status)
	assert.Equal(t, ownerID, job.OwnerID)

	// Wait for background execution to start
	select {
	case <-startedChan:
		// Verify manager shows running status and preserves OwnerID now
		j, err := manager.Get(job.ID)
		require.NoError(t, err)
		assert.Equal(t, StatusRunning, j.Status)
		assert.Equal(t, ownerID, j.OwnerID)
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for job to start")
	}

	// Unblock work
	close(blockChan)

	// Wait for completion
	require.Eventually(t, func() bool {
		j, err := manager.Get(job.ID)
		return err == nil && j.Status == StatusCompleted && j.OwnerID == ownerID
	}, 2*time.Second, 10*time.Millisecond)
}

func TestJobRunner_SuccessfulWork(t *testing.T) {
	manager := NewManager()
	runner := NewRunner(manager)

	done := make(chan struct{})
	ownerID := "223e4567-e89b-12d3-a456-426614174000"
	job, err := runner.Submit(JobTypeTranscript, ownerID, func(ctx context.Context) error {
		close(done)
		return nil
	})
	require.NoError(t, err)
	assert.Equal(t, ownerID, job.OwnerID)

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for work")
	}

	require.Eventually(t, func() bool {
		j, err := manager.Get(job.ID)
		return err == nil && j.Status == StatusCompleted && j.StartedAt != nil && j.FinishedAt != nil && j.OwnerID == ownerID
	}, 2*time.Second, 10*time.Millisecond)

	j, err := manager.Get(job.ID)
	require.NoError(t, err)
	assert.Equal(t, StatusCompleted, j.Status)
	assert.NotEmpty(t, j.ID)
	assert.Equal(t, JobTypeTranscript, j.Type)
	assert.Equal(t, ownerID, j.OwnerID)
}

func TestJobRunner_FailedWork(t *testing.T) {
	manager := NewManager()
	runner := NewRunner(manager)

	expectedErr := errors.New("processing failed")
	ownerID := "323e4567-e89b-12d3-a456-426614174000"
	job, err := runner.Submit(JobTypeSubtitle, ownerID, func(ctx context.Context) error {
		return expectedErr
	})
	require.NoError(t, err)
	assert.Equal(t, ownerID, job.OwnerID)

	require.Eventually(t, func() bool {
		j, err := manager.Get(job.ID)
		return err == nil && j.Status == StatusFailed && j.OwnerID == ownerID
	}, 2*time.Second, 10*time.Millisecond)

	j, err := manager.Get(job.ID)
	require.NoError(t, err)
	assert.Equal(t, StatusFailed, j.Status)
	assert.Equal(t, expectedErr.Error(), j.Error)
	assert.NotNil(t, j.FinishedAt)
	assert.Equal(t, ownerID, j.OwnerID)
}

func TestJobRunner_IndependentContext(t *testing.T) {
	manager := NewManager()
	runner := NewRunner(manager)

	ctxCh := make(chan context.Context, 1)
	ownerID := "423e4567-e89b-12d3-a456-426614174000"

	job, err := runner.Submit(JobTypeExport, ownerID, func(ctx context.Context) error {
		ctxCh <- ctx
		return nil
	})
	require.NoError(t, err)
	assert.Equal(t, ownerID, job.OwnerID)

	select {
	case receivedCtx := <-ctxCh:
		require.NotNil(t, receivedCtx)
		select {
		case <-receivedCtx.Done():
			t.Fatal("independent background context should not be canceled")
		default:
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for context capture")
	}

	j, err := manager.Get(job.ID)
	require.NoError(t, err)
	assert.Equal(t, ownerID, j.OwnerID)
}

func TestJobRunner_PanicRecovery(t *testing.T) {
	manager := NewManager()
	runner := NewRunner(manager)
	ownerID := "523e4567-e89b-12d3-a456-426614174000"

	job, err := runner.Submit(JobTypeExport, ownerID, func(ctx context.Context) error {
		panic("unexpected system fault")
	})
	require.NoError(t, err)
	assert.Equal(t, ownerID, job.OwnerID)

	require.Eventually(t, func() bool {
		j, err := manager.Get(job.ID)
		return err == nil && j.Status == StatusFailed && j.OwnerID == ownerID
	}, 2*time.Second, 10*time.Millisecond)

	j, err := manager.Get(job.ID)
	require.NoError(t, err)
	assert.Equal(t, StatusFailed, j.Status)
	assert.Contains(t, j.Error, "job panicked")
	assert.Contains(t, j.Error, "unexpected system fault")
	assert.NotNil(t, j.FinishedAt)
	assert.Equal(t, ownerID, j.OwnerID)
}

func TestJobRunner_MultipleConcurrentJobs(t *testing.T) {
	manager := NewManager()
	runner := NewRunner(manager)

	const numJobs = 20
	jobIDs := make([]string, numJobs)
	ownerIDs := make([]string, numJobs)
	startedCh := make(chan struct{}, numJobs)
	releaseCh := make(chan struct{})
	doneCh := make(chan struct{}, numJobs)

	for i := 0; i < numJobs; i++ {
		idx := i
		ownerIDs[i] = fmt.Sprintf("owner-uuid-%d", i)
		job, err := runner.Submit(JobTypeExport, ownerIDs[i], func(ctx context.Context) error {
			startedCh <- struct{}{}
			<-releaseCh
			defer func() {
				doneCh <- struct{}{}
			}()
			if idx%3 == 0 {
				return errors.New("even-index failure")
			}
			return nil
		})
		require.NoError(t, err)
		assert.Equal(t, ownerIDs[i], job.OwnerID)
		jobIDs[i] = job.ID
	}

	// Wait for all jobs to start
	for i := 0; i < numJobs; i++ {
		select {
		case <-startedCh:
		case <-time.After(2 * time.Second):
			t.Fatal("timed out waiting for job to start")
		}
	}

	// Release all jobs
	close(releaseCh)

	// Wait for all jobs to finish work
	for i := 0; i < numJobs; i++ {
		select {
		case <-doneCh:
		case <-time.After(3 * time.Second):
			t.Fatal("timed out waiting for job completion")
		}
	}

	// Verify uniqueness, terminal states, status correctness, owner ID preservation, timestamps from main goroutine
	seenIDs := make(map[string]bool)
	for idx, id := range jobIDs {
		assert.False(t, seenIDs[id], "job IDs should be unique")
		seenIDs[id] = true

		var j *Job
		require.Eventually(t, func() bool {
			var err error
			j, err = manager.Get(id)
			return err == nil && j.Status.IsTerminal()
		}, 2*time.Second, time.Millisecond)

		assert.NotNil(t, j.StartedAt)
		assert.NotNil(t, j.FinishedAt)
		assert.Equal(t, ownerIDs[idx], j.OwnerID)

		if idx%3 == 0 {
			assert.Equal(t, StatusFailed, j.Status)
		} else {
			assert.Equal(t, StatusCompleted, j.Status)
		}
	}
}

func TestJobRunner_RunningTransitionBeforeWork(t *testing.T) {
	manager := NewManager()
	runner := NewRunner(manager)

	jobIDCh := make(chan string, 1)
	statusCh := make(chan JobStatus, 1)
	done := make(chan struct{})
	ownerID := "623e4567-e89b-12d3-a456-426614174000"

	job, err := runner.Submit(JobTypeExport, ownerID, func(ctx context.Context) error {
		jobID := <-jobIDCh
		j, err := manager.Get(jobID)
		if err == nil {
			statusCh <- j.Status
			assert.Equal(t, ownerID, j.OwnerID)
		} else {
			statusCh <- StatusQueued
		}
		close(done)
		return nil
	})
	require.NoError(t, err)
	assert.Equal(t, ownerID, job.OwnerID)

	// Send job ID after Submit() has returned and assigned job
	jobIDCh <- job.ID

	select {
	case status := <-statusCh:
		assert.Equal(t, StatusRunning, status)
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for work execution check")
	}

	select {
	case <-done:
	case <-time.After(2 * time.Second):
	}

	j, err := manager.Get(job.ID)
	require.NoError(t, err)
	assert.Equal(t, ownerID, j.OwnerID)
}
