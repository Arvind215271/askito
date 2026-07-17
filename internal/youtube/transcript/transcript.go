// ./internal/youtube/transcript/transcript.go
package transcript

import (
	"fmt"
	"strings"
)

func (t *Transcript) ToTimelineText() string {
	var b strings.Builder

	for _, s := range t.Segments {
		b.WriteString(fmt.Sprintf(
			"[%0.2f - %0.2f] %s\n",
			s.Start,
			s.End,
			s.Text,
		))
	}

	return strings.TrimSpace(b.String())
}

func (t *Transcript) ToPlainText() string {
	var b strings.Builder

	for i, s := range t.Segments {
		if i > 0 {
			b.WriteString(" ")
		}
		b.WriteString(s.Text)
	}

	return strings.TrimSpace(b.String())
}

func (t *Transcript) GroupByDuration(windowSeconds float64) *Transcript {
	if len(t.Segments) == 0 {
		return &Transcript{Language: t.Language, Segments: nil}
	}

	var result []TranscriptSegment

	current := TranscriptSegment{
		Start: t.Segments[0].Start,
	}

	windowEnd := current.Start + windowSeconds

	for _, seg := range t.Segments {
		if seg.Start >= windowEnd {
			current.Text = strings.TrimSpace(current.Text)

			if current.Text != "" {
				result = append(result, current)
			}

			current = TranscriptSegment{
				Start: seg.Start,
			}

			windowEnd = seg.Start + windowSeconds
		}

		if current.Text != "" {
			current.Text += " "
		}

		current.Text += seg.Text
		current.End = seg.End
	}

	current.Text = strings.TrimSpace(current.Text)

	if current.Text != "" {
		result = append(result, current)
	}

	return &Transcript{
		Language: t.Language,
		Segments: result,
	}
}
