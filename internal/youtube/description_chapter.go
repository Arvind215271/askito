// ./internal/youtube/descriptoin_chapter.go

// reference: https://stackoverflow.com/questions/63821605/how-do-i-get-info-about-a-youtube-videos-chapters-from-the-api

package youtube

import (
	"regexp"
	"strconv"
	"strings"
)



var chapterRegex = regexp.MustCompile(
	`(?i)^(.*?)?(\d{1,2}:\d{2}(?::\d{2})?)(.*?)$`,
)


func ExtractChapters(description string) []Chapter {
	
	var chapters []Chapter
	
	// if descriptoin is nill... return empty chapter
	if description == "" {
		return chapters
	}

	// split lines to be processed independently.
	lines := strings.Split(description, "\n")
	

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		match := chapterRegex.FindStringSubmatch(line)
		if len(match) == 0 {
			continue
		}

		// get the timestamp for the it...
		timestamp := match[2]

		// get the title as well
		title := strings.TrimSpace(
			strings.ReplaceAll(
				strings.ReplaceAll(line, timestamp, ""),
				"()", "",
			),
		)

		title = strings.Trim(title, "-:|[]() ")

		chapters = append(chapters, Chapter{
			Title:     title,
			Timestamp: timestamp,
			Seconds:   parseTimestamp(timestamp),
		})
	}

	if !ValidateChapters(chapters) {
		return nil
	}

	return chapters
}

// function to parse different time stamp
//
// yt have timestamp as minute or hours. So we have to deal with both the cases.
func parseTimestamp(ts string) int {
	parts := strings.Split(ts, ":")

	switch len(parts) {
	case 2:
		m, _ := strconv.Atoi(parts[0])
		s, _ := strconv.Atoi(parts[1])
		return m*60 + s

	case 3:
		h, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		s, _ := strconv.Atoi(parts[2])
		return h*3600 + m*60 + s
	}

	return 0
}



func ValidateChapters(chapters []Chapter) bool {
	// YouTube requires at least 3 chapters
	if len(chapters) < 3 {
		return false
	}

	// First chapter must start at 0
	if chapters[0].Seconds != 0 {
		return false
	}

	// Timestamps must be strictly increasing
	for i := 1; i < len(chapters); i++ {
		if chapters[i].Seconds <= chapters[i-1].Seconds {
			return false
		}
	}

	return true
}


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


func ExtractChapterText(description string) string {
	chapters := ExtractChapters(description)

	if len(chapters) == 0 {
		return ""
	}

	return ChaptersToText(chapters)
}

