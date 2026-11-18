package garminmap

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

	"github.com/open-wanderer/wanderer/plugins/garmin/internal/garminapi"
)

const detailFixture = `{
 "activityId": 12345,
 "activityName": "Morning Hike",
 "description": "Nice trail",
 "startTimeGMT": "2026-06-01 10:00:00",
 "duration": 3600,
 "distance": 10000,
 "elevationGain": 500,
 "elevationLoss": 480,
 "averageSpeed": 2.7,
 "maxSpeed": 5.1,
 "calories": 600,
 "activityType": {"typeKey": "hiking"},
 "privacy": "friends",
 "difficultyRating": "moderate",
 "locationName": "Somewhere",
 "photos": [
   {"photoId": 1, "photoUrl": "https://example.com/p1.jpg", "latitude": 46.0, "longitude": 7.0},
   {"photoId": 2, "photoUrl": ""},
   {"photoId": 3, "photoUrl": "//cdn.example.com/p3.jpg"},
   {"photoId": 4, "photoUrl": "::::"}
 ],
 "trackPoints": [
   {"latitude": 46.0, "longitude": 7.0, "elevation": 500, "timestamp": "2026-06-01T10:00:00Z"},
   {"latitude": 46.001, "longitude": 7.001, "elevation": 510, "timestamp": "2026-06-01T10:05:00Z"}
 ]
}`

func loadDetail(t *testing.T) *garminapi.ActivityDetail {
	t.Helper()
	var detail garminapi.ActivityDetail
	if err := json.Unmarshal([]byte(detailFixture), &detail); err != nil {
		t.Fatal(err)
	}
	return &detail
}

func TestMapActivityToTrail(t *testing.T) {
	item, warnings, err := MapActivityToTrail(loadDetail(t), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %v", warnings)
	}

	if item.Source.Provider != "garmin" || item.Source.ExternalID != "12345" {
		t.Fatalf("unexpected source: %#v", item.Source)
	}
	if !strings.Contains(item.Source.URL, "/activity/12345") {
		t.Fatalf("unexpected source URL: %q", item.Source.URL)
	}
	if item.Kind != "completed" || item.Name != "Morning Hike" || item.Description != "Nice trail" {
		t.Fatalf("unexpected identity fields: %#v", item)
	}
	if item.StartedAt != "2026-06-01T10:00:00Z" {
		t.Fatalf("StartedAt = %q; want normalized RFC3339 UTC", item.StartedAt)
	}
	if item.ActivityType != "hiking" {
		t.Fatalf("ActivityType = %q; want hiking", item.ActivityType)
	}
	if item.Privacy == nil || *item.Privacy != "friends" {
		t.Fatalf("Privacy = %v; want friends", item.Privacy)
	}
	if item.Difficulty != "moderate" {
		t.Fatalf("Difficulty = %q; want moderate", item.Difficulty)
	}

	if item.Track.Format != "gpx" || item.Track.ContentBase64 == "" {
		t.Fatalf("track not set: %#v", item.Track)
	}
	gpx, err := base64.StdEncoding.DecodeString(item.Track.ContentBase64)
	if err != nil {
		t.Fatalf("track is not valid base64: %v", err)
	}
	if !strings.Contains(string(gpx), "<gpx") || !strings.Contains(string(gpx), "<trkpt") {
		t.Fatalf("track does not look like GPX:\n%s", gpx)
	}

	if len(item.Photos) != 2 {
		t.Fatalf("got %d photos; want 2 (empty and invalid URLs skipped)", len(item.Photos))
	}
	if item.Photos[0].Source.URL != "https://example.com/p1.jpg" {
		t.Fatalf("photo[0] URL = %q", item.Photos[0].Source.URL)
	}
	if item.Photos[0].Lat == nil || *item.Photos[0].Lat != 46.0 {
		t.Fatalf("photo[0] lat = %v", item.Photos[0].Lat)
	}
	if item.Photos[1].Source.URL != "https://cdn.example.com/p3.jpg" {
		t.Fatalf("photo[1] URL = %q; want https-prefixed", item.Photos[1].Source.URL)
	}

	if item.Metadata["distance"] != 10000.0 || item.Metadata["providerCategory"] != "hiking" {
		t.Fatalf("unexpected metadata: %#v", item.Metadata)
	}
}

func TestMapActivityToTrailUsesProvidedGPX(t *testing.T) {
	item, _, err := MapActivityToTrail(loadDetail(t), []byte("<gpx>provided</gpx>"))
	if err != nil {
		t.Fatal(err)
	}
	gpx, err := base64.StdEncoding.DecodeString(item.Track.ContentBase64)
	if err != nil || string(gpx) != "<gpx>provided</gpx>" {
		t.Fatalf("provided GPX not used: %q err=%v", gpx, err)
	}
}

func TestMapActivityToTrailBadStartTimeWarns(t *testing.T) {
	detail := loadDetail(t)
	detail.StartTime = "garbage"
	item, warnings, err := MapActivityToTrail(detail, nil)
	if err != nil {
		t.Fatal(err)
	}
	if item.StartedAt != "" {
		t.Fatalf("StartedAt = %q; want empty", item.StartedAt)
	}
	if len(warnings) != 1 || !strings.Contains(warnings[0], "startTimeGMT") {
		t.Fatalf("warnings = %v; want one startTimeGMT warning", warnings)
	}
}

func TestMapPrivacy(t *testing.T) {
	tests := map[string]string{
		"private": "private", "PUBLIC": "public", "friends": "friends", "unknown": "private", "": "private",
	}
	for in, want := range tests {
		if got := MapPrivacy(in); got != want {
			t.Errorf("MapPrivacy(%q) = %q; want %q", in, got, want)
		}
	}
}

func TestMapDifficulty(t *testing.T) {
	tests := map[string]string{
		"easy": "easy", "Moderate": "moderate", "hard": "difficult", "difficult": "difficult", "": "",
	}
	for in, want := range tests {
		if got := MapDifficulty(in); got != want {
			t.Errorf("MapDifficulty(%q) = %q; want %q", in, got, want)
		}
	}
}

func TestMapActivityType(t *testing.T) {
	tests := map[string]string{
		"running": "running", "biking": "cycling", "MOUNTAIN_BIKING": "mountain_biking", "hiking": "hiking", "something_new": "other",
	}
	for in, want := range tests {
		if got := MapActivityType(in); got != want {
			t.Errorf("MapActivityType(%q) = %q; want %q", in, got, want)
		}
	}
}

func TestMapPhotosEmpty(t *testing.T) {
	if got := MapPhotos(nil); got != nil {
		t.Fatalf("MapPhotos(nil) = %v; want nil", got)
	}
}
