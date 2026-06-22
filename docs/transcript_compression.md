# Transcript Compression Research

This document records the different transcript compression formats explored during development.

---

# Light Stop Word Reduction

## Process

- Convert transcript to lowercase.
- Remove punctuation.
- Split into words.
- Remove common stop words.
- Join remaining words.

## Output

```
This operating system course will teach you about process management.

↓

Operating System Course Teach Process Management
```

## Results

- Character Reduction: **30-40%**
- Word Reduction: **20-30%**

## Notes

**Advantages**
- Very simple and fast.
- Preserves transcript order.

**Drawbacks**
- Removes only common stop words.
- Still contains significant repetition.

---

# Heavy Stop Word Reduction

## Process

- Convert transcript to lowercase.
- Remove punctuation.
- Split into words.
- Remove words from the heavy stop-word dictionary (~1500 words).
- Join remaining words.

## Output

```
This operating system course will teach you everything about process management,
memory management and CPU scheduling.

↓

operating system course teach process management memory management cpu scheduling
```

## Results

- Character Reduction: **50-60%**
- Word Reduction: **30-40%**

## Notes

**Advantages**
- Better compression than the light version.
- Retains most technical concepts.

**Drawbacks**
- Technical words are still repeated many times.
- Long lectures remain highly redundant.


# Unique Word Reduction

## Process

- Convert transcript to lowercase.
- Remove punctuation.
- Split into words.
- Keep only the first occurrence of each word.
- Join remaining words.

## Output

```
machine learning machine learning regression regression decision tree

↓

machine learning regression decision tree
```

## Results

Typical reduction

- Character Reduction: **94–97%**
- Unique Words: **3k–5k**

## Notes

**Advantages**

- Extremely high compression.
- Captures the vocabulary covered by the transcript.

**Drawbacks**

- Loses all word frequency information.
- Loses transcript progression and timeline.


# Heavy Stop Word + Unique Words

## Process

- Convert transcript to lowercase.
- Remove punctuation.
- Remove heavy stop words.
- Keep only the first occurrence of each remaining word.
- Join remaining words.

## Output

```
this machine learning course teaches machine learning algorithms
machine learning applications and machine learning projects

↓

machine learning course teaches algorithms applications projects
```

## Results

Typical reduction

- Character Reduction: **95–98%**
- Unique Words: **2.5k–4.5k**

## Notes

**Advantages**

- Very high compression.
- Retains most technical vocabulary while removing conversational words.

**Drawbacks**

- No timeline information.
- No concept frequency.
- Relationship between concepts is largely lost.

---

# Word Statistics (Full Frequency + Duration)

## Process

- Convert transcript to lowercase.
- Remove punctuation.
- Split transcript into words.
- Track:
  - Word frequency (count)
  - Total segment duration contribution per word
- Aggregate results across full transcript.
- Sort words by total duration (descending).

## Output

```

Word, Count, Duration
the, 4068, 136975s
you, 3132, 105575s
and, 2868, 96567s
is, 2775, 93500s
...
data, 699, 23611s
model, 419, 14099s
learning, 328, 11036s

```

## Results

Typical output characteristics:

- Word list size: **5k+ unique words (large lectures)**
- Character reduction: **~90–95% (vs raw transcript export)**

## Notes

**Advantages**

- Preserves full word frequency information.
- Captures semantic importance via duration weighting.
- Useful for ranking concepts by attention/time span.

**Drawbacks**

- Still relatively large output size.
- No explicit structure (flat list only).
- Difficult to interpret at scale without grouping.

---

# Heavy Word Statistics (Stop-word Filtered)

## Process

- Convert transcript to lowercase.
- Remove punctuation.
- Split into words.
- Remove heavy stop words (~1500-word dictionary).
- Track:
  - Word frequency
  - Duration contribution
- Aggregate across transcript.
- Sort by duration (descending).

## Output

```

Word, Count, Duration
data, 699, 23611s
model, 419, 14099s
learning, 328, 11036s
theta, 262, 8844s
cluster, 240, 8079s
```

## Results

Typical reduction

- Character Reduction: **93–96%**
- Unique Words: **2.5k–4.5k**

## Notes

**Advantages**

