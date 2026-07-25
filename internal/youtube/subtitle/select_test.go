package subtitle

import (
	"testing"
)

func TestSelectSubtitle(t *testing.T) {
	meta := SubtitleMetadata{
		Manual: []SubtitleTrack{
			{LanguageCode: "en", LanguageName: "English", Formats: []string{"json3", "vtt"}},
			{LanguageCode: "fr", LanguageName: "French", Formats: []string{"vtt"}},
		},
		Automatic: []SubtitleTrack{
			{LanguageCode: "en", LanguageName: "English (auto)", Formats: []string{"json3"}},
			{LanguageCode: "es", LanguageName: "Spanish", Formats: []string{"json3"}},
		},
	}

	t.Run("Exact language match manual preferred", func(t *testing.T) {
		prefs := []SubtitlePreference{
			{Language: "en", Type: ManualPreferred},
		}
		selected, err := SelectSubtitle(&meta, prefs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if selected.Language != "en" || selected.Type != "manual" {
			t.Errorf("got language=%s type=%s, want en, manual", selected.Language, selected.Type)
		}
	})

	t.Run("Automatic preferred when manual missing", func(t *testing.T) {
		prefs := []SubtitlePreference{
			{Language: "es", Type: ManualPreferred},
		}
		selected, err := SelectSubtitle(&meta, prefs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if selected.Language != "es" || selected.Type != "automatic" {
			t.Errorf("got language=%s type=%s, want es, automatic", selected.Language, selected.Type)
		}
	})

	t.Run("Wildcard language match", func(t *testing.T) {
		prefs := []SubtitlePreference{
			{Language: "fr*", Type: ManualOnly},
		}
		selected, err := SelectSubtitle(&meta, prefs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if selected.Language != "fr" {
			t.Errorf("got language=%s, want fr", selected.Language)
		}
	})

	t.Run("No match found", func(t *testing.T) {
		prefs := []SubtitlePreference{
			{Language: "de", Type: ManualOnly},
		}
		_, err := SelectSubtitle(&meta, prefs)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("Manual only match success", func(t *testing.T) {
		prefs := []SubtitlePreference{
			{Language: "en", Type: ManualOnly},
		}
		selected, err := SelectSubtitle(&meta, prefs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if selected.Type != "manual" || selected.Language != "en" {
			t.Errorf("got %+v", selected)
		}
	})

	t.Run("Manual only fallback fails when only automatic exists", func(t *testing.T) {
		prefs := []SubtitlePreference{
			{Language: "es", Type: ManualOnly},
		}
		_, err := SelectSubtitle(&meta, prefs)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("Automatic only match success", func(t *testing.T) {
		prefs := []SubtitlePreference{
			{Language: "es", Type: AutomaticOnly},
		}
		selected, err := SelectSubtitle(&meta, prefs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if selected.Type != "automatic" || selected.Language != "es" {
			t.Errorf("got %+v", selected)
		}
	})

	t.Run("Automatic preferred picks automatic when available", func(t *testing.T) {
		prefs := []SubtitlePreference{
			{Language: "en", Type: AutomaticPreferred},
		}
		selected, err := SelectSubtitle(&meta, prefs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if selected.Type != "automatic" {
			t.Errorf("got type %s, want automatic", selected.Type)
		}
	})

	t.Run("Wildcard match with multiple fallback preferences", func(t *testing.T) {
		prefs := []SubtitlePreference{
			{Language: "it*", Type: ManualPreferred},
			{Language: "fr*", Type: ManualPreferred},
		}
		selected, err := SelectSubtitle(&meta, prefs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if selected.Language != "fr" {
			t.Errorf("got language %s, want fr", selected.Language)
		}
	})
}
