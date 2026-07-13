#!/usr/bin/env python3

import argparse
import json
import logging
import os
import struct
import sys
import traceback
import orjson
import time

import yt_dlp

# ----------------------------------------------------------------------
# Configuration
# ----------------------------------------------------------------------

WORKER_ID = "unknown"

stdin = sys.stdin.buffer
stdout = sys.stdout.buffer

# ----------------------------------------------------------------------
# Logging
# ----------------------------------------------------------------------

LOG_DIR = "debug/yt_bench/output/logs"
os.makedirs(LOG_DIR, exist_ok=True)

logger = None

def setup_logging():
    global logger
    log_path = os.path.join(LOG_DIR, "python_worker_single.log")

    logging.basicConfig(
        filename=log_path,
        level=logging.INFO,
        format="%(asctime)s [%(levelname)s] %(message)s",
    )
    logger = logging.getLogger(__name__)


def log(msg):
    logger.info(f"worker={WORKER_ID} {msg}")


def log_request_failed(video_id, error):
    log(f"REQUEST_FAILED video={video_id}")
    logger.error(error)


def log_worker_error(error):
    log("WORKER_ERROR")
    logger.error(error)


def log_worker_started():
    log("WORKER_STARTED")


def log_worker_initialized(duration):
    log(f"INITIALIZED duration={duration:.4f}")


def log_worker_warmup(duration):
    log(f"WARMUP_DONE duration={duration:.4f}")


def log_request_start(video_id):
    log(f"REQUEST_START video={video_id}")


def log_request_end(video_id):
    log(f"REQUEST_END video={video_id}")


def log_worker_shutdown():
    log("WORKER_STOPPED")


def read_exact(n: int) -> bytes:
    buf = bytearray()
    while len(buf) < n:
        chunk = stdin.read(n - len(buf))
        if not chunk:
            raise EOFError("stdin closed")
        buf.extend(chunk)
    return bytes(buf)


def recv():
    header = read_exact(4)
    length = struct.unpack(">I", header)[0]
    return orjson.loads(read_exact(length))


def send(obj):
    payload = orjson.dumps(obj)
    stdout.write(struct.pack(">I", len(payload)))
    stdout.write(payload)
    stdout.flush()


def main():
    global WORKER_ID
    parser = argparse.ArgumentParser()
    parser.add_argument("--worker-id", help="Worker ID")
    parser.add_argument("pos_worker_id", nargs="?", help="Worker ID")
    args = parser.parse_args()
    WORKER_ID = args.worker_id or args.pos_worker_id or "unknown"
    setup_logging()

    log_worker_started()
    
    start_init = time.perf_counter()
    ydl = yt_dlp.YoutubeDL(
        {
            "skip_download": True,
            "quiet": True,
            "verbose": False,
            "no_warnings": True,
            "cachedir": False,
            "check_formats": False,
            "ignore_no_formats_error": True,
            "noplaylist": True,
            "writesubtitles": False,
            "writeautomaticsub": False,
            "writeinfojson": True,
            "writethumbnail": False,
            "write_all_thumbnails": False,
            "writedescription": False,
            "writeannotations": False,
            "extract_flat": False,
            "ignoreconfig": True,
            "socket_timeout": 10,
            "extractor_retries": 0,
            "retries": 0,
            # "fragment_retries": 0,
            "noplaylist": True,
            "extractor_args": {
                "youtube": {
                    "player_client": ["android_vr"],
                }
            },
        }
    )
    init_time = time.perf_counter() - start_init
    log_worker_initialized(init_time)

    while True:
        try:
            req = recv()
        except EOFError:
            log_worker_shutdown()
            break
        except Exception:
            logger.exception("Failed reading request")
            break

        try:
            cmd = req.get("cmd", "video")
            if cmd == "warmup":
                # logger.info("Warmup requested")
                start_warmup = time.time()
                ydl.extract_info("https://www.youtube.com/watch?v=dQw4w9WgXcQ", download=False)
                warmup_time = time.time() - start_warmup
                log_worker_warmup(warmup_time)
                send({"ok": True})
                continue

            video_id = req.get("id")
            if video_id is None:
                raise ValueError("No id provided")

            log_request_start(video_id)
            
            info = ydl.extract_info(
                f"https://www.youtube.com/watch?v={video_id}",
                download=False,
            )
            
            log_request_end(video_id)
            send({"ok": True, "data": info})

        except Exception:
            err = traceback.format_exc()
            log_request_failed(req.get("id", "unknown"), err)
            send(
                {
                    "ok": False,
                    "error": err,
                }
            )


if __name__ == "__main__":
    main()
