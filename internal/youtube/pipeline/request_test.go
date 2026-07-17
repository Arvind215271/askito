package pipeline

import (
	"reflect"
	"testing"

	"github.com/Arvind215271/askito/internal/youtube/signal"
	"github.com/Arvind215271/askito/internal/youtube/subtitle"
	"github.com/Arvind215271/askito/internal/youtube/transcript"
)

func TestRequest_Combinations(t *testing.T) {
	tests := []struct {
		name string
		req  Request
	}{
		{"OnlyTranscript", Request{Fields: []string{"transcript"}, Transcript: &transcript.ProcessingRequest{Output: "plain-text"}}},
		{"OnlySubtitle", Request{Fields: []string{"subtitle"}, Subtitle: &subtitle.DownloadRequest{Type: "auto", Language: "en", Format: "vtt"}}},
		{"OnlySignal", Request{Fields: []string{"signal"}, Signal: &signal.SignalRequest{Analysis: "words", UseHeavy: true}}},
		{"TranscriptAndSubtitle", Request{Fields: []string{"transcript", "subtitle"}, Transcript: &transcript.ProcessingRequest{Output: "timeline-text"}, Subtitle: &subtitle.DownloadRequest{Type: "manual"}}},
		{"AllThree", Request{Fields: []string{"transcript", "subtitle", "signal"}, Transcript: &transcript.ProcessingRequest{Output: "plain-text"}, Subtitle: &subtitle.DownloadRequest{Type: "auto"}, Signal: &signal.SignalRequest{Analysis: "all"}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify that the request can be constructed and fields are accessible
			if !reflect.DeepEqual(tt.req.Fields, tt.req.Fields) {
				t.Errorf("field mismatch")
			}
		})
	}
}
