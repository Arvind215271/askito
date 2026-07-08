# Video API TODO

## Phase 1: Video Metadata

- [x] Video endpoint
  Accept a YouTube video URL and return the normalized video model.
  We are using our input method that we had developed to process the input but for a single URL only... that is all.

- [x] Video ID endpoint
  Accept a YouTube video ID and return the normalized video model.
  

- [x] Normalize YouTube metadata
  Convert the YouTube API response into the project's internal schema.

---

## Phase 2: Field Selection

- [ ] Support field selection
  Return only the metadata fields requested by the user.

- [ ] Support nested field selection
  Allow users to request specific nested objects.

- [ ] Validate requested fields
  Ignore or report unsupported field names.

---

## Phase 3: Description Processing

- [x] Description processing endpoint
  Process the video description and return extracted metadata.

- [x] Extract chapters
  Detect chapters from timestamp patterns.

- [x] Extract links
  Return URLs found in the description.

- [x] Extract contact information
  Return publicly available email addresses.

- [x] Return cleaned description
  Remove extracted metadata while preserving meaningful text.

---

## Phase 4: Subtitle Metadata

- [x] Subtitle endpoint
  Return subtitle information for the video.

- [x] List available subtitle languages
  Return every available subtitle language.

- [x] Return subtitle metadata
  Include language, source and availability information.

- [x] Support translated subtitles
  Return available translation languages when supported.

- [ ] Fix fixed format as youtube supports 7 formats currently. To reduce repeated data in JSON response.  
---

## Phase 5: Transcript

- [x] Transcript endpoint
  Return the transcript for a selected subtitle language.

- [x] Support transcript language selection
  Allow users to request transcripts in a specific language.

- [x] Return transcript metadata
  Include transcript language, source and availability.

- [x] Return normalized transcript
  Return transcript data using the project's internal transcript model.

---

## Phase 6: Chapter Generation

- [ ] Chapter generation endpoint
  Generate chapters from transcript content.

- [ ] Return generated chapters
  Return timestamps and titles for generated chapters.

- [ ] Preserve chapter source
  Distinguish between description chapters and generated chapters.

- [ ] Return unified chapter model
  Normalize all chapter sources into a common structure.

---

## Phase 7: Transcript Signals

- [x] Transcript signal endpoint
  Analyze transcript content and return extracted signals.

- [ ] Keyword extraction (NOT DOING) 
  Return important keywords from the transcript.

- [ ] Topic extraction (NOT DOING)
  Return detected topics discussed throughout the video.

- [x] Word statistics
  Return word frequency and occurrence statistics.

- [x] Window statistics
  Return window-based signal analysis.

- [ ] Concept extraction (NOT DOING)
  Return higher-level concepts identified from transcript content.

---

## Phase 8: Export

- [ ] JSON export

- [ ] CSV export

- [ ] Markdown export

- [ ] HTML export

- [ ] Excel export

- [ ] SQLite export

- [ ] Custom field export
  Export only the fields requested by the user.

- [ ] Batch export
  Export multiple video analyses in a single request.

---

## Phase 9: Cache

- [ ] Add cache at metadata and subtitle
  These endpoints are going to be reasked unless we maintain session ID but works well tho.

- [ ] Add file downloading at each of the endpoint and reuse them.
  Then flush those file to be automatically deleted after a period of days like a 28 days or 180 days. To reduce file size being stored.

