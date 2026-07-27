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

- [x] Export normalized data as CSV
  Convert metadata into a spreadsheet-friendly format.

- [x] Flatten nested structures
  Ensure playlist and video data can be represented as rows.

- [x] Support field selection
  Include only the requested columns in the exported file.

---

## Phase 3: Multiple Export Support

- [x] Export multiple playlists in a single request
  Support batch exports without requiring multiple downloads.

- [x] Export multiple videos in a single request
  Allow batch exports for video metadata.

- [x] Export mixed inputs
  Support exports containing both videos and playlists.

- [x] Generate ZIP archives
  Package multiple exported files into a single download.

- [x] Support merged exports
  Combine all extracted data into a single export file when requested.

---

## Phase 4: Additional Export Formats

- [x] Markdown export
  Generate readable markdown documents from extracted data.

- [x] HTML export
  Generate browser-friendly reports.

- [x] XML export
  Support structured XML exports.

- [x] YAML export
  Support YAML exports for configuration and tooling workflows.

- [x] Text export
  Generate simple text-based exports.

---

## Phase 5: Document-Based Exports

- [x] Excel export
  Generate spreadsheet exports for larger datasets.

- [x] Word export
  Generate editable document exports.

- [x] SQLite export
  Export extracted data into a portable database file.

---

## Phase 6: Enhanced Metadata Exports

- [x] Export native chapter information
  Include creator-provided chapters when available.

- [x] Export generated chapter information
  Include AI-generated chapters derived from transcripts.

- [x] Export transcript content
  Include extracted transcript data.

- [x] Export translated transcripts
  Include normalized English transcripts when available.

- [x] Export future metadata fields
  Ensure new enrichment fields can be exported without redesigning the export layer.
