// reference: https://stackoverflow.com/questions/63821605/how-do-i-get-info-about-a-youtube-videos-chapters-from-the-api

package chapter

import (
	"regexp"
	"strconv"
	"strings"
)



var chapterRegex = regexp.MustCompile(
	`(?i)^(.*?)?(\d{1,2}:\d{2}(?::\d{2})?)(.*?)$`,
)


func ExtractChapters(description string) Chapters {
	if description == "" {
		return Chapters{}
	}

	var list []Chapter

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

		timestamp := match[2]

		title := strings.TrimSpace(
			strings.ReplaceAll(
				strings.ReplaceAll(line, timestamp, ""),
				"()", "",
			),
		)

		title = strings.Trim(title, "-:|[]() ")

		list = append(list, Chapter{
			Title:     title,
			Timestamp: timestamp,
			Seconds:   parseTimestamp(timestamp),
		})
	}

	return Chapters{
		List:  list,
		Text:  ChaptersToText(list),
		Valid: ValidateChapters(list),
	}
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


func ExtractChapterText(description string) string {
	chapters := ExtractChapters(description)

	if len(chapters.Text) == 0 {
		return ""
	}

	return chapters.Text
}

