---
name: garmin-integration-test
description: Test Garmin plugin changes with fixture-first checks, a bounded live smoke test when needed, API verification, and state restoration; destructive cleanup is limited to isolated data.
---

# Garmin integration test

Use this procedure when asked to run or debug the Garmin plugin against the local Wanderer development checkout. Prefer deterministic fixture tests for parsing and filtering, followed by a small live smoke test. A live sync can create persistent data; never clear existing trails as a routine setup step. Run destructive or full-history scenarios only against an isolated disposable database or after preserving a recoverable snapshot and confirming the exact impact with the user.

Repository context: this skill is stored in `wanderer-garmin-plugin/.agents/skills/garmin-integration-test/`. The Wanderer checkout is normally the parent directory. Read `DEVELOPMENT.md` and `SCRIPTS.md` in this repository for current commands and helper behavior before running the test.

## Test sequence

Use the shortest test that provides evidence for the requested change. Do not wait for unrelated historical activities or for every background sync in the environment to finish.

1. **Read the task and choose coverage.** Read `DEVELOPMENT.md`, `SCRIPTS.md`, and the referenced messageboard item. Translate each requested behavior into an assertion. For an `after` filter, cover: missing option, empty option if behavior differs, valid ISO date, invalid date, and malformed Garmin start time if relevant. Prefer fixture-backed/unit tests for these cases. The valid fixture must include one activity before the cutoff and one exactly on or after it. A full-history import (`after` empty) is not a suitable default live test because it can take a long time.

2. **Run fast checks first.** Run the relevant Go tests and pinned TinyGo build/install checks described by the project. A build or unit test failure should stop the live test unless the live behavior is specifically needed to diagnose it.

3. **Identify the target environment and baseline.** Confirm the local development stack and identify its Compose project, DB service, API base URL, schedule source, and `POCKETBASE_CRON_SYNC_SCHEDULE`. Record the exact schedule and timezone before changing them. Never target production. Use authenticated API reads for user-visible counts, status, and trail fields; an unauthenticated zero count is not evidence that the collection is empty. Do not print resolved Compose configuration because it can contain credentials. Read only the schedule or redact secrets.

4. **Prepare safe test data and config.** Prefer a dedicated test account, isolated DB, or fixture. Record current route/activity counts and preserve a restorable snapshot before any destructive action. Do not run `scripts/clear_trails.py` on a shared or valuable dev dataset by default: it permanently removes trails and external references. Do not use `reset_wanderer_data.py --delete` for a route-only test because it also deletes activities. Keep API keys in the environment/local secret store, never in command text, history, source, or logs. Back up the exact plugin config and change only the target option, using the authenticated API if it supports the update. Use SQLite only if no supported API write exists; then use a reviewed helper that preserves unrelated JSON fields and restores in a `finally` block.

5. **Run one small live smoke test when needed.** Choose a recent cutoff or a dedicated test dataset expected to return few activities. Record the original cron schedule and config first. Set a one-off cron schedule through the existing dev Compose/environment configuration and recreate the DB service. Calculate the schedule in the DB container's timezone; if `TZ` is unset, inspect the image/runtime timezone or explicitly use UTC rather than guessing. Avoid `* * * * *`, which can overlap long imports. Do not wait for unrelated sync work to complete.

   ```sh
   docker compose up -d --force-recreate db
   ```

6. **Prove that the sync actually ran.** Follow sanitized, Garmin-filtered DB logs and inspect plugin status. A `syncing` field alone is not proof: require evidence of a real transition from baseline, such as a newer `last_sync_at`, a Garmin start/completion log, or new matching API records. Poll with a bounded timeout. If the trigger helper only changes status but produces no worker/log/API activity, stop treating it as active, recover the stale status safely, and use the documented cron path. Do not leave an instance stuck in `syncing`.

   ```sh
   docker compose logs -f db
   ```

   Completion requires the plugin to leave the running state and return to `configured`, with no nonempty last error. Treat both `null` and `{}` as empty error values. Report nested plugin errors even if the host wraps them (for example, `provider_unavailable` containing `invalid_request`). Filter unrelated provider logs out of the evidence.

7. **Verify via API; use DB only for gaps.** Use authenticated API responses to check plugin status, `last_sync_at`, and newly created Garmin trails/external references. Compare against the baseline or isolated test data. Counts alone do not prove date filtering. The exposed trail `date` may not be the original Garmin start time; verify source dates from fixtures/plugin-level assertions or explicit aggregate filter diagnostics (`received`, `eligible`, `filtered`, `invalid start time`). Use read-only SQL for diagnosis only when the API does not expose a needed field, and do not treat DB inspection as the final proof of client-visible behavior.

8. **Stop early when the result is sufficient.** If the test is taking longer than the bounded timeout or the user asks to shorten it, stop the one-off run safely. Report it as partial/stopped early, not passed. Preserve concise evidence (elapsed time, last observed status/time, authenticated API result, and relevant sanitized logs); do not spend time waiting for all historical data to sync unless full-history completion is explicitly the test target.

9. **Restore state in a `finally` step.** Restore the exact original config and cron schedule, recreate the DB service if needed, and verify both through the API/config read and effective schedule. Clear stale status only after confirming no worker is still active. Stop only services started solely for the test; leave the already-running stack alone. Verify and report any data changes; restoring config/schedule does not restore deleted or imported records. Never leave a one-off schedule active.

## Reporting

Report which fixture/unit/build/live checks ran and passed, whether a live run completed or was stopped early, the specific evidence (including API authentication and baseline/delta), data impact, errors, restored config/schedule, and any services stopped. Distinguish `passed`, `failed`, `timed out`, and `partial/stopped early`. Do not include API tokens, Garmin keys, raw secret-bearing config, or unsanitized logs.
