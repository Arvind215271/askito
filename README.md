# Askito

Askito is a high-performance backend application designed for deep YouTube video analysis, transcript extraction, subtitle downloading, statistical signaling, and batch exports. It combines Go for robust API routing and concurrency with persistent Python worker pools running `yt-dlp` and `orjson` for high-speed data retrieval.

---

## System Requirements & Prerequisites

Make sure your machine has:
- **Go** (version 1.20 or higher)
- **Python** (version 3.8 or higher) with `pip` and `python3-venv`
- **FFmpeg** (Recommended for robust media handling by `yt-dlp`)

---

## Complete Installation & Setup Guide

### 1. Clone the Repository
```bash
git clone https://github.com/your-username/askito.git
cd askito
```

### 2. Go Dependencies Setup
Download and tidy all Go module dependencies:
```bash
go mod tidy
go mod download
```

### 3. Python Virtual Environment & Dependencies Setup
Askito communicates with a persistent Python worker pool ([`internal/youtube/metadata/ytdlp/python/python_worker_single.py`](internal/youtube/metadata/ytdlp/python/python_worker_single.py)) which requires specialized Python packages (`yt-dlp` and `orjson`).

Set up your virtual environment and install the required dependencies:

```bash
# Create a python virtual environment
python3 -m venv venv

# Activate the virtual environment
# On Linux/macOS:
source venv/bin/activate
# On Windows (Command Prompt / PowerShell):
# venv\Scripts\activate

# Upgrade pip and install required packages
pip install --upgrade pip
pip install yt-dlp orjson
```

> **Important**: Ensure your environment uses the virtual environment so that `python3` resolves to the venv where `yt-dlp` and `orjson` are installed.

### 4. Configuration (`.env`)
Copy the sample environment configuration file:
```bash
cp .env.sample .env
```

Open `.env` and customize your configuration parameters:
```env
APP_ENV=development
PORT=8080

# Optional: YouTube Data API v3 key if using the YouTube Data API provider
# YOUTUBE_API_KEY=your_api_key_here

# Cache directory configuration
YTDLP_CACHE_DIR=./.cache/ytdlp
YTDLP_CACHE_TTL_DAYS=28
YTDLP_CACHE_MAX_FILES=2000

# Number of concurrent Python workers for yt-dlp processing
# Each worker consumes ~50-60 MB RAM. Adjust based on your system capacity.
PYTHON_WORKERS=16
```

### 5. Running the Application

To run the application in development mode with live logs:
```bash
go run main.go
```

To build and run a compiled binary:
```bash
go build -o askito main.go
./askito
```

---

## API Endpoints, Parameters & Options Reference

Once the server is running at `http://localhost:8080`, you can access the following REST endpoints:

---

### 1. Video Metadata API (`/videos`)

#### Get Video By ID
- **Endpoint**: `GET /videos/id`
- **Query Parameters**:
  - `id` (string, required): YouTube Video ID (`VIDEO_ID`, e.g., `dQw4w9WgXcQ`).
  - `provider` (string, optional): Metadata provider (`ytdlp` or `youtube_api`).
- **Example Request**:
  ```bash
  curl -X GET "http://localhost:8080/videos/id?id=dQw4w9WgXcQ&provider=ytdlp"
  ```

#### Get Video By URL
- **Endpoint**: `GET /videos/url`
- **Query Parameters**:
  - `url` (string, required): Full YouTube video URL (`VIDEO_URL`, e.g., `https://www.youtube.com/watch?v=dQw4w9WgXcQ`).
  - `provider` (string, optional): Metadata provider (`ytdlp` or `youtube_api`).
- **Example Request**:
  ```bash
  curl -X GET "http://localhost:8080/videos/url?url=https://www.youtube.com/watch?v=dQw4w9WgXcQ&provider=ytdlp"
  ```

---

### 2. Subtitles API (`/subtitles`)

#### Get Subtitle Options
- **Endpoint**: `POST /subtitles/options`
- **Request Body Options (`JSON`)**:
  - `url` (string, required): YouTube video URL (`VIDEO_URL`).
- **Example Request**:
  ```bash
  curl -X POST "http://localhost:8080/subtitles/options" \
    -H "Content-Type: application/json" \
    -d '{"url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"}'
  ```

#### Download Subtitle
- **Endpoint**: `POST /subtitles/download`
- **Request Body Options (`JSON`)**:
  - `url` (string, required): YouTube video URL (`VIDEO_URL`).
  - `type` (string, required): Subtitle type (`manual` or `automatic`).
  - `language` (string, required): Subtitle/track language code (e.g., `en`).
  - `format` (string, optional): Desired output subtitle format (e.g., `vtt`, `json3`).
