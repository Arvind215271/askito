package youtube

// import (
// 	"context"
// 	"time"

// 	"strconv"
	
	
// 	yt "google.golang.org/api/youtube/v3"
// 	"google.golang.org/api/option"
// )

// type Client struct {
// 	service *yt.Service
// }

 

// func NewClient(
// 	ctx context.Context,
// 	apiKey string,
// ) (*Client, error) {

// 	service, err := yt.NewService(
// 		ctx,
// 		option.WithAPIKey(apiKey),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &Client{
// 		service: service,
// 	}, nil
// }


// func parseTime(s string) time.Time {
// 	t, err := time.Parse(time.RFC3339, s)
// 	if err != nil {
// 		return time.Time{}
// 	}
// 	return t
// }

// func parseUint(s string) uint64 {
// 	v, err := strconv.ParseUint(s, 10, 64)
// 	if err != nil {
// 		return 0
// 	}
// 	return v
// }

// func getPlaylistThumbnail(
// 	playlist *yt.Playlist,
// ) string {

// 	if playlist.Snippet == nil ||
// 		playlist.Snippet.Thumbnails == nil {
// 		return ""
// 	}

// 	if playlist.Snippet.Thumbnails.High != nil {
// 		return playlist.Snippet.Thumbnails.High.Url
// 	}

// 	if playlist.Snippet.Thumbnails.Medium != nil {
// 		return playlist.Snippet.Thumbnails.Medium.Url
// 	}

// 	if playlist.Snippet.Thumbnails.Default != nil {
// 		return playlist.Snippet.Thumbnails.Default.Url
// 	}

// 	return ""
// }

// func getVideoThumbnail(
// 	video *yt.Video,
// ) string {

// 	// there are different resolutoin type per video. Simply return whatever is present... Not much needed to be true. We can simply keep those field as empty.
// 	if video.Snippet == nil ||
// 		video.Snippet.Thumbnails == nil {
// 		return ""
// 	}

// 	if video.Snippet.Thumbnails.High != nil {
// 		return video.Snippet.Thumbnails.High.Url
// 	}

// 	if video.Snippet.Thumbnails.Medium != nil {
// 		return video.Snippet.Thumbnails.Medium.Url
// 	}

// 	if video.Snippet.Thumbnails.Default != nil {
// 		return video.Snippet.Thumbnails.Default.Url
// 	}

// 	return ""
// }

// // returns metadata about the playlist, but does not contain any videoID in it 
// func (c *Client) GetPlaylist(
// 	ctx context.Context,
// 	playlistID string,
// ) (*yt.Playlist, error) {

// 	resp, err := c.service.Playlists.
// 		List([]string{
// 			"snippet",
// 			"status",
// 			"contentDetails",

// 		}).
// 		Id(playlistID).
// 		Context(ctx).
// 		Do()

// 	if err != nil {
// 		return nil, Err.Playlist.FetchFailed().Wrap(err)
// 	}

// 	if len(resp.Items) == 0 {
// 		return nil, Err.Playlist.NotFound()
// 	}

// 	item := resp.Items[0]

	

// 	return item, nil
// }





// // this returns the videoID in the playlist with some data
// func (c *Client) GetPlaylistItems(
// 	ctx context.Context,
// 	playlistID string,
// ) ([]*yt.PlaylistItem, error) {

// 	var items []*yt.PlaylistItem
// 	// items represent an array of playlist items.
// 	// the reason is that YT api returns upto 0 to 50 videoID per call.
// 	// Thus, to retrieve them all, we need to append those into the system for them to be resuable and such.
// 	// so we have to use the pagetoken for the upcoming ID if exist.
// 	var pageToken string

// 	for {
// 		resp, err := c.service.PlaylistItems.
// 			List([]string{
// 				"snippet",
// 				"contentDetails",
// 			}).
// 			PlaylistId(playlistID).
// 			MaxResults(50).
// 			PageToken(pageToken).
// 			Context(ctx).
// 			Do()

// 		if err != nil {
// 			return nil, Err.Playlist.FetchFailed().Wrap(err)
// 		}
// 		// add the recieved items to an array. 	
// 		items = append(items, resp.Items...)

// 		if resp.NextPageToken == "" {
// 			break
// 		}

// 		pageToken = resp.NextPageToken
// 	}

// 	return items, nil
// }



// // fetches actual video resources from YouTube
// func (c *Client) GetVideos(
// 	ctx context.Context,
// 	videoIDs []string,
// ) ([]*yt.Video, error) {

// 	if len(videoIDs) == 0 {
// 		return nil, nil
// 	}

// 	var videos []*yt.Video

// 	// youtube allows max 50 ids per request
// 	for start := 0; start < len(videoIDs); start += 50 {
		
// 		end := start + 50
// 		if end > len(videoIDs) {
// 			end = len(videoIDs)
// 		}

// 		chunk := videoIDs[start:end]

// 		resp, err := c.service.Videos.
// 			List([]string{
// 				"snippet",
// 				"contentDetails",
// 				"statistics",
// 				"status",
// 			}).
// 			Id(chunk...).
// 			Context(ctx).
// 			Do()

// 		if err != nil {
// 			return nil, Err.Video.FetchFailed().Wrap(err)
// 		}

// 		videos = append(videos, resp.Items...)
// 	}

// 	return videos, nil
// }

// // maps the fetched Youtube API video type to domain logic video type
// func GetPlaylistVideos(
// 	items []*yt.PlaylistItem,
// 	videos []*yt.Video,
// ) []PlaylistVideo {

// 	// create a map of yt.video and total video we have/
// 	videoMap := make(map[string]*yt.Video, len(videos))

// 	// for each video. map that in map with its video ID.
// 	// why? Because each video in Playlist Video have a position associated to it. Thus, we would need to map the video to it. Thus, it would be far easier to do so.  
// 	for _, video := range videos {
// 		videoMap[video.Id] = video
// 	}

// 	result := make([]PlaylistVideo, 0, len(items))

// 	for _, item := range items {

// 		if item == nil ||
// 			item.Snippet == nil ||
// 			item.ContentDetails == nil {
// 			continue
// 		}

// 		video := videoMap[item.ContentDetails.VideoId]

// 		if video == nil ||
// 			video.Snippet == nil ||
// 			video.ContentDetails == nil ||
// 			video.Statistics == nil ||
// 			video.Status == nil {
// 			continue
// 		}

// 		result = append(result, PlaylistVideo{
// 			Video: Video{
// 				ID:          video.Id,
// 				Title:       video.Snippet.Title,
// 				Description: video.Snippet.Description,

// 				ChannelID:    video.Snippet.ChannelId,
// 				ChannelTitle: video.Snippet.ChannelTitle,

// 				ThumbnailURL: getVideoThumbnail(video),

// 				PublishedAt: parseTime(video.Snippet.PublishedAt),

// 				Duration: video.ContentDetails.Duration,

// 				ViewCount:    video.Statistics.ViewCount,
// 				LikeCount:    video.Statistics.LikeCount,
// 				CommentCount: video.Statistics.CommentCount,

// 				Tags: video.Snippet.Tags,

// 				CategoryID: video.Snippet.CategoryId,

// 				CaptionAvailable: video.ContentDetails.Caption == "true",

// 				PrivacyStatus:       video.Status.PrivacyStatus,
// 				LiveBroadcastStatus: video.Snippet.LiveBroadcastContent,
// 			},

// 			Position: int(item.Snippet.Position),
// 			AddedAt:  parseTime(item.Snippet.PublishedAt),
// 		})
// 	}

// 	return result
// }