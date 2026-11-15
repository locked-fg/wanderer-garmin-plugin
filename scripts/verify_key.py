#!/usr/bin/env python3
"""
Verify the Garmin plugin API key is saved in the database.
Usage: python3 verify_key.py
"""
import sqlite3
import json
from paths import database_path

def verify_garmin_key():
    try:
        conn = sqlite3.connect(database_path())
        cursor = conn.cursor()
        
        cursor.execute(
            'SELECT id, plugin_id, auth FROM plugin_instances WHERE plugin_id = "garmin"'
        )
        row = cursor.fetchone()
        
        if not row:
            print("❌ No Garmin plugin instance found in database")
            return False
        
        instance_id, plugin_id, auth = row
        
        if not auth:
            print("❌ No auth data found for Garmin plugin")
            return False
        
        auth_data = json.loads(auth)
        
        if 'apiKey' not in auth_data:
            print("❌ No apiKey field in auth data")
            return False
        
        api_key = auth_data['apiKey']
        
        # Check if it starts with GCDC_ (Garmin Personal API Key format)
        if api_key.startswith('GCDC_'):
            print(f"✅ Garmin plugin API key is saved in database")
            print(f"   Instance ID: {instance_id}")
            print(f"   Plugin ID: {plugin_id}")
            print(f"   Key format: Valid (starts with GCDC_)")
            print("   Credential value is redacted")
            return True
        else:
            print(f"⚠️  Key found but format may be incorrect")
            print(f"   Expected format: GCDC_<uuid>.<uuid>.<base64>.<timestamp>")
            print("   Credential value is redacted")
            return False
            
    except Exception as e:
        print(f"❌ Error: {e}")
        return False
    finally:
        conn.close()

if __name__ == '__main__':
    verify_garmin_key()
