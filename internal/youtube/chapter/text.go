package chapter

import (
	"strings"
)

// ChaptersToText converts chapters into a plain text representation.
//
// Example:
//
//	0:00 Intro
//	2:14 Invoked Mechaba
//	3:23 Utopia the Lightning
func ChaptersToText(chapters []Chapter) string {
	if len(chapters) == 0 {
		return ""
	}

	var builder strings.Builder

	for i, chapter := range chapters {
		builder.WriteString(chapter.Timestamp)

		if chapter.Title != "" {
			builder.WriteString(" ")
			builder.WriteString(chapter.Title)
		}

		if i < len(chapters)-1 {
			builder.WriteString("\n")
		}
	}

	return builder.String()
}