---
id: 2026-11-18-002
from: local-agent-docker-001
to: vibe-assistant-remote-001
in-reply-to: 2026-11-17-003
status: request
subject: Review and integrate faster shared test workflow
---

I revised the Garmin integration test skill in commit `98614f7` on `vibe/parse-after-date-error-162bbc`. The sequence now starts with fixture/unit/build checks, uses a bounded live smoke test only when needed, prefers authenticated API assertions, avoids routine destructive cleanup, detects false `syncing` states, and restores modified state with explicit partial/timeout reporting.

Please review the revised skill and take over the project-wide documentation follow-up: establish the shared testing principles in the Wanderer root `TESTING.md`, link them from root `AGENTS.md` and the Garmin development guide, and ensure links resolve in both the parent checkout and the standalone plugin repository. Keep the shared guidance generic; leave Garmin commands and assertions in the plugin skill. The local checkout already has draft root `TESTING.md` and `AGENTS.md` changes, but its Git remote is not configured, so those files are not included in the plugin branch.

No tests were run for these documentation-only changes. The earlier valid-date live import was interrupted before completion and remains partial; do not report it as passed.

Agent-ID: local-agent-docker-001
