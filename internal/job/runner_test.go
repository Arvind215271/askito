package job

import (
	"context"
	"errors"
	"sync"
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

	job, err := runner.Submit(JobTypeExport, func(ctx context.Context) error {
		close(startedChan)
		<-blockChan
		return nil
	})

	require.NoError(t, err)
	require.NotNil(t, job)
	assert.Equal(t, StatusQueued, job.Status)

	// Wait for background execution to start
	select {
	case <-startedChan:
		// Verify manager shows running status now
		j, err := manager.Get(job.ID)
		require.NoError(t, err)
		assert.Equal(t, StatusRunning, j.Status)
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for job to start")
	}

	// Unblock work
	close(blockChan)

	// Wait for completion
	require.Eventually(t, func() bool {
		j, err := manager.Get(job.ID)
		return err == nil && j.Status == StatusCompleted
	}, 2*time.Second, 10*time.Millisecond)
}

func TestJobRunner_SuccessfulWork(t *testing.T) {
	manager := NewManager()
	runner := NewRunner(manager)

	done := make(chan struct{})
	job, err := runner.Submit(JobTypeTranscript, func(ctx context.Context) error {
		close(done)
		return nil
	})
	require.NoError(t, err)

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for work")
	}

	require.Eventually(t, func() bool {
		j, err := manager.Get(job.ID)
		return err == nil && j.Status == StatusCompleted && j.StartedAt != nil && j.FinishedAt != nil
	}, 2*time.Second, 10*time.Millisecond)

	j, err := manager.Get(job.ID)
	require.NoError(t, err)
	assert.Equal(t, StatusCompleted, j.Status)
	assert.NotEmpty(t, j.ID)
	assert.Equal(t, JobTypeTranscript, j.Type)
}

func TestJobRunner_FailedWork(t *testing.T) {
	manager := NewManager()
	runner := NewRunner(manager)

	expectedErr := errors.New("processing failed")
	job, err := runner.Submit(JobTypeSubtitle, func(ctx context.Context) error {
		return expectedErr
	})
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		j, err := manager.Get(job.ID)
		return err == nil && j.Status == StatusFailed
	}, 2*time.Second, 10*time.Millisecond)

	j, err := manager.Get(job.ID)
	require.NoError(t, err)
	assert.Equal(t, StatusFailed, j.Status)
	assert.Equal(t, expectedErr.Error(), j.Error)
	assert.NotNil(t, j.FinishedAt)
}

func TestJobRunner_IndependentContext(t *testing.T) {
	manager := NewManager()
	runner := NewRunner(manager)

	var receivedCtx context.Context
	done := make(chan struct{})

	_, err := runner.Submit(JobTypeExport, func(ctx context.Context) error {
		receivedCtx = ctx
		close(done)
		return nil
	})
	require.NoError(t, err)

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for context capture")
	}

	require.NotNil(t, receivedCtx)
	// Verify context is not already cancelled
	select {
	case <-receivedCtx.Done():
		t.Fatal("independent background context should not be canceled")
	default:
	}
}

func TestJobRunner_PanicRecovery(t *testing.T) {
	manager := NewManager()
	runner := NewRunner(manager)

	job, err := runner.Submit(JobTypeExport, func(ctx context.Context) error {
		panic("unexpected system fault")
	})
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		j, err := manager.Get(job.ID)
		return err == nil && j.Status == StatusFailed
	}, 2*time.Second, 10*time.Millisecond)

	j, err := manager.Get(job.ID)
	require.NoError(t, err)
	assert.Equal(t, StatusFailed, j.Status)
	assert.Contains(t, j.Error, "job panicked")
	assert.Contains(t, j.Error, "unexpected system fault")
	assert.NotNil(t, j.FinishedAt)
}

func TestJobRunner_MultipleConcurrentJobs(t *testing.T) {
	manager := NewManager()
	runner := NewRunner(manager)

	const numJobs = 20
	var wg sync.WaitGroup
	wg.Add(numJobs)

	jobIDs := make([]string, numJobs)
	var mu sync.Mutex

	for i := 0; i < numJobs; i++ {
		idx := i
		job, err := runner.Submit(JobTypeExport, func(ctx context.Context) error {
			time.Sleep(10 * time.Millisecond)
			if idx%3 == 0 {
				return errors.New("even-index failure")
			}
			return nil
		})
		require.NoError(t, err)
		mu.Lock()
		jobIDs[idx] = job.ID
		mu.Unlock()
	}

	// Verify all jobs finish and reach terminal state
	for _, id := range jobIDs {
		go func(jobID string) {
			defer wg.Done()
			require.Eventually(t, func() bool {
				j, err := manager.Get(jobID)
				return err == nil && j.Status.IsTerminal()
			}, 3*time.Second, 10*time.Millisecond)
		}(id)
	}

	wg.Wait()

	// Verify uniqueness and states
	seenIDs := make(map[string]bool)
	for idx, id := range jobIDs {
		assert.False(t, seenIDs[id], "job IDs should be unique")
		seenIDs[id] = true

		j, err := manager.Get(id)
		require.NoError(t, err)
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

	statusAtStartChan := make(chan JobStatus, 1)
	done := make(chan struct{})

	var jobID string
	job, err := runner.Submit(JobTypeExport, func(ctx context.Context) error {
		// Check job status in manager at the very beginning of work execution
		j, err := manager.Get(jobID)
		if err == nil {
			statusAtStartChan <- j.Status
		} else {
			statusAtStartChan <- StatusQueued
		}
		close(done)
		return nil
	})
	require.NoError(t, err)
	jobID = job.ID

	select {
	case status := <-statusAtStartChan:
		assert.Equal(t, StatusRunning, status)
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for work execution check")
	}

	select {
	case <-done:
	case <-time.After(2 * time.Second):
	}
}
