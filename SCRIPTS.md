# Helper scripts

Optional Python utilities live in `scripts/`. They locate the sibling Wanderer `data/` directory from their own location, so they work from any current directory. For example, from the plugin repository root:

```bash
python3 scripts/check_plugin_status.py
```

For a different data directory, set `WANDERER_DATA_DIR` to its path. The directory should contain `pb_data/data.db` and `plugins/garmin/`:

```bash
WANDERER_DATA_DIR="../data" python3 scripts/check_plugin_status.py
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

For the local dev checkout, the Wanderer API token is available in the `WANDERER_DEV_API` environment variable. `reset_wanderer_data.py` reads `WANDERER_API_TOKEN`, so pass it without printing the value:

```bash
WANDERER_API_TOKEN="$WANDERER_DEV_API" python3 scripts/reset_wanderer_data.py
```

Do not copy the token value into documentation or command text.

## Data cleanup

- `reset_wanderer_data.py`: report or delete routes and activities through the Wanderer API; supports a confirmation-oriented `--delete` mode.
- `clear_trails.py`, `clear_all_activity_data.py`: delete local activity/trail records from SQLite; inspect `--help` and use `--dry-run` before deletion.
- `clear_meilisearch.py`: clear the local Meilisearch index.

Cleanup scripts are destructive. Use them only against a disposable development instance and check their options before running. SQLite helpers default to `data/pb_data/data.db` beside the Wanderer checkout; some also accept `--db-path` for an explicit database file.
