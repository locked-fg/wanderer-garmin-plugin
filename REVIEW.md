# Garmin plugin source review

## Scope and overall assessment

This review focuses on readability and maintainability of the Go plugin source, its parser, tests, and build layout. It is a static review; it does not claim that every Garmin API behavior or production sync path was re-tested.

The plugin is understandable at a high level, and the Garmin response parser has already been separated into a small package with fixture-based tests. The code also uses descriptive domain names and keeps the plugin manifest and build instructions outside the implementation. Those are good foundations.

The source is serviceable, but it is not yet easy to change safely. Most behavior is concentrated in `garmin.go`, and the tests cover only a narrow part of the code. I recommend a staged cleanup of the remaining findings below, followed by smaller files and broader unit tests. This is a maintainability assessment, not a recommendation to rewrite the working integration.

## Findings

All findings from the original review have been addressed and squash-merged to `main` (`7319e27`, formerly the `code-cleanup` branch):

- Package globals / cross-account state: resolved earlier (per-instance `garminClient` state).
- Invalid `after` values silently disabling filtering: resolved earlier (`internal/garminutil` returns `ConfigError`, host reports `invalid_request`).
- `garmin.go` mixing transport, paging, logging, and mapping: split into `client.go`, `sync.go`, `map.go`.
- Duplicated retry request construction: extracted `doRequest` helper in `client.go`.
- Dead helpers/fields: removed (`parseInput`, `outputJSON`, field getters, `getActivityWorkouts`, write-only `baseURL`); routes limitation documented in `plugin.json`.
- Duplicated string-or-object decoding: shared `garminutil.UnmarshalTypeKey` with tests.
- GPX string formatting: `garmintrack.BuildGPX` uses `encoding/xml` with tests.
- Parser coordinate validation: `validLatLon` rejects out-of-range/non-finite coordinates, with tests.
- Narrating comments and `types.go` alias comments: removed; constraint/quirk comments kept.
- Makefile build dependencies: `plugin.wasm` now depends on all Go files plus `go.mod`/`go.sum`.

## Testing gaps

Covered with unit/fixture tests: date parsing and invalid-`after` handling, track parser (polyline, metrics fallback, invalid coordinates), GPX generation and XML validity, string-or-object type decoding, activity search filtering/pagination (`internal/garminapi`), and full detail mapping incl. photos and timestamps (`internal/garminmap`).

Host-dependent behavior (pagination against the live API, token exchange/retry) was validated via live sync by the local agent (messageboard `2026-11-18-004`, status `done`).

## Suggested improvement order

All steps complete: split `garmin.go`, broadened unit tests (incl. mapping/pagination fixtures via `internal/garminapi` + `internal/garminmap`), removed dead code, fixed Makefile dependencies — merged to `main` (`7319e27`), live validation passed (messageboard `2026-11-18-004`). No open follow-ups.
