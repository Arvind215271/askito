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
Askito communicates with a persistent Python worker pool [`internal/youtube/metadata/ytdlp/python/python_worker_single.py`](internal/youtube/metadata/ytdlp/python/python_worker_single.py) which requires specialized Python packages (`yt-dlp` and `orjson`).

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

## API Endpoints Reference & Examples

Once the server is running at `http://localhost:8080`, you can access the following REST endpoints:

### 1. Videos
- **Get Video By ID** (`GET /videos/id`)
  - **Query Parameters**: `id` (string, required), `provider` (`ytdlp` or `youtube_api`, optional).
  - **Example**:
    ```bash
    curl -X GET "http://localhost:8080/videos/id?id=dQw4w9WgXcQ"
    ```

- **Get Video By URL** (`GET /videos/url`)
  - **Query Parameters**: `url` (string, required), `provider` (`ytdlp` or `youtube_api`, optional).
  - **Example**:
    ```bash
    curl -X GET "http://localhost:8080/videos/url?url=https://www.youtube.com/watch?v=dQw4w9WgXcQ"
    ```

---

### 2. Subtitles
- **Get Subtitle Options** (`POST /subtitles/options`)
  - **Request Body**:
    ```json
    {
      "url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
    }
    ```
  - **Example**:
    ```bash
    curl -X POST "http://localhost:8080/subtitles/options" \
      -H "Content-Type: application/json" \
      -d '{"url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"}'
    ```

- **Download Subtitle** (`POST /subtitles/download`)
  - **Request Body**:
    ```json
    {
      "url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
      "language": "en",
      "format": "vtt"
    }
    ```
  - **Example**:
    ```bash
    curl -X POST "http://localhost:8080/subtitles/download" \
      -H "Content-Type: application/json" \
      -d '{"url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ", "language": "en", "format": "vtt"}'
    ```

---

### 3. Transcripts
- **Get Transcript** (`POST /transcripts`)
  - **Request Body**:
    ```json
    {
      "url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
      "language": "en"
    }
    ```
  - **Example**:
    ```bash
    curl -X POST "http://localhost:8080/transcripts" \
      -H "Content-Type: application/json" \
      -d '{"url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ", "language": "en"}'
    ```

---

### 4. Signals & Statistics
- **Get Video Signals** (`POST /signals`)
  - **Request Body**:
    ```json
    {
      "url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
    }
    ```
  - **Example**:
    ```bash
    curl -X POST "http://localhost:8080/signals" \
      -H "Content-Type: application/json" \
      -d '{"url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"}'
    ```

---

### 5. Export
- **Export Single Video** (`POST /export/video`)
- **Export Multiple Videos (Batch)** (`POST /export/videos`)
- **Export Playlist** (`POST /export/playlist`)
  - **Example (Playlist Export)**:
    ```bash
    curl -X POST "http://localhost:8080/export/playlist" \
      -H "Content-Type: application/json" \
      -d '{"url": "https://www.youtube.com/playlist?list=PL_EXAMPLE"}'
    ```
