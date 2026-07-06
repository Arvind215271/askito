package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type YouTubeJSON3 struct {
	Events []JSON3Event `json:"events"`
}

type JSON3Event struct {
	StartMs    int64          `json:"tStartMs"`
	DurationMs int64          `json:"dDurationMs"`
	Segments   []JSON3Segment `json:"segs"`
}

type JSON3Segment struct {
	UTF8 string `json:"utf8"`
}

func main() {
	f, err := os.Open("-.en.json3")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var yt YouTubeJSON3
	dec := json.NewDecoder(f)
	if err := dec.Decode(&yt); err != nil {
		fmt.Printf("Error decoding JSON: %v\n", err)
	} else {
		fmt.Printf("Successfully parsed %d events\n", len(yt.Events))
	}
}
