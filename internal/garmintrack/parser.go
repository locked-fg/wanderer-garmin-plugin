package garmintrack

import (
	"bytes"
	"encoding/json"
	"math"
	"strings"
)

type MetricDescriptor struct {
	MetricsIndex int    `json:"metricsIndex"`
	Key          string `json:"key"`
}

type MetricRow struct {
	Metrics []json.RawMessage `json:"metrics"`
}

type GeoPoint struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
	Elevation float64 `json:"altitude"`
	Valid     *bool   `json:"valid"`
}

type Details struct {
	MetricDescriptors     []MetricDescriptor `json:"metricDescriptors"`
	ActivityDetailMetrics []MetricRow        `json:"activityDetailMetrics"`
	GeoPolylineDTO        struct {
		Polyline []GeoPoint `json:"polyline"`
	} `json:"geoPolylineDTO"`
}

type Point struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Elevation float64 `json:"elevation"`
	Timestamp string  `json:"timestamp"`
}

// ParseActivityTrackDetails converts Garmin activity-details JSON into
// normalized track points. Polyline is preferred; metric rows are a fallback.
func ParseActivityTrackDetails(body []byte) (Details, []Point, error) {
	var details Details
	if err := json.Unmarshal(body, &details); err != nil {
		return details, nil, err
	}

	points := make([]Point, 0, len(details.GeoPolylineDTO.Polyline))
	for _, point := range details.GeoPolylineDTO.Polyline {
		if point.Valid != nil && !*point.Valid {
			continue
		}
		if !validLatLon(point.Latitude, point.Longitude) {
			continue
		}
		points = append(points, Point{Latitude: point.Latitude, Longitude: point.Longitude, Elevation: point.Elevation})
	}
	if len(points) > 0 {
		return details, points, nil
	}

	indices := map[string]int{}
	for _, descriptor := range details.MetricDescriptors {
		key := strings.ToLower(descriptor.Key)
		switch {
		case strings.Contains(key, "latitude"):
			indices["latitude"] = descriptor.MetricsIndex
		case strings.Contains(key, "longitude"):
			indices["longitude"] = descriptor.MetricsIndex
		case strings.Contains(key, "elevation") || strings.Contains(key, "altitude"):
			indices["elevation"] = descriptor.MetricsIndex
		}
	}
	latIndex, hasLat := indices["latitude"]
	lonIndex, hasLon := indices["longitude"]
	if !hasLat || !hasLon {
		return details, points, nil
	}

	for _, row := range details.ActivityDetailMetrics {
		if latIndex < 0 || lonIndex < 0 || latIndex >= len(row.Metrics) || lonIndex >= len(row.Metrics) {
			continue
		}
		lat, latOK := jsonNumber(row.Metrics[latIndex])
		lon, lonOK := jsonNumber(row.Metrics[lonIndex])
		if !latOK || !lonOK || !validLatLon(lat, lon) {
			continue
		}
		point := Point{Latitude: lat, Longitude: lon}
		if elevationIndex, ok := indices["elevation"]; ok && elevationIndex >= 0 && elevationIndex < len(row.Metrics) {
			point.Elevation, _ = jsonNumber(row.Metrics[elevationIndex])
		}
		points = append(points, point)
	}
	return details, points, nil
}

// validLatLon rejects non-finite and out-of-range coordinates
func validLatLon(lat, lon float64) bool {
	if math.IsNaN(lat) || math.IsInf(lat, 0) || math.IsNaN(lon) || math.IsInf(lon, 0) {
		return false
	}
	return lat >= -90 && lat <= 90 && lon >= -180 && lon <= 180
}

func jsonNumber(raw json.RawMessage) (float64, bool) {
	if len(raw) == 0 || bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
		return 0, false
	}
	var value float64
	if err := json.Unmarshal(raw, &value); err != nil {
		return 0, false
	}
	return value, true
}
