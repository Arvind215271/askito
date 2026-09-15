package job

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// JobManager manages jobs in memory with thread-safe operations.
type JobManager struct {
	mu   sync.RWMutex
	jobs map[string]Job
}

// NewManager creates a new instance of JobManager.
func NewManager() *JobManager {
	return &JobManager{
		jobs: make(map[string]Job),
	}
}

// Create creates a new job of the specified type and owner ID, assigns a UUID, sets status to queued,
// records CreatedAt in UTC, stores it, and returns a snapshot copy.
func (m *JobManager) Create(jobType JobType, ownerID string) (*Job, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, j := range m.jobs {
		if j.OwnerID == ownerID && (j.Status == StatusQueued || j.Status == StatusRunning) {
			return nil, ErrActiveJobExists
		}
	}

	id := uuid.New().String()
	now := time.Now().UTC()

	job := Job{
		ID:        id,
		Type:      jobType,
		Status:    StatusQueued,
		OwnerID:   ownerID,
		CreatedAt: now,
	}

	m.jobs[id] = job

	snapshot := job
	return &snapshot, nil
}

// Get retrieves a job by ID and returns a snapshot copy to prevent caller mutation and data races.
func (m *JobManager) Get(id string) (*Job, error) {
	m.mu.RLock()
	job, exists := m.jobs[id]
	m.mu.RUnlock()

	if !exists {
		return nil, ErrJobNotFound
	}

	snapshot := job
	return &snapshot, nil
}

// GetForUser retrieves a job by ID for a specific owner, returning a snapshot copy if found and owned.
func (m *JobManager) GetForUser(id string, ownerID string) (*Job, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	job, exists := m.jobs[id]
	if !exists {
		return nil, ErrJobNotFound
	}

	if job.OwnerID != ownerID {
		return nil, ErrJobNotFound
	}

	snapshot := job
	return &snapshot, nil
}

// Transition validates and applies a state transition for a job atomically.
func (m *JobManager) Transition(id string, targetStatus JobStatus, jobErr error) (*Job, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	job, exists := m.jobs[id]
	if !exists {
		return nil, ErrJobNotFound
	}

	// Validate transition rules
	if err := validateTransition(job.Status, targetStatus); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	job.Status = targetStatus

	switch targetStatus {
	case StatusRunning:
		if job.StartedAt == nil {
			job.StartedAt = &now
		}
	case StatusCompleted, StatusCancelled:
		if job.FinishedAt == nil {
			job.FinishedAt = &now
		}
	case StatusFailed:
		if job.FinishedAt == nil {
			job.FinishedAt = &now
		}
		if jobErr != nil {
			job.Error = jobErr.Error()
		}
	}

	m.jobs[id] = job
	snapshot := job
	return &snapshot, nil
}

// validateTransition checks whether transitioning from currentStatus to targetStatus is allowed.
func validateTransition(current, target JobStatus) error {
	// Terminal states cannot transition further
	if current.IsTerminal() {
		return ErrInvalidTransition
	}

	switch current {
	case StatusQueued:
		if target == StatusRunning || target == StatusCancelled {
			return nil
		}
	case StatusRunning:
		if target == StatusCompleted || target == StatusFailed || target == StatusCancelled {
			return nil
		}
	}

	return ErrInvalidTransition
}
