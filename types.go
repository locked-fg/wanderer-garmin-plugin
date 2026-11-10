//go:build tinygo

package main

import (
	"github.com/open-wanderer/wanderer/plugins/sdk"
)

type listInput = sdk.ListInput
type listOutput = sdk.ListOutput
type detailInput = sdk.DetailInput
type trailSummary = sdk.TrailSummary
type trailImport = sdk.TrailImport
type trailImportSource = sdk.TrailImportSource
type track = sdk.Track
type waypoint = sdk.Waypoint
type photo = sdk.Photo
type mediaSource = sdk.MediaSource
type mediaRef = sdk.MediaRef
type queryParam = sdk.QueryParam
type pluginError = sdk.PluginError

// Garmin API response types (verified against personal.connectapi.garmin.com)

// garminTypeRef is Garmin's {typeId, typeKey} reference object used for
// activityType, eventType, privacy (accessControlRuleDTO), fileFormat, etc.
type garminTypeRef struct {
	TypeID  int    `json:"typeId"`
	TypeKey string `json:"typeKey"`
}

// Activity list item from /activitylist-service/activities/search/activities
// (the endpoint returns a bare JSON array).
type garminActivity struct {
	ActivityID     int64         `json:"activityId"`
	ActivityName   string        `json:"activityName"`
	StartTimeGMT   string        `json:"startTimeGMT"`
	StartTimeLocal string        `json:"startTimeLocal"`
	BeginTimestamp int64         `json:"beginTimestamp"` // ms epoch
	Duration       float64       `json:"duration"`
	Distance       float64       `json:"distance"`
	ElevationGain  float64       `json:"elevationGain"`
	ElevationLoss  float64       `json:"elevationLoss"`
	AverageSpeed   float64       `json:"averageSpeed"`
	MaxSpeed       float64       `json:"maxSpeed"`
	Calories       float64       `json:"calories"`
	ActivityType   garminTypeRef `json:"activityType"`
	EventType      garminTypeRef `json:"eventType"`
	StartLat       float64       `json:"startLatitude"`
	StartLon       float64       `json:"startLongitude"`
	EndLat         float64       `json:"endLatitude"`
	EndLon         float64       `json:"endLongitude"`
	HasImages      bool          `json:"hasImages"`
	Privacy        garminTypeRef `json:"privacy"`
	LocationName   string        `json:"locationName"`
}

// Activity detail from /activity-service/activity/{id}
type garminActivityDetail struct {
	ActivityID   int64          `json:"activityId"`
	ActivityName string         `json:"activityName"`
	Description  string         `json:"description"`
	LocationName string         `json:"locationName"`
	ActivityType garminTypeRef  `json:"activityTypeDTO"`
	Privacy      garminTypeRef  `json:"accessControlRuleDTO"`
	Summary      garminSummary  `json:"summaryDTO"`
	Metadata     garminMetadata `json:"metadataDTO"`
}

type garminSummary struct {
	StartTimeGMT    string  `json:"startTimeGMT"`
	StartTimeLocal  string  `json:"startTimeLocal"`
	StartLat        float64 `json:"startLatitude"`
	StartLon        float64 `json:"startLongitude"`
	Distance        float64 `json:"distance"`
	Duration        float64 `json:"duration"`
	MovingDuration  float64 `json:"movingDuration"`
	ElapsedDuration float64 `json:"elapsedDuration"`
	ElevationGain   float64 `json:"elevationGain"`
	ElevationLoss   float64 `json:"elevationLoss"`
	AverageSpeed    float64 `json:"averageSpeed"`
	MaxSpeed        float64 `json:"maxSpeed"`
	Calories        float64 `json:"calories"`
	AverageHR       float64 `json:"averageHR"`
	MaxHR           float64 `json:"maxHR"`
}

type garminFileFormat struct {
	FormatID  int    `json:"formatId"`
	FormatKey string `json:"formatKey"`
}

type garminMetadata struct {
	FileFormat garminFileFormat `json:"fileFormat"`
}

// Track points from /activity-service/activity/{id}/details (geoPolylineDTO)
type garminActivityDetails struct {
	GeoPolyline garminGeoPolyline `json:"geoPolylineDTO"`
}

type garminGeoPolyline struct {
	Polyline []garminPolyPoint `json:"polyline"`
}

type garminPolyPoint struct {
	Lat       float64  `json:"lat"`
	Lon       float64  `json:"lon"`
	Altitude  *float64 `json:"altitude"`
	TimeMs    int64    `json:"time"`
}
