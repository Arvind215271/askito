package chapter	

type Chapter struct {
	Title     string
	Timestamp string
	Seconds   int
}


type Chapters struct {
	// list of chapters to be used as primary structural data
	List []Chapter `json:"list"`
	// optional metadata (NOT enforcement)
	Valid bool `json:"valid"`
}