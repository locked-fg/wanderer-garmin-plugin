#!/usr/bin/env python3
import sqlite3
import json

conn = sqlite3.connect('../data/pb_data/data.db')
cursor = conn.cursor()

cursor.execute('SELECT id, plugin_id, auth, config, status FROM plugin_instances')
rows = cursor.fetchall()

print("=" * 70)
print("ALL PLUGIN INSTANCES:")
print("=" * 70)

for row in rows:
    instance_id, plugin_id, auth, config, status = row
    print(f"\nID: {instance_id}")
    print(f"Plugin ID: {plugin_id}")
    print(f"Status: {status}")
    
    if auth:
        try:
            auth_data = json.loads(auth)
            api_key = auth_data.get('apiKey', '')
            print(f"API Key: {api_key[:80]}...")
            
            if api_key.startswith('GCDC_'):
                parts = api_key.split('.')
                if len(parts) > 1:
                    print(f"  Extracted UUID: {parts[0].replace('GCDC_', '')}")
            elif api_key == '1234':
                print(f"  ⚠️  KEY IS '1234' (test key)")
            else:
                print(f"  ℹ️  Encrypted/Unknown format")
        except:
            print(f"Auth: {auth}")

conn.close()
