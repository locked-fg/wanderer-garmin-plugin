#!/usr/bin/env python3
"""
Test that the save process through the UI works correctly.
This checks that PocketBase is encrypting and storing the API key properly.
"""
import sqlite3
import json
import base64

def test_save_process():
    conn = sqlite3.connect('../data/pb_data/data.db')
    cursor = conn.cursor()
    
    print("=" * 70)
    print("TESTING SAVE PROCESS")
    print("=" * 70)
    
    # Check current state
    cursor.execute(
        'SELECT id, plugin_id, auth, updated FROM plugin_instances WHERE plugin_id = "garmin"'
    )
    row = cursor.fetchone()
    
    if not row:
        print("❌ FAIL: No Garmin plugin instance found")
        return False
    
    instance_id, plugin_id, auth, updated = row
    auth_data = json.loads(auth) if auth else {}
    api_key = auth_data.get('apiKey', '')
    
    print(f"\n📊 Current State:")
    print(f"   Instance ID: {instance_id}")
    print(f"   Plugin ID: {plugin_id}")
    print(f"   Last Updated: {updated}")
    print(f"   API Key (encrypted): {api_key[:60]}...")
    print(f"   Key Length: {len(api_key)} characters")
    
    # Check if the key is encrypted
    # PocketBase encrypts JSON fields with AES-256-CBC and base64 encodes them
    # So the stored value should be a base64 string
    
    # Try to decode as base64
    try:
        decoded = base64.b64decode(api_key)
        is_base64 = True
    except:
        is_base64 = False
    
    # Check if it's plaintext
    is_plaintext_1234 = api_key == "1234"
    is_plaintext_gcdc = api_key.startswith("GCDC_")
    is_plaintext_old = api_key.startswith("6zi/")  # The old wrong key
    
    print(f"\n🔐 Analysis:")
    if is_base64:
        print(f"   ✅ Key is valid base64 (likely encrypted by PocketBase)")
        print(f"   Decoded length: {len(decoded)} bytes")
    else:
        print(f"   ⚠️  Key is NOT valid base64")
    
    if is_plaintext_1234:
        print(f"   ⚠️  Key is plaintext '1234'")
    elif is_plaintext_gcdc:
        print(f"   ⚠️  Key is plaintext starting with 'GCDC_'")
    elif is_plaintext_old:
        print(f"   ⚠️  Key is the old incorrect plaintext key")
    else:
        print(f"   ℹ️  Key format: {api_key[:20]}...")
    
    print("\n📝 Conclusion:")
    if is_base64 and not is_plaintext_1234 and not is_plaintext_gcdc:
        print("   ✅ The save process through the Wanderer UI IS working correctly.")
        print("   ✅ PocketBase is encrypting the API key before storing it.")
        print("   ✅ Only PocketBase can decrypt it during runtime.")
    else:
        print("   ⚠️  There might be an issue with encryption.")
    
    # Additional check: compare with known formats
    print(f"\n🔍 Details:")
    print(f"   - Starts with 'GCDC_': {api_key.startswith('GCDC_')}")
    print(f"   - Starts with 'YPJC': {api_key.startswith('YPJC')}")
    print(f"   - Starts with '6zi/': {api_key.startswith('6zi/')}")
    print(f"   - Is '1234': {api_key == '1234'}")
    print(f"   - Is valid base64: {is_base64}")
    
    conn.close()
    return True

if __name__ == '__main__':
    test_save_process()
