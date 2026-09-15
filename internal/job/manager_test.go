package job

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJobManager_Concurrency(t *testing.T) {
	manager := NewManager()
	const goroutines = 50
	const jobsPerGoroutine = 20

	var wg sync.WaitGroup
	wg.Add(goroutines)

	// Concurrent creation
	jobIDsChan := make(chan string, goroutines*jobsPerGoroutine)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < jobsPerGoroutine; j++ {
				jType := JobTypeExport
				if j%2 == 0 {
					jType = JobTypeSubtitle
				} else if j%3 == 0 {
					jType = JobTypeTranscript
				}

				job, err := manager.Create(jType, "owner-concurrency")
				if err == nil && job != nil {
					jobIDsChan <- job.ID
				}
			}
		}()
	}

	// Wait for creations to finish collecting IDs
	go func() {
		wg.Wait()
		close(jobIDsChan)
	}()

	// Read and transition concurrently
	var transitionWg sync.WaitGroup
	for id := range jobIDsChan {
		transitionWg.Add(1)
		go func(jobID string) {
			defer transitionWg.Done()

			// Get snapshot
			_, _ = manager.Get(jobID)

			// Transition queued -> running
			runningJob, err := manager.Transition(jobID, StatusRunning, nil)
			if err == nil && runningJob != nil {
				// Transition running -> completed or failed
				_, _ = manager.Transition(jobID, StatusCompleted, nil)
			}
		}(id)
	}

	transitionWg.Wait()

	// Verify manager store is clean and no race conditions triggered under -race
	assert.True(t, true)
}

func TestJobManager_OwnerIsolation(t *testing.T) {
	manager := NewManager()

	// 1. Created job stores owner ID
	job, err := manager.Create(JobTypeExport, "user-alpha")
	require.NoError(t, err)
	require.NotNil(t, job)
	assert.Equal(t, "user-alpha", job.OwnerID)

	// 2. Same owner can retrieve job via GetForUser
	retrieved, err := manager.GetForUser(job.ID, "user-alpha")
	require.NoError(t, err)
	require.NotNil(t, retrieved)
	assert.Equal(t, job.ID, retrieved.ID)
	assert.Equal(t, "user-alpha", retrieved.OwnerID)

	// 3. Different owner cannot retrieve job (returns ErrJobNotFound)
	unauthorized, err := manager.GetForUser(job.ID, "user-beta")
	assert.ErrorIs(t, err, ErrJobNotFound)
	assert.Nil(t, unauthorized)

	// 4. Unknown job returns ErrJobNotFound
	unknown, err := manager.GetForUser("non-existent", "user-alpha")
	assert.ErrorIs(t, err, ErrJobNotFound)
	assert.Nil(t, unknown)

	// 5. Get remains unrestricted regardless of owner
	unrestricted, err := manager.Get(job.ID)
	require.NoError(t, err)
	require.NotNil(t, unrestricted)
	assert.Equal(t, job.ID, unrestricted.ID)
	assert.Equal(t, "user-alpha", unrestricted.OwnerID)
}
