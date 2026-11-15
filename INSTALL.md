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

## Production Docker Compose

The backend must have the plugin directory mounted into its container. In the production Compose file, add this bind mount to the Wanderer backend service (usually `db`):

```yaml
services:
  db:
    volumes:
      - ./data/plugins:/data/plugins
```

Keep any existing `volumes` entries and add this line alongside them. The host path `./data/plugins` is relative to the Compose file's directory. Put the plugin bundle directly in `data/plugins/garmin/`; that directory must contain `plugin.json` and `plugin.wasm` (and may contain `icon.svg`). Do not copy the repository source directory or add another nested `garmin/` directory.

After copying or updating the bundle, restart the backend so it discovers the plugin:

```bash
docker compose restart db
```

Then reload **Settings → Plugins**. If Garmin is not listed, check that the mount points to the expected host directory and that `/data/plugins/garmin/plugin.json` and `/data/plugins/garmin/plugin.wasm` exist inside the backend container.
