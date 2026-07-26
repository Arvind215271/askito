package fields

const (
	GroupMetadata    = "Metadata"
	GroupDescription = "Description"
	GroupTranscript  = "Transcript"
	GroupSignal      = "Signal"
)

var (
	MetadataFields = []string{
		FieldID,
		FieldErrors,
		FieldTitle,
		FieldDescription,
		FieldChannelID,
		FieldChannelTitle,
		FieldThumbnails,
		FieldPublishedAt,
		FieldDuration,
		FieldDurationSeconds,
		FieldDurationMinutes,
		FieldDurationTimestamp,
		FieldViewCount,
		FieldLikeCount,
		FieldCommentCount,
		FieldTags,
		FieldCategoryID,
		FieldCaptionAvailable,
		FieldPrivacyStatus,
		FieldLiveBroadcastStatus,
	}

	DescriptionFields = []string{
		FieldDescriptionChapters,
		FieldDescriptionLinks,
		FieldDescriptionEmails,
		FieldDescriptionCleaned,
	}

	TranscriptFields = []string{
		FieldTranscriptText,
	}

	SignalFields = []string{
		FieldTranscriptSignal,
	}
)

