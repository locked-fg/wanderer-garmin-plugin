#!/usr/bin/env python3
"""
Check if Garmin API key is saved and trigger import.
Usage: python3 check_and_trigger_import.py [--trigger]
"""
import sqlite3
import json
import sys
from datetime import datetime, timezone

def check_garmin_key():
    """Check if Garmin API key is saved in the database."""
    print("=" * 70)
    print("CHECKING GARMIN PLUGIN API KEY")
    print("=" * 70)
    
    try:
        conn = sqlite3.connect('../data/pb_data/data.db')
        cursor = conn.cursor()
        
        # Get the plugin instance
        cursor.execute('SELECT id, plugin_id, auth, status FROM plugin_instances WHERE plugin_id = "garmin"')
        row = cursor.fetchone()
        
        if not row:
            print("❌ No Garmin plugin instance found in database")
            return False
        
        instance_id, plugin_id, auth, status = row
        
        print(f"\nPlugin Instance Found:")
        print(f"  ID: {instance_id}")
        print(f"  Plugin: {plugin_id}")
        print(f"  Status: {status}")
        
        if not auth:
            print("\n❌ No auth data found for Garmin plugin")
            return False
        
        auth_data = json.loads(auth)
        
        if 'apiKey' not in auth_data:
            print("\n❌ No apiKey field in auth data")
            print(f"   Available fields: {list(auth_data.keys())}")
            return False
        
        api_key = auth_data['apiKey']
        
        # Check if key is saved (non-empty)
        if api_key and len(api_key) > 0:
            print(f"\n✅ API key is SAVED in database")
            print(f"   Instance ID: {instance_id}")
            print(f"   Key length: {len(api_key)} characters")
            print(f"   Key preview: {api_key[:50]}...")
            
            # Check if it's encrypted (PocketBase encrypts auth fields)
            if api_key.startswith('mN'):
                print(f"   Note: Key appears to be encrypted (PocketBase encrypts auth fields)")
            
            return True
        else:
            print("\n❌ API key is empty")
            return False
            
    except Exception as e:
        print(f"\n❌ Error: {e}")
        return False
    finally:
        conn.close()

def trigger_import():
    """Trigger a plugin sync by updating the timestamp."""
    print("\n" + "=" * 70)
    print("TRIGGERING GARMIN PLUGIN IMPORT")
    print("=" * 70)
    
    try:
        conn = sqlite3.connect('../data/pb_data/data.db')
        cursor = conn.cursor()
        
        now = datetime.now(timezone.utc).isoformat(timespec='milliseconds') + 'Z'
        
        # Clear any errors and update the timestamp
        cursor.execute(
            'UPDATE plugin_instances SET last_error = NULL, status = "syncing", updated = ? WHERE plugin_id = "garmin"',
            (now,)
        )
        
        conn.commit()
        
        # Verify
        cursor.execute('SELECT id, status, updated FROM plugin_instances WHERE plugin_id = "garmin"')
        row = cursor.fetchone()
        
        if row:
            print(f"\n✅ Sync triggered successfully")
            print(f"   Instance ID: {row[0]}")
            print(f"   Status: {row[1]}")
            print(f"   Updated: {row[2]}")
            print(f"\n📌 PocketBase should now pick up the change and start importing")
            print(f"📌 Check the GUI for import progress")
            return True
        else:
            print("\n❌ No Garmin instance found")
            return False
            
    except Exception as e:
        print(f"\n❌ Error triggering sync: {e}")
        return False
    finally:
        conn.close()

if __name__ == '__main__':
    trigger = '--trigger' in sys.argv
    
    # Always check the key first
    key_saved = check_garmin_key()
    
    if not key_saved:
        print("\n" + "=" * 70)
        print("ABORTING: API key is not saved properly")
        print("=" * 70)
        sys.exit(1)
    
    # If --trigger flag or no flag, trigger import
    if trigger or len(sys.argv) == 1:
        trigger_import()
    else:
        print("\n✅ API key is saved. Use --trigger to start import.")
