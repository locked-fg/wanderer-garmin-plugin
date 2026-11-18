# Development

## What is implemented

The plugin implements Wanderer's `list_activities`, `get_activity_detail`, `list_routes`, and `get_route_detail` capabilities. Garmin Personal API authentication exchanges the configured `apiKey` for an access token. Activity listing uses 50-item pages and continues until Garmin returns a shorter page. The optional `after` date filters the import window. Activity details are mapped to Wanderer trail metadata and track points; Garmin detail response variations are handled by `internal/garmintrack`.

## Project structure

- `main.go`: exported WASM capability entry points.
- `garmin.go`, `types.go`: Garmin API requests, pagination, and Wanderer data mapping.
- `internal/garmintrack/`: activity-detail JSON and track parser.
- `internal/garmintrack/parser_test.go`: fixture-based parser tests.
- `testdata/`: sanitized example API responses used by those tests.
- `go.mod`: pins the Wanderer plugin SDK module version used by this plugin.
- `scripts/`: optional build and instance-debugging helpers.

## Requirements

Docker is the supported build environment; host Go/TinyGo installations are not required. Parser tests use Go's standard test runner. The tested TinyGo image is `tinygo/tinygo:0.44.0`.

## Test and build

Run the parser tests from this directory in the supported Docker environment:

```bash
docker run --rm -v "$PWD:/workspace" -w /workspace tinygo/tinygo:0.44.0 go test ./internal/garmintrack
```

The SDK version is pinned in `go.mod` and is downloaded during the build, so a separate checkout is not required to build the plugin. To inspect or edit the SDK source locally, clone the full Wanderer repository from this directory:

```bash
git clone https://github.com/open-wanderer/wanderer.git ../wanderer-sdk-source
```

This creates the sibling layout below; the SDK module is under `plugins/sdk/`:

```text
parent-directory/
├── wanderer-garmin-plugin/
└── wanderer-sdk-source/
    └── plugins/sdk/
```

Build the distributable files into `dist/garmin/`:

```bash
docker run --rm -v "$PWD:/workspace" -w /workspace tinygo/tinygo:0.44.0 make clean install
```

Install into the sibling Wanderer checkout's `data/plugins/garmin/` directory:

```bash
docker run --rm \
  -v "$PWD/..:/workspace" \
  -w "/workspace/$(basename "$PWD")" \
  tinygo/tinygo:0.44.0 make clean install-local
```

From the Wanderer repository root (the parent of the plugin directory), restart the database service to reload plugin files:

```bash
docker compose restart db
```

`make clean` removes generated build output. `/dist/` is ignored by Git. `make install-local` uses `WANDERER_DATA_DIR` if set; otherwise it installs to the sibling `data/` directory. The install command mounts the parent directory so that default location is available inside Docker.

## Git Workflow and Branching Strategy

### Branching Model

**Always use feature branches for development.** This ensures `main` remains stable and functional at all times, which is critical when multiple agents collaborate.

**Branch Naming Convention:**
```
feature/<short-description>
fix/<issue-description>
docs/<documentation-topic>
```

**Examples:**
- `feature/fix-caching`
- `feature/garmin-parser`
- `fix/token-exchange`
- `docs/install-guide`

### Standard Workflow

1. **Start from main:**
   ```bash
   git checkout main
   git pull origin main
   ```

2. **Create feature branch:**
   ```bash
   git checkout -b feature/<description>
   git push origin feature/<description>
   ```

3. **Make changes and commit:**
   ```bash
   git add .
   git commit -m "<type>: <description>"
   git push origin feature/<description>
   ```

4. **Test thoroughly**

5. **Merge to main (with squash):**
   **Always squash feature branch commits into a single commit with a comprehensive message.**
   
   ```bash
   git checkout main
   git merge --squash feature/<description>
   git commit -m "feat/<description>: Comprehensive commit message
   
   This commit addresses...
   
   ## Changes
   - Change 1
   - Change 2
   
   ## Validation
   - Test 1 passed
   - Test 2 passed
   
   Agent-ID: <your-agent-id>
   Status: Validated"
   git push origin main
   ```
   
   **Why squash?**
   - Keeps main history clean and readable
   - Combines multiple implementation steps into one logical change
   - Makes it easier to understand what was done and why
   - Reduces noise in git history

6. **Clean up:**
   ```bash
   git branch -d feature/<description>
   git push origin --delete feature/<description>
   ```

### Why Branches Are Required

- **Main stability:** Main should always be in a working state
- **Isolated experiments:** Multiple agents can work on different features simultaneously
- **Easy rollback:** If something breaks, simply delete the branch
- **Code review:** Pull requests enable review before merging
- **Upstream contribution:** GitHub PRs require branches

### When Main Development Is Acceptable

Only for:
- Tiny, safe fixes (e.g., documentation typos)
- Emergency hotfixes that must be deployed immediately
- Single-agent work with immediate testing

**For multi-agent collaboration: Always use branches!**

---

## Development checks

After installation, add/configure Garmin in Wanderer, provide the `apiKey`, and sync. Inspect plugin status and container logs if the import fails. Optional helpers and their invocation are documented in [SCRIPTS.md](SCRIPTS.md). Some cleanup helpers delete local Wanderer data; check each script's `--help` and dry-run option before using it.

For Garmin plugin validation in the local development environment, follow the [Garmin integration test skill](.agents/skills/garmin-integration-test/SKILL.md). It prioritizes fixture tests and a bounded live smoke test, and covers API verification, safe state handling, and restoring the schedule and configuration.

For shared testing principles on choosing test scope, verifying results, limiting side effects, handling credentials, and reporting partial runs, see the repository [testing guide](../TESTING.md).

### Run a one-off plugin sync in the dev instance

The dev Compose file sets `POCKETBASE_CRON_SYNC_SCHEDULE` to `0 2 * * *`. To test an import immediately, temporarily change it to a one-time cron expression for the next minute in the database container's timezone, then recreate the database service so it picks up the schedule:

```bash
docker compose up -d --force-recreate db
docker compose logs -f db
```

Wait for the plugin instance to return to `configured` and review the import logs. Restore the normal schedule (`0 2 * * *`) and recreate `db` afterward. Avoid `* * * * *` for Garmin imports: a sync that fetches many activity details can take longer than a minute, allowing another run to start before the first one finishes.

For API cleanup, `scripts/reset_wanderer_data.py` reports the Activities (`summit_logs`) and Routes (`trails`) counts by default. Its `--delete` option deletes both collections, so use it only when both should be cleared. Route-only cleanup should delete records from the `trails` collection through the Wanderer API and verify that the collection is empty afterward. The script requires `WANDERER_API_TOKEN`; do not put the token in command history or documentation.
