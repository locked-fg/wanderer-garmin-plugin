#!/usr/bin/env python3
"""
Clear Meilisearch indexes for Wanderer testing.

This script clears the 'trails' index from both Meilisearch instances:
- Port 7700 (wanderer-search, for wanderer-db on port 8090)
- Port 7701 (wanderer-dev-search-1, for wanderer-dev-db-1 on port 8091)

Usage:
    python3 clear_meilisearch.py

Note: Requires Meilisearch master key from docker-compose.yaml
"""
import requests
import sys

MEILI_MASTER_KEY = "b4aff09abab4f79e4d8af413245f66c2"

MEILI_INSTANCES = [
    {"url": "http://localhost:7700", "name": "wanderer-search (port 8090)"},
    {"url": "http://localhost:7701", "name": "wanderer-dev-search-1 (port 8091)"},
]

HEADERS = {
    "Authorization": f"Bearer {MEILI_MASTER_KEY}",
    "Content-Type": "application/json"
}


def clear_meilisearch_index(url, instance_name):
    """Clear the trails index from a Meilisearch instance."""
    print(f"\nClearing Meilisearch: {instance_name}")
    print(f"URL: {url}")
    
    try:
        # Check if trails index exists
        response = requests.get(f"{url}/indexes", headers=HEADERS)
        if response.status_code != 200:
            print(f"  ❌ Could not list indexes: {response.text}")
            return False
        
        indexes = response.json().get("results", [])
        trails_index = next((idx for idx in indexes if idx["uid"] == "trails"), None)
        
        if not trails_index:
            print(f"  ✅ No 'trails' index to delete")
            return True
        
        # Delete the trails index
        response = requests.delete(f"{url}/indexes/trails", headers=HEADERS)
        if response.status_code != 202:
            print(f"  ❌ Failed to delete index: {response.text}")
            return False
        
        task_uid = response.json().get("taskUid")
        print(f"  ✅ Deletion enqueued (task: {task_uid})")
        
        # Wait for completion
        if task_uid:
            response = requests.get(f"{url}/tasks/{task_uid}", headers=HEADERS)
            if response.status_code == 200:
                task = response.json()
                if task.get("status") == "succeeded":
                    deleted = task.get("details", {}).get("deletedDocuments", 0)
                    print(f"  ✅ Deleted {deleted} documents")
                else:
                    print(f"  ⚠️  Task status: {task.get('status')}")
        
        return True
        
    except Exception as e:
        print(f"  ❌ Error: {e}")
        return False


def main():
    print("=" * 70)
    print("CLEARING MEILISEARCH INDEXES")
    print("=" * 70)
    
    all_success = True
    for instance in MEILI_INSTANCES:
        success = clear_meilisearch_index(instance["url"], instance["name"])
        all_success = all_success and success
    
    if all_success:
        print("\n" + "=" * 70)
        print("✅ All Meilisearch indexes cleared")
        print("=" * 70)
    else:
        print("\n" + "=" * 70)
        print("⚠️  Some indexes could not be cleared")
        print("=" * 70)
    
    sys.exit(0 if all_success else 1)


if __name__ == "__main__":
    main()
