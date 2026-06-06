// ./internal/youtube/input/input_test.go
 package youtube

import "testing"


func TestParsePlaylistURL(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		wantID    string
		wantError bool
	}{
		{
			name:      "mobile youtube domain",
			url:       "https://m.youtube.com/playlist?list=PL123",
			wantID:    "PL123",
			wantError: false,
		},
		{
			name:      "youtube.com domain",
			url:       "https://youtube.com/playlist?list=PL123",
			wantID:    "PL123",
			wantError: false,
		},
		{
			name:      "playlist with extra query params",
			url:       "https://www.youtube.com/playlist?list=PL123&si=abcdef",
			wantID:    "PL123",
			wantError: false,
		},
		{
			name:      "playlist from watch url",
			url:       "https://www.youtube.com/watch?v=abc123&list=PL123",
			wantID:    "PL123",
			wantError: false,
		},
		{
			name:      "empty url",
			url:       "",
			wantError: true,
		},
		{
			name:      "whitespace url",
			url:       "   ",
			wantError: true,
		},
		{
			name:      "invalid url format",
			url:       "://broken-url",
			wantError: true,
		},
		{
			name:      "missing list query value",
			url:       "https://www.youtube.com/playlist?list=",
			wantError: true,
		},
		{
			name:      "random text",
			url:       "hello world",
			wantError: true,
		},
		{
			name:      "subdomain not youtube",
			url:       "https://music.google.com/playlist?list=PL123",
			wantError: true,
		},
		{
			name:      "fake youtube domain",
			url:       "https://youtube-fake.com/playlist?list=PL123",
			wantError: true,
		},
		{
			name:      "notyoutube domain",
			url:       "https://notyoutube.com/playlist?list=PL123",
			wantError: true,
		},
		{
			name:      "country domain uk",
			url:       "https://youtube.co.uk/playlist?list=PL123",
			wantID:    "PL123",
			wantError: false,
		},
		{
			name:      "country domain de",
			url:       "https://youtube.de/playlist?list=PL123",
			wantID:    "PL123",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			playlist, err := Inputs.ParsePlaylistURL(tt.url)

			if tt.wantError {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if playlist.ID != tt.wantID {
				t.Fatalf(
					"expected %s, got %s",
					tt.wantID,
					playlist.ID,
				)
			}
		})
	}
}