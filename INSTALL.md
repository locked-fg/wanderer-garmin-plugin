# Garmin Plugin Installation

Requires Docker and a running Wanderer instance. TinyGo runs in Docker; no host Go tools are needed.

Build and install with the project’s Makefile, then restart Wanderer’s database service:

```bash
cd /opt/wanderer-dev
docker run --rm -v "$PWD:/workspace" -w /workspace/wanderer-garmin-plugin tinygo/tinygo:0.44.0 make clean install-local
docker compose restart db
```

In Wanderer, add the Garmin plugin and enter your [Garmin Personal API token](https://connect.garmin.com/app/settings/personalApiToken) in the `apiKey` field. Then save and sync.
