//go:build tinygo

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/open-wanderer/wanderer/plugins/sdk"
)

const (
	garminBaseURL                = "https://personal.connectapi.garmin.com"
	garminTokenPath              = "/personalapitoken-service/personal-api/personal-api-key/token"
	garminDownloadMaxBytes int64 = 16 * 1024 * 1024
	defaultPageSize              = 50
)

// garminClient talks to the Garmin Personal API. The user's long-lived API key
// (GCDC_...) is exchanged for a short-lived (24h) JWT access token, which is
// then used as a Bearer token on all data endpoints.
type garminClient struct {
	accessToken string
	baseURL     string
}

// newGarminClient builds a client from the auth map. Prefers an existing
// accessToken; otherwise exchanges the apiKey for an access token.
func newGarminClient(auth map[string]any) (*garminClient, error) {
	if token := sdk.StringField(auth, "accessToken"); token != "" {
		return &garminClient{accessToken: token, baseURL: garminBaseURL}, nil
	}

	apiKey := sdk.StringField(auth, "apiKey")
	if apiKey == "" {
		return nil, fmt.Errorf("accessToken or apiKey is required")
	}

	token, err := exchangeAPIKey(apiKey)
	if err != nil {
		return nil, err
	}
	return &garminClient{accessToken: token, baseURL: garminBaseURL}, nil
}

// exchangeAPIKey trades the Personal API key for a JWT access token.
// The token endpoint returns the JWT as the plain response body.
func exchangeAPIKey(apiKey string) (string, error) {
	_, body, err := request("POST", garminTokenPath, nil, map[string]string{
		"Authorization": "ApiKey " + apiKey,
		"Content-Type":  "application/json",
		"Accept":        "text/plain, application/json, */*",
	}, "", []string{
		"text/plain",
		"application/json",
		"application/octet-stream",
	})
	if err != nil {
		return "", fmt.Errorf("token exchange request failed: %w", err)
	}
	token := strings.TrimSpace(string(body))
	if token == "" {
		return "", fmt.Errorf("token exchange returned an empty token")
	}
	return token, nil
}

// searchActivities fetches one page of activities. The endpoint returns a
// bare JSON array.
func (c *garminClient) searchActivities(page int, limit int, after time.Time) ([]garminActivity, error) {
	query := []sdk.QueryParam{
		{Name: "start", Value: strconv.Itoa((page - 1) * limit)},
		{Name: "limit", Value: strconv.Itoa(limit)},
	}
	if !after.IsZero() {
		query = append(query, sdk.QueryParam{Name: "startDate", Value: after.Format("2006-01-02")})
	}

	_, body, err := c.request("GET", "/activitylist-service/activities/search/activities", query, nil, []string{"application/json"})
	if err != nil {
		return nil, err
	}

	var activities []garminActivity
	if err := json.Unmarshal(body, &activities); err != nil {
		return nil, fmt.Errorf("failed to parse activities: %w", err)
	}
	return activities, nil
}

// getActivity fetches /activity-service/activity/{id}
func (c *garminClient) getActivity(activityID int64) (*garminActivityDetail, error) {
	var detail garminActivityDetail
	_, body, err := c.request("GET", fmt.Sprintf("/activity-service/activity/%d", activityID), nil, nil, []string{"application/json"})
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &detail); err != nil {
		return nil, fmt.Errorf("failed to parse activity detail: %w", err)
	}
	return &detail, nil
}

// getActivityTrack fetches /activity-service/activity/{id}/details and returns
// the geo polyline track points.
func (c *garminClient) getActivityTrack(activityID int64) ([]garminPolyPoint, error) {
	var details garminActivityDetails
	_, body, err := c.request("GET", fmt.Sprintf("/activity-service/activity/%d/details", activityID), nil, nil, []string{"application/json"})
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &details); err != nil {
		return nil, fmt.Errorf("failed to parse activity track: %w", err)
	}
	return details.GeoPolyline.Polyline, nil
}

// request performs an authenticated request against the "api" connector.
func (c *garminClient) request(method, path string, query []sdk.QueryParam, extraHeaders map[string]string, contentTypes []string) (sdk.HostResponse, []byte, error) {
	return request(method, path, query, extraHeaders, c.accessToken, contentTypes)
}

