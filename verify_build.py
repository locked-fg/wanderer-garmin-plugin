#!/usr/bin/env python3
"""
Compare build configuration between Garmin and reference plugins.

Usage:
    python3 verify_build.py [--reference /path/to/reference/plugin.json]

Options:
    --reference   Path to reference plugin.json (default: uses Strava plugin from SDK)
    --quiet       Suppress verbose output

Description:
    Compares the build configuration (auth context, secret fields, etc.)
    between the Garmin plugin and a reference plugin (default: Strava).
    
    This helps identify configuration differences that might cause issues.
"""

import argparse
import json
import os
import sys

GARMIN_PLUGIN_JSON = "/opt/wanderer-dev/wanderer-garmin-plugin/plugin.json"
STRAVA_PLUGIN_JSON = "/opt/wanderer-dev/wanderer-sdk-source/plugins/strava/plugin.json"


def load_plugin_json(path):
    """Load and parse a plugin.json file."""
    try:
        with open(path, 'r') as f:
            return json.load(f)
    except FileNotFoundError:
        return None
    except json.JSONDecodeError as e:
        print(f"❌ Failed to parse {path}: {e}")
        return None


def compare_auth_configs(garmin, reference, quiet=False):
    """Compare auth configurations between plugins."""
    if not quiet:
        print("=" * 60)
        print("AUTH CONFIGURATION COMPARISON")
        print("=" * 60)
    
    garmin_auth = garmin.get('auth', {})
    reference_auth = reference.get('auth', {})
    
    if not garmin_auth and not reference_auth:
        if not quiet:
            print("ℹ️  Neither plugin has auth configuration")
        return True
    
    # Compare contexts
    garmin_contexts = garmin_auth.get('contexts', {})
    reference_contexts = reference_auth.get('contexts', {})
    
    if not quiet:
        print(f"\nGarmin contexts: {list(garmin_contexts.keys())}")
        print(f"Reference contexts: {list(reference_contexts.keys())}")
    
    # Check if they have the same context type
    if garmin_contexts and reference_contexts:
        garmin_type = next(iter(garmin_contexts.values())).get('type')
        ref_type = next(iter(reference_contexts.values())).get('type')
        
        if garmin_type != ref_type:
            print(f"\n⚠️  Context types differ:")
            print(f"   Garmin: {garmin_type}")
            print(f"   Reference: {ref_type}")
    
    # Compare secret field configuration
    if not quiet:
        print("\n" + "-" * 60)
        print("SECRET FIELD CONFIGURATION")
        print("-" * 60)
    
    issues = []
    
    for context_name, context_config in garmin_contexts.items():
        if not quiet:
            print(f"\n[Garmin] Context: {context_name}")
            print(f"  Type: {context_config.get('type')}")
        
        # Check secretField
        secret_field = context_config.get('secretField')
        secret_fields = context_config.get('secretFields')
        
        if secret_field:
            if not quiet:
                print(f"  secretField: {secret_field}")
            
            # Check if reference uses secretFields
            for ref_context_name, ref_context_config in reference_contexts.items():
                ref_secret_fields = ref_context_config.get('secretFields')
                if ref_secret_fields:
                    issues.append(
                        f"Garmin uses 'secretField' (singular) but reference "
                        f"'{ref_context_name}' uses 'secretFields' (plural)"
                    )
        
        if secret_fields:
            if not quiet:
                print(f"  secretFields: {secret_fields}")
    
    for context_name, context_config in reference_contexts.items():
        if not quiet:
            print(f"\n[Reference] Context: {context_name}")
            print(f"  Type: {context_config.get('type')}")
        
        secret_field = context_config.get('secretField')
        secret_fields = context_config.get('secretFields')
        
        if secret_field:
            if not quiet:
                print(f"  secretField: {secret_field}")
        
        if secret_fields:
            if not quiet:
                print(f"  secretFields: {secret_fields}")
    
    # Print issues
    if issues:
        if not quiet:
            print("\n" + "=" * 60)
            print("⚠️  POTENTIAL ISSUES FOUND")
            print("=" * 60)
        for issue in issues:
            print(f"  • {issue}")
        return False
    else:
        if not quiet:
            print("\n✅ No configuration issues found")
        return True


def compare_build_targets(garmin, reference, quiet=False):
    """Compare build targets."""
    if not quiet:
        print("\n" + "=" * 60)
        print("BUILD TARGET RECOMMENDATION")
        print("=" * 60)
    
    # Check if reference uses wasi or wasm
    # We can't directly tell from plugin.json, but we can infer from context type
    ref_contexts = reference.get('auth', {}).get('contexts', {})
    if ref_contexts:
        ref_type = next(iter(ref_contexts.values())).get('type')
        
        # Strava uses oauth2, which typically uses wasi
        if ref_type == "oauth2":
            if not quiet:
                print("Reference plugin uses OAuth2")
                print("Recommended build: tinygo build -target wasi -scheduler=none -no-debug")
        elif ref_type == "bearer":
            if not quiet:
                print("Reference plugin uses Bearer token")
                print("Recommended build: tinygo build -target wasi -scheduler=none -no-debug")
    
    return True


def main():
    parser = argparse.ArgumentParser(
        description="Compare build configuration with reference plugin"
    )
    parser.add_argument(
        "--reference",
        default=STRAVA_PLUGIN_JSON,
        help="Path to reference plugin.json"
    )
    parser.add_argument(
        "--quiet",
        action="store_true",
        help="Suppress verbose output"
    )
    
    args = parser.parse_args()
    
    # Load plugins
    garmin = load_plugin_json(GARMIN_PLUGIN_JSON)
    if garmin is None:
        print(f"❌ Failed to load Garmin plugin.json")
        sys.exit(1)
    
    reference = load_plugin_json(args.reference)
    if reference is None:
        print(f"⚠️  Reference plugin not found at {args.reference}")
        print("   Using internal defaults...")
        reference = {"auth": {"contexts": {"oauth_access_token": {"type": "oauth2", "secretFields": ["accessToken"]}}}}
    
    # Compare
    auth_ok = compare_auth_configs(garmin, reference, args.quiet)
    compare_build_targets(garmin, reference, args.quiet)
    
    if not auth_ok:
        sys.exit(1)
    sys.exit(0)


if __name__ == "__main__":
    main()
