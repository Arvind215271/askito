package wordstats

import()



func BuildResult(
	originalWords int,
	originalChars int,
	words []*WordStats,
	buckets []Bucket,
) Result {

	// FILTERED COUNTS
	keptWords := len(words)

	filteredChars := 0
	for _, w := range words {
		filteredChars += len([]rune(w.Word))
	}

	// COMPRESSION
	wordRatio := 0.0
	if originalChars > 0 {
		wordRatio = float64(originalChars-filteredChars) / float64(originalChars) * 100
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