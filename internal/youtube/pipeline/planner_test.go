package pipeline

import (
	"testing"

	"github.com/Arvind215271/askito/internal/youtube/fields"
)

func TestPlanner_NeedsMetadata(t *testing.T) {
	tests := []struct {
		name   string
		fields []string
		want   bool
	}{
		{"Empty fields (everything)", []string{}, true},
		{"Metadata field", []string{fields.FieldTitle}, true},
		{"Non-metadata field", []string{fields.FieldTranscriptText}, false},
		{"Description requires metadata", []string{fields.FieldDescriptionCleaned}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pFields, _ := fields.NewPlanner(tt.fields)
			p := NewPlanner(pFields)
			if got := p.NeedsMetadata(); got != tt.want {
				t.Errorf("NeedsMetadata() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlanner_NeedsTranscript(t *testing.T) {
	tests := []struct {
		name   string
		fields []string
		want   bool
	}{
		{"Empty fields (everything)", []string{}, true},
		{"Transcript field", []string{fields.FieldTranscriptText}, true},
		{"Non-transcript field", []string{fields.FieldTitle}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pFields, _ := fields.NewPlanner(tt.fields)
			p := NewPlanner(pFields)
			if got := p.NeedsTranscript(); got != tt.want {
				t.Errorf("NeedsTranscript() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlanner_NeedsSignal(t *testing.T) {
	tests := []struct {
		name   string
		fields []string
		want   bool
	}{
		{"Empty fields (everything)", []string{}, true},
		{"Signal field", []string{fields.FieldTranscriptSignal}, true},
		{"Non-signal field", []string{fields.FieldTitle}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pFields, _ := fields.NewPlanner(tt.fields)
			p := NewPlanner(pFields)
			if got := p.NeedsSignal(); got != tt.want {
				t.Errorf("NeedsSignal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlanner_NeedsSubtitle(t *testing.T) {
	tests := []struct {
		name   string
		fields []string
		want   bool
	}{
		{"Empty fields (everything)", []string{}, true},
		{"Transcript field", []string{fields.FieldTranscriptText}, true},
		{"Signal field", []string{fields.FieldTranscriptSignal}, true},
		{"Neither field", []string{fields.FieldTitle}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pFields, _ := fields.NewPlanner(tt.fields)
			p := NewPlanner(pFields)
			if got := p.NeedsSubtitle(); got != tt.want {
				t.Errorf("NeedsSubtitle() = %v, want %v", got, tt.want)
			}
		})
	}
}
