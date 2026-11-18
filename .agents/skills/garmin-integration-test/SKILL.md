---
name: garmin-integration-test
description: Run and verify a Garmin sync against the local Wanderer development environment, including trail cleanup, one-off cron scheduling, log inspection, API verification, and environment restoration.
---

# Garmin integration test

Use this procedure when asked to run or debug a Garmin sync against the local Wanderer development checkout. This procedure is for the disposable development environment only. It deletes all existing trails and related trail records before the test.

Repository context: this skill is stored in `wanderer-garmin-plugin/.agents/skills/garmin-integration-test/`. The Wanderer checkout is normally the parent directory. Read `DEVELOPMENT.md` and `SCRIPTS.md` in this repository for current commands and helper behavior before running the test.

## Required outcome

Run one Garmin import and confirm that the plugin finishes successfully and the expected imported trails are visible through the Wanderer API. At the end, restore the cron schedule to its original value and stop any environment started solely for this test.

## Procedure

1. **Identify and record the target environment.** Confirm this is the local development stack and identify its Compose project, database service, API base URL, schedule source, and current `POCKETBASE_CRON_SYNC_SCHEDULE`. The documented default is `0 2 * * *`; record the actual value and timezone before changing it. Never target production. Do not print, save, or commit credentials.

2. **Prepare test configuration.** Locate the Garmin plugin instance and configure the test `apiKey` and any other required inputs using the project's existing local configuration mechanism. Keep secrets in the environment or local secret store; do not add them to source, command history, or logs. Record the exact variable/config values that must be restored.

3. **Clear prior trails.** Use `scripts/clear_trails.py` against the development database, following `SCRIPTS.md`. This removes trails and related trail records; it does not clear activities. Verify that the `trails` collection is empty through the API before starting the import. Do not use `reset_wanderer_data.py --delete` for a route-only test because it also deletes activities.

4. **Set a one-off cron schedule.** Choose a one-time expression for the next minute in the database container's timezone. Change the schedule through the existing dev Compose/environment configuration and recreate the database service so it loads the new value:

   ```sh
   docker compose up -d --force-recreate db
   ```

   Avoid `* * * * *`: a sync that imports many activity details can run longer than a minute and overlap another run.

5. **Run and observe the import.** Follow database logs and inspect plugin instance status using the helpers described in `SCRIPTS.md`:

   ```sh
   docker compose logs -f db
   ```

   Wait for the Garmin instance to leave its running state and return to `configured`. Treat a reported error or a stuck run as a failed test; capture relevant sanitized log lines and status details.

6. **Verify results through the API.** Query the Wanderer `trails` collection using the development API and confirm records were created by the Garmin import. Check Garmin origin/external references if available, and compare the records with the expected test account data. A nonzero trail count alone is insufficient if unrelated records could have been created during the test. Also confirm the plugin has no last error and the sync timestamp advanced.

7. **Restore state, even after failure.** In a cleanup/finally step, restore the exact cron schedule and configuration values recorded in steps 1–2, then recreate `db` to apply the restored schedule. Stop only services that this procedure started solely for the test; leave an already-running development stack running. Confirm the effective schedule is restored and report whether import verification passed. Never leave the one-off schedule active.

## Reporting

Report the target development stack, whether the trails collection was cleared, the sync result and verification evidence, any errors, the restored cron schedule, and whether any test-started services were stopped. Do not include API tokens, Garmin keys, or other secrets in the report.
