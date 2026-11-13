#!/usr/bin/env python3
import sqlite3
import json

conn = sqlite3.connect('../data/pb_data/data.db')
cursor = conn.cursor()

cursor.execute('SELECT * FROM plugin_instances WHERE plugin_id = "garmin"')
row = cursor.fetchone()

if row:
    # Get column names
    cursor.execute('PRAGMA table_info(plugin_instances)')
    columns = [col[1] for col in cursor.fetchall()]
    
    print("=" * 70)
    print("FULL GARMIN PLUGIN INSTANCE DATA:")
    print("=" * 70)
    
    for col, val in zip(columns, row):
        if val is None:
            print(f"{col}: NULL")
        elif col in ['auth', 'config', 'state', 'last_error']:
            try:
                parsed = json.loads(val)
                print(f"{col}:")
                print(f"  {json.dumps(parsed, indent=2)}")
            except:
                print(f"{col}: {val[:100]}...")
        else:
            print(f"{col}: {val}")
else:
    print("No Garmin instance found")

conn.close()
