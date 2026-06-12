package parser

import (
	"fmt"
	"regexp"
	"strings"
)


var (
	timestampRegex = regexp.MustCompile(`^(\d{2}:\d{2}:\d{2}\.\d+)\s-->.*$`)


	tagRegex = regexp.MustCompile(
		`<[^>]+>`,
	)
)


func isRollingDuplicate(prev, curr string) bool {
	if prev == "" {
		return false
	}

	// exact match
	if prev == curr {
		return true
	}

	// curr is suffix extension of prev (rolling behavior)
	if strings.HasSuffix(curr, prev) {
		return true
	}

	// prev is prefix of curr (stream continuation)
	if strings.HasPrefix(curr, prev) {
		return true
	}

	// small edit distance shortcut (cheap heuristic)
	if len(prev) > 0 && len(curr) > 0 {
		diff := abs(len(curr) - len(prev))
		if diff <= 3 && strings.Contains(curr, prev) {
			return true
		}
	}

	return false
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}



func ExtractCleanTranscriptText(vtt string) string {
	lines := strings.Split(vtt, "\n")

	var result []string

	var prevText string
	var prevStable string
	var buffer []string

	flush := func(time string, text string) {
		text = strings.TrimSpace(text)
		if text == "" {
			return
		}

		if text != prevStable {
			result = append(result, fmt.Sprintf("[%s] %s", time, text))
			prevStable = text
		}
	}

	var currentTime string

	for _, raw := range lines {
		line := strings.TrimSpace(raw)

		if line == "" {
			continue
		}

		// skip headers
		if line == "WEBVTT" ||
			strings.HasPrefix(line, "Kind:") ||
			strings.HasPrefix(line, "Language:") {
			continue
		}

		// timestamp
		if timestampRegex.MatchString(line) {
			// flush previous cue
			if len(buffer) > 0 {
				text := strings.Join(buffer, " ")
				text = tagRegex.ReplaceAllString(text, "")
				text = strings.TrimSpace(text)

				// 🔥 CORE LOGIC: rolling subtitle filter
				if !isRollingDuplicate(prevText, text) {
					flush(currentTime, text)
					prevText = text
				}
			}

			parts := strings.Split(line, " --> ")
			currentTime = strings.TrimSpace(parts[0])
			buffer = buffer[:0]
			continue
		}

		line = tagRegex.ReplaceAllString(line, "")
		line = strings.TrimSpace(line)

		if line != "" {
			buffer = append(buffer, line)
		}
	}

	// final flush
	if len(buffer) > 0 {
		text := strings.Join(buffer, " ")
		text = tagRegex.ReplaceAllString(text, "")
		text = strings.TrimSpace(text)

		if !isRollingDuplicate(prevText, text) {
			flush(currentTime, text)
		}
	}

	return strings.Join(result, "\n")
}



// ExtractCleanTranscriptText parses a pre-flattened VTT file (No duplication logic required!)
func ExtractCleanTranscriptTextAlreadyClean(vtt string) string {
	var result []string
	lines := strings.Split(vtt, "\n")
	var currentTime string

	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || line == "WEBVTT" || strings.HasPrefix(line, "Kind:") || strings.HasPrefix(line, "Language:") {
			continue
		}

		// Grab the timestamp
		if timestampRegex.MatchString(line) {
			parts := strings.Split(line, " --> ")
			currentTime = strings.TrimSpace(parts[0])
			continue
		}

		// This line is guaranteed to be clean, static text without rolling duplicates
		line = tagRegex.ReplaceAllString(line, "")
		line = strings.TrimSpace(line)

		if line != "" {
			result = append(result, fmt.Sprintf("[%s] %s", currentTime, line))
		}
	}

	return strings.Join(result, "\n")
}

