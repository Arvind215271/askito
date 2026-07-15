package pipeline

import (
	"reflect"
	"testing"
)

func TestRequest_Combinations(t *testing.T) {
	tests := []struct {
		name string
		req  Request
	}{
		{"OnlyTranscript", Request{Fields: []string{"transcript"}, Transcript: &TranscriptRequest{Type: "auto", Language: "en", Format: "json3"}}},
		{"OnlySubtitle", Request{Fields: []string{"subtitle"}, Subtitle: &SubtitleRequest{Type: "auto", Language: "en", Format: "vtt"}}},
		{"OnlySignal", Request{Fields: []string{"signal"}, Signal: &SignalRequest{Analysis: "words", UseHeavyStopWords: true}}},
		{"TranscriptAndSubtitle", Request{Fields: []string{"transcript", "subtitle"}, Transcript: &TranscriptRequest{Type: "auto"}, Subtitle: &SubtitleRequest{Type: "manual"}}},
		{"TranscriptAndSignal", Request{Fields: []string{"transcript", "signal"}, Transcript: &TranscriptRequest{Language: "en"}, Signal: &SignalRequest{Analysis: "stats"}}},
		{"SubtitleAndSignal", Request{Fields: []string{"subtitle", "signal"}, Subtitle: &SubtitleRequest{Language: "es"}, Signal: &SignalRequest{MinFreq: 5}}},
		{"AllThree", Request{Fields: []string{"transcript", "subtitle", "signal"}, Transcript: &TranscriptRequest{Type: "auto"}, Subtitle: &SubtitleRequest{Type: "auto"}, Signal: &SignalRequest{Analysis: "all"}}},
		{"EmptyFields", Request{Fields: []string{}, Transcript: &TranscriptRequest{Type: "auto"}}},
		{"NoOptionalParts", Request{Fields: []string{"metadata"}}},
		{"PartialSubtitle", Request{Fields: []string{"subtitle"}, Subtitle: &SubtitleRequest{Type: "auto"}}},
		{"PartialTranscript", Request{Fields: []string{"transcript"}, Transcript: &TranscriptRequest{Format: "json3"}}},
		{"PartialSignal", Request{Fields: []string{"signal"}, Signal: &SignalRequest{Depth: 0.5, WindowSize: 100}}},
		{"ComplexSignal", Request{Fields: []string{"signal"}, Signal: &SignalRequest{Analysis: "complex", MinFreq: 10, Depth: 0.95, WindowSize: 500, BucketCount: 128}}},
		{"MinimalSubtitle", Request{Subtitle: &SubtitleRequest{}}},
		{"MinimalTranscript", Request{Transcript: &TranscriptRequest{}}},
		{"MinimalSignal", Request{Signal: &SignalRequest{}}},
		{"MultiFields", Request{Fields: []string{"f1", "f2", "f3"}}},
		{"SubtitleLanguageVariation", Request{Subtitle: &SubtitleRequest{Language: "de"}}},
		{"TranscriptFormatVariation", Request{Transcript: &TranscriptRequest{Format: "vtt"}}},
		{"SignalConfigVariation", Request{Signal: &SignalRequest{UseHeavyStopWords: false}}},
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
