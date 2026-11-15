# Install

From the plugin repository root, build and install into the `data/` directory beside the Wanderer checkout:

```bash
docker run --rm \
  -v "$PWD/..:/workspace" \
  -w "/workspace/$(basename "$PWD")" \
  tinygo/tinygo:0.44.0 make clean install-local
```

This installs `plugin.wasm`, `plugin.json`, and `icon.svg` under the sibling Wanderer checkout's `data/plugins/garmin/`. Then restart Wanderer's `db` service from the Wanderer repository root:

```bash
cd ..
docker compose restart db
```

In Wanderer, add Garmin under **Settings → Plugins**, enter a Garmin Personal API token in `apiKey`, optionally set `after`, then save and sync.

To build without installing, run the same Docker command with `make clean build`. See [DEVELOPMENT.md](DEVELOPMENT.md) for tests and build details.