- **Example Request**:
  ```bash
  curl -X POST "http://localhost:8080/subtitles/download" \
    -H "Content-Type: application/json" \
    -d '{
      "url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
      "type": "manual",
      "language": "en",
      "format": "vtt"
    }'
  ```

---

### 3. Transcripts API (`/transcripts`)

#### Get Transcript
- **Endpoint**: `POST /transcripts`
- **Request Body Options (`JSON`)**:
  - `url` (string, required): YouTube video URL (`VIDEO_URL`).
  - `type` (string, required): Transcript type (`manual` or `automatic`).
  - `language` (string, required): Transcript language code (e.g., `en`).
- **Example Request**:
  ```bash
  curl -X POST "http://localhost:8080/transcripts" \
    -H "Content-Type: application/json" \
    -d '{
      "url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
      "type": "manual",
      "language": "en"
    }'
  ```

---

### 4. Signals API (`/signals`)

#### Get Video Signals
- **Endpoint**: `POST /signals`
- **Request Body Options (`JSON`)**:
  - `url` (string, required): YouTube video URL (`VIDEO_URL`).
  - `analysis` (string, required): Statistical analysis type (`word-stats` or `window-stats`).
  - `type` (string, required): Transcript type (`manual` or `automatic`).
  - `language` (string, required): Language code (e.g., `en`).
  - `use_heavy_stopwords` (boolean, optional): Whether to filter out heavy stopwords (`true` or `false`).
  - `min_freq` (integer, optional): Minimum occurrence frequency for words (`>= 0`).
  - `depth` (float, optional): Analysis depth threshold (`0.0` to `1.0`).
  - `window_size` (float, optional): Sliding window size (`> 0`).
  - `bucket_count` (integer, optional): Number of statistical buckets (`> 0`).
- **Example Request**:
  ```bash
  curl -X POST "http://localhost:8080/signals" \
    -H "Content-Type: application/json" \
    -d '{
      "url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
      "analysis": "word-stats",
      "type": "manual",
      "language": "en",
      "use_heavy_stopwords": true,
      "min_freq": 2
    }'
  ```

---

### 5. Export API (`/export`)

#### Export Single Video
- **Endpoint**: `POST /export/video`
- **Request Body Options (`JSON`)**:
  - `input` (string, required): YouTube Video URL or ID (`VIDEO_URL` or `VIDEO_ID`).
  - `format` (string, required): Export output format (`json`).
  - `fields` (array of strings, optional): Specific metadata/data fields to include.
  - `subtitle` (object, optional): Nested subtitle download options:
    - `url` (string)
    - `type` (string: `manual` | `automatic`)
    - `language` (string)
    - `format` (string)
  - `transcript` (object, optional): Nested transcript request options:
    - `url` (string)
    - `type` (string: `manual` | `automatic`)
    - `language` (string)
  - `signal` (object, optional): Nested signal analysis options:
    - `url` (string)
    - `analysis` (string: `word-stats` | `window-stats`)
    - `type` (string: `manual` | `automatic`)
    - `language` (string)
    - `use_heavy_stopwords` (boolean)
    - `min_freq` (integer)
- **Example Request**:
  ```bash
  curl -X POST "http://localhost:8080/export/video" \
    -H "Content-Type: application/json" \
    -d '{
      "input": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
      "format": "json"
    }'
  ```

#### Export Multiple Videos (Batch)
- **Endpoint**: `POST /export/videos`
- **Request Body Options (`JSON`)**:
  - `inputs` (array of strings, required): List of YouTube video URLs or IDs (`VIDEO_URL` / `VIDEO_ID`).
  - `format` (string, required): Export output format (`json`).
  - `fields`, `subtitle`, `transcript`, `signal` (optional, same nested structures as single video export).
- **Example Request**:
  ```bash
  curl -X POST "http://localhost:8080/export/videos" \
    -H "Content-Type: application/json" \
    -d '{
      "inputs": [
        "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
      ],
      "format": "json"
    }'
  ```

#### Export Playlist
- **Endpoint**: `POST /export/playlist`
- **Request Body Options (`JSON`)**:
  - `input` (string, required): YouTube Playlist URL or ID (`PLAYLIST_URL` or `PLAYLIST_ID`).
  - `format` (string, required): Export output format (`json`).
  - `fields`, `subtitle`, `transcript`, `signal` (optional, same nested structures as single video export).
- **Example Request**:
  ```bash
  curl -X POST "http://localhost:8080/export/playlist" \
    -H "Content-Type: application/json" \
    -d '{
      "input": "https://www.youtube.com/playlist?list=PL_EXAMPLE",
      "format": "json"
    }'
  ```
