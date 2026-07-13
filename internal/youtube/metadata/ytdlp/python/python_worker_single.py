#!/usr/bin/env python3

import json
import logging
import os
import struct
import sys
import traceback
import orjson
import time

import yt_dlp

stdin = sys.stdin.buffer
stdout = sys.stdout.buffer

# ----------------------------------------------------------------------
# Logging
# ----------------------------------------------------------------------

LOG_DIR = "debug/yt_bench/output/logs"
os.makedirs(LOG_DIR, exist_ok=True)

logging.basicConfig(
    filename=os.path.join(LOG_DIR, "python_worker_single.log"),
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(message)s",
)

logger = logging.getLogger(__name__)


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
    start_send = time.time()
    payload = orjson.dumps(obj)
    stdout.write(struct.pack(">I", len(payload)))
    stdout.write(payload)
    stdout.flush()
    send_time = time.time() - start_send
    logger.info("Send time: %.4f seconds", send_time)


def main():
    logger.info("Python worker started")
    
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
    logger.info("Initialization time: %.4f seconds", init_time)

    while True:
        try:
            req = recv()
        except EOFError:
            logger.info("stdin closed")
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
                logger.info("Warmup finished in %.4f seconds", warmup_time)
                send({"ok": True})
                continue

            video_id = req.get("id")
            if video_id is None:
                raise ValueError("No id provided")

            logger.info("Processing video: %s", video_id)
            
            start_fetch = time.time()
            info = ydl.extract_info(
                f"https://www.youtube.com/watch?v={video_id}",
                download=False,
            )
            fetch_time = time.time() - start_fetch
            
            start_extract = time.time()
            extract_time = time.time() - start_extract
            
            logger.info("Finished metadata for %s. Fetch: %.4f, Extract: %.4f", video_id, fetch_time, extract_time)
            send({"ok": True, "data": info})

        except Exception:
            logger.exception("Failed to process command %s", req)
            send(
                {
                    "ok": False,
                    "error": traceback.format_exc(),
                }
            )


if __name__ == "__main__":
    main()
