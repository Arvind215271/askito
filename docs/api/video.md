# Video API TODO

## Phase 1: Video Metadata

- [x] Video endpoint
  Accept a YouTube video URL and return the normalized video model.

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

- [ ] Description processing endpoint
  Process the video description and return extracted metadata.

- [ ] Extract chapters
  Detect chapters from timestamp patterns.

- [ ] Extract links
  Return URLs found in the description.

- [ ] Extract contact information
  Return publicly available email addresses.

- [ ] Return cleaned description
  Remove extracted metadata while preserving meaningful text.

---

## Phase 4: Subtitle Metadata

- [ ] Subtitle endpoint
  Return subtitle information for the video.

- [ ] List available subtitle languages
  Return every available subtitle language.

- [ ] Return subtitle metadata
  Include language, source and availability information.

- [ ] Support translated subtitles
  Return available translation languages when supported.

---

## Phase 5: Transcript

- [ ] Transcript endpoint
  Return the transcript for a selected subtitle language.

- [ ] Support transcript language selection
  Allow users to request transcripts in a specific language.

- [ ] Return transcript metadata
  Include transcript language, source and availability.

- [ ] Return normalized transcript
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

- [ ] Transcript signal endpoint
  Analyze transcript content and return extracted signals.

- [ ] Keyword extraction
  Return important keywords from the transcript.

- [ ] Topic extraction
  Return detected topics discussed throughout the video.

- [ ] Word statistics
  Return word frequency and occurrence statistics.

- [ ] Window statistics
  Return window-based signal analysis.

- [ ] Concept extraction
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