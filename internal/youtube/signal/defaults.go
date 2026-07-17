package signal

func DefaultSignalRequest(videoID string) SignalRequest {
	return SignalRequest{
		VideoID:    videoID,
		Analysis:   "word-stats",
		UseHeavy:   true,
		MinFreq:    3,
		Depth:      1,
		WindowSize: 300,
	}
}
