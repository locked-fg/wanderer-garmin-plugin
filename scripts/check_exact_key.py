#!/usr/bin/env python3
"""
Check whether a Garmin API key is present and inspect its format without printing it.
"""
import sqlite3
import json
from paths import database_path

conn = sqlite3.connect(database_path())
cursor = conn.cursor()

cursor.execute(
    'SELECT id, plugin_id, auth FROM plugin_instances WHERE plugin_id = "garmin"'
)
row = cursor.fetchone()

if row:
    instance_id, plugin_id, auth = row
    auth_data = json.loads(auth) if auth else {}
    api_key = auth_data.get('apiKey', 'NOT FOUND')
    
    print("=" * 60)
    print("CURRENT DATABASE STATE:")
    print("=" * 60)
    print(f"Instance ID: {instance_id}")
    print(f"Plugin ID: {plugin_id}")
    print(f"Auth field present: {'apiKey' in auth_data}")
    print(f"API Key length: {len(api_key)}")
    
    # Check if it's the test key or real key
    if api_key == "1234":
        print("\n⚠️  WARNING: Key is '1234' (test key)")
    elif api_key.startswith("GCDC_"):
        print("\n✅ Key starts with GCDC_ (real key format)")
    elif api_key.startswith("YPJC"):
        print("\n⚠️  Key is encrypted/hashed (PocketBase encryption)")
    else:
        print(f"\n⚠️  Key format unknown")
else:
    print("❌ No Garmin plugin instance found")

conn.close()
