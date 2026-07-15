#!/usr/bin/env python3

import csv
import json
import os
import struct
import subprocess
import threading
import time


VIDEO_ID = "dQw4w9WgXcQ"
PLAYLIST_ID = "PL3tRBEVW0hiCjpQl2LpnE7IFS4ik3A0LX"
N = 50

WORKER = (
    "internal/youtube/metadata/ytdlp/python/"
    "python_worker_single.py"
)

CACHE_DIR = os.path.join(
    os.getcwd(),
    ".cache",
    "ytdlp",
)

os.makedirs(CACHE_DIR, exist_ok=True)


def send(proc, obj):
    payload = json.dumps(obj).encode()

    proc.stdin.write(
        struct.pack(">I", len(payload))
    )

    proc.stdin.write(payload)
    proc.stdin.flush()


def recv(proc):
    hdr = proc.stdout.read(4)

    if len(hdr) != 4:
        raise EOFError("worker closed stdout")

    length = struct.unpack(">I", hdr)[0]

    data = proc.stdout.read(length)

    if len(data) != length:
        raise EOFError("worker closed stdout before full response")

    return json.loads(data)


def sampler(pid, stop_event, rows):
    status = f"/proc/{pid}/status"

    while not stop_event.is_set():
        rss = None
        vsz = None

        try:
            with open(status) as f:
                for line in f:
                    if line.startswith("VmRSS:"):
                        rss = int(line.split()[1])

                    elif line.startswith("VmSize:"):
                        vsz = int(line.split()[1])

        except FileNotFoundError:
            break

        rows.append(
            (
                time.time(),
                rss,
                vsz,
            )
        )

        time.sleep(0.02)


def main():
    proc = subprocess.Popen(
        [
            "python3",
            WORKER,
            "--worker-id",
            "bench",
        ],
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
    )

    print("Worker PID:", proc.pid)

    rows = []

    stop = threading.Event()

    t = threading.Thread(
        target=sampler,
        args=(
            proc.pid,
            stop,
            rows,
        ),
        daemon=True,
    )

    t.start()

    start = time.perf_counter()

    for i in range(N):

        # --------------------------------------------------
        # 1. Fetch metadata
        # --------------------------------------------------

        send(
            proc,
            {
                "cmd": "metadata",
                "video_id": VIDEO_ID,
            },
        )

        metadata_resp = recv(proc)

        if not metadata_resp.get("ok"):
            print(
                f"Metadata request {i + 1} failed"
            )

            print(
                metadata_resp.get("error")
            )

            break

        metadata = metadata_resp["data"]

        print(
            f"{i + 1}/{N} metadata fetched: "
            f"{metadata.get('id')}"
        )


        # --------------------------------------------------
        # 2. Fetch subtitles
        # --------------------------------------------------

        send(
            proc,
            {
                "cmd": "subtitle",
                "video_id": VIDEO_ID,
                "language": "en",
                "type": "manual",
                "format": "json3",
                "cache_dir": CACHE_DIR,
            },
        )

        subtitle_resp = recv(proc)

        if not subtitle_resp.get("ok"):
            print(
                f"Subtitle request {i + 1} failed"
            )

            print(
                subtitle_resp.get("error")
            )

            break

        filename = subtitle_resp["filename"]

        path = os.path.join(
            CACHE_DIR,
            VIDEO_ID,
            filename,
        )

        with open(
            path,
            "r",
            encoding="utf-8",
        ) as f:
            preview = "".join(
                f.readlines()[:500]
            )

        print(
            f"{i + 1}/{N} subtitle fetched"
        )

        print(preview)


        # --------------------------------------------------
        # 3. Fetch playlist
        # --------------------------------------------------

        send(
            proc,
            {
                "cmd": "playlist",
                "playlist_id": PLAYLIST_ID,
            },
        )

        playlist_resp = recv(proc)

        if not playlist_resp.get("ok"):
            print(
                f"Playlist request {i + 1} failed"
            )

            print(
                playlist_resp.get("error")
            )

            break

        playlist = playlist_resp["data"]

        entries = playlist.get(
            "entries",
            [],
        )

        print(
            f"{i + 1}/{N} playlist fetched: "
            f"{playlist.get('id')} "
            f"({len(entries)} entries)"
        )


    elapsed = time.perf_counter() - start

    stop.set()

    t.join(timeout=1)

    proc.terminate()

    try:
        proc.wait(timeout=2)

    except subprocess.TimeoutExpired:
        proc.kill()


    with open(
        "worker_memory.csv",
        "w",
        newline="",
    ) as f:

        w = csv.writer(f)

        w.writerow(
            [
                "timestamp",
                "rss_kb",
                "vmsize_kb",
            ]
        )

        w.writerows(rows)


    print(
        f"\nFinished {N} request cycles "
        f"in {elapsed:.2f}s"
    )

    print(
        "Each cycle performed:"
    )

    print(
        "  1. Metadata fetch"
    )

    print(
        "  2. Subtitle fetch"
    )

    print(
        "  3. Playlist fetch"
    )

    print(
        "Memory log written to worker_memory.csv"
    )


if __name__ == "__main__":
    main()