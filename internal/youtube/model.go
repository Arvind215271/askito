// ./internal/youtube/model.go

package youtube

import (
	"time"
	"github.com/Arvind215271/askito/internal/youtube/transcript"
	"github.com/Arvind215271/askito/internal/youtube/description"
)

type Playlist struct {
	ID string `json:"id"`

	Title       string `json:"title"`
	Description string `json:"description"`

	ChannelID    string `json:"channel_id"`
	ChannelTitle string `json:"channel_title"`

	ThumbnailURL string `json:"thumbnail_url"`

	ItemCount int `json:"item_count"`

	PrivacyStatus string `json:"privacy_status"`

	PublishedAt time.Time `json:"published_at"`

	Videos []PlaylistVideo `json:"videos,omitempty"`
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

	Transcript *transcript.Transcript `json:"transcript,omitempty"`
	

	ChannelID    string `json:"channel_id"`
	ChannelTitle string `json:"channel_title"`

	ThumbnailURL string `json:"thumbnail_url"`

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
}
