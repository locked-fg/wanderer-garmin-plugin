# Garmin Plugin Debugging Scripts

## Overview
This directory contains Python scripts for debugging and managing the Garmin plugin.

## Available Scripts

### 🔍 Check API Key
- **`verify_key.py`** - Verifies the Garmin API key is saved in the database
- **`check_auth_fields.py`** - Checks what auth fields exist in the plugin instance
- **`check_exact_key.py`** - Checks if the API key matches the exact expected value

### 🔄 Trigger Import
- **`trigger_sync.py`** - Triggers a plugin sync by updating the timestamp
- **`check_and_trigger_import.py`** - Combined script: checks API key + triggers import

### 🧹 Database Cleanup
- **`clear_trails.py`** - Deletes all trails from the database (for testing)

### 📊 Status Checks
- **`check_plugin_status.py`** - Checks the plugin instance status
- **`check_all_instances.py`** - Lists all plugin instances
- **`check_current_error.py`** - Checks for current errors in plugin instances
- **`check_full_instance.py`** - Shows full plugin instance data

### 🔐 WASM Verification
- **`check_wasm.py`** - Verifies the WASM module has correct exports
- **`verify_build.py`** - Verifies the plugin build

### ⚡ Build Scripts
- **`build_plugin.py`** - Builds the plugin WASM using Docker

### 🗃️ Database
- **`decrypt_test.py`** - Tests decryption of encrypted auth fields

---

## Usage Examples

### Check if API Key is Saved
```bash
python3 check_auth_fields.py
python3 verify_key.py
```

### Trigger Import
```bash
python3 trigger_sync.py
# or
python3 check_and_trigger_import.py
```

### Check Plugin Status
```bash
python3 check_plugin_status.py
```

### Build Plugin
```bash
python3 build_plugin.py
```

### Clear Trails (for testing)
```bash
python3 clear_trails.py
```

---

## Common Workflows

### 1. Verify Setup and Trigger Import
```bash
# Check if API key is saved
python3 check_auth_fields.py

# Trigger import
python3 trigger_sync.py
```

### 2. Full Debug Cycle
```bash
# Check plugin status
python3 check_plugin_status.py

# Verify API key
python3 verify_key.py

# Trigger sync
python3 trigger_sync.py

# Check status again
python3 check_plugin_status.py
```

### 3. Clean Database for Testing
```bash
# Clear all trails
python3 clear_trails.py

# Trigger fresh import
python3 trigger_sync.py
```

---

## Notes
- All scripts use relative path `../data/pb_data/data.db` to access the PocketBase database
- Scripts are designed to be run from this directory: `wanderer-garmin-plugin/`
- Auth fields (like API keys) are encrypted by PocketBase, so they appear as encrypted strings in the database
- The `mN` prefix in encrypted values indicates PocketBase's encryption format
