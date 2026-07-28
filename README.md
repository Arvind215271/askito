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

## API Endpoints, Parameters & Models Reference

Once the server is running at `http://localhost:8080`, you can access the following REST endpoints exposed by [`internal/api`](internal/api):

---

### 1. Export API (`/export`)

#### Export Resources
- **Endpoint**: `POST /export`
- **Request Body Model**: [`internal/api/export/models.go`](internal/api/export/models.go)
  ```json
  {
    "inputs": ["string"],
    "subtitle": {
      "video_id": "string",
      "language": "string",
      "type": "string",
      "format": "string"
    },
    "preferences": [
      {
        "language": "string",
        "type": "string" // options: "manual", "automatic", "manual>automatic", "automatic>manual"
      }
    ],
    "transcript": {
      "language": "string",
      "type": "string"
    },
    "signal": {
      "analysis": "string", // e.g. "word-stats", "window-stats"
      "language": "string",
      "type": "string",
      "use_heavy_stopwords": true,
      "min_freq": 0
    },
    "format": "string", // options: "json", "csv", "markdown", "excel", "yaml", "xml"
    "fields": ["string"]
  }
  ```
- **Available Fields**:
  - Metadata: `id`, `errors`, `title`, `description`, `channel_id`, `channel_title`, `thumbnails`, `published_at`, `duration`, `duration_seconds`, `duration_minutes`, `duration_timestamp`, `view_count`, `like_count`, `comment_count`, `tags`, `category_id`, `caption_available`, `privacy_status`, `live_broadcast_status`
  - Description: `description_chapters`, `description_links`, `description_emails`, `description_cleaned`
  - Transcript/Subtitle: `transcript_text`, `transcript_signal`, `subtitle_metadata`

- **Example Request**:
  ```bash
  curl -X POST "http://localhost:8080/export" \
    -H "Content-Type: application/json" \
    -d '{
      "inputs": ["https://www.youtube.com/watch?v=dQw4w9WgXcQ"],
      "format": "json",
      "fields": ["id", "title", "duration"]
    }'
  ```

---

### 2. Playlist API (`/playlist`)

#### Expand Playlist Videos
- **Endpoint**: `POST /playlist/videos`
- **Request Body Model**: [`internal/api/playlist/models.go`](internal/api/playlist/models.go)
  ```json
  {
    "url": "string",
    "provider": "string", // options: "ytdlp", "youtube_api"
    "output": "string"    // options: "id", "url", "both"
  }
  ```
- **Response Model**:
  ```json
  {
    "playlist_id": "string",
    "total": 0,
    "videos": [
      {
        "id": "string",
        "url": "string",
        "position": 0,
        "added_at": "string"
      }
    ]
  }
  ```
- **Example Request**:
  ```bash
  curl -X POST "http://localhost:8080/playlist/videos" \
    -H "Content-Type: application/json" \
    -d '{
      "url": "https://www.youtube.com/playlist?list=PL_EXAMPLE",
      "output": "both"
    }'
  ```

---

### 3. Subtitle API (`/subtitle`)

#### Get Subtitle Options
- **Endpoint**: `POST /subtitle/options`
- **Request Body Model**: [`internal/api/subtitle/models.go`](internal/api/subtitle/models.go)
  ```json
  {
    "inputs": ["string"]
  }
  ```
- **Example Request**:
  ```bash
  curl -X POST "http://localhost:8080/subtitle/options" \
    -H "Content-Type: application/json" \
    -d '{
      "inputs": ["https://www.youtube.com/watch?v=dQw4w9WgXcQ"]
    }'
  ```

#### Download Subtitle
- **Endpoint**: `POST /subtitle/download`
- **Request Body Model**: [`internal/api/subtitle/models.go`](internal/api/subtitle/models.go)
  ```json
  {
    "inputs": ["string"],
    "preferences": [
      {
        "language": "string",
        "type": "string" // options: "manual", "automatic", "manual>automatic", "automatic>manual"
      }
    ],
    "format": "string" // options: "json3", "vtt", "srt", "srv1", "srv2", "srv3", "ttml"
  }
  ```
- **Example Request**:
  ```bash
  curl -X POST "http://localhost:8080/subtitle/download" \
    -H "Content-Type: application/json" \
    -d '{
      "inputs": ["https://www.youtube.com/watch?v=dQw4w9WgXcQ"],
      "format": "vtt"
    }'
  ```

---

### 4. Transcript API (`/transcript`)

#### Get Transcript
- **Endpoint**: `POST /transcript`
- **Request Body Model**: [`internal/api/transcript/models.go`](internal/api/transcript/models.go)
  ```json
  {
    "inputs": ["string"],
    "preferences": [
      {
        "language": "string",
        "type": "string" // options: "manual", "automatic", "manual>automatic", "automatic>manual"
      }
    ]
  }
  ```
- **Example Request**:
  ```bash
  curl -X POST "http://localhost:8080/transcript" \
    -H "Content-Type: application/json" \
    -d '{
      "inputs": ["https://www.youtube.com/watch?v=dQw4w9WgXcQ"]
    }'
  ```
