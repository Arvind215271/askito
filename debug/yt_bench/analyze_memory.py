#!/usr/bin/env python3
"""
Memory analysis script for yt_bench log files.
Parses memory events and tracemalloc allocations from python_worker_single.log
"""

import re
import csv
from pathlib import Path

PATH = "debug/yt_bench/debug/yt_bench/output/logs_memory_stats/"
LOG_PATH = Path("debug/yt_bench/debug/yt_bench/output/logs_memory_stats/python_worker_single.log")


def load_log_lines():
    """Load all lines from the log file."""
    with open(LOG_PATH, "r") as f:
        lines = f.readlines()
    print(f"Loaded {len(lines)} log lines")
    return lines


def parse_memory_events(lines):
    """
    Parse memory events from log lines.
    Returns list of dicts: {request, stage, rss_mb}
    """
    events = []
    worker_request_counts = {}  # Track request count per worker
    
    # Pattern for REQUEST_COUNT lines: worker=N REQUEST_COUNT M
    req_count_pattern = re.compile(r'worker=(\d+)\s+REQUEST_COUNT\s+(\d+)')
    # Pattern for MEMORY lines: worker=N MEMORY stage=X rss_mb=Y
    memory_pattern = re.compile(r'worker=(\d+)\s+MEMORY\s+stage=(\S+)\s+rss_mb=([\d.]+)')
    # Pattern for REQUEST_MEMORY_START: worker=N REQUEST_MEMORY_START rss_mb=X
    req_mem_start_pattern = re.compile(r'worker=(\d+)\s+REQUEST_MEMORY_START\s+rss_mb=([\d.]+)')
    # Pattern for REQUEST_MEMORY_AFTER_GC: worker=N REQUEST_MEMORY_AFTER_GC rss_mb=X
    req_mem_gc_pattern = re.compile(r'worker=(\d+)\s+REQUEST_MEMORY_AFTER_GC\s+rss_mb=([\d.]+)')
    
    for line in lines:
        # Track request count per worker
        req_match = req_count_pattern.search(line)
        if req_match:
            worker_id = int(req_match.group(1))
            req_num = int(req_match.group(2))
            worker_request_counts[worker_id] = req_num
            continue
        
        # Parse MEMORY stage events
        mem_match = memory_pattern.search(line)
        if mem_match:
            worker_id = int(mem_match.group(1))
            stage = mem_match.group(2)
            rss_mb = float(mem_match.group(3))
            request_num = worker_request_counts.get(worker_id, 0)
            events.append({
                "request": request_num,
                "stage": stage,
                "rss_mb": rss_mb,
            })
            continue
        
        # Parse REQUEST_MEMORY_START as before_extract stage
        req_start_match = req_mem_start_pattern.search(line)
        if req_start_match:
            worker_id = int(req_start_match.group(1))
            rss_mb = float(req_start_match.group(2))
            request_num = worker_request_counts.get(worker_id, 0)
            events.append({
                "request": request_num,
                "stage": "before_extract",
                "rss_mb": rss_mb,
            })
            continue
        
        # Parse REQUEST_MEMORY_AFTER_GC as after_gc stage
        req_gc_match = req_mem_gc_pattern.search(line)
        if req_gc_match:
            worker_id = int(req_gc_match.group(1))
            rss_mb = float(req_gc_match.group(2))
            request_num = worker_request_counts.get(worker_id, 0)
            events.append({
                "request": request_num,
                "stage": "after_gc",
                "rss_mb": rss_mb,
            })
            continue
    
    return events


def parse_allocations(lines):
    """
    Parse tracemalloc allocation entries from log lines.
    Returns list of dicts: {request, file, line, size_kb, count}
    """
    allocations = []
    current_request = 0
    current_worker = None
    worker_request_counts = {}
    
    # Pattern for REQUEST_COUNT
    req_count_pattern = re.compile(r'worker=(\d+)\s+REQUEST_COUNT\s+(\d+)')
    # Pattern for TOP 20 ALLOCATIONS marker
    alloc_marker_pattern = re.compile(r'worker=(\d+)\s+--- TOP 20 ALLOCATIONS ---')
    # Pattern for tracemalloc stat lines: file.py:123: size=184 KiB, count=57
    alloc_line_pattern = re.compile(r'^(.+?):(\d+):\s+size=([\d.]+)\s*(\w+)B,\s*count=(\d+)')
    
    in_allocations = False
    
    for line in lines:
        # Track request count per worker
        req_match = req_count_pattern.search(line)
        if req_match:
            worker_id = int(req_match.group(1))
            req_num = int(req_match.group(2))
            worker_request_counts[worker_id] = req_num
            continue
        
        # Detect allocation marker
        alloc_marker = alloc_marker_pattern.search(line)
        if alloc_marker:
            current_worker = int(alloc_marker.group(1))
            current_request = worker_request_counts.get(current_worker, 0)
            in_allocations = True
            continue
        
        # Parse allocation lines
        if in_allocations:
            alloc_match = alloc_line_pattern.search(line)
            if alloc_match:
                file_path = alloc_match.group(1)
                line_num = int(alloc_match.group(2))
                size = float(alloc_match.group(3))
                unit = alloc_match.group(4)
                count = int(alloc_match.group(5))
                
                # Convert to KB
                if unit == "Ki":
                    size_kb = size
                elif unit == "Mi":
                    size_kb = size * 1024
                elif unit == "Gi":
                    size_kb = size * 1024 * 1024
                else:
                    size_kb = size / 1024  # bytes to KB
                
                allocations.append({
                    "request": current_request,
                    "file": file_path,
                    "line": line_num,
                    "size_kb": round(size_kb, 2),
                    "count": count,
                })
            else:
                # End of allocation block (empty line or non-matching line)
                if line.strip() == "" or not line.strip().startswith(" "):
                    in_allocations = False
    
    return allocations


def export_memory_csv(events):
    """Export memory events to CSV."""
    output_path = Path(PATH, "memory.csv")
    with open(output_path, "w", newline="") as f:
        writer = csv.DictWriter(f, fieldnames=["request", "stage", "rss_mb"])
        writer.writeheader()
        writer.writerows(events)
    print(f"Exported {len(events)} memory events to {output_path}")


def export_allocations_csv(allocations):
    """Export allocation events to CSV."""
    output_path = Path(PATH, "allocations.csv")
    with open(output_path, "w", newline="") as f:
        writer = csv.DictWriter(f, fieldnames=["request", "file", "line", "size_kb", "count"])
        writer.writeheader()
        writer.writerows(allocations)
    print(f"Exported {len(allocations)} allocation records to {output_path}")


def main():
    print("Phase 1: Loading log file...")
    lines = load_log_lines()
    
    print("\nPhase 2: Parsing memory events...")
    events = parse_memory_events(lines)
    print(f"Parsed {len(events)} memory events")
    
    print("\nPhase 3: Parsing tracemalloc allocations...")
    allocations = parse_allocations(lines)
    print(f"Parsed {len(allocations)} allocation records")
    
    print("\nPhase 4: Exporting CSV files...")
    export_memory_csv(events)
    export_allocations_csv(allocations)
    
    print("\nDone!")


if __name__ == "__main__":
    main()