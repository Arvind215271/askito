import json
import re
import statistics
from datetime import datetime


LOG_FILE = "debug/yt_bench/debug/yt_bench/output/logs/kaggle_worker.log"

OUTPUT_FILE = (
    "debug/yt_bench/debug/yt_bench/output/logs/"
    "kaggle_stats_worker.json"
)


class BenchmarkRun:

    def __init__(self, start):

        self.start = start
        self.end = None

        self.workers = set()

        self.init_times = []
        self.warmup_times = []

        self.fetch_start = None
        self.fetch_end = None

        self.fetch_times = []
        self.videos = set()


    def runtime(self):

        if self.end is None:
            return None

        return (
            self.end - self.start
        ).total_seconds()


    def fetch_runtime(self):

        if (
            self.fetch_start is None
            or
            self.fetch_end is None
        ):
            return None

        return (
            self.fetch_end -
            self.fetch_start
        ).total_seconds()



def parse_time(line):

    return datetime.strptime(
        line[:23],
        "%Y-%m-%d %H:%M:%S,%f",
    )



def percentile(values, p):

    if not values:
        return None

    values = sorted(values)

    k = (
        len(values) - 1
    ) * (
        p / 100
    )

    f = int(k)

    c = min(
        f + 1,
        len(values) - 1,
    )

    if f == c:
        return values[f]

    return (
        values[f]
        +
        (
            values[c] -
            values[f]
        )
        *
        (
            k - f
        )
    )



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


worker_started = re.compile(
    r"worker=(\d+).*WORKER_STARTED"
)

worker_initialized = re.compile(
    r"INITIALIZED duration=([0-9.]+)"
)

worker_warmup = re.compile(
    r"WARMUP_DONE duration=([0-9.]+)"
)

request_start = re.compile(
    r"REQUEST_START video=([^\s]+)"
)

request_end = re.compile(
    r"REQUEST_END video=([^\s]+)"
)

worker_stopped = re.compile(
    r"worker=(\d+).*WORKER_STOPPED"
)


with open(LOG_FILE, "r") as f:

    for line in f:

        timestamp = parse_time(line)

        #
        # New benchmark starts
        #

        match = worker_started.search(line)

        if match:

            if current is None:
                current = BenchmarkRun(timestamp)

            current.workers.add(
                int(match.group(1))
            )

            continue

        if current is None:
            continue

        #
        # Worker initialization
        #

        match = worker_initialized.search(line)

        if match:

            current.init_times.append(
                float(match.group(1))
            )

            continue

        #
        # Warmup
        #

        match = worker_warmup.search(line)

        if match:

            current.warmup_times.append(
                float(match.group(1))
            )

            continue

        #
        # Request started
        #

        match = request_start.search(line)

        if match:

            if current.fetch_start is None:
                current.fetch_start = timestamp

            current.videos.add(
                match.group(1)
            )

            continue

        #
        # Request finished
        #

        match = request_end.search(line)

        if match:

            current.fetch_end = timestamp

            continue

        #
        # Worker stopped
        #

        match = worker_stopped.search(line)

        if match:

            current.workers.discard(
                int(match.group(1))
            )

            if not current.workers:

                current.end = timestamp

                runs.append(current)

                current = None


#
# Build fetch durations by matching start/end timestamps
#

request_start_times = {}

with open(LOG_FILE, "r") as f:

    for line in f:

        timestamp = parse_time(line)

        match = request_start.search(line)

        if match:

            request_start_times[
                match.group(1)
            ] = timestamp

            continue

        match = request_end.search(line)

        if match:

            video = match.group(1)

            start = request_start_times.pop(
                video,
                None,
            )

            if (
                start is not None
                and
                runs
            ):

                duration = (
                    timestamp -
                    start
                ).total_seconds()

                for run in runs:

                    if (
                        run.fetch_start
                        and
                        run.fetch_end
                        and
                        run.fetch_start
                        <= start
                        <= run.fetch_end
                    ):

                        run.fetch_times.append(
                            duration
                        )

                        break


output = []


for index, run in enumerate(runs, 1):

    fetch_runtime = run.fetch_runtime()

    video_count = len(run.videos)

    avg_fetch = (
        statistics.mean(run.fetch_times)
        if run.fetch_times
        else None
    )

    output.append({

        "run":
            index,

        "workers":
            len(run.init_times),

        "start":
            run.start.isoformat(),

        "end":
            run.end.isoformat()
            if run.end
            else None,

        "total_runtime_seconds":
            run.runtime(),

        "fetch_cycle": {

            "start":
                run.fetch_start.isoformat()
                if run.fetch_start
                else None,

            "end":
                run.fetch_end.isoformat()
                if run.fetch_end
                else None,

            "duration_seconds":
                fetch_runtime,

            "videos_processed":
                video_count,

            "throughput_videos_per_second":
                (
                    video_count /
                    fetch_runtime
                )
                if fetch_runtime
                else None,
        },

        "parallelism": {

            "ideal_worker_seconds":
                (
                    video_count *
                    avg_fetch
                )
                if avg_fetch
                else None,

            "effective_parallelism":
                (
                    (
                        video_count *
                        avg_fetch
                    )
                    /
                    fetch_runtime
                )
                if (
                    avg_fetch
                    and
                    fetch_runtime
                )
                else None,
        },

        "initialization":
            stats(
                run.init_times
            ),

        "warmup":
            stats(
                run.warmup_times
            ),

        "metadata_fetch":
            stats(
                run.fetch_times
            ),
    })


with open(
    OUTPUT_FILE,
    "w",
) as f:

    json.dump(
        output,
        f,
        indent=4,
    )


print(
    f"Saved stats to {OUTPUT_FILE}"
)