package planner

// ResourceInfo contains information regarding requested resources.
type ResourceInfo struct {
	VideosRequested    bool
	PlaylistsRequested bool
	MixedResources     bool
}

// ExecutionPlan represents the ordered list of execution stages and resource info.
type ExecutionPlan struct {
	Stages    []Stage
	Resources ResourceInfo
}

// HasStage checks if a specific stage is present in the execution plan.
type StageFlags map[Stage]bool

// Flags returns a map of stage flags for easy checks.
func (p *ExecutionPlan) Flags() StageFlags {
	flags := make(StageFlags)
	for _, s := range p.Stages {
		flags[s] = true
	}
	return flags
}

// NeedsStage checks if the execution plan includes the given stage.
func (p *ExecutionPlan) NeedsStage(stage Stage) bool {
	for _, s := range p.Stages {
		if s == stage {
			return true
		}
	}
	return false
}

// Convenience helper methods as requested

func (p *ExecutionPlan) NeedsPlaylistItems() bool {
	return p.NeedsStage(StagePlaylistItems)
}

func (p *ExecutionPlan) NeedsMetadata() bool {
	return p.NeedsStage(StageMetadata)
}

func (p *ExecutionPlan) NeedsDescription() bool {
	return p.NeedsStage(StageDescription)
}

func (p *ExecutionPlan) NeedsSubtitle() bool {
	return p.NeedsStage(StageSubtitle)
}

func (p *ExecutionPlan) NeedsTranscript() bool {
	return p.NeedsStage(StageTranscript)
}

func (p *ExecutionPlan) NeedsSignal() bool {
	return p.NeedsStage(StageSignal)
}
