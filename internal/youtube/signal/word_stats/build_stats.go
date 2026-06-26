package wordstats

import (
	"regexp"
	"strings"
	"math"

	transcript "github.com/Arvind215271/askito/internal/youtube/transcript"
)

var cleanRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

// RAW Word Stats
func BuildWordStats(t *transcript.Transcript) map[string]*WordStats {
	stats := make(map[string]*WordStats)

	if t == nil {
		return stats
	}

	for _, seg := range t.Segments {
		clean := cleanRegex.ReplaceAllString(strings.ToLower(seg.Text), "")
		words := strings.Fields(clean)

		duration := seg.End - seg.Start

		for _, w := range words {
			if w == "" {
				continue
			}

			if _, ok := stats[w]; !ok {
				stats[w] = &WordStats{Word: w}
			}

			stats[w].Count++
			stats[w].Duration += duration
		}
	}

	return stats
}





// ADD SCORE
func ScoreWords(stats map[string]*WordStats) []*WordStats {
	var list []*WordStats

	for _, w := range stats {
		w.Score = math.Log1p(float64(w.Count))*2 +
			math.Log1p(w.Duration)

		list = append(list, w)
	}

	return list
}


