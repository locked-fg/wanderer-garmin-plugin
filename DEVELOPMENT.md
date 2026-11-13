# Development Guide - Wanderer Garmin Plugin

This document describes the development and build process for the Garmin plugin for Wanderer.

---

## 🛠️ Prerequisites

- **Docker** - Required for building without host Go installation
- **Git** - For version control
- **Wanderer Instance** - For testing the plugin

**No Go installation required on the host system!** ✅

---

## 🚀 Quick Start

### Clone & Build

```bash
# Clone the repository (if not already done)
git clone <repository-url> wanderer-garmin-plugin
cd wanderer-garmin-plugin

# Build using Docker
docker run --rm -v $(pwd):/src -w /src tinygo/tinygo:0.44.0 make build

# Install to local Wanderer
mkdir -p /opt/wanderer-dev/data/plugins/garmin
cp dist/garmin/plugin.json /opt/wanderer-dev/data/plugins/garmin/
cp dist/garmin/plugin.wasm /opt/wanderer-dev/data/plugins/garmin/
cp assets/icon.svg /opt/wanderer-dev/data/plugins/garmin/icon.svg

# Restart Wanderer
docker compose restart wanderer-db
```

---

## 📦 Project Structure

```
wanderer-garmin-plugin/
├── main.go           # Plugin entry point and exported functions
├── garmin.go         # Garmin API client and data mapping
├── types.go          # Type definitions
├── go.mod            # Go module definition
├── go.sum            # Go dependencies checksum
├── Makefile          # Build scripts
├── plugin.json       # Plugin manifest
├── README.md         # User documentation
├── INSTALL.md        # Installation guide
├── DEVELOPMENT.md    # This file
├── assets/
│   └── icon.svg      # Plugin icon
└── sdk/              # Local copy of Wanderer Plugin SDK
    ├── go.mod
    ├── types.go
    ├── host_http.go
    └── ...
```

---

## 🔨 Build Methods

### Method 1: Using TinyGo in wanderer-db Container

**Pros:**
- Simple integration with existing infrastructure
- No separate containers to manage

**Cons:**
- Modifies the production container
- Build dependencies mixed with runtime

```bash
# Install TinyGo and Make in the container
docker exec -it wanderer-db apt-get update && apt-get install -y tinygo make

# Build and install
docker exec -it wanderer-db bash -c "cd /app/plugins/garmin && make build && make install-local"
docker compose restart wanderer-db
```

### Method 2: Using Separate Dev Container (Recommended)

**Pros:**
- ✅ No modification to production containers
- ✅ Clean separation of build and runtime
- ✅ Easier to maintain and update
- ✅ Better for CI/CD pipelines

**Cons:**
- Requires additional container

```bash
# Build using standalone TinyGo container
docker run --rm -v $(pwd):/src -w /src tinygo/tinygo:0.44.0 make build

# Install manually
make install-local
```

Or with explicit Docker command:

```bash
# Build
docker run --rm \
  -v $(pwd):/src \
  -v /opt/wanderer-dev/data:/opt/wanderer-dev/data \
  -w /src \
  tinygo/tinygo:0.44.0 make install-local

# Restart
docker compose restart wanderer-db
```

---

## 📝 Makefile Targets

| Target | Description |
|--------|-------------|
| `make build` | Build plugin.wasm only |
| `make install` | Build and copy to dist/garmin/ |
| `make install-local` | Build and install to ../../../data/plugins/garmin/ |
| `make clean` | Remove build artifacts |
| `make package` | Create tar.gz package for distribution |

### Custom Build Options

```bash
# Specify TinyGo version
TINYGO_VERSION=0.44.0 make build

# Specify target
TINYGO_TARGET=wasm make build

# Specify flags (use with caution)
TINYGO_FLAGS="-opt=0" make build
```

**Note:** The default flags (`-opt=2 -scheduler=none -gc=none`) were removed due to WASM linking errors with TinyGo 0.44.0.

