package youtubeurl

import (
	"strings"
	"regexp"
)

func normalizeInput(raw string) string {
	raw = strings.TrimSpace(raw)

	// unify commas and newlines into spaces
	raw = strings.ReplaceAll(raw, ",", " ")
	raw = strings.ReplaceAll(raw, "\n", " ")

	// collapse multiple spaces later via Fields()
	return raw
}

var urlRegex = regexp.MustCompile(`https?://[^\s,]+`)

func extractURLs(raw string) []string {
	return urlRegex.FindAllString(raw, -1)
}

func cleanURL(u string) string {
	return strings.Trim(u, ",.()[]{}<>\"'")
}

func ParseMany(raw string) []ParseItemResult {
	candidates := extractURLs(raw)

	results := make([]ParseItemResult, 0, len(candidates))

	for _, c := range candidates {
		input, err := Parse(cleanURL(c))

		results = append(results, ParseItemResult{
			Input: input,
			Error: err,
			Raw:   c,
		})
	}

	return results
}