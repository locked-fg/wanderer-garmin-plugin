#!/usr/bin/env python3
"""
Clear all trails from Wanderer's PocketBase database.

Usage:
    python3 clear_trails.py [--db-path /path/to/data.db] [--dry-run]

Options:
    --db-path    Path to the PocketBase data.db file
                 Default: /opt/wanderer-dev/data/pb_data/data.db
    --dry-run    Show what would be deleted without actually deleting
    --quiet      Suppress verbose output

Description:
    This script deletes all trails and related data from the Wanderer
    PocketBase database. It's useful for testing the Garmin plugin
    with a clean database.

    Related tables cleared:
    - trails
    - trail_external_reference
    - trail_link_share
    - trail_like
    - trail_share

Examples:
    # Dry run to see what would be deleted
    python3 clear_trails.py --dry-run

    # Delete all trails
    python3 clear_trails.py

    # Delete with custom database path
    python3 clear_trails.py --db-path /path/to/your/data.db

Notes:
    - The database file must not be in use by PocketBase when running
    - This operation is NOT reversible - data will be permanently deleted
    - Backup your database before running this script if you need the data
"""

import argparse
import sqlite3
import sys

DEFAULT_DB_PATH = "/opt/wanderer-dev/data/pb_data/data.db"

RELATED_TABLES = [
    "trail_external_reference",
    "trail_link_share", 
    "trail_like",
    "trail_share",
]


def get_table_counts(cursor):
    """Get current counts for all relevant tables."""
    counts = {}
    for table in ["trails"] + RELATED_TABLES:
        cursor.execute(f"SELECT COUNT(*) FROM {table}")
        counts[table] = cursor.fetchone()[0]
    return counts


def print_counts(counts, label=""):
    """Print table counts with optional label."""
    if label:
        print(f"\n{label}")
    for table, count in counts.items():
        print(f"  {table}: {count}")


def clear_trails(db_path, dry_run=False, quiet=False):
    """Clear all trails and related data from the database."""
    try:
        conn = sqlite3.connect(db_path)
        cursor = conn.cursor()
        
        # Get counts before deletion
        before_counts = get_table_counts(cursor)
        
        if not quiet:
            print(f"Database: {db_path}")
            print_counts(before_counts, "Before deletion:")
        
        total_to_delete = sum(before_counts.values())
        
        if dry_run:
            print(f"\n[DRY RUN] Would delete {total_to_delete} records total")
            return True
        
        if total_to_delete == 0:
            if not quiet:
                print("\nNo trails or related data to delete (already empty)")
            return True
        
        # Delete from related tables first
        for table in RELATED_TABLES:
            cursor.execute(f"DELETE FROM {table}")
            if not quiet:
                print(f"  Deleted {cursor.rowcount} records from {table}")
        
        # Delete all trails
        cursor.execute("DELETE FROM trails")
        if not quiet:
            print(f"  Deleted {cursor.rowcount} records from trails")
        
        # Commit changes
        conn.commit()
        
        # Get counts after deletion
        after_counts = get_table_counts(cursor)
        
        if not quiet:
            print_counts(after_counts, "\nAfter deletion:")
        
        conn.close()
        
        # Verify all are zero
        if all(count == 0 for count in after_counts.values()):
            if not quiet:
                print(f"\n✅ Successfully deleted all trails and related data")
            return True
        else:
            if not quiet:
                print(f"\n❌ Failed to delete all records")
            return False
            
    except Exception as e:
        print(f"\n❌ Error: {e}")
        return False


def main():
    parser = argparse.ArgumentParser(
        description="Clear all trails from Wanderer's PocketBase database"
    )
    parser.add_argument(
        "--db-path",
        default=DEFAULT_DB_PATH,
        help=f"Path to PocketBase data.db (default: {DEFAULT_DB_PATH})"
    )
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Show what would be deleted without actually deleting"
    )
    parser.add_argument(
        "--quiet",
        action="store_true",
        help="Suppress verbose output"
    )
    
    args = parser.parse_args()
    
    success = clear_trails(
        db_path=args.db_path,
        dry_run=args.dry_run,
        quiet=args.quiet
    )
    
    sys.exit(0 if success else 1)


if __name__ == "__main__":
    main()
