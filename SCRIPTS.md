# Helper scripts

Optional Python utilities live in `scripts/`. Run them from the plugin repository root so their paths to the sibling Wanderer `data/` directory resolve correctly:

```bash
cd /opt/wanderer-dev/wanderer-garmin-plugin
python3 scripts/check_plugin_status.py
```

## Build and verification

- `build_plugin.py`: build, optionally install, and verify the WASM using Docker (`--help` for options).
- `check_wasm.py`: inspect a WASM file's imports and size.
- `verify_build.py`: check plugin manifest/auth configuration.

## Diagnostics

- `check_plugin_status.py`, `check_all_instances.py`, `check_current_error.py`, `check_full_instance.py`: inspect plugin instance state and errors.
- `check_auth_fields.py`, `verify_key.py`, `check_exact_key.py`, `check_and_trigger_import.py`: inspect saved auth fields. These scripts do not print a stored credential in full.
- `trigger_sync.py`, `check_and_trigger_import.py`: trigger a sync.
- `decrypt_test.py`, `test_save_process.py`: inspect/test local save and encryption behavior.

## Data cleanup

- `reset_wanderer_data.py`: report or delete routes and activities through the Wanderer API; supports a confirmation-oriented `--delete` mode.
- `clear_trails.py`, `clear_all_activity_data.py`: delete local activity/trail records from SQLite; inspect `--help` and use `--dry-run` before deletion.
- `clear_meilisearch.py`: clear the local Meilisearch index.

Cleanup scripts are destructive. Use them only against a disposable development instance and check their options before running. Scripts that query SQLite expect the development checkout at `/opt/wanderer-dev` by default; scripts with a `--db-path` option can target a different database.
