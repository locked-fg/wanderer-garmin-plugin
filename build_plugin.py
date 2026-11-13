#!/usr/bin/env python3
"""
Build script for the Garmin plugin.

Builds the plugin WASM inside a Docker container with TinyGo,
installs it to the local wanderer data directory, and verifies the build.

Usage:
    python3 build_plugin.py [--clean] [--no-install] [--verbose]

Options:
    --clean       Remove build artifacts before building
    --no-install  Build but don't install to data directory
    --verbose     Show detailed output
    --help        Show this help message
"""

import argparse
import os
import shutil
import subprocess
import sys
from pathlib import Path

# Configuration
PLUGIN_DIR = Path(__file__).parent.resolve()
WANDERER_ROOT = PLUGIN_DIR.parent.resolve()
DIST_DIR = PLUGIN_DIR / "dist" / "garmin"
INSTALL_DIR = WANDERER_ROOT / "data" / "plugins" / "garmin"

# Docker image for TinyGo
TINYGO_DOCKER_IMAGE = "tinygo/tinygo:0.44.0"


def run_command(cmd, cwd=None, verbose=False, check=True):
    """Run a shell command and return the result."""
    if verbose:
        print(f"\n[CMD] {cmd}")
        print(f"[CWD] {cwd}")
    
    result = subprocess.run(
        cmd,
        shell=True,
        cwd=cwd,
        capture_output=True,
        text=True
    )
    
    if verbose:
        if result.stdout:
            print(f"[STDOUT]\n{result.stdout}")
        if result.stderr:
            print(f"[STDERR]\n{result.stderr}")
    
    if check and result.returncode != 0:
        print(f"\n❌ Command failed: {cmd}")
        print(f"stdout: {result.stdout}")
        print(f"stderr: {result.stderr}")
        sys.exit(1)
    
    return result


def clean(verbose=False):
    """Remove build artifacts."""
    print("🧹 Cleaning build artifacts...")
    
    # Remove dist directory
    if DIST_DIR.exists():
        shutil.rmtree(DIST_DIR)
        if verbose:
            print(f"  Removed: {DIST_DIR}")
    
    # Remove plugin.wasm in plugin root
    wasm_file = PLUGIN_DIR / "plugin.wasm"
    if wasm_file.exists():
        wasm_file.unlink()
        if verbose:
            print(f"  Removed: {wasm_file}")
    
    print("✅ Cleaned build artifacts")


def build_in_docker(verbose=False):
    """Build the plugin WASM inside a Docker container."""
    print("\n🔨 Building Garmin plugin in Docker...")
    
    # Create dist directory
    DIST_DIR.mkdir(parents=True, exist_ok=True)
    
    # Docker build command
    # Mount current directory into container and run tinygo build
    docker_cmd = f"""
docker run --rm \\
    -v {PLUGIN_DIR}:/src \\
    -w /src \\
    {TINYGO_DOCKER_IMAGE} \\
    tinygo build -target=wasi -scheduler=none -no-debug -o /src/dist/garmin/plugin.wasm .
"""
    
    result = run_command(docker_cmd, cwd=PLUGIN_DIR, verbose=verbose, check=True)
    
    # Check if WASM was created
    wasm_path = DIST_DIR / "plugin.wasm"
    if not wasm_path.exists():
        print(f"\n❌ Build failed: {wasm_path} not created")
        sys.exit(1)
    
    # Check file size
    file_size = wasm_path.stat().st_size
    print(f"✅ Built plugin.wasm ({file_size:,} bytes)")
    
    return wasm_path


def copy_assets(verbose=False):
    """Copy plugin assets (json, icon) to dist directory."""
    print("\n📁 Copying plugin assets...")
    
    # Copy plugin.json
    plugin_json = PLUGIN_DIR / "plugin.json"
    if plugin_json.exists():
        shutil.copy2(plugin_json, DIST_DIR)
        if verbose:
            print(f"  Copied: {plugin_json} -> {DIST_DIR}")
    
    # Copy icon.svg
    icon_svg = PLUGIN_DIR / "assets" / "icon.svg"
    if icon_svg.exists():
        shutil.copy2(icon_svg, DIST_DIR / "icon.svg")
        if verbose:
            print(f"  Copied: {icon_svg} -> {DIST_DIR}")
    else:
        # Try alternative location
        icon_svg = PLUGIN_DIR / "icon.svg"
        if icon_svg.exists():
            shutil.copy2(icon_svg, DIST_DIR / "icon.svg")
            if verbose:
                print(f"  Copied: {icon_svg} -> {DIST_DIR}")
    
    print("✅ Copied plugin assets")


