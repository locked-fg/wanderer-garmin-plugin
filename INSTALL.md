# Garmin Plugin Installation

Requires Docker and a running Wanderer instance. TinyGo runs in Docker; no host Go tools are needed.

From the project directory, build and install the plugin, then restart Wanderer’s database service:

```bash
cd /opt/wanderer-dev/wanderer-garmin-plugin
docker run --rm -v "$PWD:/src" -w /src tinygo/tinygo:0.44.0 make clean build
rm -rf ../data/plugins/garmin
mkdir -p ../data/plugins/garmin
cp dist/garmin/plugin.wasm plugin.json assets/icon.svg ../data/plugins/garmin/
cd ..
docker compose restart db
```

In Wanderer, add the Garmin plugin and enter your [Garmin Personal API token](https://connect.garmin.com/app/settings/personalApiToken) in the `apiKey` field. Then save and sync.
