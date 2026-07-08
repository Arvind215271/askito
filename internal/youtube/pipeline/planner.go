package pipeline

import (
	"github.com/Arvind215271/askito/internal/youtube/fields"
)

// Planner handles pipeline execution decisions based on requested fields.
type Planner struct {
	fields *fields.Planner
}

// NewPlanner creates a new pipeline planner.
func NewPlanner(p *fields.Planner) *Planner {
	return &Planner{fields: p}
}

// NeedsMetadata returns true if every field or any metadata field is requested.
func (p *Planner) NeedsMetadata() bool {
	return p.fields.HasAny(fields.MetadataFields) || p.NeedsDescription()
}

// NeedsDescription returns true if every field or any description export field is requested.
func (p *Planner) NeedsDescription() bool {
	return p.fields.HasAny(fields.DescriptionFields)
}

// NeedsTranscript returns true if every field or any transcript export field is requested.
func (p *Planner) NeedsTranscript() bool {
	return p.fields.HasAny(fields.TranscriptFields)
}

// NeedsSignal returns true if every field or only transcript signal is requested.
func (p *Planner) NeedsSignal() bool {
	return p.fields.HasAny(fields.SignalFields)
}

// NeedsSubtitle returns true if transcript or signal are requested.
func (p *Planner) NeedsSubtitle() bool {
	return p.NeedsTranscript() || p.NeedsSignal()
}
