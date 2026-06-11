# Description Processing TODO

## Phase 1: Description Extraction

- [x] Extract raw description content
  Store the original video description without modification.

- [x] Preserve original formatting
  Retain line breaks and structure for future processing.

---

## Phase 2: Chapter Detection

- [x] Detect timestamp patterns
  Identify timestamps that may represent chapter markers.
  There is a page on stackoverflow. Which is what i have used for reference as there were a lot of peps guiding on how they solved this problem. Most of it was extracting stuff using regex normally. 
  Because, it was done by peps. They must have already tested most of it as well with it.     

- [x] Extract chapter candidates
  Pair timestamps with surrounding text to create chapter entries.
  For this, we can simply extract the video line by line instead.
  We might miss detail that might be present if the timeline had breaks and spaces or like multiline.
  But currently the simplest thing is to use regex per line. 

- [x] Validate chapter sequence
  Ensure extracted chapters follow a valid chronological order.
  Aka the time must be incresaing.
  Also other stuff we can add is that atlesat 3 chapter must exist per youtube video.
  We cannot be sure if we extracted actual chapter or what. But whatever... 

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