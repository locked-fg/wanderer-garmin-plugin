#!/usr/bin/env python3
"""Reset imported Wanderer activities/routes through the Wanderer API.

By default this only reports record counts. Pass --delete to remove summit
logs (Activities) and trails (Routes), then verify both collections are empty.
Provide WANDERER_API_TOKEN in the environment or enter it at the prompt.
"""

import argparse
import getpass
import json
import os
import sys
import urllib.error
import urllib.parse
import urllib.request


def request_json(base_url, path, token=None, method="GET", body=None):
    data = None if body is None else json.dumps(body).encode()
    headers = {"Accept": "application/json"}
    if data is not None:
        headers["Content-Type"] = "application/json"
    if token:
        headers["Authorization"] = "Bearer " + token
    request = urllib.request.Request(
        base_url.rstrip("/") + path, data=data, headers=headers, method=method
    )
    try:
        with urllib.request.urlopen(request, timeout=30) as response:
            payload = response.read()
            return response.status, json.loads(payload) if payload else None
    except urllib.error.HTTPError as error:
        detail = error.read().decode("utf-8", errors="replace")[:500]
        raise RuntimeError(f"API {method} {path} returned {error.code}: {detail}") from error
    except urllib.error.URLError as error:
        raise RuntimeError(f"Could not reach Wanderer API: {error.reason}") from error


def authenticate(base_url, api_token):
    _, result = request_json(base_url, "/auth/token", body={"api_token": api_token}, method="POST")
    return result["token"]


def list_record_ids(base_url, token, collection):
    ids = []
    page = 1
    while True:
        query = urllib.parse.urlencode({"perPage": 500, "page": page})
        _, result = request_json(
            base_url, f"/api/collections/{collection}/records?{query}", token
        )
        ids.extend(record["id"] for record in result.get("items", []))
        if page >= result.get("totalPages", 1):
            return ids
        page += 1


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--base-url", default="http://127.0.0.1:8091")
    parser.add_argument(
        "--delete", action="store_true", help="delete Activities and Routes (default is report-only)"
    )
    args = parser.parse_args()

    api_token = os.environ.get("WANDERER_API_TOKEN") or getpass.getpass("Wanderer API token: ")
    try:
        token = authenticate(args.base_url, api_token)
        collections = (("summit_logs", "Activities"), ("trails", "Routes"))
        initial = {}
        for collection, label in collections:
            initial[collection] = list_record_ids(args.base_url, token, collection)
            print(f"{label}: {len(initial[collection])}")

        if not args.delete:
            print("Report only. Pass --delete to remove these records.")
            return 0

        for collection, label in collections:
            deleted = 0
            for record_id in initial[collection]:
                path = f"/api/collections/{collection}/records/{record_id}"
                request_json(args.base_url, path, token, method="DELETE")
                deleted += 1
            print(f"Deleted {label}: {deleted}")

        remaining = {}
        for collection, label in collections:
            remaining[collection] = list_record_ids(args.base_url, token, collection)
            print(f"Remaining {label}: {len(remaining[collection])}")
        if any(remaining.values()):
            return 1
        print("API verification passed: Activities and Routes are empty.")
        return 0
    except (KeyError, RuntimeError, json.JSONDecodeError) as error:
        print(f"Reset failed: {error}", file=sys.stderr)
        return 1


if __name__ == "__main__":
    sys.exit(main())
