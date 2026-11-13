# Garmin Plugin - Build & Development Notes

## Build Instructions

**No Go installation required on host!** All compilation happens inside Docker containers.

### Quick Build (Recommended)

```bash
# Build using TinyGo Docker container
docker run --rm -v $(pwd):/src -w /src tinygo/tinygo:0.44.0 make build

# Install to local Wanderer
mkdir -p ../../../data/plugins/garmin
cp dist/garmin/plugin.wasm ../../../data/plugins/garmin/
cp dist/garmin/plugin.json ../../../data/plugins/garmin/
cp assets/icon.svg ../../../data/plugins/garmin/

# Restart Wanderer
docker compose restart wanderer-db
```

### Full Build from Source

```bash
cd /opt/wanderer-dev/wanderer-garmin-plugin
docker run --rm -v $(pwd):/src -w /src tinygo/tinygo:0.44.0 make clean build install
```

### Install Only

```bash
cd /opt/wanderer-dev/wanderer-garmin-plugin
make install-local
```

### Makefile Targets

| Target | Description |
|--------|-------------|
| `make build` | Build plugin.wasm |
| `make install` | Build and copy to dist/garmin/ |
| `make install-local` | Build and install to ../../../data/plugins/garmin/ |
| `make clean` | Remove build artifacts |

---

## Critical Fix Applied

### Issue 1: API Key Not Recognized (Setup-Fehler)

**Problem:** The plugin.json was configured to look for `accessToken` in the auth map, but users save the credential as `apiKey` in the Wanderer UI.

**Fix:** Changed plugin.json auth context from `secretField: "accessToken"` to `secretField: "apiKey"`.

**File:** `plugin.json`
```json
"auth": {
  "contexts": {
    "bearer_token": {
      "type": "bearer",
      "secretField": "apiKey"  // Changed from "accessToken"
    }
  }
}
```

### Issue 2: Wrong API Key in Database

**Problem:** The database contained an incorrect/old API key, causing Garmin to return:
```
"Personal API Key with ID aeeb944f-774e-46c1-9004-9b25180cdba5 does not exist."
```

**Fix:** Updated the API key in the database for plugin instance `z3ahsdr1i0pfty9` to the correct value provided by the user:
```
GCDC_62a8961e-0f89-448a-8780-8b89c2738095.c80cd6d4-0cd3-4ee9-bbe4-1759aa037b55.wOmMuxj/lqjrZwo9RI8fw/PWKuZyAPKLk5NFWTmiziQ.1612962821
```

---

## Troubleshooting

### Check if API Key is Saved in Database

Use the provided `verify_key.py` script:

```bash
python3 verify_key.py
```

Or use this one-liner:

```bash
python3 -c "import sqlite3,json;conn=sqlite3.connect('/opt/wanderer-dev/data/pb_data/data.db');c=conn.cursor();c.execute('SELECT auth FROM plugin_instances WHERE plugin_id=\"garmin\"');r=c.fetchone();print('✅ Key exists' if r and 'apiKey' in json.loads(r[0]) else '❌ No key')"
```

### View Full Plugin Instance Details

```bash
python3 << 'EOF'
import sqlite3, json
conn = sqlite3.connect('/opt/wanderer-dev/data/pb_data/data.db')
cursor = conn.cursor()
cursor.execute('SELECT id, status, enabled, auth, last_error FROM plugin_instances WHERE plugin_id = "garmin"')
instance = cursor.fetchone()
if instance:
    print(f"ID: {instance[0]}, Status: {instance[1]}, Enabled: {instance[2]}")
    auth = json.loads(instance[3]) if instance[3] else {}
    print(f"Auth: {json.dumps(auth, indent=2)}")
    if instance[4]:
        error = json.loads(instance[4])
        print(f"Last Error: {json.dumps(error, indent=2)}")
else:
    print("No garmin instance found")
conn.close()
EOF
```

### Verify plugin.json Configuration

```bash
grep -A 2 '"bearer_token"' /opt/wanderer-dev/data/plugins/garmin/plugin.json
```
Should show: `"secretField": "apiKey"`

### Update API Key in Database

If you need to update the Garmin API key manually:

```bash
python3 << 'EOF'
import sqlite3, json

conn = sqlite3.connect('/opt/wanderer-dev/data/pb_data/data.db')
cursor = conn.cursor()

new_key = "YOUR_NEW_API_KEY_HERE"
cursor.execute('UPDATE plugin_instances SET auth = ?, last_error = NULL WHERE plugin_id = "garmin"',
               (json.dumps({"apiKey": new_key}),))
conn.commit()
conn.close()
print("✅ API key updated")
EOF
```

### Restart Wanderer After Changes

```bash
cd /opt/wanderer-dev
docker compose restart db
```

---

## Code Fixes Applied During Initial Build

1. **Variable Shadowing in garmin.go**
   - Changed `var track track` to `var trackData track`

2. **InputJSON Misuse in main.go**
   - Changed `pdk.InputJSON(input)` to `pdk.InputJSON(&input)`

3. **TinyGo Version & Flags**
   - Updated from TinyGo 0.32.0 to 0.44.0
   - Removed problematic flags: `-opt=2 -scheduler=none -gc=none`

4. **SDK Integration**
   - Added local SDK copy at `sdk/`
   - Updated go.mod replace directive: `=> ./sdk`

---

## Project Conventions

- **No CHANGES.md**: Keep change history in git commits only
- Use descriptive commit messages that explain the "what" and "why"

---

## Requirements

- **TinyGo**: 0.44.0 (includes Go 1.27.3)
- **Docker**: For containerized builds
- **Wanderer**: Running instance for testing
- **Garmin Personal API Key**: From https://connect.garmin.com/app/settings/personalApiToken

---

## Useful Commands

### Build & Install in One Command
```bash
cd /opt/wanderer-dev/wanderer-garmin-plugin && \
docker run --rm -v $(pwd):/src -w /src tinygo/tinygo:0.44.0 make build && \
mkdir -p /opt/wanderer-dev/data/plugins/garmin && \
cp dist/garmin/plugin.wasm /opt/wanderer-dev/data/plugins/garmin/ && \
cp plugin.json /opt/wanderer-dev/data/plugins/garmin/ && \
cp assets/icon.svg /opt/wanderer-dev/data/plugins/garmin/ && \
cd /opt/wanderer-dev && docker compose restart db
```

### Test Garmin API Key
```bash
curl -H "Authorization: Bearer YOUR_API_KEY" \
  "https://personal.connectapi.garmin.com/activitylist-service/activities/search/activities?start=0&limit=1"
```