- Removes conversational noise.
- Preserves domain-specific technical vocabulary.
- Better signal-to-noise ratio than raw stats.

**Drawbacks**

- Still flat representation.
- No hierarchical structure.
- Duration still spread across noisy aggregation.

---

# Bucketed Word Statistics (Importance-Based Compression)

## Process

- Convert transcript to lowercase.
- Remove punctuation.
- Split into words.
- Remove heavy stop words.
- Compute per-word:
  - Count
  - Duration
  - Importance score:
    ```
    score = log(1 + count) * 2 + log(1 + duration)
    ```
- Normalize scores and assign words into **32 buckets**.
- Aggregate words inside each bucket.
- Sort buckets from high importance → low importance.

---

## Output Format

```

# BUCKETED TRANSCRIPT WORD STATISTICS (HIGH → LOW IMPORTANCE)

# FORMAT: bucketId,countMin-countMax,durationMin-durationMax,words(space separated)

31,699-699,23611-23611,data
28,419-419,14099-14099,model
27,262-328,8844-11036,theta learning
26,201-240,6782-8079,equals regression cluster set
25,158-196,5309-6609,entropy called decision training function algorithm
...
19,35-44,1176-1505,gain based plane divide learners kmeans overfitting

```

---

## Results

Typical reduction

- Character Reduction: **95–98%**
- Output size: **Highly compact structured representation**
- Bucket count: **32 fixed groups**

## Notes

**Advantages**

- Strong structural compression.
- Groups semantically similar importance levels.
- Great for LLM-friendly summarization.
- Reduces hallucination risk by ordering importance.

**Drawbacks**

- Loses exact word frequency detail inside buckets.
- No strict temporal sequence.
- Interpretation requires understanding bucket semantics.

---

# Window Word Statistics (300s Windowing)

## Process

- Segment the transcript into fixed 300-second windows.
- Clean text, lowercase, and split into individual words within each time bound.
- Calculate frequency counts and exact duration offsets for terms inside each window boundary independently.

## Output Format

Each window:
```text
WINDOW <id> (start - end)
word,count,duration
```

## Example Output

```text
WINDOW 0 (0s - 300s)
learning,29,11
machine,22,8
syllabus,9,3
start,7,2
applications,5,2
deep,5,2
youtube,4,2
absolutely,4,1

WINDOW 1 (300s - 600s)
house,11,4
input,11,4
data,8,3
definition,9,3
machine,9,3
learning,9,3
variable,7,3
experience,7,3

WINDOW 2 (600s - 900s)
learning,19,7
house,11,4
machine,8,3
algorithm,7,3
prediction,6,2
detection,6,2
price,6,2
```

## Results

- Character Reduction: **92-96%**
- Unique Words: **4k-8k** (Repeated in different window)

## Notes

### Advantages
- Strong temporal segmentation
- Preserves topic drift over time
- Good for lecture-style videos

### Drawbacks
- No concept-level grouping
- Repetition across windows
- Formatting overhead is high

---

# Concept Distillation (Heuristic Clustering)

## Process

- Scan the full transcript to find highly co-occurring terms.
- Extract descriptive text tags based on shared structural distributions.
- Group complementary semantic keywords together using adaptive threshold heuristics to establish localized conceptual markers.

## Output Format

```csv
concept,score,freq,topWords
```

## Example Output

```csv
learning,328,11036,learingin
machine,290,9800,machine
house,210,7400,house
theta,180,6200,theta
```

## Results

- Character Reduction: **90-92%**
- Concepts Extracted: **All but can be reduced by freq (varies by transcript)**

## Notes

### Advantages


### Drawbacks
- Same Results like Unique Word Frequency
- Heuristic-based (not true NLP clustering)
- Concept overlap issues
- No strict hierarchy

---

# Window Concept Fingerprint (v1-5)

## Process

- Set up a fixed window scale parameter (WS=300).
- Extract terms dynamically per temporal section.
- Stratify terms cleanly into a three-tier hierarchical frequency layer layout (L3 for core, L2 for mid-tier, L1 for specific accents).

## Formats

