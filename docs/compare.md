# Comparison TODO

## Phase 1: Topic Extraction

- [ ] Extract topics from available chapter data
  Use native chapters, description chapters, or generated chapters as the primary source of topics.

- [ ] Extract topics from transcript when chapters are unavailable
  Fall back to transcript analysis when chapter information does not exist.

- [ ] Normalize extracted topics
  Merge similar topics into a consistent representation.

- [ ] Attach topics to videos
  Store extracted topics as part of the normalized video structure.

---

## Phase 2: Video Comparison

- [ ] Compare topics between two videos
  Identify common and unique topics across both videos.

- [ ] Identify missing topics
  Determine which topics exist in one video but not the other.

- [ ] Calculate topic overlap
  Measure how similar two videos are based on covered topics.

- [ ] Generate comparison summary
  Provide a concise overview of the major similarities and differences.

---

## Phase 3: Playlist Comparison

- [ ] Aggregate topics across all videos in a playlist
  Build a playlist-level topic representation.

- [ ] Compare topic coverage between playlists
  Identify overlap and gaps between learning paths.

- [ ] Identify unique topics per playlist
  Highlight content covered exclusively by one playlist.

- [ ] Calculate playlist similarity score
  Estimate how closely two playlists align in subject matter.

- [ ] Generate playlist comparison summary
  Provide a high-level overview of the comparison results.

---

## Phase 4: Coverage Analysis

- [ ] Identify beginner topics
  Detect foundational concepts covered within a playlist.

- [ ] Identify intermediate topics
  Detect concepts that build upon foundational knowledge.

- [ ] Identify advanced topics
  Detect more specialized or advanced material.

- [ ] Estimate learning depth
  Determine how deeply a playlist covers its topics.

- [ ] Identify coverage gaps
  Highlight areas that may be missing from a learning path.

---

## Phase 5: Learning Path Analysis

- [ ] Analyze topic ordering
  Evaluate how topics are introduced throughout a playlist.

- [ ] Detect prerequisite relationships
  Identify concepts that should be learned before others.

- [ ] Identify learning progression
  Determine how knowledge develops throughout the playlist.

- [ ] Compare learning structures
  Evaluate differences in teaching flow between playlists.

---

## Phase 6: Recommendation Engine

- [ ] Recommend complementary playlists
  Suggest playlists that fill identified knowledge gaps.

- [ ] Recommend missing topics
  Highlight important concepts not covered by the selected playlist.

- [ ] Recommend alternative playlists
  Suggest playlists with similar coverage and learning goals.

- [ ] Generate learning recommendations
  Provide actionable next steps based on comparison results.
