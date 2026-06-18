package chapter	

type Chapter struct {
	Title     string
	Timestamp string
	Seconds   int
}


type Chapters struct {
	List []Chapter `json:"list"`

	// AI / export friendly representation
	Text string `json:"text"`

	// optional metadata (NOT enforcement)
	Valid bool `json:"valid"`
}