---

## 🧰 Dependencies

### Go Modules

```
module github.com/open-wanderer/wanderer/plugins/garmin

go 1.25.0

require (
    github.com/extism/go-pdk v1.1.3
    github.com/open-wanderer/wanderer/plugins/sdk v0.0.0
)

replace github.com/open-wanderer/wanderer/plugins/sdk => ./sdk
```

### SDK Integration

The plugin uses the **Wanderer Plugin SDK** which is included as a local module via the `replace` directive in `go.mod`.

The SDK provides:
- Type definitions (`ListInput`, `DetailInput`, `TrailSummary`, etc.)
- Host communication functions (`HostRequest`, `InputJSON`, `OutputJSON`)
- Utility functions (`StringField`, `IntState`, `LogWarn`, etc.)

**SDK Location:** `wanderer-garmin-plugin/sdk/`

**To update the SDK:**

```bash
# Clone the latest version from the main repository
cd wanderer-garmin-plugin
git clone --depth 1 https://github.com/open-wanderer/wanderer.git wanderer-sdk-source
cp -r wanderer-sdk-source/plugins/sdk .
rm -rf wanderer-sdk-source

# Test the build
docker run --rm -v $(pwd):/src -w /src tinygo/tinygo:0.44.0 make build
```

---

## 🐛 Bug Fixes Applied

During the initial build process, several code issues were identified and fixed:

### 1. Variable Shadowing in garmin.go

**File:** `garmin.go` (Lines ~337 & ~344)

**Problem:**
```go
var track track  // ❌ Variable name shadows type
if gpxData != nil && len(gpxData) > 0 {
    track = track{  // ❌ Using variable as type
        Format: "gpx",
        ContentBase64: base64.StdEncoding.EncodeToString(gpxData),
    }
}
```

**Solution:**
```go
var trackData track  // ✅ Different variable name
if gpxData != nil && len(gpxData) > 0 {
    trackData = track{  // ✅ Using type correctly
        Format: "gpx",
        ContentBase64: base64.StdEncoding.EncodeToString(gpxData),
    }
}

// ... later in the return statement
trackData,  // ✅ Use the renamed variable
```

### 2. InputJSON Misuse in main.go

**File:** `main.go` (Line ~139)

**Problem:**
```go
func parseInput(input any) error {
    data, err := pdk.InputJSON(input)  // ❌ Returns 1 value, not 2
    if err != nil {
        return err
    }
    if data == nil {
        return json.Unmarshal([]byte("{}"), &input)
    }
    return nil
}
```

**Solution:**
```go
func parseInput(input any) error {
    err := pdk.InputJSON(&input)  // ✅ Correct: takes pointer, returns error
    if err != nil {
        return err
    }
    return nil
}
```

**Note:** Also removed unused `encoding/json` import.

### 3. TinyGo Flags Causing WASM Linker Errors

**File:** `Makefile`

**Problem:**
```makefile
TINYGO_FLAGS ?= -opt=2 -scheduler=none -gc=none
```

This caused:
```
wasm-ld: error: main.lto..o: undefined symbol: runtime.alloc
```

**Solution:**
```makefile
TINYGO_VERSION ?= 0.44.0
TINYGO ?= tinygo
TINYGO_TARGET ?= wasm
TINYGO_FLAGS ?=  # ✅ Empty - no problematic flags
```

**Rationale:** TinyGo 0.44.0 has improved WASM support and doesn't need these specific optimization flags. Removing them resolved the linking errors.

---

## 📊 Testing

### Manual Testing

1. **Build the plugin** (as described above)
2. **Configure the plugin** in Wanderer settings:
   - Add the Garmin plugin
   - Enter your Garmin Personal API Token in the `accessToken` field
   - Set an optional start date (e.g., `2024-01-01`)
   - Enable activities sync
