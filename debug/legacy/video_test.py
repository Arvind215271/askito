import requests
import time
import json
from collections import Counter


BASE_URL = "http://127.0.0.1:8080/videos/url"


VIDEOS = [
    "nH9R0Jpqeqc",
    "dQw4w9WgXcQ",
    # Add more video IDs here
]


results = []

all_errors = Counter()
error_by_video = {}

successful_videos = 0
failed_videos = 0


for index, video_id in enumerate(VIDEOS, start=1):

    url = f"https://www.youtube.com/watch?v={video_id}"

    print("\n" + "=" * 80)
    print(f"{index}/{len(VIDEOS)}")
    print("Video ID:", video_id)
    print("URL:", url)
    print("=" * 80)

    start = time.perf_counter()

    try:

        response = requests.get(
            BASE_URL,
            json={
                "url": url,
            },
            timeout=300,
        )

        elapsed = time.perf_counter() - start

        try:
            body = response.json()

        except Exception:

            body = response.text

        result = {
            "video_id": video_id,
            "url": url,
            "status_code": response.status_code,
            "time_seconds": round(elapsed, 2),
            "response": body,
        }

        results.append(result)

        print("Status:", response.status_code)
        print("Time:", round(elapsed, 2), "seconds")

        video_errors = []

        # ---------------------------------------------------------
        # Extract errors from the API response
        # ---------------------------------------------------------

        if isinstance(body, dict):

            response_errors = body.get("errors", [])

            if isinstance(response_errors, list):

                video_errors.extend(response_errors)

            elif response_errors:

                video_errors.append(str(response_errors))

        # ---------------------------------------------------------
        # HTTP errors are also recorded
        # ---------------------------------------------------------

        if response.status_code != 200:

            http_error = (
                f"HTTP request failed with status "
                f"{response.status_code}"
            )

            video_errors.append(http_error)

        # ---------------------------------------------------------
        # Record and print errors
        # ---------------------------------------------------------

        if video_errors:

            failed_videos += 1

            error_by_video[video_id] = video_errors

            print("\nERRORS:")

            for error in video_errors:

                print("  -", error)

                all_errors[str(error)] += 1

        else:

            successful_videos += 1

            print("Errors: none")

        # ---------------------------------------------------------
        # Print response structure
        # ---------------------------------------------------------

        if isinstance(body, dict):

            print("\nResponse keys:")

            for key in body.keys():

                print("  -", key)

        print()

    except Exception as error:

        elapsed = time.perf_counter() - start

        error_string = (
            f"Request exception: {type(error).__name__}: {error}"
        )

        result = {
            "video_id": video_id,
            "url": url,
            "status_code": None,
            "time_seconds": round(elapsed, 2),
            "request_error": str(error),
        }

        results.append(result)

        failed_videos += 1

        error_by_video[video_id] = [error_string]

        all_errors[error_string] += 1

        print("REQUEST ERROR:", error_string)

    # Small delay between requests.
    time.sleep(2)


# ---------------------------------------------------------------------
# Save detailed result for every video
# ---------------------------------------------------------------------

with open("video_test_results.json", "w") as f:

    json.dump(
        results,
        f,
        indent=2,
    )


# ---------------------------------------------------------------------
# Save error summary
# ---------------------------------------------------------------------

error_summary = {
    "total_videos": len(VIDEOS),
    "successful_videos": successful_videos,
    "failed_videos": failed_videos,
    "total_errors": sum(all_errors.values()),
    "error_counts": dict(all_errors),
    "errors_by_video": error_by_video,
}


with open("video_test_error_summary.json", "w") as f:

    json.dump(
        error_summary,
        f,
        indent=2,
    )


# ---------------------------------------------------------------------
# Final summary
# ---------------------------------------------------------------------

print("\n\n" + "=" * 80)
print("FINAL SUMMARY")
print("=" * 80)

print("Total videos:", len(VIDEOS))
print("Successful:", successful_videos)
print("Failed:", failed_videos)
print("Total errors:", sum(all_errors.values()))


print("\n" + "=" * 80)
print("ERROR COUNTS")
print("=" * 80)

if all_errors:

    for error, count in all_errors.most_common():

        print(f"{count}x {error}")

else:

    print("No errors found.")


print("\n" + "=" * 80)
print("ERRORS BY VIDEO")
print("=" * 80)

if error_by_video:

    for video_id, errors in error_by_video.items():

        print("\n", video_id)

        for error in errors:

            print("  -", error)

else:

    print("No failed videos.")


print("\n")
print("Saved:")
print("  video_test_results.json")
print("  video_test_error_summary.json")