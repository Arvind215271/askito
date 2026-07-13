#!/usr/bin/env python3

import argparse
import json
import logging
import os
import struct
import sys
import traceback

import yt_dlp

stdin = sys.stdin.buffer
stdout = sys.stdout.buffer

# ----------------------------------------------------------------------
# Logging
# ----------------------------------------------------------------------

# Using a directory relative to the workspace root for the logs
LOG_DIR = "debug/yt_bench/output/logs"
os.makedirs(LOG_DIR, exist_ok=True)

logger = None

def setup_logging():
    global logger
    logging.basicConfig(
        filename=os.path.join(LOG_DIR, "python_worker_batch.log"),
        level=logging.INFO,
        format="%(asctime)s [%(levelname)s] %(message)s",
    )
    logger = logging.getLogger(__name__)


def log(msg):
    logger.info("EVENT=" + msg)


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
    payload = read_exact(length)
    return json.loads(payload)


def send(obj):
    payload = json.dumps(obj, default=str).encode("utf-8")
    stdout.write(struct.pack(">I", len(payload)))
    stdout.write(payload)
    stdout.flush()

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--worker-id", help="Worker ID")
    args = parser.parse_args()
    worker_id = args.worker_id or "unknown"
    setup_logging()

    log(f"WORKER_STARTED worker={worker_id}")

    ydl = yt_dlp.YoutubeDL(
        {
            # "skip_download": True,
            # "quiet": True,
            "verbose": False,
            # "no_warnings": True,
            "cachedir": False,
            "check_formats": False,
            "ignore_no_formats_error": True,
            "extractor_args": {
                "youtube": {
                    "player_client": ["android_vr"],
                }
            },
        }
    )

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
            ids = req.get("ids")
            if ids is None:
                raise ValueError("No ids provided")

            logger.info("Received batch (%d videos)", len(ids))
            logger.info("IDs: %s", ",".join(ids))

            urls = " ".join(
                f"https://www.youtube.com/watch?v={vid}"
                for vid in ids
)

# info = ydl.extract_info(url, download=False)

            logger.info("Calling yt-dlp once for %d URLs", len(urls))

            info = ydl.extract_info(urls, download=False)

            logger.info("Returned object type: %s", type(info).__name__)

            if isinstance(info, dict):
                logger.info("Top-level keys: %s", list(info.keys()))

            results = []

            # Playlist-style return
            if isinstance(info, dict) and "entries" in info:
                entries = list(info.get("entries") or [])

                logger.info("Found %d entries", len(entries))

                for entry in entries:
                    if entry is None:
                        continue

                    results.append(
                        {
                            "ok": True,
                            "data": entry,
                        }
                    )

            # Single video
            elif isinstance(info, dict):
                logger.info("Single metadata object returned")
                results.append(
                    {
                        "ok": True,
                        "data": info,
                    }
                )

            # Unknown return type
            else:
                logger.warning("Unexpected return type: %r", type(info))
                logger.warning("Returned value: %r", info)

            returned = {
                r["data"]["id"]
                for r in results
                if r.get("ok") and "id" in r["data"]
            }

            missing = [video_id for video_id in ids if video_id not in returned]

            logger.info(
                "yt-dlp completed. Returned %d/%d metadata objects",
                len(returned),
                len(ids),
            )

            if missing:
                logger.warning("Missing videos: %s", ",".join(missing))

            send(results)

        except Exception:
            err = traceback.format_exc()
            log(f"REQUEST_FAILED worker={worker_id} error={err}")
            logger.error("Batch failed: %s", err)

            send(
                {
                    "ok": False,
                    "error": err,
                }
            )


if __name__ == "__main__":
    main()