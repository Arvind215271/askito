# Export TODO

## Phase 1: JSON Export

- [x] Export normalized data as JSON
  Use JSON as the first export format since it directly maps to the internal data structure.

- [x] Support single playlist exports
  Allow users to export metadata from a single playlist.

- [x] Support single video exports
  Allow users to export metadata from a single video.

- [x] Support field selection
  Export only the fields requested by the user.

---

## Phase 2: CSV Export

- [ ] Export normalized data as CSV
  Convert metadata into a spreadsheet-friendly format.

- [ ] Flatten nested structures
  Ensure playlist and video data can be represented as rows.

- [ ] Support field selection
  Include only the requested columns in the exported file.

---

## Phase 3: Multiple Export Support

- [ ] Export multiple playlists in a single request
  Support batch exports without requiring multiple downloads.

- [ ] Export multiple videos in a single request
  Allow batch exports for video metadata.

- [ ] Export mixed inputs
  Support exports containing both videos and playlists.

- [ ] Generate ZIP archives
  Package multiple exported files into a single download.

- [ ] Support merged exports
  Combine all extracted data into a single export file when requested.

---

## Phase 4: Additional Export Formats

- [ ] Markdown export
  Generate readable markdown documents from extracted data.

- [ ] HTML export
  Generate browser-friendly reports.

- [ ] XML export
  Support structured XML exports.

- [ ] YAML export
  Support YAML exports for configuration and tooling workflows.

- [ ] Text export
  Generate simple text-based exports.

---

## Phase 5: Document-Based Exports

- [ ] Excel export
  Generate spreadsheet exports for larger datasets.

- [ ] Word export
  Generate editable document exports.

- [ ] SQLite export
  Export extracted data into a portable database file.

---

## Phase 6: Enhanced Metadata Exports

- [ ] Export native chapter information
  Include creator-provided chapters when available.

- [ ] Export generated chapter information
  Include AI-generated chapters derived from transcripts.

- [ ] Export transcript content
  Include extracted transcript data.

- [ ] Export translated transcripts
  Include normalized English transcripts when available.

- [ ] Export future metadata fields
  Ensure new enrichment fields can be exported without redesigning the export layer.
