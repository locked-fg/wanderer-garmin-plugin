//go:build tinygo

package main

import (
	"github.com/open-wanderer/wanderer/plugins/garmin/internal/garmintrack"
	"github.com/open-wanderer/wanderer/plugins/garmin/internal/garminutil"
	"github.com/open-wanderer/wanderer/plugins/sdk"
)

// Plugin-facing aliases for SDK types
type listInput = sdk.ListInput
type listOutput = sdk.ListOutput

type detailInput = sdk.DetailInput

type trailSummary = sdk.TrailSummary

type trailImportSource = sdk.TrailImportSource

type trailImport = sdk.TrailImport

type track = sdk.Track

type waypoint = sdk.Waypoint

type photo = sdk.Photo

type mediaSource = sdk.MediaSource

type mediaRef = sdk.MediaRef

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
	key, err := garminutil.UnmarshalTypeKey(data)
	if err != nil {
		return err
	}
	activityType.TypeKey = key
	return nil
}

type garminPrivacy struct {
	TypeKey string `json:"typeKey"`
}

func (privacy *garminPrivacy) UnmarshalJSON(data []byte) error {
	key, err := garminutil.UnmarshalTypeKey(data)
	if err != nil {
		return err
	}
	privacy.TypeKey = key
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
