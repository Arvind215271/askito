package wordstats

import(
	transcript "github.com/Arvind215271/askito/internal/youtube/transcript"

	"strings"

)

func BuildResult(t *transcript.Transcript, words []*WordStats, buckets []Bucket) Result {

	// ---- RAW TRANSCRIPT COUNTS (SOURCE OF TRUTH) ----
	originalWords := 0
	originalChars := 0

	for _, seg := range t.Segments {
		originalChars += len([]rune(seg.Text))

		// word count from raw transcript
		// (same logic used in BuildWordStats)
		originalWords += len(strings.Fields(seg.Text))
	}

	// ---- FILTERED COUNTS ----
	keptWords := len(words)

	filteredChars := 0
	for _, w := range words {
		filteredChars += len([]rune(w.Word)) 
	}

	// ---- COMPRESSION ----
	wordRatio := 0.0
	if originalChars > 0 {
		wordRatio = float64(originalChars - filteredChars) / float64(originalChars) * 100
	}

	return Result{
		Buckets: buckets,

		OriginalWords: originalWords,
		KeptWords:     keptWords,

		OriginalChars: originalChars,
		OutputChars:   filteredChars,

		CompressionRatio: wordRatio,
	}
}