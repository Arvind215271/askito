#!/usr/bin/env python3

import json
import struct
import sys
import traceback

import yt_dlp

stdin = sys.stdin.buffer
stdout = sys.stdout.buffer
stderr = sys.stderr


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
    ydl = yt_dlp.YoutubeDL(
        {
            "skip_download": True,
            "quiet": True,
            "verbose": False,
            "no_warnings": True,
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
            break
        except Exception:
            traceback.print_exc(file=stderr)
            break

        try:
            if "ids" in req:
                results = []
                for video_id in req["ids"]:
                    info = ydl.extract_info(
                        f"https://www.youtube.com/watch?v={video_id}",
                        download=False,
                    )
                    results.append({"ok": True, "data": info})
                send(results)
            elif "id" in req:
                video_id = req["id"]
                info = ydl.extract_info(
                    f"https://www.youtube.com/watch?v={video_id}",
                    download=False,
                )
                send({"ok": True, "data": info})
            else:
                raise ValueError("Neither id nor ids provided")

        except Exception:
            send(
                {
                    "ok": False,
                    "error": traceback.format_exc(),
                }
            )


if __name__ == "__main__":
    main()