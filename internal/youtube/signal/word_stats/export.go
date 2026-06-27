package wordstats

import (
	"strings"
	"fmt"
)


type ExportMode int

const (
	ExportCompact ExportMode = iota
	ExportVerbose
)


func Export(result Result, mode ExportMode) string {
	switch mode {

	case ExportVerbose:
		return exportVerbose(result)

	default:
		return exportCompact(result)
	}
}


func exportCompact(result Result) string {
	var b strings.Builder

	b.WriteString("# COMPRESSED TRANSCRIPT WORD SIGNAL (BUCKET MODEL)\n")
	b.WriteString("# This is a compressed representation of transcript word importance.\n")
	b.WriteString("# Each line represents a bucket ordered by importance (highest → lowest).\n")
	b.WriteString("#\n")
	b.WriteString("# FORMAT:\n")
	b.WriteString("# countMin-countMax, timeSpentMin-timeSpentMax(s), words\n")
	b.WriteString("#\n")
	b.WriteString("# FIELD MEANING:\n")
	b.WriteString("# count = how many times word appears in transcript\n")
	b.WriteString("# timeSpent = total duration (sum of timestamps where word occurs)\n\n")

	for _, bucket := range result.Buckets {
		if len(bucket.Words) == 0 {
			continue
		}

		b.WriteString(fmt.Sprintf(
			"%d-%d,%.0f-%.0f(s),%s\n",
			bucket.CountMin,
			bucket.CountMax,
			bucket.DurationMin,
			bucket.DurationMax,
			strings.Join(bucket.Words, " "),
		))
	}

	return b.String()
}


func exportVerbose(result Result) string {
	var b strings.Builder

	b.WriteString("# COMPRESSED TRANSCRIPT SIGNAL (BUCKET IMPORTANCE MODEL)\n")
	b.WriteString("#\n")
	b.WriteString("# This output is a hierarchical compression of a transcript.\n")
	b.WriteString("# Words are grouped into buckets ordered by importance.\n")
	b.WriteString("# Higher position = more important semantic signal.\n")
	b.WriteString("#\n")

	b.WriteString("# FORMAT:\n")
	b.WriteString("# count <min>-<max> | timeSpent <min>-<max>(s) | words\n")
	b.WriteString("#\n")

	b.WriteString("# FIELD DEFINITIONS:\n")
	b.WriteString("# count = number of occurrences of a word in transcript\n")
	b.WriteString("# timeSpent = total accumulated duration where word appears (NOT timeline)\n")
	b.WriteString("# words = aggregated token list in each bucket\n\n")

	for _, bucket := range result.Buckets {
		if len(bucket.Words) == 0 {
			continue
		}

		b.WriteString(fmt.Sprintf(
			"count:%d-%d | timeSpent:%.0f-%.0fs | words: %s\n",
			bucket.CountMin,
			bucket.CountMax,
			bucket.DurationMin,
			bucket.DurationMax,
			strings.Join(bucket.Words, " "),
		))
	}

	return b.String()
}


func ExportWindowVerbose(
	results []Result,
	cfg AnalysisConfig,
) string {

	var b strings.Builder

	// ---------------- HEADER ----------------
	b.WriteString("# WINDOW CONCEPT FINGERPRINT (BUCKET MODEL)\n")
	b.WriteString("# Each window shows word importance grouped by frequency buckets.\n")
	b.WriteString("# Format is fully data-driven (no semantic labels).\n")
	b.WriteString("#\n")

	b.WriteString(fmt.Sprintf("# WS=%.0f | BUCKETS=%d\n\n",
		cfg.WindowSize,
		cfg.BucketCount,
	))

	b.WriteString("# FORMAT:\n")
	b.WriteString("# start-ends\n")
	b.WriteString("# min-max: words\n")
	b.WriteString("# min-max: words\n")
	b.WriteString("# ...\n")
	b.WriteString("#\n\n")

	// ---------------- WINDOWS ----------------
	for i, res := range results {

		start := float64(i) * cfg.WindowSize
		end := start + cfg.WindowSize

		b.WriteString(fmt.Sprintf("%.0f-%.0fs\n", start, end))

		// ---------------- BUCKETS ----------------
		for _, bucket := range res.Buckets {

			if len(bucket.Words) == 0 {
				continue
			}

			b.WriteString(fmt.Sprintf(
				"%d-%d: %s\n",
				bucket.CountMin,
				bucket.CountMax,
				strings.Join(bucket.Words, " "),
			))
		}

		b.WriteByte('\n')
	}

	return b.String()
}