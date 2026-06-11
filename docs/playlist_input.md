# Playlist Input TODO

## Phase 1: Single Playlist Support

- [x] Accept a single YouTube playlist URL  
  Start with the simplest case where input is only one playlist link.
  A Playlist consist of ?list= query params. So, we must use that ig.
  Also sometimes a video might be playing from a playlist. so we might end up needing to extract that as well (but do we extract the video or playlist? IDK) 

- [x] Validate YouTube playlist domain and format  
  Ensure the URL is a valid YouTube link and contains a playlist identifier.
  There are a lot of youtube domain. Some of them have been listed in the txt file in this same folder. Make sure the top level domain is youtube and or youtu.be that is the case.

- [x] Detect playlist type from URL  
  Confirm the input is a playlist using URL structure (`list=` parameter).

- [x] Extract playlist ID  
  Pull out the playlist identifier for downstream processing.

- [x] Normalize playlist input  
  Convert raw URL into a consistent internal format used by the system.
  Make them all have youtube.com instead of the country specific domain. That is all.

---

## Phase 2: Add Single Video Support

- [x] Accept a single YouTube video URL  
  Extend input handling to support individual video links.

- [x] Detect video type from URL  
  Identify video URLs using patterns like `watch?v=` or `youtu.be`.
  Also one more thing, YT shorts can also be used with watch?v=.

- [x] Extract video ID  
  Get the unique video identifier for processing.

- [x] Extend normalization structure  
  Support both video and playlist in the same internal format.

- [x] Keep type distinction intact  
  Ensure system knows whether input is a video or playlist after normalization.

---

## Phase 3: Multiple URL Support

- [ ] Accept multiple YouTube URLs in one request  
  Allow users to submit multiple inputs at once.

- [ ] Support newline-separated input format  
  Each line represents one URL for simple parsing.

- [ ] Process mixed inputs (video + playlist)  
  Handle combinations of both types in a single request.

- [ ] Apply same validation rules per input  
  Each URL is validated and processed independently.

- [ ] Normalize all inputs into a single list  
  Final output becomes a unified structured list of video and playlist items.

- [ ] Preserve type metadata for each item  
  Ensure pipeline can distinguish how each item should be processed.