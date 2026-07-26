package planner

// Stage represents an execution stage capability.
type Stage string

const (
	StagePlaylistItems Stage = "playlist_items"
	StageMetadata      Stage = "metadata"
	StageSubtitle      Stage = "subtitle"
	StageDescription   Stage = "description"
	StageTranscript    Stage = "transcript"
	StageSignal        Stage = "signal"
)

// OrderedStages defines the exact execution sequence order.
var OrderedStages = []Stage{
	StagePlaylistItems,
	StageMetadata,
	StageDescription,
	StageSubtitle,
	StageTranscript,
	StageSignal,
}
