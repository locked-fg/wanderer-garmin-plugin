# Garmin plugin for Wanderer

Imports Garmin Connect activities into Wanderer trails, including route geometry, activity metrics, dates, and available photos. It authenticates with a Garmin Personal API token and provides activity and route capabilities through Wanderer's plugin SDK.

## Develop and install

- [Development guide](DEVELOPMENT.md): architecture, tests, and build process.
- [Install guide](INSTALL.md): build, install, restart, and configure the plugin.
- [Helper scripts](SCRIPTS.md): diagnostics and local development utilities in `scripts/`.

## Implementation notes

Activities are requested in pages of 50. The plugin continues while Garmin returns a full page and stops on a shorter page. The optional `after` date filters older activities; it does not cap the number imported. Activity track geometry is parsed from Garmin's activity details response. The parser fixtures in `testdata/` capture representative response shapes.

## Repository layout

- `main.go`, `garmin.go`, `types.go`: plugin entry points, Garmin client, and response mapping.
- `internal/garmintrack/`: parser for activity detail and track responses.
- `go.mod`: pins the Wanderer plugin SDK version used by the build.
- `testdata/`: sanitized Garmin response fixtures used by parser tests.
- `scripts/`: optional build, verification, and Wanderer debugging helpers.

The build downloads the pinned SDK module from `go.mod`; a separate SDK checkout is not required. To look up or edit the SDK source locally, clone the full Wanderer repository beside this plugin repository:

```bash
git clone https://github.com/open-wanderer/wanderer.git ../wanderer-sdk-source
```

The resulting layout is:

```text
parent-directory/
├── wanderer-garmin-plugin/
└── wanderer-sdk-source/
    └── plugins/sdk/
```
