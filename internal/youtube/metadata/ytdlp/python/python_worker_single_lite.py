#!/usr/bin/env python3

import gc
import struct
import sys

import orjson
import yt_dlp

stdin = sys.stdin.buffer
stdout = sys.stdout.buffer

PACK = struct.Struct(">I")
pack = PACK.pack
unpack = PACK.unpack

YT_URL = "https://www.youtube.com/watch?v="
WARMUP_URL = YT_URL + "dQw4w9WgXcQ"

YDL_OPTS = {
    "skip_download": True,
    "quiet": True,
    "no_warnings": True,
    "ignoreconfig": True,
    "cachedir": False,
    "check_formats": False,
    "ignore_no_formats_error": True,
    "socket_timeout": 10,
    "extractor_retries": 0,
    "retries": 0,
    "lazy_playlist": True,
    "noprogress": True,
    "extractor_args": {
        "youtube": {
            "player_client": ["android_vr"],
        }
    },
}


def recv():
    hdr = bytearray(4)
    if stdin.readinto(hdr) != 4:
        raise EOFError

    size = unpack(hdr)[0]

    buf = bytearray(size)
    view = memoryview(buf)

    while view:
        n = stdin.readinto(view)
        if not n:
            raise EOFError
        view = view[n:]

    return orjson.loads(buf)


def send(obj):
    payload = orjson.dumps(obj)
    stdout.write(pack(len(payload)))
    stdout.write(payload)
    stdout.flush()


def main():
    ydl = yt_dlp.YoutubeDL(YDL_OPTS)

    gc.collect()
    gc.freeze()

    requests = 0

    while True:
        try:
            req = recv()
        except EOFError:
            break
        except Exception:
            break

        try:
            if req.get("cmd") == "warmup":
                ydl.extract_info(WARMUP_URL, download=False)
                send({"ok": True})
                continue

            info = ydl.extract_info(
                YT_URL + req["id"],
                download=False,
            )

            send({
                "ok": True,
                "data": info,
            })

            del info

        except Exception as exc:
            send({
                "ok": False,
                "error": str(exc),
            })

        finally:
            del req

            requests += 1
            if requests >= 50:
                gc.collect()
                requests = 0


if __name__ == "__main__":
    main()