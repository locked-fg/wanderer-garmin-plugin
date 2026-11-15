# Garmin plugin source review

## Scope and overall assessment

This review focuses on readability and maintainability of the Go plugin source, its parser, tests, and build layout. It is a static review; it does not claim that every Garmin API behavior or production sync path was re-tested.

The plugin is understandable at a high level, and the Garmin response parser has already been separated into a small package with fixture-based tests. The code also uses descriptive domain names and keeps the plugin manifest and build instructions outside the implementation. Those are good foundations.

The source is serviceable, but it is not yet easy to change safely. Most behavior is concentrated in `garmin.go`, mutable package globals carry request state, and the tests cover only a narrow part of the code. I recommend a staged cleanup, starting with the state and date-handling issues below, followed by smaller files and broader unit tests. This is a maintainability assessment, not a recommendation to rewrite the working integration.

## Findings

### High priority

1. **Package globals hold authentication and activity state** (`garmin.go`: `cachedAPIKey`, `cachedAccessToken`, `cachedActivities`, `trackShapeLogged`).

   Their lifetime is the WASM instance, while logical work is per plugin call or sync. That makes the behavior depend on whether the host reuses an instance and whether list/detail calls for different plugin instances share it. `cachedActivities` is not associated with an API key, so a detail request can use a summary cached by an earlier account if activity IDs collide. `trackShapeLogged` also suppresses diagnostics across later calls. Keep this state in an explicit sync/session object or remove caches where the host already supplies the data. If token caching is needed, key it by the credential and give the cache an explicit lifetime.

2. **Invalid `after` values silently disable filtering** (`garmin.go`: `parseAfterDate`).

   The function returns a zero time both when no date was configured and when a configured value cannot be parsed. The caller then treats both cases as “no filter”, which can import older activities unexpectedly. Return `(time.Time, error)` or distinguish absent from invalid configuration, and report invalid values clearly. Add tests for empty, supported, and malformed values.

3. **The Garmin client mixes transport, parsing, logging, paging, and Wanderer mapping** (`garmin.go`, currently the main implementation file).

   Authentication, HTTP request construction, response-shape diagnostics, pagination, filtering, detail retrieval, GPX generation, metadata conversion, and photo conversion all live together. This makes local changes harder to isolate and tests awkward. Split by responsibility, for example into client/auth and request methods, activity sync/detail mapping, and GPX output. Keep the package small and avoid adding layers without a concrete testing or comprehension benefit.

### Medium priority

4. **Retry request construction is duplicated** (`garmin.go`: `request`).

   The first call and token-activation retry repeat the full `HostRequestSpec`. Future header, limit, or connector changes must be kept in sync. Extract a small helper that builds/performs the request, leaving the retry policy visible in `request`.

5. **Several helpers and fields appear unused or redundant** (`main.go`, `garmin.go`).

   `parseInput`, `outputJSON`, `getStringField`, `getIntField`, `getBoolField`, `getTimeField`, `getActivityWorkouts`, and `garminClient.baseURL` have no apparent callers in the reviewed source. Remove dead code after confirming no generated or build-tag-specific caller exists. `newGarminClient` is also created for the routes capability even though Garmin routes immediately returns an empty list; the intended capability behavior should be made explicit so users do not infer that saved courses are imported.

6. **The response model contains repeated custom decoding logic** (`types.go`).

   `garminActivityType` and `garminPrivacy` implement nearly identical `UnmarshalJSON` methods to accept either a string or an object. A shared small decoder/helper or a focused test for the two accepted shapes would make this compatibility behavior easier to understand and less likely to diverge. Avoid generalizing it unless the API really uses both forms.

7. **GPX is assembled with string formatting and raw timestamps** (`garmin.go`: `buildGPXFromTrackPoints`).

   Garmin timestamps are inserted directly into XML. A malformed or unexpected timestamp containing XML metacharacters can produce invalid GPX. Use `encoding/xml` or validate/escape element content, and test empty points, elevation, timestamp, and special characters. Also consider naming the generated track from the activity rather than a fixed label if that is useful to users.

8. **The parser accepts questionable coordinates without validation** (`internal/garmintrack/parser.go`).

   It skips invalid metric indices and null/non-numeric values, but accepts out-of-range latitude/longitude and potentially non-finite numbers. Define the intended policy and validate coordinates before producing a track. Tests should cover invalid indices, missing descriptors, nulls, invalid coordinates, and a polyline where every point is marked invalid.

### Lower priority

9. **Some comments narrate implementation rather than explain constraints** (`garmin.go`, `main.go`).

   Comments such as “Create a new Garmin client” and “Helper function to fail with an error” add little. Keep comments for Garmin-specific API quirks, host lifecycle assumptions, and non-obvious compatibility behavior; let function names describe routine operations.

10. **The `types.go` aliases obscure ownership slightly**.

    Most plugin-facing types are direct aliases to SDK types. This is compact, but the repeated comments (“use SDK type”) make the file read like a list of wrappers. Either use SDK types directly where that stays readable, or retain aliases only where they add domain meaning or simplify signatures.

11. **Build dependencies are not fully represented in the Makefile target** (`Makefile`).

    The `plugin.wasm` target lists Go files but omits `go.mod` and `go.sum`; changing the pinned SDK dependency may not trigger a rebuild when a stale artifact exists. Add module files and any other build inputs to the prerequisites, or make the build target phony if that better matches the expected workflow.

## Testing gaps

The parser tests are a useful start and use sanitized fixtures. The following important behavior currently has no visible focused tests:

- `parseAfterDate` and `parseActivityStartTime` supported formats and failures.
- Garmin activity type/privacy string-or-object decoding.
- Pagination state and filtering when a page contains old or invalid dates.
- Mapping of activity metadata, photos, timestamps, and empty/missing fields.
- GPX generation and XML validity.
- Token reuse/retry decisions and response status/error handling.

Prioritize pure functions first: date parsing, response decoding, mapping, and GPX generation can be tested without a live Garmin account or a Wanderer container. Add an injectable request function only if transport behavior needs unit tests; avoid building an elaborate mock framework for this small plugin.

## Suggested improvement order

1. Add tests for date parsing and change invalid configured dates to an explicit error.
2. Make sync/cache state explicit and prevent cached activity data from crossing accounts or syncs.
3. Split `garmin.go` along the client, sync, mapping, and GPX responsibilities.
4. Add unit tests for mapping, GPX, and the flexible Garmin JSON fields.
5. Remove confirmed dead helpers and simplify comments/type aliases.
6. Fix the Makefile dependency list and document that the routes capability currently returns no routes.

These steps preserve the current integration shape while making future changes easier to reason about and safer to verify.
