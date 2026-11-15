#!/usr/bin/env python3
import sqlite3
import json

conn = sqlite3.connect('../data/pb_data/data.db')
cursor = conn.cursor()

cursor.execute('SELECT last_error, status, updated FROM plugin_instances WHERE plugin_id = "garmin"')
row = cursor.fetchone()

if row:
    last_error, status, updated = row
    
    print("=" * 70)
    print("CURRENT PLUGIN STATE:")
    print("=" * 70)
    print(f"Status: {status}")
    print(f"Last Updated: {updated}")
    
    if last_error:
        print(f"\n❌ Last Error:")
        try:
            error_data = json.loads(last_error)
            print(json.dumps(error_data, indent=2))
        except:
            print(last_error)
    else:
        print(f"\n✅ No errors (last_error is NULL)")
else:
    print("No Garmin instance found")

conn.close()
