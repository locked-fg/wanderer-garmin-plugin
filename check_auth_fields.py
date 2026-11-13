#!/usr/bin/env python3
"""
Check what fields are in the auth JSON.
"""
import sqlite3
import json

conn = sqlite3.connect('../data/pb_data/data.db')
cursor = conn.cursor()

cursor.execute('SELECT auth FROM plugin_instances WHERE plugin_id = "garmin"')
row = cursor.fetchone()

if row and row[0]:
    auth_str = row[0]
    try:
        auth_data = json.loads(auth_str)
        
        print("=" * 70)
        print("AUTH FIELDS:")
        print("=" * 70)
        
        # Check which fields exist
        has_accessToken = 'accessToken' in auth_data
        has_apiKey = 'apiKey' in auth_data
        
        print(f"Has 'accessToken': {has_accessToken}")
        print(f"Has 'apiKey': {has_apiKey}")
        
        # Check all keys
        print(f"\nAll keys in auth: {list(auth_data.keys())}")
        
        # Important: if both exist, the plugin will use accessToken first!
        if has_accessToken and has_apiKey:
            print("\n⚠️  WARNING: Both 'accessToken' and 'apiKey' exist!")
            print("   The plugin will use 'accessToken' (first in code) and ignore 'apiKey'")
        elif has_accessToken:
            print("\n⚠️  WARNING: Only 'accessToken' exists!")
            print("   The plugin will use 'accessToken'")
        elif has_apiKey:
            print("\n✅ CORRECT: Only 'apiKey' exists")
            print("   The plugin will use 'apiKey'")
        else:
            print("\n❌ ERROR: No auth fields found!")
        
    except Exception as e:
        print(f"Error: {e}")
else:
    print("No auth data found")

conn.close()
