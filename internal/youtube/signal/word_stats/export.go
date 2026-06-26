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

	b.WriteString("# COMPRESSED TRANSCRIPT WORD SIGNAL (32 BUCKET MODEL)\n")
	b.WriteString("# This is a compressed representation of transcript word importance.\n")
	b.WriteString("# Each bucket groups words by importance score (32 = highest, 1 = lowest).\n")
	b.WriteString("#\n")
	b.WriteString("# FORMAT:\n")
	b.WriteString("# bucketId, countMin-countMax, timeSpentMin-timeSpentMax(s), words(space separated)\n")
	b.WriteString("#\n")
	b.WriteString("# FIELD MEANING:\n")
	b.WriteString("# count = how many times word appears in transcript\n")
	b.WriteString("# timeSpent = total duration (sum of timestamps where word occurs)\n\n")

	for _, bucket := range result.Buckets {
		if len(bucket.Words) == 0 {
			continue
		}

		b.WriteString(fmt.Sprintf(
			"%d,%d-%d,%.0f-%.0f(s),%s\n",
			bucket.ID,
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

	b.WriteString("# COMPRESSED TRANSCRIPT SIGNAL (32 BUCKET IMPORTANCE MODEL)\n")
	b.WriteString("#\n")
	b.WriteString("# This output is a hierarchical compression of a transcript.\n")
	b.WriteString("# Words are grouped into 32 importance buckets.\n")
	b.WriteString("# Higher bucket ID = more important semantic signal.\n")
	b.WriteString("#\n")

	b.WriteString("# FORMAT:\n")
	b.WriteString("# Bucket <id> | count <min>-<max> | timeSpent <min>-<max>(s) | words\n")
	b.WriteString("#\n")

	b.WriteString("# FIELD DEFINITIONS:\n")
	b.WriteString("# count = number of occurrences of a word in transcript\n")
	b.WriteString("# timeSpent = total accumulated duration where word appears (NOT timeline)\n")
	b.WriteString("# words = aggregated token list in that bucket\n\n")

	for _, bucket := range result.Buckets {
		if len(bucket.Words) == 0 {
			continue
		}

		b.WriteString(fmt.Sprintf(
			"Bucket %d | count:%d-%d | timeSpent:%.0f-%.0fs | words: %s\n",
			bucket.ID,
			bucket.CountMin,
			bucket.CountMax,
			bucket.DurationMin,
			bucket.DurationMax,
			strings.Join(bucket.Words, " "),
		))
	}

	return b.String()
}