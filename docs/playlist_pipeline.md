# Playlist Pipeline TODO

## Phase 1: Basic YouTube Metadata

- [x] Fetch playlist metadata
  Retrieve basic playlist information directly from the YouTube API.
  Youtube API have 3 main types of data that we can use ie Playlist, PlaylistItems, and Videos.
  Playlist contain information about the playlist only (no Video ID). It contain channel name, channel ID, Channel Title, Playlist Descritption, + locallisted content (same title and descirption but in local language), status (public, private, etc).
  Playlist item is what we need here. It contain videoID + publishtedAt. It contain other info like channel ID, channel title, videoTitle, video Descriptoin, playlistID, position, etc.

   
    

- [x] Fetch video metadata
  Retrieve the videos contained within a playlist.
  Again, as explained. The first would be getting the Data about the playlist, second would be getting the playlistItems list which contain all the ID of the video (youtube API can get us 0 to 50 video per call. So, we might have to use for loop or recursively call to get all the video ID). Then we need to call the for getting data for the video (again 0 to 50 ID can be sent per call). 

- [x] Extract basic video information
  Gather fields such as title, video URL, description, thumbnail, channel name, duration, views, and upload date.
  We do not really need everything from the call. We can give what fields to give (but they are combined so... we will get unnecssary data anyway). 
  So get those data and use only the once that is required by us

- [x] Normalize metadata format
  Convert API responses into a consistent internal structure.
  We will be using all these struct everywhere.
  In case we might go and end up using other service like yt-dlp. That would get us different data type.
  Thus, having a consistent data type for the internal working of the system is required and needed.

- [x] Pass normalized data to processing layer
  Hand off the collected metadata for formatting and export preparation.

---

## Phase 2: Chapter Extraction

- [x] Research available chapter sources
  Determine whether chapters can be obtained from the YouTube API or require alternative extraction methods.
  Currently, there is no option to get Chapters from Youtube API.
  We can use yt-dlp to get those... which be a little expensive to call for each video. (But single user can do that with their device)
  So... We might not be able to extract chapter for a full playlist. But we can do it for single videos. Because we might need more data to process a single video than it is for full video. That is all. 
  The simplest thing would be to give ADMIN the right to do everything they want... They can perform the most expensive operation.
  For normal user... they won't have everything they want here. That is all.    

- [x] Extract chapter information
  Retrieve chapter titles and timestamps for videos that provide them. 
  We are going to use description to get those chapter. 

- [x] Attach chapters to video metadata
  Merge chapter information into the existing video structure.

- [x] Handle videos without chapters
  Ensure missing chapter data does not break the pipeline.

- [x] Pass enriched metadata to processing layer
  Continue using the same processing flow with additional chapter information.

---

## Phase 3: Transcript Extraction

- [ ] Research transcript availability
  Determine the most reliable method for obtaining subtitles or transcripts.

- [ ] Extract transcript content
  Retrieve subtitle or transcript text for supported videos.

- [ ] Extract transcript metadata
  Capture information such as transcript language and availability.

- [ ] Attach transcripts to video metadata
  Include transcript content within the normalized structure.

- [ ] Handle unavailable transcripts
  Gracefully continue processing when transcripts are disabled or unavailable.

- [ ] Pass enriched metadata to processing layer
  Continue pipeline processing with transcript support enabled.

---

## Phase 4: Transcript Normalization

- [ ] Detect transcript language
  Identify the original language of the extracted transcript.

- [ ] Preserve original transcript
  Store the original transcript without modification.

- [ ] Translate transcript to English
  Generate an English version to create a consistent internal format.

- [ ] Store translation metadata
  Keep track of source language and translation status.

- [ ] Attach normalized transcript data
  Merge translated and original transcript data into the existing structure.

- [ ] Pass normalized transcript to processing layer
  Ensure all downstream systems can operate on a consistent language.

---

## Phase 5: Chapter Generation

- [ ] Research chapter generation approaches
  Evaluate methods for generating chapters from transcript content.

- [ ] Generate chapter boundaries
  Identify logical section breaks based on transcript context.

- [ ] Generate chapter titles
  Create meaningful titles for each generated chapter.

- [ ] Attach generated chapters to video metadata
  Merge generated chapter information into the normalized structure.

- [ ] Handle videos without transcripts
  Ensure chapter generation degrades gracefully when transcripts are unavailable.

- [ ] Pass enriched metadata to processing layer
  Continue processing with generated chapter support.

## Phase 6: Additional Metadata Research

- [ ] Investigate AI-generated summaries
  Determine whether YouTube-generated summaries can be accessed reliably.

- [ ] Investigate additional video insights
  Explore other metadata sources that may provide useful information.

- [ ] Evaluate usefulness versus extraction cost
  Ensure new fields provide enough value to justify pipeline complexity.

- [ ] Integrate approved metadata fields
  Extend the existing structure without breaking previous phases.