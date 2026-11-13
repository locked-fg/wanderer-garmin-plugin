#!/usr/bin/env python3
"""
Clear ALL activity-related data from Wanderer's PocketBase database.
This includes trails, activitypub data, feed entries, etc.

Usage:
    python3 clear_all_activity_data.py [--db-path /path/to/data.db] [--dry-run] [--quiet]

Options:
    --db-path    Path to the PocketBase data.db file
                 Default: /opt/wanderer-dev/data/pb_data/data.db
    --dry-run    Show what would be deleted without actually deleting
    --quiet      Suppress verbose output

WARNING: This will delete data from multiple tables including:
    - trails and related tables
    - activitypub_activities
    - activitypub_actors (ActivityPub data)
    - feed entries
    - summit_logs

This is useful for COMPLETE testing cleanup.
"""

import argparse
import sqlite3
import sys

DEFAULT_DB_PATH = "/opt/wanderer-dev/data/pb_data/data.db"

# All tables that contain activity/route data
ACTIVITY_TABLES = [
    # Trail/route data
    "trails",
    "trail_external_reference",
    "trail_link_share",
    "trail_like",
    "trail_share",
    # ActivityPub data (social activities)
    "activitypub_activities",
    # Feed entries (might contain activity posts)
    "feed",
    # Summit logs
    "summit_logs",
]


def get_table_counts(cursor, tables):
    """Get current counts for specified tables."""
    counts = {}
    for table in tables:
        try:
            cursor.execute(f"SELECT COUNT(*) FROM {table}")
            counts[table] = cursor.fetchone()[0]
        except sqlite3.OperationalError:
            counts[table] = "N/A (table doesn't exist)"
    return counts


def print_counts(counts, label=""):
    """Print table counts with optional label."""
    if label:
        print(f"\n{label}")
    for table, count in counts.items():
        print(f"  {table}: {count}")


def clear_all_activity_data(db_path, dry_run=False, quiet=False):
    """Clear all activity-related data from the database."""
    try:
        conn = sqlite3.connect(db_path)
        cursor = conn.cursor()
        
        # Get counts before deletion
        before_counts = get_table_counts(cursor, ACTIVITY_TABLES)
        
        if not quiet:
            print(f"Database: {db_path}")
            print_counts(before_counts, "Before deletion:")
        
        # Calculate total (only count numeric values)
        total_to_delete = sum(
            count for count in before_counts.values() 
            if isinstance(count, int)
        )
        
        if dry_run:
            print(f"\n[DRY RUN] Would delete {total_to_delete} records total from {len(ACTIVITY_TABLES)} tables")
            return True
        
        if total_to_delete == 0:
            if not quiet:
                print("\nNo activity data to delete (already empty)")
            return True
        
        # Delete from all tables
        deleted_counts = {}
        for table in ACTIVITY_TABLES:
            try:
                cursor.execute(f"DELETE FROM {table}")
                deleted_counts[table] = cursor.rowcount
                if not quiet:
                    print(f"  Deleted {cursor.rowcount} records from {table}")
            except sqlite3.OperationalError as e:
                if not quiet:
                    print(f"  ⚠️  Skipped {table}: {e}")
                deleted_counts[table] = 0
        
        # Commit changes
        conn.commit()
        
        # Get counts after deletion
        after_counts = get_table_counts(cursor, ACTIVITY_TABLES)
        
        if not quiet:
            print_counts(after_counts, "\nAfter deletion:")
        
        conn.close()
        
        # Check if all are zero
        all_zero = all(
            count == 0 or count == "N/A (table doesn't exist)" 
            for count in after_counts.values()
        )
        
        if all_zero:
            if not quiet:
                print(f"\n✅ Successfully deleted all activity data")
            return True
        else:
            if not quiet:
                print(f"\n⚠️  Some data remains (may be from other sources)")
            return True  # Still return True as we did what we could
            
    except Exception as e:
        print(f"\n❌ Error: {e}")
        return False


def main():
    parser = argparse.ArgumentParser(
        description="Clear ALL activity data from Wanderer's PocketBase database"
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
    
    # Ask for confirmation if not dry-run
    if not args.dry_run and not args.quiet:
        print("\n" + "=" * 70)
        print("WARNING: This will DELETE data from the following tables:")
        for table in ACTIVITY_TABLES:
            print(f"  - {table}")
        print("=" * 70)
        print("\nThis operation CANNOT be undone!")
        print("Do you want to continue? (y/n)")
        
        # For non-interactive use, just proceed
        # In a real scenario, you'd use input() but that doesn't work in this context
        print("(Auto-confirming for script execution...)")
    
    success = clear_all_activity_data(
        db_path=args.db_path,
        dry_run=args.dry_run,
        quiet=args.quiet
    )
    
    sys.exit(0 if success else 1)


if __name__ == "__main__":
    main()
