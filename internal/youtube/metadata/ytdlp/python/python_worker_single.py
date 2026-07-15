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

                # Check for CacheDir passed in the request
                if "cache_dir" not in req:
                    raise RuntimeError("cache_dir missing in subtitle request")

                cache_dir = req["cache_dir"]
                video_dir = os.path.join(cache_dir, video_id)
                os.makedirs(video_dir, exist_ok=True)

                # old_write = ydl.params["writesubtitles"]
                # old_auto = ydl.params["writeautomaticsub"]
                # old_langs = ydl.params["subtitleslangs"]
                # old_fmt = ydl.params["subtitlesformat"]
                old_outtmpl = dict(ydl.params["outtmpl"])

                outtmpl = dict(old_outtmpl)
                outtmpl["default"] = os.path.join(video_dir, "%(id)s")

                try:
                    ydl.params["writesubtitles"] = subtitle_type == "manual"
                    ydl.params["writeautomaticsub"] = subtitle_type == "automatic"
                    ydl.params["subtitleslangs"] = [language]
                    ydl.params["subtitlesformat"] = fmt
                    ydl.params["outtmpl"] = outtmpl

                    ydl.download([YT_URL + video_id])

                finally:
                    # ydl.params["writesubtitles"] = old_write
                    # ydl.params["writeautomaticsub"] = old_auto
                    # ydl.params["subtitleslangs"] = old_langs
                    # ydl.params["subtitlesformat"] = old_fmt
                    ydl.params["outtmpl"] = old_outtmpl

                downloaded = None

                for name in os.listdir(video_dir):
                    if (
                        name.startswith(video_id + ".")
                        and name.endswith("." + fmt)
                    ):
                        downloaded = os.path.join(video_dir, name)
                        break

                if downloaded is None:
                    raise RuntimeError("subtitle file not created")

                final_name = f"subtitles.{language}.{fmt}"
                final_path = os.path.join(video_dir, final_name)

                if downloaded != final_path:
                    os.replace(downloaded, final_path)

                send({
                    "ok": True,
                    "filename": final_name,
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
