#!/usr/bin/env python3
"""
wait_for_signal.py
==================
Polls a GitHub Gist every POLL_INTERVAL seconds, waiting for Jenkins to write
a signal for the expected stage.  Exits 0 on success, 1 on error or timeout.

Environment variables (all required unless noted):
  GH_TOKEN          – GitHub PAT with gist scope
  GIST_ID           – ID of the signal Gist
  RUN_ID            – GitHub Actions run ID (used as a safety key)
  STAGE_NAME        – the stage we are waiting for (e.g. "env-setup")
  TIMEOUT_SECONDS   – hard timeout (default 600)
"""

import json
import os
import sys
import time
import urllib.error
import urllib.request
from datetime import datetime, timezone

# ── config ────────────────────────────────────────────────────────────────────
GH_TOKEN        = os.environ["GH_TOKEN"]
GIST_ID         = os.environ["GIST_ID"]
RUN_ID          = os.environ["RUN_ID"]
STAGE_NAME      = os.environ["STAGE_NAME"]
TIMEOUT_SECONDS = int(os.environ.get("TIMEOUT_SECONDS", "600"))
POLL_INTERVAL   = 8   # seconds between Gist reads


def ts() -> str:
    return datetime.now(timezone.utc).strftime("%H:%M:%S")


def fetch_signal() -> dict | None:
    """Fetch and parse the current Gist content. Returns None on transient errors."""
    url = f"https://api.github.com/gists/{GIST_ID}"
    req = urllib.request.Request(
        url,
        headers={
            "Authorization": f"token {GH_TOKEN}",
            "Accept": "application/vnd.github+json",
            "X-GitHub-Api-Version": "2022-11-28",
        },
    )
    try:
        with urllib.request.urlopen(req, timeout=15) as resp:
            gist = json.load(resp)
    except urllib.error.HTTPError as exc:
        print(f"[{ts()}] ⚠  GitHub API HTTP {exc.code} — will retry")
        return None
    except Exception as exc:
        print(f"[{ts()}] ⚠  Network error ({exc}) — will retry")
        return None

    # The file is named ci-signal-{run_id}.json
    files = gist.get("files", {})
    if not files:
        return None

    raw_content = next(iter(files.values()), {}).get("content", "")
    if not raw_content:
        return None

    try:
        return json.loads(raw_content)
    except json.JSONDecodeError:
        print(f"[{ts()}] ⚠  Gist content is not valid JSON yet — will retry")
        return None


def die(message: str, annotation: str = "error") -> None:
    print(f"::{annotation}::{message}")
    sys.exit(1)


# ── main ──────────────────────────────────────────────────────────────────────
deadline = time.monotonic() + TIMEOUT_SECONDS
print(f"[{ts()}] Waiting for Jenkins stage '{STAGE_NAME}' "
      f"(timeout {TIMEOUT_SECONDS}s, polling every {POLL_INTERVAL}s) …")

last_seen_stage = ""

while True:
    # ── hard timeout ──────────────────────────────────────────────────────────
    remaining = deadline - time.monotonic()
    if remaining <= 0:
        die(
            f"Timed out after {TIMEOUT_SECONDS}s waiting for Jenkins stage "
            f"'{STAGE_NAME}'. The pipeline may be hung or unreachable. "
            f"Check Jenkins logs for run-id {RUN_ID}.",
        )

    signal = fetch_signal()

    if signal is None:
        time.sleep(POLL_INTERVAL)
        continue

    # Safety check: make sure this signal is for *our* run
    if str(signal.get("run_id")) != str(RUN_ID):
        die(
            f"Signal Gist contains run_id '{signal.get('run_id')}' "
            f"but we expected '{RUN_ID}'. Possible Gist collision.",
        )

    status  = signal.get("status", "")
    stage   = signal.get("stage", "")
    message = signal.get("message", "")

    # ── pipeline-level error from Jenkins ────────────────────────────────────
    if status == "error":
        print(f"\n[{ts()}] ❌  Jenkins reported a pipeline error:")
        print(f"       Stage:   {stage or '(unknown)'}")
        print(f"       Message: {message or '(no message)'}")
        die(
            f"Jenkins pipeline failed at stage '{stage}': {message}",
        )

    # ── pipeline-level abort ──────────────────────────────────────────────────
    if status == "aborted":
        die(f"Jenkins pipeline was aborted. Stage: '{stage}'. Message: {message}")

    # ── stage completed successfully ──────────────────────────────────────────
    if status == "success" and stage == STAGE_NAME:
        elapsed = TIMEOUT_SECONDS - remaining
        print(f"\n[{ts()}] ✅  Stage '{STAGE_NAME}' completed successfully "
              f"({elapsed:.0f}s elapsed)")
        if message:
            print(f"       Note: {message}")
        sys.exit(0)

    # ── progress update (different stage just finished) ───────────────────────
    if stage and stage != last_seen_stage:
        print(f"[{ts()}] ⏳  Jenkins is at stage '{stage}' (status: {status}) "
              f"— still waiting for '{STAGE_NAME}' …")
        last_seen_stage = stage

    time.sleep(POLL_INTERVAL)
