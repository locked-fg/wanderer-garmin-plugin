package garmintrack

import (
	"os"
	"testing"
)

func TestParseActivityTrackDetailsPolylineFixture(t *testing.T) {
	body, err := os.ReadFile("../../testdata/garmin_activity_details_polyline.json")
	if err != nil {
		t.Fatal(err)
	}
	details, points, err := ParseActivityTrackDetails(body)
	if err != nil {
		t.Fatal(err)
	}
	if len(details.ActivityDetailMetrics) != 2 || len(points) != 2 {
		t.Fatalf("got %d metric rows and %d usable points; want 2 and 2", len(details.ActivityDetailMetrics), len(points))
	}
	if points[0].Latitude != 46.0001 || points[0].Longitude != 7.0001 || points[0].Elevation != 1200.5 {
		t.Fatalf("unexpected first point: %#v", points[0])
	}
}

func TestParseActivityTrackDetailsMetricsFallbackFixture(t *testing.T) {
	body, err := os.ReadFile("../../testdata/garmin_activity_details_metrics_fallback.json")
	if err != nil {
		t.Fatal(err)
	}
	_, points, err := ParseActivityTrackDetails(body)
	if err != nil {
		t.Fatal(err)
	}
	if len(points) != 2 {
		t.Fatalf("got %d fallback points; want 2", len(points))
	}
	if points[1].Latitude != 46.0002 || points[1].Longitude != 7.0002 || points[1].Elevation != 1201.0 {
		t.Fatalf("unexpected fallback point: %#v", points[1])
	}
}

func TestParseActivityTrackDetailsRejectsMalformedJSON(t *testing.T) {
	if _, _, err := ParseActivityTrackDetails([]byte("{")); err == nil {
		t.Fatal("malformed JSON unexpectedly parsed")
	}
}
