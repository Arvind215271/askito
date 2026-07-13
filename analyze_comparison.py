import re
import statistics
from datetime import datetime
import json


LOG_FILE = "debug/yt_bench/debug/yt_bench/output/logs/kaggle.log"

OUTPUT_FILE = (
    "debug/yt_bench/debug/yt_bench/output/logs/"
    "kaggle_stats.json"
)


class BenchmarkRun:

    def __init__(self, start):

        self.start = start
        self.end = None

        self.fetch_start = None
        self.fetch_end = None

        self.workers = 0

        self.init_times = []
        self.warmup_times = []

        self.videos = []
        self.fetch_times = []


    def runtime(self):

        if self.end:
            return (
                self.end - self.start
            ).total_seconds()

        return None


    def fetch_runtime(self):

        if self.fetch_start and self.fetch_end:

            return (
                self.fetch_end -
                self.fetch_start
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

    k = (
        len(values) - 1
    ) * (
        p / 100
    )

    f = int(k)

    c = min(
        f + 1,
        len(values)-1
    )

    if f == c:
        return values[f]

    return (
        values[f]
        +
        (
            values[c]
            -
            values[f]
        )
        *
        (
            k-f
        )
    )



runs = []

current = None



with open(LOG_FILE, "r") as f:

    for line in f:


        t = parse_time(line)



        #
        # Worker start
        #

        if "Python worker started" in line:


            if current is None:

                current = BenchmarkRun(t)


            current.workers += 1



        elif current is not None:



            #
            # Initialization
            #

            if "Initialization time:" in line:


                value = float(
                    re.search(
                        r"Initialization time: ([0-9.]+)",
                        line
                    ).group(1)
                )


                current.init_times.append(value)



            #
            # Warmup
            #

            elif "Warmup finished in" in line:


                value = float(
                    re.search(
                        r"Warmup finished in ([0-9.]+)",
                        line
                    ).group(1)
                )


                current.warmup_times.append(value)



            #
            # Video started
            #

            elif "Processing video:" in line:


                video = (
                    line
                    .split("Processing video:")[1]
                    .strip()
                )


                current.videos.append(video)



            #
            # Metadata finished
            #

            elif "Finished metadata for" in line:



                if current.fetch_start is None:

                    current.fetch_start = t



                current.fetch_end = t



                match = re.search(
                    r"Fetch: ([0-9.]+)",
                    line
                )


                if match:

                    current.fetch_times.append(
                        float(
                            match.group(1)
                        )
                    )



            #
            # Worker shutdown
            #

            elif "stdin closed" in line:


                current.end = t


                runs.append(current)


                current = None





if current:

    runs.append(current)





def stats(values):

    return {

        "count":
            len(values),


        "avg":
            statistics.mean(values)
            if values else None,


        "p50":
            percentile(values,50),


        "p95":
            percentile(values,95)

    }





output = []



for i, run in enumerate(runs, 1):


    fetch_cycle = run.fetch_runtime()


    avg_fetch = (

        statistics.mean(
            run.fetch_times
        )

        if run.fetch_times

        else None

    )



    data = {


        "run":

            i,



        "workers":

            run.workers,



        "start":

            run.start.isoformat(),



        "end":

            run.end.isoformat()
            if run.end
            else None,



        #
        # Whole benchmark
        #

        "total_runtime_seconds":

            run.runtime(),



        #
        # Actual fetch phase
        #

        "fetch_cycle_seconds":

            fetch_cycle,



        "videos_processed":

            len(run.videos),



        "fetch_throughput_videos_per_sec":

            (
                len(run.videos)
                /
                fetch_cycle
            )

            if fetch_cycle

            else None,



        #
        # Parallel efficiency
        #

        "fetch_parallelism_efficiency":

            (
                (
                    len(run.fetch_times)
                    *
                    avg_fetch
                )
                /
                fetch_cycle
            )

            if avg_fetch and fetch_cycle

            else None,



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
            )

    }


    output.append(data)





with open(
    OUTPUT_FILE,
    "w"
) as f:


    json.dump(
        output,
        f,
        indent=4
    )



print(
    f"Saved stats to {OUTPUT_FILE}"
)