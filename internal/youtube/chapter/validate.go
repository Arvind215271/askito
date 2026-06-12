package chapter

func ValidateChapters(chapters []Chapter) bool {
	// YouTube requires at least 3 chapters
	if len(chapters) < 3 {
		return false
	}

	// First chapter must start at 0
	if chapters[0].Seconds != 0 {
		return false
	}

	// Timestamps must be strictly increasing
	for i := 1; i < len(chapters); i++ {
		if chapters[i].Seconds <= chapters[i-1].Seconds {
			return false
		}
	}

	return true
}