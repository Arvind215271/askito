package wordstats

type BucketConfig struct {
	Count int
}

// CreateBuckets assigns words into importance buckets
func CreateBuckets(words []*WordStats, cfg BucketConfig) []Bucket {
	if len(words) == 0 {
		return nil
	}

	var maxScore float64
	for _, w := range words {
		if w.Score > maxScore {
			maxScore = w.Score
		}
	}

	buckets := make([]Bucket, cfg.Count)

	for _, w := range words {
		bid := 0
		if maxScore > 0 {
			bid = int((w.Score / maxScore) * float64(cfg.Count-1))
		}

		b := &buckets[bid]

		b.ID = bid
		b.Words = append(b.Words, w.Word)

		if b.CountMin == 0 || w.Count < b.CountMin {
			b.CountMin = w.Count
		}
		if w.Count > b.CountMax {
			b.CountMax = w.Count
		}

		if b.DurationMin == 0 || w.Duration < b.DurationMin {
			b.DurationMin = w.Duration
		}
		if w.Duration > b.DurationMax {
			b.DurationMax = w.Duration
		}
	}

	// compact output (remove empty buckets)
	out := make([]Bucket, 0, cfg.Count)
	for i := cfg.Count - 1; i >= 0; i-- {
		if len(buckets[i].Words) > 0 {
			out = append(out, buckets[i])
		}
	}

	return out
}