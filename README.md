# Garmin Plugin for Wanderer

Imports your Garmin activities and routes into Wanderer.

This plugin connects [Wanderer](https://wanderer.to) with the [Garmin Personal API](https://developer.garmin.com/gc-developer-program/personal/personal-api/) to automatically synchronize your activities, routes, and associated data (GPX tracks, photos, metadata).

---

## Features

- Import Garmin activities as Wanderer trails
- Full GPX track data with elevation, timestamps, and coordinates
- Activity photos with geotags
- Metadata: distance, duration, elevation gain/loss, calories, average/max speed
- Automatic category mapping (Running, Cycling, Hiking, etc.)
- Privacy settings preserved
- Difficulty ratings
- Start date filter for selective imports

---

## Requirements

- Running Wanderer instance in Docker
- Docker installed on your system
- Garmin Connect account with at least one activity
- [Garmin Personal API Token](https://connect.garmin.com/app/settings/personalApiToken)

---

## Build Requirements

To build this plugin, you need:

- **TinyGo** (version 0.44.0 or later) - for WASM compilation
- **Make** - for build automation
- **Go** (optional, only needed if building outside TinyGo container)

**Note:** None of these tools need to be installed on your host system. All compilation happens inside Docker containers.

---

## Installation

### Docker-only Build (Recommended)

This is the only supported method. All compilation happens in Docker containers.

#### Build the plugin:

```bash
# Navigate to plugin directory
cd /opt/wanderer-dev/wanderer-garmin-plugin

# Build using TinyGo Docker container
docker run --rm -v $(pwd):/src -w /src tinygo/tinygo:0.44.0 make clean build
```

#### Install the plugin:

```bash
# Create plugin directory
mkdir -p /opt/wanderer-dev/data/plugins/garmin

# Copy plugin files
cp dist/garmin/plugin.wasm /opt/wanderer-dev/data/plugins/garmin/
cp plugin.json /opt/wanderer-dev/data/plugins/garmin/
cp assets/icon.svg /opt/wanderer-dev/data/plugins/garmin/

# Restart Wanderer
docker compose restart db
```

#### One-liner:

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

## Configuration

1. Open Wanderer: `http://localhost:8080`
2. Go to: Settings -> Plugins
3. Find: "Garmin" and click + Add
4. Configure:
   - **apiKey**: Your Garmin Personal API Token (from Garmin Connect settings)
   - **Start date (`after`)**: Optional earliest activity date, e.g. `2024-01-01`. Leave it unset to import all activities Garmin makes available.
   - Enable: Activities (required), Routes (optional)
5. Save and Activate

### Importing older activities

The start date limits which activities are eligible for import; it does not set a maximum number of activities. The plugin requests activities in pages of 50 using Garmin's `start` and `limit` parameters, then requests subsequent pages while the previous page is full. It finishes when Garmin returns a page with fewer than 50 activities. For example, a sync that receives 50 activities followed by 7 continues through both pages.

Set `after` to the earliest date you want included, or leave it unset for the full history available through Garmin's API. Allow the sync to finish; a failed API request can interrupt a run, so check the plugin sync status and logs if the import appears incomplete.

---

## Usage

- Click "Sync Now" in the plugin settings or wait for automatic sync
- Activities will appear in "Trails" after 1-5 minutes
- The plugin automatically exchanges your API key for an access token
- Token activation may take 10-30 seconds (handled automatically with retry)

---

## Authentication

The Garmin Personal API uses two-step authentication:

1. **API Key Exchange**: Your Personal API Key is exchanged for an Access Token
2. **API Calls**: The Access Token is used for all Garmin API requests

The plugin handles this automatically. You only need to provide your Personal API Key in the Wanderer UI.

**Important:** Access tokens may need 10-30 seconds to activate after exchange. The plugin includes automatic retry logic to handle this.

---

## Supported Activity Types

The plugin maps Garmin activity types to Wanderer categories:

- **Running**: Running, Trail Running
- **Biking**: Cycling, Mountain Biking, Gravel Cycling, E-Bike
- **Hiking**: Hiking, Walking
- **Skiing**: Skiing, Snowboarding, Cross-Country Skiing, Alpine Skiing, Backcountry Skiing
- **Canoeing**: Kayaking, Rowing, Paddle Boarding
- **Climbing**: Climbing, Rock Climbing, Ice Climbing
- **Other**: Swimming, Other

---

## Troubleshooting

### Token Activation Issues

If you see "Token is not active" errors:
- This is normal for Garmin Personal API tokens
- The plugin automatically waits 15 seconds and retries
- If it still fails, wait 30 seconds and try again

### Build Failures

If the build fails:
- Ensure Docker is running
- Check that you have internet access for downloading the TinyGo image
- Verify the TinyGo image tag (0.44.0) is available

### Plugin Not Loading

- Check that all files are in the correct directory: `/opt/wanderer-dev/data/plugins/garmin/`
- Verify the plugin.json, plugin.wasm, and icon.svg files exist
- Restart the Wanderer container

---

## Technical Details

### Plugin Structure

- `main.go`: WASM export functions (entry points)
- `garmin.go`: Garmin API client implementation
- `types.go`: Type definitions for API responses
- `plugin.json`: Plugin manifest with capabilities and auth
- `Makefile`: Build automation
- `assets/icon.svg`: Plugin icon

### API Endpoints

- Token Exchange: `POST /personalapitoken-service/personal-api/personal-api-key/token`
- Activities Search: `GET /activitylist-service/activities/search/activities`
- Activity Details: `GET /activity-service/activity/{id}`
- Activity Download: `GET /download-service/files/activity/{id}`
- User Profile: `GET /userprofile-service/userprofile`

### Observed API Response Format

A direct request to the Personal API token endpoint returned HTTP `200` with the bearer token as `text/plain`. The activities search endpoint returned HTTP `200` with a top-level JSON array, including when requesting one item. Relevant fields in an activity object:

| Field | JSON type | Notes |
|---|---|---|
| `activityId` | integer | Activity identifier |
| `activityType` | object | Contains `typeId`, `typeKey`, and `parentTypeId` |
| `startTimeGMT` | string | UTC start time; this is the field used for date filtering |
| `startTimeLocal` | string | Local start time |
| `calories` | number | May be fractional (for example, JSON numeric values are not always integers) |
| `privacy` | object | Contains `typeId` and `typeKey` |
| `distance`, `duration`, `elevationGain`, `elevationLoss`, `averageSpeed`, `maxSpeed` | number | Activity metrics |
| `startLatitude`, `startLongitude`, `endLatitude`, `endLongitude` | number | Coordinates |

The list response contains additional fields. Do not assume `activityType` or `privacy` are strings, or that the activity list is wrapped in an object.

`startTimeGMT` is a string and may arrive as `YYYY-MM-DD HH:MM:SS` (for example, `2026-07-20 16:07:41`), which is not RFC3339. The plugin parses this GMT value and emits `startedAt` as RFC3339 UTC for Wanderer's host, such as `2026-10-12T16:09:45Z`. Appending `T00:00:00Z` to this timestamp makes it invalid and causes the host to reject the detail response.

For the tested account, `GET /activity-service/activity/{id}/details` returned a JSON object with `activityDetailMetrics`, `metricDescriptors`, and `geoPolylineDTO`. The Garmin Python client also exposes the activity download endpoint; its response is a ZIP archive containing a `.fit` file for the tested activity, not a GPX file. The importer therefore gets track geometry from the activity details response instead of treating the downloaded ZIP as GPX.

Garmin's [official Python client](https://github.com/garmin/connect-personal-api-python) currently wraps activity summary, workouts, and ZIP download. The details endpoint can be called through its generic `get()` method, but it has no dedicated wrapper.

### Authentication Flow

```
User provides API Key
    ↓
Plugin exchanges API Key for Access Token (POST with ApiKey header)
    ↓
Wait 15s for token activation
    ↓
Plugin uses Access Token for API calls (GET with Bearer header)
    ↓
Data imported into Wanderer
```

---

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

---

## License

This plugin is open source and available under the same license as Wanderer.

---

## Plugin Comparison

| Feature | Garmin | Strava | Komoot | Hammerhead |
|---------|--------|--------|--------|------------|
| Activities | ✅ | ✅ | ✅ | ✅ |
| Routes | ✅ | ✅ | ✅ | ✅ |
| Photos | ✅ | ✅ | ✅ | ❌ |
| GPX Tracks | ✅ | ✅ | ✅ | ✅ |
| Auto Sync | ✅ | ✅ | ✅ | ✅ |
| OAuth | ❌ | ✅ | ✅ | ✅ |
| API Key | ✅ | ❌ | ❌ | ❌ |

---

## Notes

- This plugin uses the Garmin Personal API, not the Garmin Connect API
- Personal API tokens are long-lived (typically 2 years)
- No OAuth flow required - simple API key authentication
- All sensitive data is encrypted at rest by Wanderer/PocketBase
- Docker-only deployment ensures no host system dependencies
