//go:build tinygo

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/open-wanderer/wanderer/plugins/garmin/internal/garminutil"
	"github.com/open-wanderer/wanderer/plugins/sdk"
)

func syncActivities(client *garminClient, input listInput) (listOutput, error) {
	page := sdk.IntState(input.State, "page", 1)
	if page <= 0 {
		page = 1
	}

	limit := sdk.SyncLimit(input)
	if limit <= 0 {
		limit = defaultPageSize
	}

	// Parse after date from options
	after, err := garminutil.ParseAfterDate(input.Options)
	if err != nil {
		logHelper(sdk.LogLevelError, "Garmin: "+err.Error())
		return listOutput{}, err
	}

	activities, err := client.searchActivities(page, limit, after)
	if err != nil {
		return listOutput{}, err
	}

	items := make([]trailSummary, 0, len(activities))
	filteredByDate := 0
	invalidStartTime := 0
	for _, activity := range activities {
		// Filter by after date if not already filtered by API
		if !after.IsZero() {
			startTime, parseErr := garminutil.ParseActivityStartTime(activity.StartTime)
			if parseErr != nil && activity.StartTimeLocal != "" {
				startTime, parseErr = garminutil.ParseActivityStartTime(activity.StartTimeLocal)
			}
			if parseErr != nil {
				invalidStartTime++
				logHelper(sdk.LogLevelWarn, fmt.Sprintf("Garmin: skipping activity %d: unparseable start time %q / %q: %v", activity.ActivityID, activity.StartTime, activity.StartTimeLocal, parseErr))
				continue
			}
			if startTime.Before(after) {
				filteredByDate++
				continue
			}
		}

		externalID := strconv.FormatInt(activity.ActivityID, 10)
		items = append(items, trailSummary{
			Source: trailImportSource{
				Provider:   "garmin",
				ExternalID: externalID,
			},
			Kind: "completed",
		})

		if len(items) >= limit {
			break
		}
	}
	logHelper(sdk.LogLevelInfo, fmt.Sprintf("Garmin: activities received=%d eligible=%d filtered_by_date=%d invalid_start_time=%d", len(activities), len(items), filteredByDate, invalidStartTime))
	logActivitySummaries(items, page, len(activities), len(activities) >= limit)

	nextPage := page + 1
	hasMore := len(activities) >= limit

	return listOutput{
		Items:   items,
		State:   sdk.NextPageState(nextPage, hasMore),
		HasMore: hasMore,
	}, nil
}

// Garmin Personal API has no separate routes endpoint; syncRoutes returns an empty list
// but we can treat saved activities as routes
func syncRoutes(client *garminClient, input listInput) (listOutput, error) {
	// For Garmin Personal API, we'll use the same activity search
	// but filter for activities that might be routes
	// Note: Garmin Personal API doesn't distinguish between activities and routes
	// in the same way as Strava. We'll treat all activities as completed for now.
	// If Garmin adds route support, we can update this.

	// For now, return empty list for routes
	// Garmin Personal API doesn't have a separate routes endpoint like Strava
	return listOutput{
		Items:   []trailSummary{},
		State:   map[string]any{},
		HasMore: false,
	}, nil
}

func getActivityDetail(client *garminClient, input detailInput, kind string) (trailImport, error) {
	logHelper(sdk.LogLevelInfo, fmt.Sprintf("Garmin: detail request received kind=%q provider=%q external_id=%q summary_kind=%q", kind, input.Summary.Source.Provider, input.Summary.Source.ExternalID, input.Summary.Kind))
	activityID, err := strconv.ParseInt(input.Summary.Source.ExternalID, 10, 64)
	if err != nil {
		return trailImport{}, fmt.Errorf("invalid activity ID: %w", err)
	}

	// Always fetch activity details directly from API (no caching to prevent cross-account issues)
	detail, err := client.getActivity(activityID)
	if err != nil {
		return trailImport{}, fmt.Errorf("failed to get activity detail: %w", err)
	}
	points, err := client.getActivityTrackPoints(activityID)
	if err != nil {
		return trailImport{}, fmt.Errorf("failed to get activity track: %w", err)
	}
	if len(points) == 0 {
		logHelper(sdk.LogLevelWarn, "Garmin: activity track details contained no usable GPS points")
		return trailImport{}, fmt.Errorf("Garmin activity has no GPS track points")
	}
	detail.Track = points
	item, err := mapActivityToTrail(detail, nil)
	if err == nil {
		logActivityImportItem(item, points)
	}
	return item, err
}

func logActivitySummaries(items []trailSummary, page, received int, hasMore bool) {
	logHelper(sdk.LogLevelInfo, fmt.Sprintf("Garmin: list capability output page=%d received=%d returned_items=%d has_more=%t", page, received, len(items), hasMore))
}

func logActivityImportItem(item trailImport, points []trackPoint) {
	gpxBytes, decodeErr := base64.StdEncoding.DecodeString(item.Track.ContentBase64)
	if decodeErr != nil {
		logHelper(sdk.LogLevelWarn, fmt.Sprintf("Garmin: unable to measure generated GPX activity_id=%q started_at=%q points=%d: %v", item.Source.ExternalID, item.StartedAt, len(points), decodeErr))
		return
	}
	logHelper(sdk.LogLevelInfo, fmt.Sprintf("Garmin: activity ready for import id=%q started_at=%q points=%d gpx_bytes=%d gpx_base64_bytes=%d", item.Source.ExternalID, item.StartedAt, len(points), len(gpxBytes), len(item.Track.ContentBase64)))
}

func logActivityTrackShape(body []byte) {
	var top map[string]json.RawMessage
	if err := json.Unmarshal(body, &top); err != nil {
		logHelper(sdk.LogLevelWarn, "Garmin: track details response was not a JSON object")
		return
	}
	var geo map[string]json.RawMessage
	if err := json.Unmarshal(top["geoPolylineDTO"], &geo); err != nil {
		logHelper(sdk.LogLevelWarn, "Garmin: track details response had no readable geoPolylineDTO")
		return
	}
	var polyline []json.RawMessage
	_ = json.Unmarshal(geo["polyline"], &polyline)
	firstPointShape := "empty"
	if len(polyline) > 0 {
		firstPointShape = jsonValueShape(polyline[0], 0)
	}
	logHelper(sdk.LogLevelInfo, "Garmin: track details field shape: "+jsonValueShape(top["geoPolylineDTO"], -1)+"; first polyline point: "+firstPointShape)
}
