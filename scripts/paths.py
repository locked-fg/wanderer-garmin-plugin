"""Shared path defaults for Garmin development helpers."""

import os
from pathlib import Path


PLUGIN_ROOT = Path(__file__).resolve().parent.parent
WANDERER_DATA_DIR = Path(
    os.environ.get("WANDERER_DATA_DIR", PLUGIN_ROOT.parent / "data")
).expanduser().resolve()


def database_path():
    return WANDERER_DATA_DIR / "pb_data" / "data.db"


def installed_plugin_dir():
    return WANDERER_DATA_DIR / "plugins" / "garmin"