3. **Run synchronization** and verify:
   - Activities appear in the Trails list
   - GPX tracks are imported correctly
   - Metadata (distance, elevation, etc.) is populated
   - Photos are imported with geotags (if available)

### Test Data

For development, you can use:

```bash
# Create a test directory
mkdir -p test/data

# Test with sample Garmin API responses
# (Create test JSON files matching Garmin API structure)
```

### Test Helper Scripts

Several Python scripts are provided in this directory to assist with testing:

| Script | Purpose | Usage |
|--------|---------|-------|
| `build_plugin.py` | **BUILD SCRIPT**: Build plugin in Docker, copy assets, install to data directory | `python3 build_plugin.py` or `python3 build_plugin.py --clean --verbose` |
| `clear_trails.py` | **Delete all trails** from database for clean testing | `python3 clear_trails.py --dry-run` or `python3 clear_trails.py` |
| `check_auth_fields.py` | Check if API key/access token is saved correctly | `python3 check_auth_fields.py` |
| `check_current_error.py` | View the last error from the plugin instance | `python3 check_current_error.py` |
| `check_full_instance.py` | Display full plugin instance details | `python3 check_full_instance.py` |
| `check_plugin_status.py` | **NEW**: Full plugin status, errors, auth config, state | `python3 check_plugin_status.py` |
| `check_wasm.py` | **NEW**: Verify WASM has no gojs dependency | `python3 check_wasm.py` |
| `verify_build.py` | **NEW**: Compare build config with Strava plugin | `python3 verify_build.py` |
| `decrypt_test.py` | Test PocketBase encryption/decryption | `python3 decrypt_test.py` |
| `trigger_sync.py` | Manually trigger plugin sync | `python3 trigger_sync.py` |
| `verify_key.py` | Verify the saved API key matches expected | `python3 verify_key.py` |
| `check_all_instances.py` | List all plugin instances | `python3 check_all_instances.py` |
| `check_exact_key.py` | Check if exact API key is saved | `python3 check_exact_key.py` |
| `test_save_process.py` | Test the save process | `python3 test_save_process.py` |

#### build_plugin.py

**Build and install the plugin with a single command:**

```bash
# Build and install
python3 build_plugin.py

# Build with clean (remove old artifacts first)
python3 build_plugin.py --clean

# Build but don't install
python3 build_plugin.py --no-install

# Build with detailed output
python3 build_plugin.py --verbose

# Clean, build, install with verbose output
python3 build_plugin.py --clean --verbose

# Help
python3 build_plugin.py --help
```

**What it does:**
1. Builds plugin.wasm inside a Docker container with TinyGo
2. Copies plugin.json and icon.svg to dist directory
3. Installs all files to `/opt/wanderer-dev/data/plugins/garmin/`
4. Verifies the build (file size, no gojs dependency)

**Docker Image:** Uses `tinygo/tinygo:0.44.0` by default

**Build Command:**
```
tinygo build -target=wasi -scheduler=none -no-debug -o /src/dist/garmin/plugin.wasm .
```

This is the **recommended** way to build the plugin as it:
- ✅ Uses Docker (no host Go required)
- ✅ Handles all build steps automatically
- ✅ Verifies the build was successful
- ✅ Can be used in CI/CD pipelines

#### clear_trails.py

**Most useful for testing**: Clears all trails and related data from the database.

```bash
# Dry run (shows what would be deleted)
python3 clear_trails.py --dry-run

# Actually delete all trails
python3 clear_trails.py

# With custom database path
python3 clear_trails.py --db-path /path/to/your/data.db

# Quiet mode (only errors)
python3 clear_trails.py --quiet

# Help
python3 clear_trails.py --help
```

**What it deletes:**
- All records from `trails` table
- All records from `trail_external_reference` table
- All records from `trail_link_share` table
- All records from `trail_like` table
- All records from `trail_share` table

**Important:** This operation is **NOT reversible**. Make sure you have a backup if you need the data.