func request(method, path string, query []sdk.QueryParam, extraHeaders map[string]string, bearer string, contentTypes []string) (sdk.HostResponse, []byte, error) {
	headers := map[string]string{
		sdk.AuthHeaderAuthorization: sdk.AuthSchemeBearer + " " + bearer,
		"Accept":                    "application/json",
	}
	if len(contentTypes) > 0 {
		headers["Accept"] = contentTypes[0]
	}
	for k, v := range extraHeaders {
		headers[k] = v
	}

	var body *sdk.HostRequestBody
	if method != "GET" && method != "HEAD" {
		body = &sdk.HostRequestBody{Type: sdk.HostRequestBodyTypeJSON, JSON: map[string]any{}}
	}

	return sdk.HostRequest(sdk.HostRequestSpec{
		Method: method,
		Target: sdk.RequestTarget{
			Type:      "connector",
			Connector: "api",
			Path:      path,
			Query:     query,
		},
		Headers: headers,
		Body:    body,
		Expect: sdk.ResponseExpect{
			ContentTypes: contentTypes,
			MaxBytes:     garminDownloadMaxBytes,
		},
	})
}

// syncActivities lists activities page by page.
func syncActivities(client *garminClient, input listInput) (listOutput, error) {
	page := sdk.IntState(input.State, "page", 1)
	if page <= 0 {
		page = 1
	}

	limit := sdk.SyncLimit(input)
	if limit <= 0 {
		limit = defaultPageSize
	}

	after := parseAfterDate(input.Options)

	activities, err := client.searchActivities(page, limit, after)
	if err != nil {
		return listOutput{}, err
	}

	items := make([]trailSummary, 0, len(activities))
	for _, activity := range activities {
		// Client-side date filter (the API startDate param is inclusive by day).
		if !after.IsZero() && activity.BeginTimestamp > 0 {
			if time.UnixMilli(activity.BeginTimestamp).Before(after) {
				continue
			}
		}

		items = append(items, trailSummary{
			Source: trailImportSource{
				Provider:   "garmin",
				ExternalID: strconv.FormatInt(activity.ActivityID, 10),
			},
			Kind: "completed",
		})

		if len(items) >= limit {
			break
		}
	}

	hasMore := len(activities) >= limit
	return listOutput{
		Items:   items,
		State:   sdk.NextPageState(page+1, hasMore),
		HasMore: hasMore,
	}, nil
}

// syncRoutes is a no-op: the Garmin Personal API has no separate routes endpoint.
func syncRoutes(client *garminClient, input listInput) (listOutput, error) {
	return listOutput{
		Items:   []trailSummary{},
		State:   map[string]any{},
		HasMore: false,
	}, nil
}

// parseAfterDate parses the "after" option (YYYY-MM-DD or RFC3339).
func parseAfterDate(options map[string]any) time.Time {
	afterStr := sdk.StringOption(options, "after")
	if afterStr == "" {
		return time.Time{}
	}
	for _, layout := range []string{"2006-01-02", time.RFC3339, "2006-01-02T15:04:05.000Z", "2006-01-02T15:04:05.0"} {
		if date, err := time.Parse(layout, afterStr); err == nil {
			return date
		}
	}
	return time.Time{}
}

// getActivityDetail fetches detail + track and maps them to a trailImport.
func getActivityDetail(client *garminClient, input detailInput, kind string) (trailImport, error) {
	activityID, err := strconv.ParseInt(input.Summary.Source.ExternalID, 10, 64)
	if err != nil {
		return trailImport{}, fmt.Errorf("invalid activity ID: %w", err)
	}

	detail, err := client.getActivity(activityID)
	if err != nil {
		return trailImport{}, fmt.Errorf("failed to get activity detail: %w", err)
	}

	points, err := client.getActivityTrack(activityID)
	if err != nil {
		sdk.LogWarn(fmt.Sprintf("Failed to get activity track, importing metadata only: %v", err))
		points = nil
	}

	return mapActivityToTrail(detail, points), nil
}

