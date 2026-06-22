package youtube

import (
	"strconv"
	"strings"
	
)

// ParseYouTubeDuration converts ISO8601 duration (PT#H#M#S) into seconds.
func ParseYouTubeDuration(d string) float64 {
	if d == "" {
		return 0
	}

	// Must start with PT
	if !strings.HasPrefix(d, "PT") {
		return 0
	}

	d = strings.TrimPrefix(d, "PT")

	var hours, minutes, seconds int

	var num strings.Builder

	for _, r := range d {
		switch {
		case r >= '0' && r <= '9':
			num.WriteRune(r)

		case r == 'H':
			hours = mustAtoi(num.String())
			num.Reset()

		case r == 'M':
			minutes = mustAtoi(num.String())
			num.Reset()

		case r == 'S':
			seconds = mustAtoi(num.String())
			num.Reset()
		}
	}

	return float64(hours*3600 + minutes*60 + seconds)
}

func mustAtoi(s string) int {
	if s == "" {
		return 0
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return n
}


func (v *Video) Seconds() float64 {
	return ParseYouTubeDuration(v.Duration)
}