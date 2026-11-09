//go:build tinygo

package main

import (
	"github.com/open-wanderer/wanderer/plugins/sdk"
)

// Input types for list capabilities
type listInput struct {
	Instance sdk.InstanceRef `json:"instance"`
	Auth     map[string]any  `json:"auth,omitempty"`
	State    map[string]any  `json:"state,omitempty"`
	Options  map[string]any  `json:"options,omitempty"`
	Limits   sdk.SyncLimits  `json:"limits,omitempty"`
}

// Output types for list capabilities
type listOutput struct {
	Items   []trailSummary   `json:"items"`
	State   map[string]any   `json:"state,omitempty"`
	HasMore bool             `json:"hasMore"`
	Error   *sdk.PluginError `json:"error,omitempty"`
}

// Input types for detail capabilities
type detailInput struct {
	Instance sdk.InstanceRef `json:"instance"`
	Auth     map[string]any  `json:"auth,omitempty"`
	Options  map[string]any  `json:"options,omitempty"`
	Summary  trailSummary    `json:"summary"`
}

// Trail summary structure
type trailSummary struct {
	Source trailImportSource `json:"source"`
	Kind   string            `json:"kind,omitempty"`
}

// Trail import source
type trailImportSource struct {
	Provider   string `json:"provider"`
	ExternalID string `json:"externalId"`
	URL        string `json:"url,omitempty"`
}

// Trail import structure
type trailImport struct {
	Source       trailImportSource `json:"source"`
	Kind         string            `json:"kind,omitempty"`
	Name         string            `json:"name"`
	Description  string            `json:"description,omitempty"`
	StartedAt    string            `json:"startedAt,omitempty"`
	ActivityType string            `json:"activityType,omitempty"`
	Privacy      *string           `json:"privacy,omitempty"`
	Track        track             `json:"track"`
	Waypoints    []waypoint        `json:"waypoints,omitempty"`
	Photos       []photo           `json:"photos,omitempty"`
	Metadata     map[string]any    `json:"metadata,omitempty"`
	Difficulty   string            `json:"difficulty,omitempty"`
}

// Track structure
type track struct {
	Format        string `json:"format"`
	ContentBase64 string `json:"contentBase64"`
}

// Waypoint structure
type waypoint struct {
	ExternalID  string   `json:"externalId,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Lat         float64  `json:"lat"`
	Lon         float64  `json:"lon"`
	Ele         *float64 `json:"ele,omitempty"`
	Icon        string   `json:"icon,omitempty"`
	Photos      []photo  `json:"photos,omitempty"`
}

// Photo structure
type photo struct {
	ExternalID  string      `json:"externalId,omitempty"`
	Filename    string      `json:"filename,omitempty"`
	ContentType string      `json:"contentType,omitempty"`
	Lat         *float64    `json:"lat,omitempty"`
	Lon         *float64    `json:"lon,omitempty"`
	Source      mediaSource `json:"source"`
}

// Media source structure
type mediaSource struct {
	Type     string    `json:"type"`
	URL      string    `json:"url,omitempty"`
	MediaRef *mediaRef `json:"mediaRef,omitempty"`
}

// Media reference structure
type mediaRef struct {
	Connector string       `json:"connector"`
	Auth      string       `json:"auth,omitempty"`
	Path      string       `json:"path,omitempty"`
	Query     []queryParam `json:"query,omitempty"`
	AssetID   string       `json:"assetId,omitempty"`
}

// Query parameter structure
type queryParam struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

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
