package fields

// ValidFields maps allowed top-level field names to their JSON tag
var ValidFields = map[string]struct{}{
	FieldID:                  {},
	FieldTitle:               {},
	FieldDescription:         {},
	FieldDescriptionChapters: {},
	FieldDescriptionLinks:    {},
	FieldDescriptionEmails:   {},
	FieldDescriptionCleaned:  {},
	FieldTranscriptText:      {},
	FieldTranscriptSignal:    {},
	FieldChannelID:           {},
	FieldChannelTitle:        {},
	FieldThumbnailURL:        {},
	FieldPublishedAt:         {},
	FieldDuration:            {},
	FieldDurationSeconds:     {},
	FieldDurationMinutes:     {},
	FieldDurationTimestamp:   {},
	FieldViewCount:           {},
	FieldLikeCount:           {},
	FieldCommentCount:        {},
	FieldTags:                {},
	FieldCategoryID:          {},
	FieldCaptionAvailable:    {},
	FieldPrivacyStatus:       {},
	FieldLiveBroadcastStatus: {},
	FieldErrors:            {},
}

// Metadata fields
const (
	FieldID                  = "id"
	FieldErrors              = "errors"
	FieldTitle               = "title"
	FieldDescription         = "description"
	FieldChannelID           = "channel_id"
	FieldChannelTitle        = "channel_title"
	FieldThumbnailURL        = "thumbnail_url"
	FieldPublishedAt         = "published_at"
	FieldDuration            = "duration"
	FieldDurationSeconds     = "duration_seconds"
	FieldDurationMinutes     = "duration_minutes"
	FieldDurationTimestamp   = "duration_timestamp"
	FieldViewCount           = "view_count"
	FieldLikeCount           = "like_count"
	FieldCommentCount        = "comment_count"
	FieldTags                = "tags"
	FieldCategoryID          = "category_id"
	FieldCaptionAvailable    = "caption_available"
	FieldPrivacyStatus       = "privacy_status"
	FieldLiveBroadcastStatus = "live_broadcast_status"
)

// Description fields
const (
	FieldDescriptionChapters = "description_chapters"
	FieldDescriptionLinks    = "description_links"
	FieldDescriptionEmails   = "description_emails"
	FieldDescriptionCleaned  = "description_cleaned"
)

// Transcript fields
const (
	FieldTranscriptText   = "transcript_text"
	FieldTranscriptSignal = "transcript_signal"
)
