//go:build tinygo

package main

import (
	"github.com/open-wanderer/wanderer/plugins/garmin/internal/garminapi"
	"github.com/open-wanderer/wanderer/plugins/garmin/internal/garmintrack"
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
type garminActivitySearch = garminapi.ActivitySearch
type garminActivity = garminapi.Activity
type garminActivityDetail = garminapi.ActivityDetail
type garminPhoto = garminapi.Photo
type garminActivityType = garminapi.ActivityType
type garminPrivacy = garminapi.Privacy

type garminActivityTrackDetails = garmintrack.Details
type trackPoint = garmintrack.Point
