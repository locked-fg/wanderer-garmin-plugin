#!/usr/bin/env python3
"""
Try to decrypt the PocketBase encrypted field.
PocketBase uses AES-256-CBC with PKCS7 padding.
"""
import base64
from Crypto.Cipher import AES
from Crypto.Util.Padding import unpad
import hashlib

# PocketBase encryption key from docker-compose
encryption_key = "fde406459dc1f6ca6f348e1f44a9a2af"

# The encrypted value from database
encrypted_b64 = "YPJC1g+U3jxIRmevRPtpUUFTOcAZiDIRgvaPhrRBmHKS8JqsnKjmv3snzEjyap5JaTGpMgHceDQpy9eBRxMA5eEU4gmRZdCT2vA8b76I37I7yiTiXi5X62DkvHp1lMGkt9PF4KmyhBTiUpkxcXmJBLSya8ma/iD8Q49qZHMuyITFA4b8g1QPpP5N+7hPMUQ9A1tu0u/+rT1VgjWGcngcNnI="

# Decode from base64
data = base64.b64decode(encrypted_b64)
print(f"Decoded length: {len(data)} bytes")
print(f"Hex: {data[:32].hex()}")

# PocketBase format: first 16 bytes might be IV, rest is ciphertext
# Or it might use a derived key

# Try various approaches
# Approach 1: AES-256-CBC with first 16 bytes as IV
if len(data) > 16:
    iv = data[:16]
    ciphertext = data[16:]
    
    # Derive key from encryption_key
    key = hashlib.sha256(encryption_key.encode()).digest()
    
    try:
        cipher = AES.new(key, AES.MODE_CBC, iv)
        decrypted = unpad(cipher.decrypt(ciphertext), AES.block_size)
        print(f"\n✅ Decrypted (SHA256 key): {decrypted.decode('utf-8')}")
    except Exception as e:
        print(f"\n❌ Failed with SHA256 key: {e}")

# Approach 2: Use encryption_key directly (32 bytes)
if len(encryption_key) == 32:
    key_bytes = encryption_key.encode()
    if len(data) > 16:
        iv = data[:16]
        ciphertext = data[16:]
        try:
            cipher = AES.new(key_bytes, AES.MODE_CBC, iv)
            decrypted = unpad(cipher.decrypt(ciphertext), AES.block_size)
            print(f"\n✅ Decrypted (direct key): {decrypted.decode('utf-8')}")
        except Exception as e:
            print(f"\n❌ Failed with direct key: {e}")

# Approach 3: Use MD5 hash of encryption_key
key_md5 = hashlib.md5(encryption_key.encode()).digest()
if len(data) > 16:
    iv = data[:16]
    ciphertext = data[16:]
    try:
        cipher = AES.new(key_md5, AES.MODE_CBC, iv)
        decrypted = unpad(cipher.decrypt(ciphertext), AES.block_size)
        print(f"\n✅ Decrypted (MD5 key): {decrypted.decode('utf-8')}")
    except Exception as e:
        print(f"\n❌ Failed with MD5 key: {e}")

print("\nNote: PocketBase may use a custom encryption scheme.")
