#!/usr/bin/env python3
"""
Trigger a plugin sync by updating the plugin instance timestamp.
This should cause PocketBase to re-sync.
"""
import sqlite3
import json
from datetime import datetime

conn = sqlite3.connect('../data/pb_data/data.db')
cursor = conn.cursor()

# Clear any errors and update the timestamp
now = datetime.utcnow().isoformat() + 'Z'

cursor.execute(
    'UPDATE plugin_instances SET last_error = NULL, status = "syncing", updated = ? WHERE plugin_id = "garmin"',
    (now,)
)

conn.commit()

# Verify
cursor.execute('SELECT id, status, updated FROM plugin_instances WHERE plugin_id = "garmin"')
row = cursor.fetchone()

if row:
    print("=" * 70)
    print("SYNC TRIGGERED:")
    print("=" * 70)
    print(f"Instance ID: {row[0]}")
    print(f"Status: {row[1]}")
    print(f"Updated: {row[2]}")
    print("\n✅ Plugin instance updated - PocketBase should pick up the change")
    print("✅ The plugin should now attempt to sync with the new API key")
else:
    print("❌ No Garmin instance found")

conn.close()
