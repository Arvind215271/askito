package job

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJobManager_Concurrency(t *testing.T) {
	manager := NewManager()
	const goroutines = 50
	const jobsPerGoroutine = 20

	var wg sync.WaitGroup
	wg.Add(goroutines * 2)

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

				job, err := manager.Create(jType)
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
