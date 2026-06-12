package youtubeurl

import (
	"testing"

	"github.com/Arvind215271/askito/internal/api"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		wantID         string
		wantType       InputType
		wantNormalized string
		wantErrorCode  string
	}{
		{
			name:           "mobile youtube domain",
			url:            "https://m.youtube.com/playlist?list=PL123",
			wantID:         "PL123",
			wantType:       InputTypePlaylist,
			wantNormalized: "https://www.youtube.com/playlist?list=PL123",
		},
		{
			name:           "youtube.com domain",
			url:            "https://youtube.com/playlist?list=PL123",
			wantID:         "PL123",
			wantType:       InputTypePlaylist,
			wantNormalized: "https://www.youtube.com/playlist?list=PL123",
		},
		{
			name:           "playlist with extra query params",
			url:            "https://www.youtube.com/playlist?list=PL123&si=abcdef",
			wantID:         "PL123",
			wantType:       InputTypePlaylist,
			wantNormalized: "https://www.youtube.com/playlist?list=PL123",
		},
		{
			name:           "playlist from watch url",
			url:            "https://www.youtube.com/watch?v=abc123&list=PL123",
			wantID:         "abc123",
			wantType:       InputTypeVideo,
			wantNormalized: "https://www.youtube.com/watch?v=abc123",
		},
		{
			name:           "standard video url",
			url:            "https://www.youtube.com/watch?v=abc123",
			wantID:         "abc123",
			wantType:       InputTypeVideo,
			wantNormalized: "https://www.youtube.com/watch?v=abc123",
		},
		{
			name:           "short url",
			url:            "https://youtu.be/abc123",
			wantID:         "abc123",
			wantType:       InputTypeVideo,
			wantNormalized: "https://www.youtube.com/watch?v=abc123",
		},
		{
			name:           "shorts url",
			url:            "https://www.youtube.com/shorts/abc123",
			wantID:         "abc123",
			wantType:       InputTypeVideo,
			wantNormalized: "https://www.youtube.com/watch?v=abc123",
		},
		{
			name:          "empty url",
			url:           "",
			wantErrorCode: "YOUTUBE_EMPTY_URL",
		},
		{
			name:          "whitespace url",
			url:           "   ",
			wantErrorCode: "YOUTUBE_EMPTY_URL",
		},
		{
			name:          "invalid url format",
			url:           "://broken-url",
			wantErrorCode: "YOUTUBE_INVALID_URL",
		},
		{
			name:          "missing list query value",
			url:           "https://www.youtube.com/playlist?list=",
			wantErrorCode: "YOUTUBE_MISSING_ID",
		},
		{
			name:          "random text",
			url:           "hello world",
			wantErrorCode: "YOUTUBE_INVALID_DOMAIN",
		},
		{
			name:          "subdomain not youtube",
			url:           "https://music.google.com/playlist?list=PL123",
			wantErrorCode: "YOUTUBE_INVALID_DOMAIN",
		},
		{
			name:          "fake youtube domain",
			url:           "https://youtube-fake.com/playlist?list=PL123",
			wantErrorCode: "YOUTUBE_INVALID_DOMAIN",
		},
		{
			name:          "notyoutube domain",
			url:           "https://notyoutube.com/playlist?list=PL123",
			wantErrorCode: "YOUTUBE_INVALID_DOMAIN",
		},
		{
			name:           "country domain uk",
			url:            "https://youtube.co.uk/playlist?list=PL123",
			wantID:         "PL123",
			wantType:       InputTypePlaylist,
			wantNormalized: "https://www.youtube.com/playlist?list=PL123",
		},
		{
			name:           "country domain de",
			url:            "https://youtube.de/playlist?list=PL123",
			wantID:         "PL123",
			wantType:       InputTypePlaylist,
			wantNormalized: "https://www.youtube.com/playlist?list=PL123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.url)

			if tt.wantErrorCode != "" {
				if err == nil {
					t.Fatalf("expected error %q", tt.wantErrorCode)
				}

				appErr, ok := err.(*api.AppError)
				if !ok {
					t.Fatalf("expected *api.AppError, got %T", err)
				}

				if appErr.Code != tt.wantErrorCode {
					t.Fatalf(
						"expected error code %q, got %q",
						tt.wantErrorCode,
						appErr.Code,
					)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.ID != tt.wantID {
				t.Fatalf(
					"expected id %q, got %q",
					tt.wantID,
					got.ID,
				)
			}

			if got.InputType != tt.wantType {
				t.Fatalf(
					"expected type %q, got %q",
					tt.wantType,
					got.InputType,
				)
			}

			if got.NormalizedURL != tt.wantNormalized {
				t.Fatalf(
					"expected normalized url %q, got %q",
					tt.wantNormalized,
					got.NormalizedURL,
				)
			}
		})
	}
}