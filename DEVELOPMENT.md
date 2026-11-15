# Development

## What is implemented

The plugin implements Wanderer's `list_activities`, `get_activity_detail`, `list_routes`, and `get_route_detail` capabilities. Garmin Personal API authentication exchanges the configured `apiKey` for an access token. Activity listing uses 50-item pages and continues until Garmin returns a shorter page. The optional `after` date filters the import window. Activity details are mapped to Wanderer trail metadata and track points; Garmin detail response variations are handled by `internal/garmintrack`.

## Project structure

- `main.go`: exported WASM capability entry points.
- `garmin.go`, `types.go`: Garmin API requests, pagination, and Wanderer data mapping.
- `internal/garmintrack/`: activity-detail JSON and track parser.
- `internal/garmintrack/parser_test.go`: fixture-based parser tests.
- `testdata/`: sanitized example API responses used by those tests.
- `sdk/`: local Wanderer plugin SDK module referenced from `go.mod`.
- `scripts/`: optional build and instance-debugging helpers.

## Requirements

Docker is the supported build environment; host Go/TinyGo installations are not required. Parser tests use Go's standard test runner. The tested TinyGo image is `tinygo/tinygo:0.44.0`.

## Test and build

Run Go tests for the parser and local SDK from this directory:

```bash
go test ./internal/garmintrack
(cd sdk && go test ./...)
```

Build the distributable files into `dist/garmin/`:

```bash
docker run --rm -v "$PWD:/workspace" -w /workspace/wanderer-garmin-plugin tinygo/tinygo:0.44.0 make clean install
```

Install into the sibling Wanderer checkout's `data/plugins/garmin/` directory:

```bash
docker run --rm -v "$PWD:/workspace" -w /workspace/wanderer-garmin-plugin tinygo/tinygo:0.44.0 make clean install-local
```

From the Wanderer repository root, restart the database service to reload plugin files:

```bash
docker compose restart db
```

`make clean` removes generated build output. `/dist/` is ignored by Git. `make install-local` uses `WANDERER_DATA_DIR` if set; otherwise it installs to the sibling `data/` directory.

## Development checks

After installation, add/configure Garmin in Wanderer, provide the `apiKey`, and sync. Inspect plugin status and container logs if the import fails. Optional helpers and their invocation are documented in [SCRIPTS.md](SCRIPTS.md). Some cleanup helpers delete local Wanderer data; check each script's `--help` and dry-run option before using it.
