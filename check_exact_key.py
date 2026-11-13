#!/usr/bin/env python3
"""
Check the exact API key value in the database.
"""
import sqlite3
import json

conn = sqlite3.connect('../data/pb_data/data.db')
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
    print(f"Auth JSON: {auth}")
    print(f"\nAPI Key value: {api_key[:100]}...")
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
