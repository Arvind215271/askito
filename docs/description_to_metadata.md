# Description Processing TODO

## Phase 1: Description Extraction

- [ ] Extract raw description content
  Store the original video description without modification.

- [ ] Preserve original formatting
  Retain line breaks and structure for future processing.

---

## Phase 2: Chapter Detection

- [ ] Detect timestamp patterns
  Identify timestamps that may represent chapter markers.

- [ ] Extract chapter candidates
  Pair timestamps with surrounding text to create chapter entries.

- [ ] Validate chapter sequence
  Ensure extracted chapters follow a valid chronological order.

- [ ] Attach detected chapters to video metadata
  Include description-derived chapters in the normalized structure.

---

## Phase 3: Link Extraction

- [ ] Extract URLs from description
  Identify and collect all links present in the description.

- [ ] Remove duplicate links
  Prevent duplicate entries from appearing in extracted metadata.

- [ ] Preserve link context
  Store nearby text that may describe the purpose of each link.

- [ ] Attach extracted links to video metadata
  Include link information within the normalized structure.

---

## Phase 4: Contact Information Extraction

- [ ] Extract email addresses
  Identify publicly shared contact information.

- [ ] Validate extracted email addresses
  Filter obvious invalid matches.

- [ ] Attach contact information to video metadata
  Include extracted contact details in the normalized structure.

---

## Phase 5: Description Cleanup

- [ ] Remove extracted metadata from description
  Separate chapters, links, and contact information from the remaining content.

- [ ] Generate cleaned description content
  Produce a reduced version focused on meaningful text.

- [ ] Preserve original description
  Ensure the raw description remains available for future processing.