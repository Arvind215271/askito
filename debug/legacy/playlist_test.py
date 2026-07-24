import requests
import time
import json
from collections import Counter


BASE_URL = "http://127.0.0.1:8080/export/playlist"


PLAYLISTS = [

    # (
    #     "100 videos",
    #     "PLakjEKwjPZDYZMtzRfB26Un77KI2xI_oB",
    # ),

    # (
    #     "200 videos",
    #     "PLiy0XOfUv4hFH2HbflPOARBXA6qN90mHt",
    # ),

    # (
    #     "500 videos",
    #     "PLtJCksbabIPyfutcI3Kx9mg6h87p2DXzE",
    # ),

    # (
    #     "584 videos",
    #     "PL-CA1f0J88gDsSdoXFz3vIefI8SnBan1s",
    # ),

    # (
    #     "1000 videos",
    #     "PLoaTDHRuxwgxk0WTXRJJJlXUYuHBSrbcF",
    # ),

    # (
    #     "1000 videos",
    #     "PL5fRL6A4m-DFddDPJU5Ugr3NRhxrm5GtR",
    # ),

    (
        "2198 videos",
        "PL0dTxWJ6ngUKlBw5eDv7qYylA3xP3Asef",
    ),

    # (
    #     "5000 videos",
    #     "PLdSukIYrTISE",
    # ),
]


results = []


