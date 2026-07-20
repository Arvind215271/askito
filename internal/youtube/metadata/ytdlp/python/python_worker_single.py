#!/usr/bin/env python3

import gc
import os
import struct
import sys

import orjson
import yt_dlp
import traceback


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
    "cachedir": True,
    "check_formats": False,
    "ignore_no_formats_error": True,
    "socket_timeout": 10,
    "extractor_retries": 0,
    "retries": 0,
    "lazy_playlist": True,
    "noprogress": True,
    # "cookiefile": "/home/arvind-saini/Projects/askito/cookies.txt",    
    "extractor_args": {
        "youtube": {
            "player_client": ["android"],
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
            traceback.print_exc(file=sys.stderr)
            break

        try:

            cmd = req.get("cmd")

            if cmd == "playlist":

                playlist_id = req["playlist_id"]

                playlist_url = (
                    "https://www.youtube.com/playlist?list="
                    + playlist_id
                )

                old_extract_flat = ydl.params.get("extract_flat")

                try:
                    ydl.params["extract_flat"] = True

                    info = ydl.extract_info(
                        playlist_url,
                        download=False,
                    )

                finally:
                    if old_extract_flat is None:
                        ydl.params.pop("extract_flat", None)
                    else:
                        ydl.params["extract_flat"] = old_extract_flat

                send({
                    "ok": True,
                    "data": info,
                })

                del info

                continue

            if cmd == "warmup":
                ydl.extract_info(WARMUP_URL, download=False)
                send({"ok": True})
                continue

            if cmd == "subtitle":

                video_id = req["video_id"]
                language = req.get("language", "en")
                subtitle_type = req.get("type", "manual")
                fmt = req.get("format", "json3")
                output_path = req.get("output_path")

                if not output_path:
                    raise RuntimeError("output_path missing in subtitle request")

                output_dir = os.path.dirname(output_path)
                os.makedirs(output_dir, exist_ok=True)

                old_outtmpl = dict(ydl.params["outtmpl"])

                outtmpl = dict(old_outtmpl)
                outtmpl["subtitle"] = output_path

                try:
                    ydl.params["writesubtitles"] = subtitle_type == "manual"
                    ydl.params["writeautomaticsub"] = subtitle_type == "automatic"
                    ydl.params["subtitleslangs"] = [language]
                    ydl.params["subtitlesformat"] = fmt
                    ydl.params["outtmpl"] = outtmpl

                    ydl.download([YT_URL + video_id])

                finally:
                    ydl.params["outtmpl"] = old_outtmpl

                send({
                    "ok": True,
                })

                continue

            info = ydl.extract_info(
                YT_URL + req["video_id"],
                download=False,
            )

            send({
                "ok": True,
                "data": info,
            })

            del info

        except Exception:
            send({
                "ok": False,
                "error": traceback.format_exc(),
            })

        finally:

            del req

            requests += 1

            if requests >= 50:
                gc.collect()
                requests = 0


if __name__ == "__main__":
    main()
