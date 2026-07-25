package subtitle

import (
	"fmt"
	"strings"
)

// SelectSubtitle evaluates preferences against available subtitle metadata
// and selects exactly one subtitle track according to wildcard matching and priority rules.
func SelectSubtitle(
	metadata *SubtitleMetadata,
	preferences []SubtitlePreference,
) (*SelectedSubtitle, error) {
	if metadata == nil {
		return nil, fmt.Errorf("subtitle metadata is nil")
	}
	if len(preferences) == 0 {
		return nil, fmt.Errorf("no subtitle preferences provided")
	}

	for _, pref := range preferences {
		if selected := matchPreference(metadata, pref); selected != nil {
			return selected, nil
		}
	}

	return nil, fmt.Errorf("no subtitle matched the given preferences")
}

func matchPreference(metadata *SubtitleMetadata, pref SubtitlePreference) *SelectedSubtitle {
	switch pref.Type {
	case ManualOnly:
		if track, ok := findMatchingTrack(metadata.Manual, pref.Language); ok {
			return &SelectedSubtitle{
				Language: track.LanguageCode,
				Type:     "manual",
				Format:   defaultFormat(track.Formats),
			}
		}
	case AutomaticOnly:
		if track, ok := findMatchingTrack(metadata.Automatic, pref.Language); ok {
			return &SelectedSubtitle{
				Language: track.LanguageCode,
				Type:     "automatic",
				Format:   defaultFormat(track.Formats),
			}
		}
	case ManualPreferred:
		if track, ok := findMatchingTrack(metadata.Manual, pref.Language); ok {
			return &SelectedSubtitle{
				Language: track.LanguageCode,
				Type:     "manual",
				Format:   defaultFormat(track.Formats),
			}
		}
		if track, ok := findMatchingTrack(metadata.Automatic, pref.Language); ok {
			return &SelectedSubtitle{
				Language: track.LanguageCode,
				Type:     "automatic",
				Format:   defaultFormat(track.Formats),
			}
		}
	case AutomaticPreferred:
		if track, ok := findMatchingTrack(metadata.Automatic, pref.Language); ok {
			return &SelectedSubtitle{
				Language: track.LanguageCode,
				Type:     "automatic",
				Format:   defaultFormat(track.Formats),
			}
		}
		if track, ok := findMatchingTrack(metadata.Manual, pref.Language); ok {
			return &SelectedSubtitle{
				Language: track.LanguageCode,
				Type:     "manual",
				Format:   defaultFormat(track.Formats),
			}
		}
	}

	return nil
}

func findMatchingTrack(tracks []SubtitleTrack, pattern string) (SubtitleTrack, bool) {
	// Exact match first
	for _, track := range tracks {
		if track.LanguageCode == pattern {
			return track, true
		}
	}

	// Wildcard matching (e.g. "en*", "*")
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		for _, track := range tracks {
			if strings.HasPrefix(track.LanguageCode, prefix) {
				return track, true
			}
		}
	}

	return SubtitleTrack{}, false
}

func defaultFormat(formats []string) string {
	for _, f := range formats {
		if f == "json3" {
			return "json3"
		}
		if f == "vtt" {
			return "vtt"
		}
	}
	if len(formats) > 0 {
		return formats[0]
	}
	return "json3"
}
