// ./internal/youtube/model.go

package youtube

import (
	"time"

	"github.com/Arvind215271/askito/internal/youtube/description"
	"github.com/Arvind215271/askito/internal/youtube/subtitle"
	"github.com/Arvind215271/askito/internal/youtube/transcript"
)

type Playlist struct {
	ID string `json:"id"`

	Title       string `json:"title"`
	Description string `json:"description"`

	ChannelID    string `json:"channel_id"`
	ChannelTitle string `json:"channel_title"`

	Thumbnails []Thumbnail `json:"thumbnails,omitempty"`

	Tags []string `json:"tags,omitempty"`

	ItemCount int `json:"item_count"`

	PrivacyStatus string `json:"privacy_status"`

	PublishedAt time.Time `json:"published_at"`
	ModifiedAt  time.Time `json:"modified_at"`

	// Lightweight playlist entries.
	Items []PlaylistItem `json:"items,omitempty"`

	// Fully processed playlist videos.
	Videos []PlaylistVideo `json:"videos,omitempty"`
}

type PlaylistItem struct {
	VideoID string `json:"video_id"`

	// Playlist metadata
	Position int       `json:"position"`
	AddedAt  time.Time `json:"added_at"`

	// Video snapshot
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`

	Duration  time.Duration `json:"duration,omitempty"`
	ViewCount uint64        `json:"view_count,omitempty"`

	// Channel snapshot
	ChannelID    string `json:"channel_id,omitempty"`
	ChannelTitle string `json:"channel_title,omitempty"`

	// URLs
	URL string `json:"url,omitempty"`

	// Thumbnails
	Thumbnails []Thumbnail `json:"thumbnails,omitempty"`

	// Playlist/video status
	PrivacyStatus string `json:"privacy_status,omitempty"`
}

type Thumbnail struct {
	URL    string `json:"url"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

type PlaylistVideo struct {
	Video

	Position int       `json:"position"`
	AddedAt  time.Time `json:"added_at"`
}

type Video struct {
	ID string `json:"id"`

	Title       string `json:"title"`
	Description string `json:"description"`

	DescriptionMetadata description.Metadata `json:"description_metadata,omitempty"`

	// Export fields
	DescriptionChapters string   `json:"description_chapters"`
	DescriptionLinks    []string `json:"description_links"`
	DescriptionEmails   []string `json:"description_emails"`
	DescriptionCleaned  string   `json:"description_cleaned"`

	Transcript *transcript.Transcript `json:"transcript,omitempty"`
	// TranscriptText stores the final transcript representation.
	TranscriptText   string `json:"transcript_text"`
	TranscriptSignal string `json:"transcript_signal"`

	SubtitleMetadata subtitle.SubtitleMetadata `json:"subtitle_metadata"`

	ChannelID    string `json:"channel_id"`
	ChannelTitle string `json:"channel_title"`

	Thumbnails []Thumbnail `json:"thumbnails,omitempty"`

	PublishedAt time.Time `json:"published_at"`

	Duration          string  `json:"duration"`
	DurationSeconds   int64   `json:"duration_seconds"`
	DurationMinutes   float64 `json:"duration_minutes"`
	DurationTimestamp string  `json:"duration_timestamp"`

	ViewCount    uint64 `json:"view_count"`
	LikeCount    uint64 `json:"like_count"`
	CommentCount uint64 `json:"comment_count"`

	Tags []string `json:"tags,omitempty"`

	CategoryID string `json:"category_id"`

	CaptionAvailable bool `json:"caption_available"`

	PrivacyStatus       string `json:"privacy_status"`
	LiveBroadcastStatus string `json:"live_broadcast_status"`

	Errors []string `json:"errors"`

}
