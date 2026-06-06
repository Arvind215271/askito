// ./internal/youtube/input/input_test.go
package youtube

import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		wantID    string
		wantType  InputType
		wantError bool
	}{
		{
			name:      "mobile youtube domain",
			url:       "https://m.youtube.com/playlist?list=PL123",
			wantID:    "PL123",
			wantType:  InputTypePlaylist,
			wantError: false,
		},
		{
			name:      "youtube.com domain",
			url:       "https://youtube.com/playlist?list=PL123",
			wantID:    "PL123",
			wantType:  InputTypePlaylist,
			wantError: false,
		},
		{
			name:      "playlist with extra query params",
			url:       "https://www.youtube.com/playlist?list=PL123&si=abcdef",
			wantID:    "PL123",
			wantType:  InputTypePlaylist,
			wantError: false,
		},
		{
			name:      "playlist from watch url",
			url:       "https://www.youtube.com/watch?v=abc123&list=PL123",
			wantID:    "abc123",
			wantType:  InputTypeVideo,
			wantError: false,
		},
		{
			name:      "standard video url",
			url:       "https://www.youtube.com/watch?v=abc123",
			wantID:    "abc123",
			wantType:  InputTypeVideo,
			wantError: false,
		},
		{
			name:      "short url",
			url:       "https://youtu.be/abc123",
			wantID:    "abc123",
			wantType:  InputTypeVideo,
			wantError: false,
		},
		{
			name:      "shorts url",
			url:       "https://www.youtube.com/shorts/abc123",
			wantID:    "abc123",
			wantType:  InputTypeVideo,
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
			wantType:  InputTypePlaylist,
			wantError: false,
		},
		{
			name:      "country domain de",
			url:       "https://youtube.de/playlist?list=PL123",
			wantID:    "PL123",
			wantType:  InputTypePlaylist,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input, err := Inputs.ParseURL(tt.url)

			if tt.wantError {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if input.ID != tt.wantID {
				t.Fatalf(
					"expected id %s, got %s",
					tt.wantID,
					input.ID,
				)
			}

			if input.InputType != tt.wantType {
				t.Fatalf(
					"expected type %s, got %s",
					tt.wantType,
					input.InputType,
				)
			}
		})
	}
}