#### check_plugin_status.py

**Check plugin instance status, errors, and configuration:**

```bash
# Check Garmin plugin status
python3 check_plugin_status.py

# With custom database path
python3 check_plugin_status.py --db-path /path/to/data.db

# Quiet mode (summary only)
python3 check_plugin_status.py --quiet
```

**Shows:**
- Plugin instance ID and status
- Whether plugin is enabled
- Last error details (code, message)
- Last sync timestamp
- Auth configuration (encrypted/decrypted status)
- Plugin state (page, etc.)

#### check_wasm.py

**Verify WASM file has no Go runtime dependencies:**

```bash
# Check installed plugin.wasm
python3 check_wasm.py

# Check a specific file
python3 check_wasm.py --file /path/to/plugin.wasm

# Quiet mode
python3 check_wasm.py --quiet
```

**What it checks:**
- File size (WASI builds are typically < 1MB, wasm builds are > 2MB)
- Presence of `gojs` imports (should be 0 for WASI builds)
- Provides fix recommendations if issues found

#### verify_build.py

**Compare build configuration with reference plugin (Strava):**

```bash
# Compare with Strava plugin (default)
python3 verify_build.py

# Compare with custom reference
python3 verify_build.py --reference /path/to/plugin.json

# Quiet mode
python3 verify_build.py --quiet
```

**Checks:**
- Auth context types (bearer, oauth2, session)
- Secret field configuration (`secretField` vs `secretFields`)
- Recommends correct build flags

### Debugging

**View logs:**
```bash
docker logs wanderer-db
```

**Enter the container:**
```bash
docker exec -it wanderer-db /bin/bash
```

**Test connectivity:**
```bash
# Test if the plugin is loaded
curl http://localhost:8091/api/plugins

# Test Garmin API connectivity (requires valid token)
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "https://personal.connectapi.garmin.com/activitylist-service/activities/search/activities?start=0&limit=10"
```

---

## 🔄 Update Process

### Updating TinyGo Version

1. Check the latest TinyGo version:
   ```bash
   curl -s https://api.github.com/repos/tinygo-org/tinygo/releases/latest | grep tag_name
   ```

2. Update the Makefile:
   ```makefile
   TINYGO_VERSION ?= <new-version>
   ```

3. Test the build:
   ```bash
   docker run --rm -v $(pwd):/src -w /src tinygo/tinygo:<new-version> make build
   ```

### Updating the SDK

