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



// ExtractScreenBufferTranscript uses a 3-state overlap buffer to stitch rolling captions.
// ExtractScreenBufferTranscript uses a 3-state overlap buffer with safety flushes.
// ExtractScreenBufferTranscript uses a line-level comparison matrix to perfectly stitch rolling subtitles.
func ExtractScreenBufferTranscript(vtt string) string {
	var result []string
	var bufferLines []string
	var currentStart string

	lines := strings.Split(vtt, "\n")
	var cueTextBuilder strings.Builder
	var cueTime string

	// Helper to calculate total character count currently in our buffer
	getBufferLen := func() int {
		total := 0
		for _, l := range bufferLines {
			total += len(l)
		}
		return total
	}

	flush := func() {
		if len(bufferLines) > 0 {
			cleanText := strings.TrimSpace(strings.Join(bufferLines, " "))
			if cleanText != "" {
				result = append(result, fmt.Sprintf("[%s] %s", currentStart, cleanText))
			}
		}
		bufferLines = bufferLines[:0]
	}

	processCue := func(time string, text string) {
		text = strings.TrimSpace(text)
		if text == "" {
			return
		}

		// Split the new text into individual lines
		incomingLines := strings.Split(text, "\n")
		var validNewLines []string

		for _, incLine := range incomingLines {
			incLine = strings.TrimSpace(incLine)
			if incLine == "" {
				continue
			}

			// Check if this incoming line already exists inside our current screen buffer
			isDuplicate := false
			for _, bufLine := range bufferLines {
				if bufLine == incLine {
					isDuplicate = true
					break
				}
			}

			if !isDuplicate {
				validNewLines = append(validNewLines, incLine)
			}
		}

		// If nothing new was added, skip entirely
		if len(validNewLines) == 0 {
			return
		}

		// Initialize start timestamp if the buffer was empty
		if len(bufferLines) == 0 {
			currentStart = time
		}

		// Safety Paragraph Break: If the buffer exceeds 180 characters, flush it out
		if getBufferLen() > 180 {
			flush()
			currentStart = time
		}

		// Append only the genuinely new text elements to the screen buffer
		bufferLines = append(bufferLines, validNewLines...)
	}

	// Main Parsing Engine Loop
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || line == "WEBVTT" || strings.HasPrefix(line, "Kind:") || strings.HasPrefix(line, "Language:") {
			continue
		}

		if matches := timestampRegex.FindStringSubmatch(line); matches != nil {
			if cueTextBuilder.Len() > 0 {
				processCue(cueTime, tagRegex.ReplaceAllString(cueTextBuilder.String(), ""))
				cueTextBuilder.Reset()
			}
			if len(matches) > 1 {
				cueTime = matches[1]
			}
			continue
		}

		// Preserve literal text newlines instead of squashing them early
		if cueTextBuilder.Len() > 0 {
			cueTextBuilder.WriteString("\n")
		}
		cueTextBuilder.WriteString(line)
	}

	if cueTextBuilder.Len() > 0 {
		processCue(cueTime, tagRegex.ReplaceAllString(cueTextBuilder.String(), ""))
	}
	flush()

	return strings.Join(result, "\n")
}



var (
	// Extracts the contents inside <c>...</c> tags
	cTagRegex = regexp.MustCompile(`<c[^>]*>([^<]+)</c>`)
	// Used for fallback cleaning
	generalTagRegex = regexp.MustCompile(`<[^>]+>`)
)

// ExtractTokenStreamTranscript extracts YouTube's DOM deltas by targeting <c> tags.
func ExtractTokenStreamTranscript(vtt string) string {
	var result []string
	var buffer strings.Builder
	var currentStart string

	lines := strings.Split(vtt, "\n")
	var cueTime string

	flush := func() {
		text := strings.TrimSpace(buffer.String())
		if text != "" {
			result = append(result, fmt.Sprintf("[%s] %s", currentStart, text))
		}
		buffer.Reset()
	}

	for _, raw := range lines {
		line := strings.TrimSpace(raw)

		if line == "" || line == "WEBVTT" || strings.HasPrefix(line, "Kind:") || strings.HasPrefix(line, "Language:") {
			continue
		}

		// Timestamp detection
		if matches := timestampRegex.FindStringSubmatch(line); matches != nil {
			cueTime = matches[1]
			if buffer.Len() == 0 {
				currentStart = cueTime
			}
			continue
		}

		// If the line contains <c> tags, YouTube is telling us exactly what the NEW words are.
		if strings.Contains(line, "<c>") || strings.Contains(line, "</c>") {
			matches := cTagRegex.FindAllStringSubmatch(line, -1)
			for _, match := range matches {
				if len(match) > 1 {
					buffer.WriteString(match[1]) // Append only the new tokens
				}
			}
		} else {
			// If there are no <c> tags, this is either a standard non-rolling VTT line,
			// or the start of a new rolling block.
			cleanLine := generalTagRegex.ReplaceAllString(line, "")
			cleanLine = strings.TrimSpace(cleanLine)
			
			// If our buffer already contains this text, ignore it (it's a static repeat)
			if !strings.Contains(buffer.String(), cleanLine) {
				// Flush the old sentence and start a new one
				flush()
				currentStart = cueTime
				buffer.WriteString(cleanLine)
			}
		}
	}

	flush() // Final flush

	// Clean up duplicate spaces that token merging might create
	finalOutput := strings.Join(result, "\n")
	spaceRegex := regexp.MustCompile(`\s+`)
	return spaceRegex.ReplaceAllString(finalOutput, " ")
}