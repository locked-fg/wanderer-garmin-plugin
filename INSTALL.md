# Garmin Plugin - Installation Guide (Docker-only)

**All compilation happens in Docker containers - nothing is installed locally!**

---

## Requirements

- Docker
- Running Wanderer instance
- [Garmin Personal API Token](https://connect.garmin.com/app/settings/personalApiToken)

---

## Build & Install

```bash
# 1. Navigate to plugin directory
cd /opt/wanderer-dev/wanderer-garmin-plugin

# 2. Build plugin in Docker container
docker run --rm -v $(pwd):/src -w /src tinygo/tinygo:0.44.0 make clean build

# 3. Install plugin files
mkdir -p /opt/wanderer-dev/data/plugins/garmin
cp dist/garmin/plugin.wasm /opt/wanderer-dev/data/plugins/garmin/
cp plugin.json /opt/wanderer-dev/data/plugins/garmin/
cp assets/icon.svg /opt/wanderer-dev/data/plugins/garmin/

# 4. Restart Wanderer
docker compose restart db
```

**One-liner:**
```bash
cd /opt/wanderer-dev/wanderer-garmin-plugin && \
docker run --rm -v $(pwd):/src -w /src tinygo/tinygo:0.44.0 make clean build && \
mkdir -p /opt/wanderer-dev/data/plugins/garmin && \
cp dist/garmin/plugin.wasm /opt/wanderer-dev/data/plugins/garmin/ && \
cp plugin.json /opt/wanderer-dev/data/plugins/garmin/ && \
cp assets/icon.svg /opt/wanderer-dev/data/plugins/garmin/ && \
docker compose restart db
```

---

## Configure in Wanderer

1. Open Wanderer: `http://localhost:8080`
2. Go to: Settings -> Plugins
3. Find: "Garmin" and click + Add
4. Configure:
   - **apiKey**: Your Garmin Personal API Token (from Step 1)
   - Start date: (Optional, e.g., 2024-01-01)
   - Enable: Activities (required), Routes (optional)
5. Save and Activate

---

## Sync

- Click "Sync Now" or wait for automatic sync
- Activities will appear in "Trails" after 1-5 minutes

---

## Notes

- **Docker-only**: TinyGo runs in `tinygo/tinygo:0.44.0` container
- **No local installation**: Go, TinyGo, Make not needed on host
- **Auth field**: Use `apiKey` (not `accessToken`)
