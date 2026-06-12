package youtube

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"

// 	"github.com/Arvind215271/askito/internal/logger"
	
// )

// func DebugYouTube(
// 	ctx context.Context,
// 	logger *logger.Logger,
// 	client *Client,
// ) {

// 	const playlistID = "PLgUwDviBIf0oF6QL8m22w1hIDC1vJ_BHz"

// 	logger.Info(
// 		"starting youtube debug",
// 		"playlist_id",
// 		playlistID,
// 	)

// 	playlist, err := client.GetPlaylist(
// 		ctx,
// 		playlistID,
// 	)
// 	if err != nil {
// 		logger.Error(
// 			"failed to fetch playlist",
// 			"error",
// 			err,
// 		)
// 		return
// 	}

// 	printSize("playlist", playlist)

// 	items, err := client.GetPlaylistItems(
// 		ctx,
// 		playlistID,
// 	)
// 	if err != nil {
// 		logger.Error(
// 			"failed to fetch playlist items",
// 			"error",
// 			err,
// 		)
// 		return
// 	}

// 	printSize("playlist items", items)

// 	videoIDs := make(
// 		[]string,
// 		0,
// 		len(items),
// 	)

// 	for _, item := range items {

// 		if item == nil ||
// 			item.ContentDetails == nil {
// 			continue
// 		}

// 		videoIDs = append(
// 			videoIDs,
// 			item.ContentDetails.VideoId,
// 		)
// 	}

// 	printSize("video ids", videoIDs)

// 	videos, err := client.GetVideos(
// 		ctx,
// 		videoIDs,
// 	)
// 	if err != nil {
// 		logger.Error(
// 			"failed to fetch videos",
// 			"error",
// 			err,
// 		)
// 		return
// 	}

// 	printSize("videos", videos)

// 	playlistVideos := GetPlaylistVideos(
// 		items,
// 		videos,
// 	)

// 	printSize(
// 		"playlist videos",
// 		playlistVideos,
// 	)

// 	fmt.Println()

// 	fmt.Printf(
// 		"playlist items : %d\n",
// 		len(items),
// 	)

// 	fmt.Printf(
// 		"video ids      : %d\n",
// 		len(videoIDs),
// 	)

// 	fmt.Printf(
// 		"videos         : %d\n",
// 		len(videos),
// 	)

// 	fmt.Printf(
// 		"playlistVideos : %d\n",
// 		len(playlistVideos),
		
// 	)

// 	fmt.Printf("\n================ SAMPLE =================\n")

// 	if len(playlistVideos) > 0 {

// 		first := playlistVideos[0]

// 		fmt.Println("Title      :", first.Title)
// 		fmt.Println("Position   :", first.Position)
// 		fmt.Println("Duration   :", first.Duration)
// 		fmt.Println("Views      :", first.ViewCount)
// 		fmt.Println("Published  :", first.PublishedAt)
// 		fmt.Println("Added To PL:", first.AddedAt)

// 		if len(first.Tags) > 0 {
// 			fmt.Println("Tags:", first.Tags[:min(5, len(first.Tags))])
// 		}

// 		fmt.Println()
// 	}

// 	// ------------------------------------
// 	// export test
// 	// ------------------------------------

// 	playlistModel := Playlist{
// 		ID:            playlist.Id,
// 		Title:         playlist.Snippet.Title,
// 		Description:   playlist.Snippet.Description,
// 		ChannelID:     playlist.Snippet.ChannelId,
// 		ChannelTitle:  playlist.Snippet.ChannelTitle,
// 		ThumbnailURL:  getPlaylistThumbnail(playlist),
// 		ItemCount:     int(playlist.ContentDetails.ItemCount),
// 		PrivacyStatus: playlist.Status.PrivacyStatus,
// 		PublishedAt:   parseTime(playlist.Snippet.PublishedAt),
// 		Videos:        playlistVideos,
// 	}

// 	exportData, err := BuildPlaylistExport(
// 		playlistModel,
// 		[]string{
// 			"title",
// 			"view_count",
// 			"duration",
// 		},
// 	)

// 	if err != nil {
// 		logger.Error(
// 			"failed to build export",
// 			"error",
// 			err,
// 		)
// 	} else {

// 		printSize("export data", exportData)

// 		// -----------------------------
// 		// VERIFY PLAYLIST FIELDS
// 		// -----------------------------
// 		fmt.Printf("\n================ PLAYLIST FIELDS =================\n")

// 		fmt.Println("Playlist Title      :", exportData["title"])
// 		fmt.Println("Playlist Channel    :", exportData["channel_title"])
// 		fmt.Println("Playlist Item Count :", exportData["item_count"])

// 		// -----------------------------
// 		// VERIFY VIDEO FILTERING
// 		// -----------------------------
// 		fmt.Printf("\n================ VIDEO SAMPLE =================\n")

// 		if vids, ok := exportData["videos"].([]any); ok && len(vids) > 0 {

// 			firstVideo, _ := vids[0].(map[string]any)

// 			fmt.Println("Title      :", firstVideo["title"])
// 			fmt.Println("Views      :", firstVideo["view_count"])
// 			fmt.Println("Duration   :", firstVideo["duration"])

// 			// these MUST ALWAYS exist (not part of filter)
// 			fmt.Println("Position   :", firstVideo["position"])
// 			fmt.Println("Added At   :", firstVideo["added_at"])
// 		}

// 		// -----------------------------
// 		// EXPORT JSON
// 		// -----------------------------
// 		fmt.Printf("\n================ EXPORT JSON =================\n")

// 		jsonExporter := export.NewJSONExporter()

// 		exportJSON, err := jsonExporter.Export(exportData)
// 		if err != nil {
// 			logger.Error(
// 				"failed to export json",
// 				"error",
// 				err,
// 			)
// 		} else {

// 			fmt.Printf("\n================ EXPORT =================\n")

// 			fmt.Println(string(exportJSON))

// 			fmt.Printf(
// 				"\nexport json size : %.2f KB\n",
// 				float64(len(exportJSON))/1024,
// 			)
// 		}
// 	}

// 	fmt.Println()
// }

// func printSize(
// 	name string,
// 	v any,
// ) {

// 	b, _ := json.Marshal(v)

// 	fmt.Printf(
// 		"%-20s %10.2f KB\n",
// 		name,
// 		float64(len(b))/1024,
// 	)
// }