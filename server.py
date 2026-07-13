#!/usr/bin/env python3

import json
import time

from yt_dlp import YoutubeDL

opts = {
    "source_address": "0.0.0.0",
    "skip_download": True,
    "simulate": True,
    "quiet": False,
    "no_warnings": False,
    "verbose": True,

    "extractor_args": {
        "youtube": {
            "player_client": ["android_vr"],
        }
    },

    "check_formats": False,
    "ignore_no_formats_error": True,
}

with YoutubeDL(opts) as ydl:
    info = ydl.extract_info("https://youtu.be/dQw4w9WgXcQ", download=False)


print("Ready", flush=True)

while True:
    try:
        url = input("> ").strip()
    except EOFError:
        break

    start = time.perf_counter()

    try:
        info = ydl.extract_info(url, download=False)
        elapsed = time.perf_counter() - start

        print(
            json.dumps(
                {
                    "title": info.get("title"),
                    "elapsed": elapsed,
                },
                indent=2,
            ),
            flush=True,
        )

    except Exception as e:
        print(e)