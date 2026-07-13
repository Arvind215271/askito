import re
import statistics
from datetime import datetime
import json


LOG_FILE = "debug/yt_bench/debug/yt_bench/output/logs/kaggle.log"


class BenchmarkRun:
    def __init__(self, start):
        self.start = start
        self.end = None

        self.workers = 0

        self.init_times = []
        self.warmup_times = []

        self.videos = []
        self.fetch_times = []

        # Actual fetch window
        self.fetch_started_at = None
        self.fetch_finished_at = None


    def runtime(self):
        if self.end:
            return (self.end - self.start).total_seconds()

        return None


    def fetch_cycle_time(self):
        if self.fetch_started_at and self.fetch_finished_at:
            return (
                self.fetch_finished_at -
                self.fetch_started_at
            ).total_seconds()

        return None



def parse_time(line):
    return datetime.strptime(
        line[:23],
        "%Y-%m-%d %H:%M:%S,%f"
    )


def percentile(values, p):

    if not values:
        return None

    values = sorted(values)

    k = (len(values)-1) * (p/100)

    f = int(k)
    c = min(
        f+1,
        len(values)-1
    )

    if f == c:
        return values[f]

    return values[f] + (
        values[c]-values[f]
    ) * (k-f)



def stats(values):

    if not values:
        return {
            "count": 0,
            "avg": None,
            "p50": None,
            "p95": None
        }

    return {
        "count": len(values),
        "avg": statistics.mean(values),
        "p50": percentile(values, 50),
        "p95": percentile(values, 95)
    }



runs = []
current = None


with open(LOG_FILE, "r") as f:

    for line in f:

        t = parse_time(line)


        if "Python worker started" in line:

            if current is None:
                current = BenchmarkRun(t)

            current.workers += 1



        elif current is not None:


            if "Initialization time:" in line:

                value = float(
                    re.search(
                        r"Initialization time: ([0-9.]+)",
                        line
                    ).group(1)
                )

                current.init_times.append(value)



            elif "Warmup finished in" in line:

                value = float(
                    re.search(
                        r"Warmup finished in ([0-9.]+)",
                        line
                    ).group(1)
                )

                current.warmup_times.append(value)



            elif "Processing video:" in line:

                if current.fetch_started_at is None:
                    current.fetch_started_at = t


                video = line.split(
                    "Processing video:"
                )[1].strip()


                current.videos.append(video)



            elif "Finished metadata for" in line:


                current.fetch_finished_at = t


                match = re.search(
                    r"Fetch: ([0-9.]+)",
                    line
                )


                if match:

                    current.fetch_times.append(
                        float(match.group(1))
                    )



            elif "stdin closed" in line:

                current.end = t

                runs.append(current)

                current = None



if current:
    runs.append(current)



output = []


for i, run in enumerate(runs, 1):

    fetch_cycle = run.fetch_cycle_time()

    videos = len(run.videos)


    data = {

        "run": i,

        "workers": run.workers,


        "runtime_seconds":
            run.runtime(),


        "videos_processed":
            videos,


        #
        # Entire process
        #
        "overall_throughput":
            videos / run.runtime()
            if run.runtime()
            else None,


        #
        # Actual fetch window
        #
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
                fetch_cycle,


            "videos_per_second":
                videos / fetch_cycle
                if fetch_cycle
                else None
        },


        #
        # Per video fetch latency
        #
        "fetch_latency":
            stats(run.fetch_times),



        "initialization":
            stats(run.init_times),



        "warmup":
            stats(run.warmup_times)
    }


    output.append(data)



with open(
    "debug/yt_bench/debug/yt_bench/output/logs/kaggle_stats.json",
    "w"
) as f:

    json.dump(
        output,
        f,
        indent=4
    )


print(
    "Saved stats to kaggle_stats.json"
)