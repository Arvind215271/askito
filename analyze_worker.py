import os
import csv
import statistics


NUM_WORKERS = 17
FILE_PATTERN = "debug/yt_bench/output/lite_worker/worker_memory_{}.csv"


def read_memory_csv(filename):
    """
    Reads CSV until Python code appears at the end of the file.
    Expected CSV format:
    timestamp,rss_kb,vmsize_kb
    """

    rows = []

    with open(filename, "r", encoding="utf-8") as f:
        reader = csv.reader(f)

        header_found = False

        for row in reader:

            # Stop when Python code starts appearing
            if not row:
                continue

            first = row[0].strip()

            # Detect python code
            python_markers = [
                "import ",
                "from ",
                "def ",
                "class ",
                "if __name__",
                "#",
                "print(",
            ]

            if any(first.startswith(x) for x in python_markers):
                break

            # Skip header
            if first == "timestamp":
                header_found = True
                continue

            if not header_found:
                continue

            # Only accept CSV rows
            if len(row) != 3:
                break

            try:
                timestamp = float(row[0])
                rss = int(row[1])
                vmsize = int(row[2])

                rows.append(
                    {
                        "timestamp": timestamp,
                        "rss": rss,
                        "vmsize": vmsize,
                    }
                )

            except ValueError:
                break

    return rows


def analyze_worker(rows):

    rss_values = [r["rss"] for r in rows]
    vmsize_values = [r["vmsize"] for r in rows]

    if not rss_values:
        return None

    start_rss = rss_values[0]
    end_rss = rss_values[-1]

    growth = end_rss - start_rss

    return {
        "samples": len(rows),
        "avg_rss_mb": statistics.mean(rss_values) / 1024,
        "peak_rss_mb": max(rss_values) / 1024,
        "final_rss_mb": end_rss / 1024,
        "rss_growth_mb": growth / 1024,
        "rss_stddev_mb": statistics.stdev(rss_values) / 1024
        if len(rss_values) > 1
        else 0,
        "peak_vmsize_mb": max(vmsize_values) / 1024,
    }


def calculate_score(stats):
    """
    Higher score = better

    Rewards:
    - low memory
    - low growth
    - stability
    """

    score = 100

    score -= stats["avg_rss_mb"] * 0.5
    score -= stats["peak_rss_mb"] * 0.2
    score -= abs(stats["rss_growth_mb"]) * 2
    score -= stats["rss_stddev_mb"] * 1

    return round(score, 2)


def main():

    results = {}

    for i in range(1, NUM_WORKERS + 1):

        filename = FILE_PATTERN.format(i)

        if not os.path.exists(filename):
            print(f"Missing: {filename}")
            continue

        rows = read_memory_csv(filename)

        stats = analyze_worker(rows)

        if stats:
            stats["score"] = calculate_score(stats)
            results[f"Worker {i}"] = stats


    print("\n========== Worker Memory Summary ==========\n")

    ranking = sorted(
        results.items(),
        key=lambda x: x[1]["score"],
        reverse=True
    )

    for rank, (worker, stats) in enumerate(ranking, 1):

        print(f"{rank}. {worker}")
        print(f"   Score:          {stats['score']}")
        print(f"   Samples:        {stats['samples']}")
        print(f"   Avg RSS:        {stats['avg_rss_mb']:.2f} MB")
        print(f"   Peak RSS:       {stats['peak_rss_mb']:.2f} MB")
        print(f"   Final RSS:      {stats['final_rss_mb']:.2f} MB")
        print(f"   RSS Growth:     {stats['rss_growth_mb']:.2f} MB")
        print(f"   RSS Stability:  {stats['rss_stddev_mb']:.2f} MB")
        print()


    if ranking:
        best = ranking[0]

        print("========== Winner ==========")
        print(best[0])
        print()
        print(
            "Selected because it has the best balance of "
            "low memory usage, low growth, and stability."
        )


if __name__ == "__main__":
    main()