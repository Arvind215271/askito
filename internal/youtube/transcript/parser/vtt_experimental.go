package parser

import(
	"fmt"
	"regexp"
	"strings"
)

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