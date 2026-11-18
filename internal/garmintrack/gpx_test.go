package garmintrack

import (
	"encoding/xml"
	"strings"
	"testing"
)

func TestBuildGPXEmpty(t *testing.T) {
	if got := BuildGPX(nil); got != nil {
		t.Fatalf("BuildGPX(nil) = %q, want nil", got)
	}
}

func TestBuildGPXProducesValidXML(t *testing.T) {
	points := []Point{
		{Latitude: 46.0001, Longitude: 7.0001, Elevation: 1200.5, Timestamp: "2026-11-15T10:00:00Z"},
		{Latitude: 46.0002, Longitude: 7.0002},
	}
	got := BuildGPX(points)
	if got == nil {
		t.Fatal("BuildGPX returned nil")
	}
	if !strings.HasPrefix(string(got), `<?xml version="1.0" encoding="UTF-8"?>`) {
		t.Fatalf("missing XML header: %q", got[:40])
	}

	var doc gpxDocument
	if err := xml.Unmarshal(got, &doc); err != nil {
		t.Fatalf("output is not valid XML: %v\n%s", err, got)
	}
	if doc.Version != "1.1" || doc.Xmlns != gpxNamespace {
		t.Fatalf("unexpected gpx attrs: %#v", doc)
	}
	if len(doc.Track.Seg.Points) != 2 {
		t.Fatalf("got %d trkpt; want 2", len(doc.Track.Seg.Points))
	}
	if doc.Track.Seg.Points[0].Lat != 46.0001 || doc.Track.Seg.Points[0].Lon != 7.0001 {
		t.Fatalf("unexpected first point: %#v", doc.Track.Seg.Points[0])
	}
	if doc.Track.Seg.Points[0].Ele == nil || *doc.Track.Seg.Points[0].Ele != 1200.5 {
		t.Fatalf("elevation not preserved: %#v", doc.Track.Seg.Points[0])
	}
	if doc.Track.Seg.Points[0].Time != "2026-11-15T10:00:00Z" {
		t.Fatalf("timestamp not preserved: %#v", doc.Track.Seg.Points[0])
	}
	if doc.Track.Seg.Points[1].Ele != nil || doc.Track.Seg.Points[1].Time != "" {
		t.Fatalf("zero elevation / empty time should be omitted: %#v", doc.Track.Seg.Points[1])
	}
}

func TestBuildGPXEscapesTimestamp(t *testing.T) {
	points := []Point{{Latitude: 46.0, Longitude: 7.0, Timestamp: `2026-11-15</time><injected>`}}
	got := string(BuildGPX(points))
	if strings.Contains(got, "<injected>") {
		t.Fatalf("timestamp was not escaped:\n%s", got)
	}
	var doc gpxDocument
	if err := xml.Unmarshal([]byte(got), &doc); err != nil {
		t.Fatalf("output with metacharacter timestamp is not valid XML: %v", err)
	}
	if doc.Track.Seg.Points[0].Time != `2026-11-15</time><injected>` {
		t.Fatalf("timestamp round-trip mismatch: %q", doc.Track.Seg.Points[0].Time)
	}
}
