package stopwords

import (
	"bufio"
	"os"
	"strings"
	"sync"
)

const defaultPath = "internal/youtube/signal/stopwords/heavy_stopwords.txt"

var (
	once  sync.Once
	heavy map[string]struct{}
	err   error
)

// load initializes stopwords once
func load() {
	heavy = make(map[string]struct{})

	file, e := os.Open(defaultPath)
	if e != nil {
		err = e
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		w := strings.TrimSpace(scanner.Text())

		if w == "" || strings.HasPrefix(w, "#") {
			continue
		}

		heavy[strings.ToLower(w)] = struct{}{}
	}

	err = scanner.Err()
}

// IsHeavy checks if word is a heavy stopword (lazy-loaded)
func IsHeavy(word string) bool {
	once.Do(load)

	if err != nil || heavy == nil {
		return false
	}

	_, ok := heavy[strings.ToLower(word)]
	return ok
}