for index, (name, playlist_id) in enumerate(
    PLAYLISTS,
    start=1,
):

    playlist_url = (
        f"https://www.youtube.com/playlist?list={playlist_id}"
    )

    print("\n" + "=" * 100)
    print(f"PLAYLIST {index}/{len(PLAYLISTS)}")
    print("Name:", name)
    print("Playlist ID:", playlist_id)
    print("URL:", playlist_url)
    print("=" * 100)

    start = time.perf_counter()

    try:

        response = requests.post(
            BASE_URL,
            json={
                "input": playlist_url,
                "video_fields": [],
                "format": "json",

                "subtitle": {
                    "type": "automatic",
                    "language": "en",
                    "format": "json3",
                },
            },
            timeout=300,
        )

        elapsed = time.perf_counter() - start

        try:

            body = response.json()

        except Exception:

            body = response.text

        result = {
            "name": name,
            "playlist_id": playlist_id,
            "playlist_url": playlist_url,
            "status_code": response.status_code,
            "time_seconds": round(elapsed, 2),
        }

        # -------------------------------------------------------------
        # Request-level failure
        # -------------------------------------------------------------

        if response.status_code != 200:

            error = (
                f"HTTP request failed with status "
                f"{response.status_code}"
            )

            result["request_error"] = error
            result["response"] = body

            results.append(result)

            print("STATUS:", response.status_code)
            print("TIME:", round(elapsed, 2), "seconds")
            print("ERROR:", error)

            continue

        # -------------------------------------------------------------
        # Validate response
        # -------------------------------------------------------------

        if not isinstance(body, dict):

            error = "Response was not a JSON object"

            result["request_error"] = error
            result["response"] = body

            results.append(result)

            print("ERROR:", error)

            continue

        # -------------------------------------------------------------
        # Extract videos
        # -------------------------------------------------------------

        videos = body.get(
            "videos",
            [],
        )

        result["videos"] = len(videos)

        # Keep the complete API response.
        result["response"] = body

        # -------------------------------------------------------------
        # Error collection
        # -------------------------------------------------------------

        all_errors = Counter()
        errors_by_video = {}

        videos_with_errors = 0

        # -------------------------------------------------------------
        # Subtitle statistics
        # -------------------------------------------------------------

        videos_with_manual_subtitles = 0
        videos_with_automatic_subtitles = 0
        videos_with_any_subtitles = 0
        videos_with_both_subtitle_types = 0
        videos_without_subtitles = 0

        # -------------------------------------------------------------
        # Language statistics
        # -------------------------------------------------------------

        manual_language_counts = Counter()
        automatic_language_counts = Counter()

        # -------------------------------------------------------------
        # Format statistics
        # -------------------------------------------------------------

        manual_format_counts = Counter()
        automatic_format_counts = Counter()

        # -------------------------------------------------------------
        # Language + format statistics
        # -------------------------------------------------------------

        manual_language_format_counts = Counter()
        automatic_language_format_counts = Counter()

        # -------------------------------------------------------------
        # Video-level format statistics
        #
        # Count a format at most once per video.
        #
        # Example:
        #
        # Video A:
        #   en -> json3, srt
        #
        # Video B:
        #   fr -> json3
        #
        # json3 = 2 videos
        # srt   = 1 video
        #
        # -------------------------------------------------------------

        manual_format_video_counts = Counter()
        automatic_format_video_counts = Counter()

        # -------------------------------------------------------------
        # Process videos
        # -------------------------------------------------------------

        for video_index, video in enumerate(
            videos,
            start=1,
        ):

            # =========================================================
            # ERRORS
            # =========================================================

            errors = video.get(
                "errors",
                [],
            )

            if errors:

                videos_with_errors += 1

                video_id = video.get(
                    "id",
                    f"video_{video_index}",
                )

                errors_by_video[video_id] = errors

                for error in errors:

                    all_errors[str(error)] += 1

            # =========================================================
            # SUBTITLE METADATA
            # =========================================================

            # subtitle_metadata = video.get(
            #     "subtitle_metadata"
            # )

            # if not isinstance(
            #     subtitle_metadata,
            #     dict,
            # ):

            #     videos_without_subtitles += 1

            #     continue

            # manual = subtitle_metadata.get(
            #     "manual"
            # )

            # automatic = subtitle_metadata.get(
            #     "automatic"
            # )

            # has_manual = (
            #     isinstance(manual, list)
            #     and bool(manual)
            # )

            # has_automatic = (
            #     isinstance(automatic, list)
            #     and bool(automatic)
            # )

            # # ---------------------------------------------------------
            # # Overall subtitle availability
            # # ---------------------------------------------------------

            # if has_manual:

            #     videos_with_manual_subtitles += 1

            # if has_automatic:

            #     videos_with_automatic_subtitles += 1

            # if has_manual or has_automatic:

            #     videos_with_any_subtitles += 1

            # else:

            #     videos_without_subtitles += 1

            # if has_manual and has_automatic:

            #     videos_with_both_subtitle_types += 1

            # # =========================================================
            # # MANUAL SUBTITLES
            # # =========================================================

            # if has_manual:

            #     seen_manual_languages = set()
            #     seen_manual_formats = set()
            #     seen_manual_combinations = set()

            #     for track in manual:

            #         if not isinstance(
            #             track,
            #             dict,
            #         ):

            #             continue

            #         language = track.get(
            #             "languageCode"
            #         )

            #         if not language:

            #             language = track.get(
            #                 "languageName"
            #             )

            #         if not language:

            #             continue

            #         # -------------------------------------------------
            #         # Language
            #         # -------------------------------------------------

            #         if language not in seen_manual_languages:

            #             manual_language_counts[
            #                 language
            #             ] += 1

            #             seen_manual_languages.add(
            #                 language
            #             )

            #         formats = track.get(
            #             "formats",
            #             [],
            #         )

            #         if not isinstance(
            #             formats,
            #             list,
            #         ):

            #             continue

            #         # -------------------------------------------------
            #         # Formats
            #         # -------------------------------------------------

            #         for subtitle_format in formats:

            #             if not subtitle_format:

            #                 continue

            #             if subtitle_format not in seen_manual_formats:

            #                 manual_format_counts[
            #                     subtitle_format
            #                 ] += 1

            #                 seen_manual_formats.add(
            #                     subtitle_format
            #                 )

            #             combination = (
            #                 language,
            #                 subtitle_format,
            #             )

            #             if (
            #                 combination
            #                 not in seen_manual_combinations
            #             ):

            #                 manual_language_format_counts[
            #                     combination
            #                 ] += 1

            #                 seen_manual_combinations.add(
            #                     combination
            #                 )

            #     # -----------------------------------------------------
            #     # Count each format once per video
            #     # -----------------------------------------------------

            #     for subtitle_format in seen_manual_formats:

            #         manual_format_video_counts[
            #             subtitle_format
            #         ] += 1

            # # =========================================================
            # # AUTOMATIC SUBTITLES
            # # =========================================================

            # if has_automatic:

            #     seen_automatic_languages = set()
            #     seen_automatic_formats = set()
            #     seen_automatic_combinations = set()

            #     for track in automatic:

            #         if not isinstance(
            #             track,
            #             dict,
            #         ):

            #             continue

            #         language = track.get(
            #             "languageCode"
            #         )

            #         if not language:

            #             language = track.get(
            #                 "languageName"
            #             )

            #         if not language:

            #             continue

            #         # -------------------------------------------------
            #         # Language
            #         # -------------------------------------------------

            #         if language not in seen_automatic_languages:

            #             automatic_language_counts[
            #                 language
            #             ] += 1

            #             seen_automatic_languages.add(
            #                 language
            #             )

            #         formats = track.get(
            #             "formats",
            #             [],
            #         )

            #         if not isinstance(
            #             formats,
            #             list,
            #         ):

            #             continue

            #         # -------------------------------------------------
            #         # Formats
            #         # -------------------------------------------------

            #         for subtitle_format in formats:

            #             if not subtitle_format:

            #                 continue

            #             if subtitle_format not in seen_automatic_formats:

            #                 automatic_format_counts[
            #                     subtitle_format
            #                 ] += 1

            #                 seen_automatic_formats.add(
            #                     subtitle_format
            #                 )

            #             combination = (
            #                 language,
            #                 subtitle_format,
            #             )

            #             if (
            #                 combination
            #                 not in seen_automatic_combinations
            #             ):

            #                 automatic_language_format_counts[
            #                     combination
            #                 ] += 1

            #                 seen_automatic_combinations.add(
            #                     combination
            #                 )

            #     # -----------------------------------------------------
            #     # Count each format once per video
            #     # -----------------------------------------------------

            #     for subtitle_format in seen_automatic_formats:

            #         automatic_format_video_counts[
            #             subtitle_format
            #         ] += 1

        # =============================================================
        # Store results
        # =============================================================

        result["videos_with_errors"] = (
            videos_with_errors
        )

        result["total_errors"] = sum(
            all_errors.values()
        )

        result["error_counts"] = dict(
            all_errors
        )

        result["errors_by_video"] = (
            errors_by_video
        )

        result["subtitle_availability"] = {

            "total_videos": len(videos),

            "videos_with_any_subtitles": (
                videos_with_any_subtitles
            ),

            "videos_without_subtitles": (
                videos_without_subtitles
            ),

            "videos_with_both_manual_and_automatic": (
                videos_with_both_subtitle_types
            ),

            "manual": {

                "videos": (
                    videos_with_manual_subtitles
                ),

                "languages": dict(
                    manual_language_counts
                    .most_common()
                ),

                "formats": dict(
                    manual_format_counts
                    .most_common()
                ),

                "format_video_counts": dict(
                    manual_format_video_counts
                    .most_common()
                ),

                "language_format_combinations": [

                    {
                        "language": language,
                        "format": subtitle_format,
                        "videos": count,
                    }

                    for (
                        language,
                        subtitle_format
                    ), count in (
                        manual_language_format_counts
                        .most_common()
                    )

                ],

            },

            "automatic": {

                "videos": (
                    videos_with_automatic_subtitles
                ),

                "languages": dict(
                    automatic_language_counts
                    .most_common()
                ),

                "formats": dict(
                    automatic_format_counts
                    .most_common()
                ),

                "format_video_counts": dict(
                    automatic_format_video_counts
                    .most_common()
                ),

                "language_format_combinations": [

                    {
                        "language": language,
                        "format": subtitle_format,
                        "videos": count,
                    }

                    for (
                        language,
                        subtitle_format
                    ), count in (
                        automatic_language_format_counts
                        .most_common()
                    )

                ],

            },

        }

        results.append(result)

        # =============================================================
        # Print result
        # =============================================================

        print("STATUS:", response.status_code)
        print(
            "TIME:",
            round(elapsed, 2),
            "seconds",
        )

        print(
            "VIDEOS:",
            len(videos),
        )

        print(
            "VIDEOS WITH ERRORS:",
            videos_with_errors,
        )

        print(
            "TOTAL ERRORS:",
            sum(all_errors.values()),
        )

        if all_errors:

            print("\nTOP ERRORS:")

            for error, count in (
                all_errors.most_common(20)
            ):

                print(
                    f"  {count}x {error}"
                )

        # =============================================================
        # SUBTITLE COVERAGE
        # =============================================================

        print("\n")
        print("SUBTITLE COVERAGE")
        print("-" * 100)

        print(
            "Total videos:",
            len(videos),
        )

        print(
            "Videos with any subtitles:",
            videos_with_any_subtitles,
        )

        print(
            "Videos without subtitles:",
            videos_without_subtitles,
        )

        print(
            "Videos with manual subtitles:",
            videos_with_manual_subtitles,
        )

        print(
            "Videos with automatic subtitles:",
            videos_with_automatic_subtitles,
        )

        print(
            "Videos with both manual + automatic:",
            videos_with_both_subtitle_types,
        )

        # =============================================================
        # MANUAL SUBTITLES
        # =============================================================

        print("\n")
        print("MANUAL SUBTITLE AVAILABILITY")
        print("-" * 100)

        print(
            "Videos with manual subtitles:",
            videos_with_manual_subtitles,
        )

        print("\nTop 20 manual languages:")

        if manual_language_counts:

            for language, count in (
                manual_language_counts
                .most_common(20)
            ):

                print(
                    f"  {count:>4} videos | "
                    f"{language}"
                )

        else:

            print(
                "  No manual subtitle languages found."
            )

        print(
            "\nTop 20 manual language + format combinations:"
        )

        if manual_language_format_counts:

            for (
                language,
                subtitle_format,
            ), count in (
                manual_language_format_counts
                .most_common(20)
            ):

                print(
                    f"  {count:>4} videos | "
                    f"{language:<30} | "
                    f"{subtitle_format}"
                )

        else:

            print(
                "  No manual subtitle metadata found."
            )

        print(
            "\nManual formats by video count:"
        )

        if manual_format_video_counts:

            for subtitle_format, count in (
                manual_format_video_counts
                .most_common()
            ):

                print(
                    f"  {count:>4} videos | "
                    f"{subtitle_format}"
                )

        else:

            print(
                "  No manual subtitle formats found."
            )

        # =============================================================
        # AUTOMATIC SUBTITLES
        # =============================================================

        print("\n")
        print("AUTOMATIC SUBTITLE AVAILABILITY")
        print("-" * 100)

        print(
            "Videos with automatic subtitles:",
            videos_with_automatic_subtitles,
        )

        print("\nTop 20 automatic languages:")

        if automatic_language_counts:

            for language, count in (
                automatic_language_counts
                .most_common(20)
            ):

                print(
                    f"  {count:>4} videos | "
                    f"{language}"
                )

        else:

            print(
                "  No automatic subtitle languages found."
            )

        print(
            "\nTop 20 automatic language + format combinations:"
        )

        if automatic_language_format_counts:

            for (
                language,
                subtitle_format,
            ), count in (
                automatic_language_format_counts
                .most_common(20)
            ):

                print(
                    f"  {count:>4} videos | "
                    f"{language:<30} | "
                    f"{subtitle_format}"
                )

        else:

            print(
                "  No automatic subtitle metadata found."
            )

        print(
            "\nAutomatic formats by video count:"
        )

        if automatic_format_video_counts:

            for subtitle_format, count in (
                automatic_format_video_counts
                .most_common()
            ):

                print(
                    f"  {count:>4} videos | "
                    f"{subtitle_format}"
                )

        else:

            print(
                "  No automatic subtitle formats found."
            )

        print()

    except Exception as error:

        elapsed = time.perf_counter() - start

        result = {
            "name": name,
            "playlist_id": playlist_id,
            "playlist_url": playlist_url,
            "status_code": None,
            "time_seconds": round(elapsed, 2),
            "request_error": (
                f"{type(error).__name__}: {error}"
            ),
        }

        results.append(result)

        print(
            "REQUEST ERROR:",
            result["request_error"],
        )


# ---------------------------------------------------------------------
# Save detailed results
# ---------------------------------------------------------------------

with open(
    "playlist_test_results.json",
    "w",
) as f:

    json.dump(
        results,
        f,
        indent=2,
    )


# ---------------------------------------------------------------------
# Final summary
# ---------------------------------------------------------------------

print("\n\n" + "=" * 100)
print("FINAL SUMMARY")
print("=" * 100)


for result in results:

    subtitle_availability = result.get(
        "subtitle_availability",
        {},
    )

    print(
        f"{result['name']:>15} | "
        f"{result.get('videos', 0):>6} videos | "
        f"{result['time_seconds']:>10.2f}s | "
        f"{result.get('videos_with_errors', 0):>6} failed | "
        f"{result.get('total_errors', 0):>6} errors | "
        f"{subtitle_availability.get('videos_with_any_subtitles', 0):>6} "
        f"with subtitles | "
        f"{subtitle_availability.get('videos_without_subtitles', 0):>6} "
        f"without subtitles"
    )


print()
print("Saved:")
print("  playlist_test_results.json")