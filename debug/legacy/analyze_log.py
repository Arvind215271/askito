#!/usr/bin/env python3

import json
import re
import statistics
from datetime import datetime


LOG_FILE = "debug/yt_bench/output/logs/python_worker_single.log"
OUTPUT_FILE = "debug/yt_bench/output/logs/python_worker_stats.json"


class BenchmarkRun:
    def __init__(self, start):
        self.start = start
        self.end = None

        self.workers = set()

        self.init_times = []
        self.warmup_times = []

        self.fetch_times = []

        self.fetch_started_at = None
        self.fetch_finished_at = None

        self.request_start = {}  # video -> timestamp
        self.videos = set()


def parse_time(line):
    return datetime.strptime(
        line[:23],
        "%Y-%m-%d %H:%M:%S,%f",
    )


def percentile(values, p):
    if not values:
        return None

    values = sorted(values)

    k = (len(values) - 1) * p / 100
    f = int(k)
    c = min(f + 1, len(values) - 1)

    if f == c:
        return values[f]

    return values[f] + (values[c] - values[f]) * (k - f)


def stats(values):
    if not values:
        return {
            "count": 0,
            "avg": None,
            "min": None,
            "max": None,
            "p50": None,
            "p95": None,
        }

    return {
        "count": len(values),
        "avg": statistics.mean(values),
        "min": min(values),
        "max": max(values),
        "p50": percentile(values, 50),
        "p95": percentile(values, 95),
    }


runs = []
current = None

worker_re = re.compile(r"worker=(\d+)")
duration_re = re.compile(r"duration=([0-9.]+)")
video_re = re.compile(r"video=([A-Za-z0-9_-]+)")

with open(LOG_FILE) as f:

    for line in f:

        timestamp = parse_time(line)

        worker_match = worker_re.search(line)
        worker = worker_match.group(1) if worker_match else None

        #
        # Start of benchmark
        #
        if "WORKER_STARTED" in line:

            if current is None:
                current = BenchmarkRun(timestamp)

            current.workers.add(worker)
            continue

        if current is None:
            continue

        #
        # Initialization
        #
        if "INITIALIZED" in line:

            m = duration_re.search(line)
            if m:
                current.init_times.append(float(m.group(1)))

            continue

        #
        # Warmup
        #
        if "WARMUP_DONE" in line:

            m = duration_re.search(line)
            if m:
                current.warmup_times.append(float(m.group(1)))

            continue

        #
        # Request started
        #
        if "REQUEST_START" in line:

            m = video_re.search(line)
            if not m:
                continue

            video = m.group(1)

            current.request_start[video] = timestamp
            current.videos.add(video)

            if current.fetch_started_at is None:
                current.fetch_started_at = timestamp

            continue

        #
        # Request finished
        #
        if "REQUEST_END" in line:

            m = video_re.search(line)
            if not m:
                continue

            video = m.group(1)

            current.fetch_finished_at = timestamp

            if video in current.request_start:

                latency = (
                    timestamp - current.request_start.pop(video)
                ).total_seconds()

                current.fetch_times.append(latency)

            continue

        #
        # Ignore failures for latency statistics
        #
        if "REQUEST_FAILED" in line:

            m = video_re.search(line)

            if m:
                current.request_start.pop(m.group(1), None)

            continue

        #
        # End benchmark once every worker exited
        #
        if "WORKER_STOPPED" in line:

            current.workers.discard(worker)

            if not current.workers:

                current.end = timestamp
                runs.append(current)
                current = None


if current:
    runs.append(current)


output = []

for i, run in enumerate(runs, 1):

    runtime = (
        (run.end - run.start).total_seconds()
        if run.end
        else None
    )

    fetch_duration = None

    if run.fetch_started_at and run.fetch_finished_at:

        fetch_duration = (
            run.fetch_finished_at -
            run.fetch_started_at
        ).total_seconds()

    videos = len(run.videos)

    output.append({

        "run": i,

        "workers": len(run.init_times),

        "runtime_seconds": runtime,

        "videos_processed": videos,

        "overall_throughput":
            videos / runtime
            if runtime
            else None,

        "fetch_cycle": {

            "start":
                run.fetch_started_at.isoformat()
                if run.fetch_started_at
                else None,

            "end":
                run.fetch_finished_at.isoformat()
                if run.fetch_finished_at
                else None,

            "duration_seconds":
                fetch_duration,

            "videos_per_second":
                videos / fetch_duration
                if fetch_duration
                else None,
        },

        "fetch_latency":
            stats(run.fetch_times),

        "initialization":
            stats(run.init_times),

        "warmup":
            stats(run.warmup_times),
    })


with open(OUTPUT_FILE, "w") as f:
    json.dump(output, f, indent=4)

print(f"Saved stats to {OUTPUT_FILE}")