```
# WINDOW CONCEPT FINGERPRINT (COMPRESSED v2)
# WS=300 (window size in seconds)
# FORMAT:
#   #<window_id> [t=range]
#     L3: freq>=8
#     L2: freq=4-7
#     L1: freq=2-3
# NOTE: freq<=1 removed


#0 [0-300s]
L3:learning,machine,syllabus
L2:start,applications,deep,youtube,absolutely
L1:understanding,sections,lot,regression,algorithms,domains,basics,teach,theory,projects

```

```
# WINDOW CONCEPT FINGERPRINT (COMPRESSED)
# FORMAT: window -> top concepts only (NO RAW WORD LISTS)


W0 (0-300s)
learning(29) machine(22) syllabus(9) start(7) applications(5) deep(5) 
```

```
# WINDOW CONCEPT FINGERPRINT v3 (ULTRA-COMPRESSED)
# WS=300 seconds per window
#
# FORMAT:
#   #<window_id>
#     B<bucket>: word list
#
# BUCKETS (log2 based importance score):
#   B4 = core concepts (very frequent / dominant ideas)
#   B3 = strong context
#   B2 = medium context
#   B1 = weak signal
#   B0 = noise floor (rare signals, often dropped)
#
# NOTE:
# - No timestamps
# - No raw transcript
# - Words only (comma-separated)

#1
B4:learning
B3:machine,syllabus,start,applications,deep
B2:youtube,absolutely,understanding,sections,lot,regression,algorithms,domains
```

```
WS=300 seconds
FORMAT:
#<window_id>
<freq_group_1>
<freq_group_2>
<freq_group_3>

FREQ RULES:
- group1: freq>=8 (core concepts)
- group2: freq 4-7 (mid importance)
- group3: freq 2-3 (low signal)
- freq<=1 removed
- no labels in output (only raw word lines)

#1
learning,machine,syllabus
start,applications,deep,youtube,absolutely
understanding,sections,lot,regression,algorithms,domains,teach,theory,basics,projects,ayush,scientist,data,talk,yeah,channel,content,divided,friends,level,advanced,basic,topics,teaching,assess,assignments,fundamentals,types,cover,visiting,hope,discuss,description,box,introduction,platform,shot
```

```
WS=300
FORMAT:
#<window_id>
line1 = freq>=8 (core)
line2 = freq 4-7 (mid)
line3 = freq 2-3 (low)
NOTE: comma-separated words only

#1
learning,machine,syllabus
start,applications,deep,youtube,absolutely
understanding,sections,lot,regression,algorithms,domains,basics,theory,teach,projects,data,ayush,scientist,talk,yeah,channel,content,divided,friends,level,advanced,basic,topics
```

## Example Output

```text
#1
learning,machine,syllabus
start,applications,deep,youtube
understanding,sections,regression,algorithms

#2
house,input,definition,machine,learning,data
variable,experience,function,size,price
```

## Results

- Character Reduction: **94–98%**
- Strong semantic compression
- Maintains window-level interpretability

## Notes

### Advantages
- Very compact representation
- Hierarchical signal preserved (L1/L2/L3)
- Best balance of compression + readability

### Drawbacks
- Still heuristic-based grouping
- Loses exact frequency precision
- Slight concept drift across windows

---

# System Logs

```json
{
  "playlist_id": "PLKnIA16_Rmvbr7zKYQuBfsVkjoLcJgxHH",
  "video_count": 134,
  "export_size_kb": 24.32
}
```


### Video NWONeJKn6kc
- **Windows:** 119
- **Reduction:** 95.37%

### Video yK1uBHPdp30
- **Windows:** 299
- **Reduction:** 94.30%



# Window Fingerprint CodeMap + Stream Encoding (v9)

Introduces dictionary encoding + binary/base64 compression.

## Core Idea

- Build global dictionary (freq ≥ 2)
- Assign uint16 IDs to words
- Encode windows using IDs
- Serialize using:
  - binary (little-endian)
  - base64 encoding
- Add repetition cap per word

## Output Structure


```
WS=<window_size>
D=<dictionary_size>

CODEMAP:
id=word; id=word; ...

WINDOWS:
#<window_id>
<base64_high>|<base64_mid>
```


## Improvements

- Compression (but usually same size as without compression)
- Efficient repetition handling

## Limitations

- Requires decoding step
- Codebook overhead
- Less LLM-friendly