def install_plugin(verbose=False):
    """Install the plugin to the local wanderer data directory."""
    print("\n🚀 Installing plugin to wanderer data directory...")
    
    # Create install directory
    INSTALL_DIR.mkdir(parents=True, exist_ok=True)
    
    # Copy files
    files_to_copy = ["plugin.json", "plugin.wasm", "icon.svg"]
    for filename in files_to_copy:
        src = DIST_DIR / filename
        dst = INSTALL_DIR / filename
        if src.exists():
            shutil.copy2(src, dst)
            if verbose:
                print(f"  Installed: {src} -> {dst}")
        else:
            print(f"  ⚠️  Missing file: {src}")
    
    print(f"✅ Plugin installed to: {INSTALL_DIR}")


def verify_build(verbose=False):
    """Verify the build was successful."""
    print("\n🔍 Verifying build...")
    
    wasm_path = INSTALL_DIR / "plugin.wasm"
    if not wasm_path.exists():
        print(f"❌ Verification failed: {wasm_path} not found")
        return False
    
    # Check file size is reasonable (> 100KB)
    file_size = wasm_path.stat().st_size
    if file_size < 100000:
        print(f"⚠️  WASM file seems small: {file_size:,} bytes")
    else:
        print(f"✅ WASM file size: {file_size:,} bytes")
    
    # Check for gojs dependency (should not exist)
    try:
        result = subprocess.run(
            ["grep", "-c", "gojs", str(wasm_path)],
            capture_output=True,
            text=True
        )
        if result.stdout.strip() != "0":
            print(f"⚠️  WASM still has gojs dependency!")
            return False
        else:
            print("✅ No gojs dependency found")
    except Exception as e:
        print(f"⚠️  Could not check for gojs: {e}")
    
    # Verify all required files exist
    required_files = ["plugin.json", "plugin.wasm", "icon.svg"]
    all_exist = True
    for filename in required_files:
        filepath = INSTALL_DIR / filename
        if filepath.exists():
            print(f"✅ {filename} exists")
        else:
            print(f"❌ {filename} missing")
            all_exist = False
    
    return all_exist


def main():
    parser = argparse.ArgumentParser(
        description="Build and install the Garmin plugin",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog=__doc__
    )
    parser.add_argument(
        "--clean",
        action="store_true",
        help="Remove build artifacts before building"
    )
    parser.add_argument(
        "--no-install",
        action="store_true",
        help="Build but don't install to data directory"
    )
    parser.add_argument(
        "--verbose", "-v",
        action="store_true",
        help="Show detailed output"
    )
    args = parser.parse_args()
    
    verbose = args.verbose
    
    print("=" * 60)
    print("Garmin Plugin Build Script")
    print("=" * 60)
    print(f"Plugin directory: {PLUGIN_DIR}")
    print(f"Install directory: {INSTALL_DIR}")
    print(f"Docker image: {TINYGO_DOCKER_IMAGE}")
    print("-" * 60)
    
    # Step 1: Clean if requested
    if args.clean:
        clean(verbose)
    
    # Step 2: Build
    build_in_docker(verbose)
    
    # Step 3: Copy assets
    copy_assets(verbose)
    
    # Step 4: Install (unless --no-install)
    if not args.no_install:
        install_plugin(verbose)
        
        # Step 5: Verify
        success = verify_build(verbose)
        if success:
            print("\n" + "=" * 60)
            print("✅ Build and installation completed successfully!")
            print(f"Plugin installed at: {INSTALL_DIR}")
            print("=" * 60)
        else:
            print("\n" + "=" * 60)
            print("⚠️  Build completed but verification had issues")
            print("=" * 60)
            sys.exit(1)
    else:
        print("\n" + "=" * 60)
        print("✅ Build completed (not installed)")
        print(f"Build artifacts at: {DIST_DIR}")
        print("=" * 60)


if __name__ == "__main__":
    main()
