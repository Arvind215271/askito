package transcript

import (
	"testing"
)

func TestValidateProcessingRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *ProcessingRequest
		wantErr bool
	}{
		{"nil", nil, false},
		{"missing output", &ProcessingRequest{Output: ""}, true},
		{"valid segments", &ProcessingRequest{Output: "segments"}, false},
		{"valid plain-text", &ProcessingRequest{Output: "plain-text"}, false},
		{"valid timeline-text", &ProcessingRequest{Output: "timeline-text"}, false},
		{"unsupported output", &ProcessingRequest{Output: "xml"}, true},
		{"negative window size", &ProcessingRequest{Output: "segments", WindowSize: -1}, true},
		{"zero window size", &ProcessingRequest{Output: "segments", WindowSize: 0}, false},
		{"positive window size", &ProcessingRequest{Output: "segments", WindowSize: 30}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateProcessingRequest(tt.req); (err != nil) != tt.wantErr {
				t.Errorf("ValidateProcessingRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGroupByDuration(t *testing.T) {
	transcript := &Transcript{
		Segments: []TranscriptSegment{
			{Start: 0, End: 10, Text: "hello"},
			{Start: 15, End: 25, Text: "world"},
			{Start: 40, End: 50, Text: "test"},
		},
	}

	t.Run("empty transcript", func(t *testing.T) {
		tr := &Transcript{Segments: nil}
		got := tr.GroupByDuration(30)
		if len(got.Segments) != 0 {
			t.Errorf("got %d segments, want 0", len(got.Segments))
		}
	})

	t.Run("grouping", func(t *testing.T) {
		got := transcript.GroupByDuration(30)
		if len(got.Segments) != 2 {
			t.Errorf("got %d segments, want 2", len(got.Segments))
		}
		if got.Segments[0].Text != "hello world" {
			t.Errorf("got text %q, want 'hello world'", got.Segments[0].Text)
		}
		if got.Segments[1].Text != "test" {
			t.Errorf("got text %q, want 'test'", got.Segments[1].Text)
		}
	})
}
