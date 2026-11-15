//go:build tinygo

package main

import (
	"encoding/json"

	"github.com/open-wanderer/wanderer/plugins/garmin/internal/garmintrack"
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

type garminActivityType struct {
	TypeKey string `json:"typeKey"`
}

func (activityType *garminActivityType) UnmarshalJSON(data []byte) error {
	var key string
	if err := json.Unmarshal(data, &key); err == nil {
		activityType.TypeKey = key
		return nil
	}

	var value struct {
		TypeKey string `json:"typeKey"`
	}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	activityType.TypeKey = value.TypeKey
	return nil
}

type garminPrivacy struct {
	TypeKey string `json:"typeKey"`
}

func (privacy *garminPrivacy) UnmarshalJSON(data []byte) error {
	var key string
	if err := json.Unmarshal(data, &key); err == nil {
		privacy.TypeKey = key
		return nil
	}

	var value struct {
		TypeKey string `json:"typeKey"`
	}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	privacy.TypeKey = value.TypeKey
	return nil
}

type garminActivity struct {
	ActivityID     int64              `json:"activityId"`
	ActivityName   string             `json:"activityName"`
	StartTime      string             `json:"startTimeGMT"`
	StartTimeLocal string             `json:"startTimeLocal"`
	Duration       float64            `json:"duration"`
	Distance       float64            `json:"distance"`
	ElevationGain  float64            `json:"elevationGain"`
	ElevationLoss  float64            `json:"elevationLoss"`
	AverageSpeed   float64            `json:"averageSpeed"`
	MaxSpeed       float64            `json:"maxSpeed"`
	Calories       float64            `json:"calories"`
	ActivityType   garminActivityType `json:"activityType"`
	StartLat       float64            `json:"startLatitude"`
	StartLon       float64            `json:"startLongitude"`
	EndLat         float64            `json:"endLatitude"`
	EndLon         float64            `json:"endLongitude"`
	HasTrack       bool               `json:"hasTrack"`
	HasPhotos      bool               `json:"hasPhotos"`
	Privacy        garminPrivacy      `json:"privacy"`
}

type garminActivityTrackDetails = garmintrack.Details

type garminActivityDetail struct {
	ActivityID    int64              `json:"activityId"`
	ActivityName  string             `json:"activityName"`
	Description   string             `json:"description"`
	StartTime     string             `json:"startTimeGMT"`
	Duration      float64            `json:"duration"`
	Distance      float64            `json:"distance"`
	ElevationGain float64            `json:"elevationGain"`
	ElevationLoss float64            `json:"elevationLoss"`
	AverageSpeed  float64            `json:"averageSpeed"`
	MaxSpeed      float64            `json:"maxSpeed"`
	Calories      float64            `json:"calories"`
	ActivityType  garminActivityType `json:"activityType"`
	StartLat      float64            `json:"startLatitude"`
	StartLon      float64            `json:"startLongitude"`
	EndLat        float64            `json:"endLatitude"`
	EndLon        float64            `json:"endLongitude"`
	Privacy       garminPrivacy      `json:"privacy"`
	Difficulty    string             `json:"difficultyRating"`
	Location      string             `json:"locationName"`
	Photos        []garminPhoto      `json:"photos"`
	Track         []trackPoint       `json:"trackPoints"`
}

type garminPhoto struct {
	PhotoID   int64   `json:"photoId"`
	PhotoURL  string  `json:"photoUrl"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp string  `json:"timestamp"`
}

type trackPoint = garmintrack.Point
