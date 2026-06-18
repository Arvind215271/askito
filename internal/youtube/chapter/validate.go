package chapter

func (c Chapters) IsValid() bool {
	// YouTube requires at least 3 chapters
	if len(c.List) < 3 {
		return false
	}

	// First chapter must start at 0
	if c.List[0].Seconds != 0 {
		return false
	}

	// Timestamps must be strictly increasing
	for i := 1; i < len(c.List); i++ {
		if c.List[i].Seconds <= c.List[i-1].Seconds {
			return false
		}
	}

	return true
}