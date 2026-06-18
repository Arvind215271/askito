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
func (c Chapters) Text() string {
	if len(c.List) == 0 {
		return ""
	}

	var builder strings.Builder

	for i, chapter := range c.List {
		builder.WriteString(chapter.Timestamp)

		if chapter.Title != "" {
			builder.WriteString(" ")
			builder.WriteString(chapter.Title)
		}

		if i < len(c.List)-1 {
			builder.WriteString("\n")
		}
	}

	return builder.String()
}


func (c Chapters) Titles() string {
	if len(c.List) == 0 {
		return ""
	}

	var builder strings.Builder

	for i, chapter := range c.List {
		if chapter.Title != "" {
			builder.WriteString(chapter.Title)
		}

		if i < len(c.List)-1 {
			builder.WriteString("\n")
		}
	}

	return builder.String()
}



