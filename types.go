//go:build tinygo

package main

import (
	"github.com/open-wanderer/wanderer/plugins/sdk"
)

// Input types for list capabilities - use SDK types directly
type listInput = sdk.ListInput
type listOutput = sdk.ListOutput

// Input types for detail capabilities - use SDK types directly
type detailInput = sdk.DetailInput

// Trail summary structure - use SDK type
type trailSummary = sdk.TrailSummary

// Trail import source - use SDK type
type trailImportSource = sdk.TrailImportSource

// Trail import structure - use SDK type
type trailImport = sdk.TrailImport

// Track structure - use SDK type
type track = sdk.Track

// Waypoint structure - use SDK type
type waypoint = sdk.Waypoint

// Photo structure - use SDK type
type photo = sdk.Photo

// Media source structure - use SDK type
type mediaSource = sdk.MediaSource

// Media reference structure - use SDK type
type mediaRef = sdk.MediaRef

// Query parameter structure - use SDK type
type queryParam = sdk.QueryParam

// Garmin API response types
type garminActivitySearch struct {
	TotalCount int              `json:"totalCount"`
	Start      int              `json:"start"`
	Limit      int              `json:"limit"`
	Activities []garminActivity `json:"activities"`
}

type garminActivity struct {
	ActivityID    int64   `json:"activityId"`
	ActivityName  string  `json:"activityName"`
	StartTime     string  `json:"startTime"`
	Duration      float64 `json:"duration"`
	Distance      float64 `json:"distance"`
	ElevationGain float64 `json:"elevationGain"`
	ElevationLoss float64 `json:"elevationLoss"`
	AverageSpeed  float64 `json:"averageSpeed"`
	MaxSpeed      float64 `json:"maxSpeed"`
	Calories      int     `json:"calories"`
	ActivityType  string  `json:"activityType"`
	StartLat      float64 `json:"startLatitude"`
	StartLon      float64 `json:"startLongitude"`
	EndLat        float64 `json:"endLatitude"`
	EndLon        float64 `json:"endLongitude"`
	HasTrack      bool    `json:"hasTrack"`
	HasPhotos     bool    `json:"hasPhotos"`
	Privacy       string  `json:"privacy"`
}

type garminActivityDetail struct {
	ActivityID    int64         `json:"activityId"`
	ActivityName  string        `json:"activityName"`
	Description   string        `json:"description"`
	StartTime     string        `json:"startTime"`
	Duration      float64       `json:"duration"`
	Distance      float64       `json:"distance"`
	ElevationGain float64       `json:"elevationGain"`
	ElevationLoss float64       `json:"elevationLoss"`
	AverageSpeed  float64       `json:"averageSpeed"`
	MaxSpeed      float64       `json:"maxSpeed"`
	Calories      int           `json:"calories"`
	ActivityType  string        `json:"activityType"`
	StartLat      float64       `json:"startLatitude"`
	StartLon      float64       `json:"startLongitude"`
	EndLat        float64       `json:"endLatitude"`
	EndLon        float64       `json:"endLongitude"`
	Privacy       string        `json:"privacy"`
	Difficulty    string        `json:"difficultyRating"`
	Location      string        `json:"locationName"`
	Photos        []garminPhoto `json:"photos"`
	Track         []trackPoint  `json:"trackPoints"`
}

type garminPhoto struct {
	PhotoID   int64   `json:"photoId"`
	PhotoURL  string  `json:"photoUrl"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp string  `json:"timestamp"`
}

type trackPoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Elevation float64 `json:"elevation"`
	Timestamp string  `json:"timestamp"`
}
