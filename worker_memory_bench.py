#!/usr/bin/env python3
import csv
import json
import struct
import subprocess
import threading
import time


VIDEO_ID = "dQw4w9WgXcQ"
N = 50

WORKER = "internal/youtube/metadata/ytdlp/python/python_worker_single.py"
import os

CACHE_DIR = os.path.join(os.getcwd(), ".cache", "ytdlp")
os.makedirs(CACHE_DIR, exist_ok=True)

def send(proc, obj):
    payload = json.dumps(obj).encode()
    proc.stdin.write(struct.pack(">I", len(payload)))
    proc.stdin.write(payload)
    proc.stdin.flush()


def recv(proc):
    hdr = proc.stdout.read(4)
    if len(hdr) != 4:
        raise EOFError("worker closed stdout")
    length = struct.unpack(">I", hdr)[0]
    data = proc.stdout.read(length)
    return json.loads(data)


def sampler(pid, stop_event, rows):
    status = f"/proc/{pid}/status"
    while not stop_event.is_set():
        rss = vsz = None
        try:
            with open(status) as f:
                for line in f:
                    if line.startswith("VmRSS:"):
                        rss = int(line.split()[1])
                    elif line.startswith("VmSize:"):
                        vsz = int(line.split()[1])
        except FileNotFoundError:
            break

        rows.append((time.time(), rss, vsz))
        time.sleep(0.02)


def main():
    proc = subprocess.Popen(
        ["python3", WORKER, "--worker-id", "bench"],
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
    )

    print("Worker PID:", proc.pid)

    rows = []
    stop = threading.Event()
    t = threading.Thread(
        target=sampler,
        args=(proc.pid, stop, rows),
        daemon=True,
    )
    t.start()

    start = time.perf_counter()

    for i in range(N):

        # -------------------------
        # 1. Fetch metadata
        # -------------------------

        send(proc, {
            "cmd": "metadata",
            "video_id": VIDEO_ID,
        })

        metadata_resp = recv(proc)

        if not metadata_resp.get("ok"):
            print(f"Metadata request {i + 1} failed")
            print(metadata_resp.get("error"))
            break

        metadata = metadata_resp["data"]

        print(
            f"{i + 1}/{N} metadata fetched: "
            f"{metadata.get('id')}"
        )


        # -------------------------
        # 2. Fetch subtitles
        # -------------------------

        send(proc, {
            "cmd": "subtitle",
            "video_id": VIDEO_ID,
            "language": "en",
            "type": "manual",
            "format": "json3",
            "cache_dir": CACHE_DIR,
        })

        subtitle_resp = recv(proc)

        if not subtitle_resp.get("ok"):
            print(f"Subtitle request {i + 1} failed")
            print(subtitle_resp.get("error"))
            break

        filename = subtitle_resp["filename"]

        path = os.path.join(
            CACHE_DIR,
            VIDEO_ID,
            filename,
        )

        with open(path, "r", encoding="utf-8") as f:
            preview = "".join(f.readlines()[:500])

        print(f"{i + 1}/{N} subtitle fetched")
        print(preview)

    elapsed = time.perf_counter() - start

    stop.set()
    t.join(timeout=1)

    proc.terminate()
    try:
        proc.wait(timeout=2)
    except subprocess.TimeoutExpired:
        proc.kill()

    with open("worker_memory.csv", "w", newline="") as f:
        w = csv.writer(f)
        w.writerow(["timestamp", "rss_kb", "vmsize_kb"])
        w.writerows(rows)

    print(f"\nFinished {N} requests in {elapsed:.2f}s")
    print("Memory log written to worker_memory.csv")


if __name__ == "__main__":
    main()