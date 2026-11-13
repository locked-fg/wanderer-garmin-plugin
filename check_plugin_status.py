#!/usr/bin/env python3
"""
Check Garmin plugin instance status, errors, and configuration.

Usage:
    python3 check_plugin_status.py [--db-path /path/to/data.db]

Options:
    --db-path    Path to PocketBase data.db (default: /opt/wanderer-dev/data/pb_data/data.db)
    --quiet      Suppress verbose output
"""

import argparse
import sqlite3
import json
import sys

DEFAULT_DB_PATH = "/opt/wanderer-dev/data/pb_data/data.db"


def check_plugin_status(db_path, quiet=False):
    """Check the Garmin plugin instance status."""
    try:
        conn = sqlite3.connect(db_path)
        cursor = conn.cursor()
        
        # Get the garmin plugin instance
        cursor.execute('SELECT id, status, enabled, last_error, last_sync_at, updated FROM plugin_instances WHERE plugin_id = "garmin"')
        instance = cursor.fetchone()
        
        if not instance:
            print("❌ No Garmin plugin instance found in database")
            return False
        
        instance_id, status, enabled, last_error, last_sync_at, updated = instance
        
        if not quiet:
            print("=" * 60)
            print("GARMIN PLUGIN INSTANCE STATUS")
            print("=" * 60)
        
        print(f"Instance ID: {instance_id}")
        print(f"Status: {status}")
        print(f"Enabled: {bool(enabled)}")
        print(f"Last Sync: {last_sync_at or 'Never'}")
        print(f"Last Updated: {updated}")
        
        if last_error:
            print(f"\n❌ LAST ERROR:")
            try:
                error_data = json.loads(last_error)
                print(f"  Code: {error_data.get('code', 'unknown')}")
                print(f"  Message: {error_data.get('message', 'N/A')}")
                if 'details' in error_data:
                    print(f"  Details: {error_data['details']}")
            except:
                print(f"  Raw: {last_error}")
        else:
            print("\n✅ No errors recorded")
        
        # Check auth configuration
        if not quiet:
            print("\n" + "=" * 60)
            print("AUTH CONFIGURATION")
            print("=" * 60)
        
        cursor.execute('SELECT auth FROM plugin_instances WHERE plugin_id = "garmin"')
        auth_row = cursor.fetchone()
        
        if auth_row and auth_row[0]:
            try:
                auth = json.loads(auth_row[0])
                if not quiet:
                    print(f"Auth fields: {list(auth.keys())}")
                
                if 'apiKey' in auth:
                    key_len = len(auth['apiKey'])
                    # Check if it looks encrypted (PocketBase encrypted values are long base64 strings)
                    is_encrypted = key_len > 100 and auth['apiKey'].startswith('/')
                    
                    if is_encrypted:
                        print(f"✅ apiKey is SAVED (encrypted, length: {key_len})")
                        print(f"   Starts with: {auth['apiKey'][:40]}...")
                    else:
                        print(f"✅ apiKey is SAVED (plain, length: {key_len})")
                        print(f"   Value: {auth['apiKey']}")
                elif 'accessToken' in auth:
                    print(f"✅ accessToken is SAVED (length: {len(auth['accessToken'])})")
                else:
                    print("❌ No apiKey or accessToken found in auth!")
            except json.JSONDecodeError:
                print(f"❌ Failed to parse auth JSON")
        else:
            print("❌ No auth configuration found")
        
        # Check state
        if not quiet:
            print("\n" + "=" * 60)
            print("PLUGIN STATE")
            print("=" * 60)
        
        cursor.execute('SELECT state FROM plugin_instances WHERE plugin_id = "garmin"')
        state_row = cursor.fetchone()
        
        if state_row and state_row[0]:
            try:
                state = json.loads(state_row[0])
                print(f"State: {json.dumps(state, indent=2)}")
            except:
                print(f"State (raw): {state_row[0][:200]}")
        else:
            print("No state recorded")
        
        conn.close()
        
        # Summary
        if not quiet:
            print("\n" + "=" * 60)
            print("SUMMARY")
            print("=" * 60)
        
        if status == "configured":
            print("✅ Plugin is configured")
        elif status == "needs_auth":
            print("⚠️  Plugin needs authentication")
        elif status == "error":
            print("❌ Plugin has errors")
        else:
            print(f"ℹ️  Plugin status: {status}")
        
        return True
        
    except Exception as e:
        print(f"❌ Error: {e}")
        return False


def main():
    parser = argparse.ArgumentParser(description="Check Garmin plugin instance status")
    parser.add_argument("--db-path", default=DEFAULT_DB_PATH, help="Path to PocketBase data.db")
    parser.add_argument("--quiet", action="store_true", help="Suppress verbose output")
    
    args = parser.parse_args()
    success = check_plugin_status(args.db_path, args.quiet)
    sys.exit(0 if success else 1)


if __name__ == "__main__":
    main()