// mapActivityToTrail maps a Garmin activity detail + track points to a trail import.
func mapActivityToTrail(activity *garminActivityDetail, points []garminPolyPoint) trailImport {
	privacy := mapPrivacy(activity.Privacy.TypeKey)

	var activityTrack track
	if len(points) > 0 {
		gpx := buildGPXFromPoints(points)
		activityTrack = track{
			Format:        "gpx",
			ContentBase64: base64.StdEncoding.EncodeToString(gpx),
		}
	}

	s := activity.Summary
	metadata := map[string]any{
		"distance":         s.Distance,
		"elevationGain":    s.ElevationGain,
		"elevationLoss":    s.ElevationLoss,
		"duration":         s.Duration,
		"movingDuration":   s.MovingDuration,
		"calories":         s.Calories,
		"averageSpeed":     s.AverageSpeed,
		"maxSpeed":         s.MaxSpeed,
		"averageHR":        s.AverageHR,
		"maxHR":            s.MaxHR,
		"location":         activity.LocationName,
		"providerCategory": activity.ActivityType.TypeKey,
		"fileFormat":       activity.Metadata.FileFormat.FormatKey,
	}

	return trailImport{
		Source: trailImportSource{
			Provider:   "garmin",
			ExternalID: strconv.FormatInt(activity.ActivityID, 10),
			URL:        fmt.Sprintf("https://connect.garmin.com/modern/activity/%d", activity.ActivityID),
		},
		Kind:         "completed",
		Name:         strings.TrimSpace(activity.ActivityName),
		Description:  activity.Description,
		StartedAt:    normalizeTime(s.StartTimeGMT),
		ActivityType: mapActivityType(activity.ActivityType.TypeKey),
		Privacy:      &privacy,
		Track:        activityTrack,
		Metadata:     metadata,
	}
}

// normalizeTime converts Garmin timestamps ("2026-10-12T16:09:45.0" or
// "2026-10-12 16:09:45", GMT) to RFC3339 UTC.
func normalizeTime(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	for _, layout := range []string{"2006-01-02T15:04:05.999", "2006-01-02T15:04:05.000Z07:00", time.RFC3339, "2006-01-02 15:04:05"} {
		if t, err := time.Parse(layout, value); err == nil {
			return t.UTC().Format(time.RFC3339)
		}
	}
	return value
}

// mapPrivacy maps Garmin accessControlRule typeKey to wanderer privacy.
func mapPrivacy(garminPrivacy string) string {
	switch strings.ToLower(garminPrivacy) {
	case "public":
		return "public"
	case "friends", "friends_only":
		return "friends"
	default:
		return "private"
	}
}

// mapActivityType maps a Garmin activity typeKey to the wanderer activity type.
func mapActivityType(garminType string) string {
	switch strings.ToLower(garminType) {
	case "running":
		return "running"
	case "trail_running":
		return "trail_running"
	case "cycling", "biking":
		return "cycling"
	case "mountain_biking":
		return "mountain_biking"
	case "gravel_cycling":
		return "gravel_cycling"
	case "hiking":
		return "hiking"
	case "walking":
		return "walking"
	case "swimming":
		return "swimming"
	case "skiing":
		return "skiing"
	case "snowboarding":
		return "snowboarding"
	case "alpine_skiing":
		return "alpine_skiing"
	case "backcountry_skiing":
		return "backcountry_skiing"
	case "cross_country_skiing":
		return "cross_country_skiing"
	case "kayaking":
		return "kayaking"
	case "rowing":
		return "rowing"
	case "paddle_boarding":
		return "paddle_boarding"
	case "climbing":
		return "climbing"
	case "rock_climbing":
		return "rock_climbing"
	case "ice_climbing":
		return "ice_climbing"
	case "e_bike_mountain", "e_bike_road", "e_bike_cargo":
		return garminType
	default:
		return "other"
	}
}

// buildGPXFromPoints builds a GPX 1.1 document from polyline track points.
func buildGPXFromPoints(points []garminPolyPoint) []byte {
	if len(points) == 0 {
		return nil
	}

	var gpx strings.Builder
	gpx.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<gpx version="1.1" creator="Garmin Personal API" xmlns="http://www.topografix.com/GPX/1/1">
  <trk>
    <name>Activity Track</name>
    <trkseg>
`)

	for _, point := range points {
		gpx.WriteString(fmt.Sprintf("      <trkpt lat=\"%.7f\" lon=\"%.7f\">\n", point.Lat, point.Lon))
		if point.Altitude != nil {
			gpx.WriteString(fmt.Sprintf("        <ele>%.1f</ele>\n", *point.Altitude))
		}
		if point.TimeMs > 0 {
			gpx.WriteString(fmt.Sprintf("        <time>%s</time>\n", time.UnixMilli(point.TimeMs).UTC().Format(time.RFC3339)))
		}
		gpx.WriteString("      </trkpt>\n")
	}

	gpx.WriteString(`    </trkseg>
  </trk>
</gpx>
`)

	return []byte(gpx.String())
}
