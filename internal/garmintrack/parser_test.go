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

func TestParseActivityTrackDetailsSkipsInvalidCoordinates(t *testing.T) {
	body := []byte(`{"geoPolylineDTO":{"polyline":[
		{"lat": 46.0, "lon": 7.0, "altitude": 500.0},
		{"lat": 91.0, "lon": 7.0},
		{"lat": 46.0, "lon": -181.0},
		{"lat": 46.1, "lon": 7.1, "valid": false}
	]}}`)
	_, points, err := ParseActivityTrackDetails(body)
	if err != nil {
		t.Fatal(err)
	}
	if len(points) != 1 {
		t.Fatalf("got %d points; want 1 (only the valid in-range point)", len(points))
	}
	if points[0].Latitude != 46.0 || points[0].Longitude != 7.0 {
		t.Fatalf("unexpected surviving point: %#v", points[0])
	}
}

func TestParseActivityTrackDetailsMetricsFallbackSkipsOutOfRange(t *testing.T) {
	body := []byte(`{"metricDescriptors":[
		{"metricsIndex": 0, "key": "latitude"},
		{"metricsIndex": 1, "key": "longitude"},
		{"metricsIndex": 2, "key": "elevation"}
	],"activityDetailMetrics":[
		{"metrics": [46.0, 7.0, 500.0]},
		{"metrics": [95.0, 7.0, 500.0]},
		{"metrics": [46.0, null, 500.0]}
	]}`)
	_, points, err := ParseActivityTrackDetails(body)
	if err != nil {
		t.Fatal(err)
	}
	if len(points) != 1 {
		t.Fatalf("got %d fallback points; want 1 (out-of-range and null rows skipped)", len(points))
	}
	if points[0].Latitude != 46.0 || points[0].Longitude != 7.0 || points[0].Elevation != 500.0 {
		t.Fatalf("unexpected surviving point: %#v", points[0])
	}
}