See [SDK Integration](#sdk-integration) above.

### Updating Dependencies

```bash
# Using the TinyGo container
docker run --rm -v $(pwd):/src -w /src tinygo/tinygo:0.44.0 go mod tidy
```

---

## 🔍 Troubleshooting & Debugging

### Checking Saved API Key in Database

If you see a "Setup-Fehler" (Setup Error) and need to verify that your API key was saved to the database:

**Method 1: Using Python (Recommended)**

```bash
# Navigate to the PocketBase data directory
cd /opt/wanderer-dev/data/pb_data

# Query the database with Python's built-in sqlite3 module
python3 << 'EOF'
import sqlite3
import json

conn = sqlite3.connect('data.db')
cursor = conn.cursor()

# Get the garmin plugin instance
cursor.execute('SELECT id, plugin_id, auth FROM plugin_instances WHERE plugin_id = "garmin"')
instance = cursor.fetchone()

if instance:
    auth_json = instance[2]
    try:
        auth = json.loads(auth_json)
        print("✅ Garmin Plugin Instance Found")
        print(f"   ID: {instance[0]}")
        print(f"   Plugin ID: {instance[1]}")
        print(f"   Auth Config: {json.dumps(auth, indent=2)}")
        if 'apiKey' in auth:
            print(f"\n✅ API Key is saved (length: {len(auth['apiKey'])})")
            print(f"   Key starts with: {auth['apiKey'][:20]}...")
        elif 'accessToken' in auth:
            print(f"\n✅ Access Token is saved (length: {len(auth['accessToken'])})")
        else:
            print("\n❌ No API key or access token found in auth!")
    except json.JSONDecodeError:
        print(f"❌ Failed to parse auth JSON: {auth_json}")
else:
    print("❌ No garmin plugin instance found in database")

conn.close()
EOF
```

**Method 2: Using Docker (if python3 not available)**

```bash
# Use the wanderer-db container (if it has python3)
docker exec wanderer-dev-db-1 python3 -c "
import sqlite3, json
conn = sqlite3.connect('/pb_data/data.db')
cursor = conn.cursor()
cursor.execute('SELECT auth FROM plugin_instances WHERE plugin_id = \"garmin\"')
result = cursor.fetchone()
if result:
    auth = json.loads(result[0]) if result[0] else {}
    if 'apiKey' in auth:
        print('API Key found, length:', len(auth['apiKey']))
    elif 'accessToken' in auth:
        print('Access Token found, length:', len(auth['accessToken']))
    else:
        print('No auth credentials found')
else:
    print('No garmin instance found')
conn.close()
"
```

### Common Issues

#### Issue 1: "Setup-Fehler" with Token Exchange Error

**Symptom:** Red "Setup-Fehler" in plugin overview
**Error:** `token exchange request failed: garmin POST /personalapitoken-service/personal-api/personal-api-key/token failed (404)`

**Causes:**
1. Invalid or incorrect Garmin Personal API Key
2. Mismatch between saved field name (`apiKey`) and manifest expectation (`accessToken`)

**Solution:**
- Verify your API key is valid using the database lookup above
- Ensure the plugin.json `secretField` matches what you saved (should be `"apiKey"`)
- Generate a new API key from [Garmin Developer Portal](https://developer.garmin.com/gc-developer-program/personal/personal-api/)

#### Issue 2: API Key Not Being Saved

**Symptom:** Database query shows empty or missing auth field

**Solution:**
- Re-enter the API key in Wanderer's plugin configuration
- Check browser console for JavaScript errors during save
- Ensure you have write permissions to the plugin instance

#### Issue 3: Plugin Status Stuck on "Needs Auth"

**Symptom:** Plugin instance status is `needs_auth` instead of `configured`

**Solution:**
```bash
# Check instance status
python3 << 'EOF'
import sqlite3
conn = sqlite3.connect('/opt/wanderer-dev/data/pb_data/data.db')
cursor = conn.cursor()
cursor.execute('SELECT id, status, auth FROM plugin_instances WHERE plugin_id = "garmin"')
for row in cursor.fetchall():
    print(f"ID: {row[0]}, Status: {row[1]}, Auth: {row[2][:100] if row[2] else 'empty'}...")
conn.close()
EOF
```

If auth exists but status is `needs_auth`, the plugin manifest may have the wrong auth context.

### Viewing Full Plugin Instance Details

```bash
python3 << 'EOF'
import sqlite3
import json

conn = sqlite3.connect('/opt/wanderer-dev/data/pb_data/data.db')
cursor = conn.cursor()

# Get all columns
cursor.execute('PRAGMA table_info(plugin_instances)')
cols = [c[1] for c in cursor.fetchall()]

# Get the garmin instance
cursor.execute('SELECT * FROM plugin_instances WHERE plugin_id = "garmin"')
instance = cursor.fetchone()

if instance:
    print("=== GARMIN PLUGIN INSTANCE DETAILS ===")
    for i, col in enumerate(cols):
        val = instance[i]
        if col == 'auth' or col == 'config' or col == 'state':
            try:
                parsed = json.loads(val) if val else {}
                print(f"\n{col}:")
                print(json.dumps(parsed, indent=2))
            except:
                print(f"\n{col}: {val[:200] if val else 'empty'}...")
        elif col == 'last_error':
            try:
                parsed = json.loads(val) if val else {}
                print(f"\n{col}:")
                print(json.dumps(parsed, indent=2))
            except:
                print(f"\n{col}: {val}")
        else:
            print(f"{col}: {val}")
else:
    print("No garmin plugin instance found")

conn.close()
EOF
```

### Checking Sync Errors

```bash
python3 << 'EOF'
import sqlite3
import json

conn = sqlite3.connect('/opt/wanderer-dev/data/pb_data/data.db')
cursor = conn.cursor()

cursor.execute('SELECT last_error FROM plugin_instances WHERE plugin_id = "garmin"')
error = cursor.fetchone()

if error and error[0]:
    try:
        error_data = json.loads(error[0])
        print("=== LAST ERROR ===")
        print(json.dumps(error_data, indent=2))
    except:
        print(f"Raw error: {error[0]}")
else:
    print("✅ No errors recorded")

conn.close()
EOF
```

---

## 📚 Resources

### Documentation

- [Wanderer Documentation](https://github.com/open-wanderer/wanderer)
- [Garmin Personal API](https://developer.garmin.com/gc-developer-program/personal/personal-api/)
- [TinyGo Documentation](https://tinygo.org/docs/)
- [Extism PDK (Go)](https://github.com/extism/go-pdk)

### API Endpoints Used

| Endpoint | Purpose |
|----------|---------|
| `/activitylist-service/activities/search/activities` | List activities |
| `/activity-service/activity/{id}` | Get activity details |
| `/activity-service/activity/{id}/workouts` | Get workouts |
| `/download-service/files/activity/{id}` | Download activity ZIP (contains GPX) |

**Base URL:** `https://personal.connectapi.garmin.com`

**Authentication:** Bearer Token (Personal API Key)

---

## 🎯 Architecture

### Plugin Flow

```
Wanderer → Plugin Manager → Garmin Plugin (WASM)
                                    ↓
                              garminClient
                                    ↓
                  [Garmin Personal API]
                                    ↓
              Activities → Mapping → Trail Import
                                    ↓
                          Wanderer Database
```

### Main Components

1. **Exported Functions** (main.go):
   - `listRoutesV1()` - List routes/courses
   - `listActivitiesV1()` - List activities
   - `getRouteDetailV1()` - Get route details
   - `getActivityDetailV1()` - Get activity details

2. **Garmin Client** (garmin.go):
   - `newGarminClient()` - Create authenticated client
   - `searchActivities()` - Search for activities
   - `getActivity()` - Get activity details
   - `downloadActivity()` - Download activity ZIP
   - `syncActivities()` - Sync activities to Wanderer format
   - `syncRoutes()` - Sync routes to Wanderer format

3. **Type Definitions** (types.go):
   - SDK type aliases
   - Garmin API response types
   - Track point structures

---

## 💡 Tips & Best Practices

### 1. Always Test with Real Data

Garmin API responses can vary. Always test with your own Garmin account to catch edge cases.

### 2. Handle Rate Limits

Garmin has rate limits. Implement:
- Exponential backoff on rate limit errors
- Configurable delay between requests
- Batch processing for large datasets

### 3. Error Handling

- Log errors with context (activity ID, user ID, etc.)
- Distinguish between:
  - Authentication errors (invalid token)
  - Rate limit errors (429)
  - API errors (5xx)
  - Network errors

### 4. Performance

- Use pagination for large datasets
- Cache responses when possible
- Minimize API calls (fetch only needed data)

### 5. Code Style

- Follow Go conventions
- Use meaningful variable names
- Add comments for complex logic
- Keep functions focused and small

---

## 🙏 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/your-feature`)
3. Commit your changes (`git commit -m 'Add some feature'`)
4. Push to the branch (`git push origin feature/your-feature`)
5. Open a Pull Request

---

## 📜 License

This plugin is part of [Wanderer](https://github.com/open-wanderer/wanderer) and is licensed under the same terms.
