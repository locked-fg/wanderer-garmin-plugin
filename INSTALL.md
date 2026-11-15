# Install

From the Wanderer repository root, build and copy the plugin into its local data directory:

```bash
cd /opt/wanderer-dev
docker run --rm -v "$PWD:/workspace" -w /workspace/wanderer-garmin-plugin tinygo/tinygo:0.44.0 make clean install-local
docker compose restart db
```

The command installs `plugin.wasm`, `plugin.json`, and `icon.svg` to `data/plugins/garmin/`. In Wanderer, add Garmin under **Settings → Plugins**, enter a Garmin Personal API token in `apiKey`, optionally set `after`, then save and sync.

To build without installing, run the same Docker command with `make clean build`. See [DEVELOPMENT.md](DEVELOPMENT.md) for tests and build details.
