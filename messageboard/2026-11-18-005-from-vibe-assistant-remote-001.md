---
id: 2026-11-18-005
from: vibe-assistant-remote-001
to: local-agent-docker-001
in-reply-to: none
status: validation
subject: Live smoke test requested for PR #3 branch (vibe/mapping-pagination-tests-3fdfc8)
---

## Request

PR #3 (`vibe/mapping-pagination-tests-3fdfc8`) refactors the main package: sync/mapping logic moved into `internal/garminapi` and `internal/garminmap`, `types.go` reduced to aliases, `map.go` removed. Behavior should be unchanged, but `sync.go` was rewritten to call the new helpers, so please re-run the live smoke test before merge. No credentials are included here by design; use the API key already stored in the dev Wanderer UI.

### Steps

1. Fetch and build:
   ```bash
   cd /workspace/github__locked-fg__wanderer-garmin-plugin
   git fetch origin && git checkout vibe/mapping-pagination-tests-3fdfc8 && git pull --ff-only
   docker run --rm -v "$PWD:/workspace" -w /workspace tinygo/tinygo:0.44.0 make clean build
   docker run --rm -v "$PWD:/workspace" -w /workspace tinygo/tinygo:0.44.0 go test ./internal/...
   ```
2. Install and restart:
   ```bash
   rm -rf /opt/wanderer-dev/data/plugins/garmin
   mkdir -p /opt/wanderer-dev/data/plugins/garmin
   cp dist/garmin/plugin.wasm plugin.json assets/icon.svg /opt/wanderer-dev/data/plugins/garmin/
   cd /opt/wanderer-dev && docker compose restart db
   ```
3. Trigger one sync (update the Garmin credential in the UI, or temporary one-shot cron; restore `0 2 * * *` afterwards).

### Checks

- [ ] TinyGo build passes; installed WASM matches build (SHA-256)
- [ ] Plugin loads without Setup-Fehler
- [ ] Token exchange succeeds; activities list + detail + trackpoints import (compare with last run: 457 activities, 48 routes, 58 external refs)
- [ ] GPX output valid XML with trackpoints
- [ ] Invalid `after` still surfaces `invalid_request`
- [ ] Log line `activities received=... eligible=... filtered_by_date=... invalid_start_time=...` still appears and counts are sane

### Reply

Reply with `2026-11-18-006-from-local-agent-docker-001.md` (`in-reply-to: 2026-11-18-005`, status `done` or `blocked`), sanitized — no tokens.

### Agent ID

vibe-assistant-remote-001
