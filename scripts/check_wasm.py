#!/usr/bin/env python3
"""
Check if a WASM file has Go runtime (gojs) dependencies.

Usage:
    python3 check_wasm.py [--file /path/to/plugin.wasm]

Options:
    --file    Path to WASM file (default: /opt/wanderer-dev/data/plugins/garmin/plugin.wasm)
    --quiet   Suppress verbose output

Description:
    WASM files built for server-side (WASI) should NOT import from 'gojs'.
    WASM files built for browser (wasm target) typically DO import from 'gojs'.
    
    For Extism/PocketBase plugins, you want WASI target with NO gojs imports.
    
    Fix: Build with `tinygo build -target wasi -scheduler=none -no-debug`
"""

import argparse
import os
import subprocess
import sys

DEFAULT_WASM_PATH = "/opt/wanderer-dev/data/plugins/garmin/plugin.wasm"


def check_wasm_file(filepath, quiet=False):
    """Check if WASM file has gojs dependencies."""
    try:
        if not os.path.exists(filepath):
            print(f"❌ File not found: {filepath}")
            return False
        
        size = os.path.getsize(filepath)
        size_kb = size / 1024
        
        if not quiet:
            print(f"Checking: {filepath}")
            print(f"Size: {size_kb:.1f} KB ({size} bytes)")
            print()
        
        # Check for gojs imports
        try:
            result = subprocess.run(
                ["strings", filepath],
                capture_output=True,
                text=True,
                timeout=5
            )
            
            gojs_count = result.stdout.lower().count("gojs")
            
            if gojs_count > 0:
                print(f"❌ WASM file has {gojs_count} 'gojs' imports")
                print("   This indicates it was built with -target=wasm (browser)")
                print("   Extism/PocketBase requires -target=wasi (server-side)")
                print()
                print("   FIX: Rebuild with:")
                print("   tinygo build -target wasi -scheduler=none -no-debug -o plugin.wasm .")
                return False
            else:
                print("✅ WASM file has NO 'gojs' imports")
                print("   This is correct for Extism/PocketBase plugins")
                
                # Additional checks
                if size_kb < 500:
                    print(f"   Note: File is small ({size_kb:.1f} KB) - likely WASI build")
                elif size_kb > 2000:
                    print(f"   Note: File is large ({size_kb:.1f} KB) - might be wasm (browser) build")
                
                return True
                
        except FileNotFoundError:
            print("❌ 'strings' command not found. Install binutils.")
            return None
            
    except Exception as e:
        print(f"❌ Error: {e}")
        return False


def main():
    parser = argparse.ArgumentParser(
        description="Check WASM file for Go runtime dependencies"
    )
    parser.add_argument(
        "--file", 
        default=DEFAULT_WASM_PATH,
        help="Path to WASM file"
    )
    parser.add_argument(
        "--quiet",
        action="store_true",
        help="Suppress verbose output"
    )
    
    args = parser.parse_args()
    result = check_wasm_file(args.file, args.quiet)
    
    if result is None:
        sys.exit(2)  # Command not found
    sys.exit(0 if result else 1)


if __name__ == "__main__":
    main()
