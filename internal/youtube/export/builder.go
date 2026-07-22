package export

import (
	youtube "github.com/Arvind215271/askito/internal/youtube"
	"github.com/Arvind215271/askito/internal/youtube/fields"
)

// this is common function that we will be using to export playlist and convert it to a simple format that can be used by an export TYPE like JSON, CSV, etc.
//
// It is the filter layer actually. We are already getting the data in our Domain Model.
// The only thing left is to filter what is needed from Video ONLY.
func BuildPlaylistExport(
	playlist youtube.Playlist,
	planner *fields.Planner,
) (ExportData, error) {

	videos := make([]any, 0, len(playlist.Videos))

	for _, v := range playlist.Videos {

		// ONLY Video struct is filtered
		videoData, err := exportStruct(v.Video, planner)
		if err != nil {
			return nil, youtube.Err.Export.MarshalFailed().Wrap(err)
		}

		// PlaylistVideo metadata is ALWAYS preserved
		videoData["position"] = v.Position
		videoData["added_at"] = v.AddedAt
		videoData["id"] = v.Video.ID

		videos = append(videos, videoData)
	}

	// Playlist itself is NOT filtered
	return ExportData{
		"id":             playlist.ID,
		"title":          playlist.Title,
		"description":    playlist.Description,
		"channel_id":     playlist.ChannelID,
		"channel_title":  playlist.ChannelTitle,
		"thumbnail_url":  playlist.ThumbnailURL,
		"item_count":     playlist.ItemCount,
		"privacy_status": playlist.PrivacyStatus,
		"published_at":   playlist.PublishedAt,

		"videos": videos,
	}, nil
}

// this is common function that we will be using to export video and convert it to a simple format that can be used by an export TYPE like JSON, CSV, etc.
//
// It is the filter layer actually. We are already getting the data in our Domain Model.
// The only thing left is to filter what is needed from Video ONLY.
func BuildVideoExport(
	video youtube.Video,
	planner *fields.Planner,
) (ExportData, error) {

	// ONLY Video is filterable
	data, err := exportStruct(video, planner)
	if err != nil {
		return nil, youtube.Err.Export.MarshalFailed().Wrap(err)
	}

	return data, nil
}

// this is common function that we will be using to export multiple videos and convert it to a simple format that can be used by an export TYPE like JSON, CSV, etc.
//
// It is the filter layer actually. We are already getting the data in our Domain Model.
// The only thing left is to filter what is needed from Video ONLY.
func BuildBatchVideoExport(
	videos []youtube.Video,
	planner *fields.Planner,
) (ExportData, error) {
	exportedVideos := make([]any, 0, len(videos))

	for _, v := range videos {
		videoData, err := exportStruct(v, planner)
		if err != nil {
			return nil, youtube.Err.Export.MarshalFailed().Wrap(err)
		}
		exportedVideos = append(exportedVideos, videoData)
	}

	return ExportData{
		"videos": exportedVideos,
	}, nil
}
