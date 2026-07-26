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

	result := ExportData{
		"id":             playlist.ID,
		"title":          playlist.Title,
		"description":    playlist.Description,
		"channel_id":     playlist.ChannelID,
		"channel_title":  playlist.ChannelTitle,
		"thumbnails":     playlist.Thumbnails,
		"tags":           playlist.Tags,
		"item_count":     playlist.ItemCount,
		"privacy_status": playlist.PrivacyStatus,
		"published_at":   playlist.PublishedAt,
		"modified_at":    playlist.ModifiedAt,
	}

	if len(playlist.Items) > 0 {
		items := make([]any, 0, len(playlist.Items))
		for _, item := range playlist.Items {
			items = append(items, item)
		}
		result["items"] = items
	}

	if len(playlist.Videos) > 0 {
		videos := make([]any, 0, len(playlist.Videos))
		for _, v := range playlist.Videos {
			videoData, err := exportStruct(v.Video, planner)
			if err != nil {
				return nil, youtube.Err.Export.MarshalFailed().Wrap(err)
			}

			videoData["position"] = v.Position
			videoData["added_at"] = v.AddedAt
			videoData["id"] = v.Video.ID

			videos = append(videos, videoData)
		}
		result["videos"] = videos
	}

	return result, nil
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
