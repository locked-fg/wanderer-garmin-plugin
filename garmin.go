//go:build tinygo

package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/open-wanderer/wanderer/plugins/garmin/internal/garmintrack"
	"github.com/open-wanderer/wanderer/plugins/sdk"
)

const (
	garminBaseURL                = "https://personal.connectapi.garmin.com"
	garminDownloadMaxBytes int64 = 16 * 1024 * 1024
	defaultPageSize              = 50
)

// Garmin client structure
type garminClient struct {
	accessToken    string
	apiKey         string
	baseURL        string
	trackShapeLogged bool // Track if shape was logged for this client instance
}

// Token exchange response
type tokenExchangeResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// logHelper sends a log message to the host
func logHelper(level sdk.LogLevel, message string) {
	sdk.Log(level, message)
}

// Create a new Garmin client
func newGarminClient(auth map[string]any) (*garminClient, error) {
	// Try to get access token first (pre-exchanged)
	accessToken := sdk.StringField(auth, "accessToken")
	apiKey := sdk.StringField(auth, "apiKey")

	client := &garminClient{
		accessToken: accessToken,
		apiKey:      apiKey,
		baseURL:     garminBaseURL,
	}

	// If we have an API key but no access token, exchange it
	// Note: Token exchange happens per client instance (no global caching)
	// This ensures each sync has its own token and prevents cross-account issues
	if apiKey != "" && accessToken == "" {
		logHelper(sdk.LogLevelInfo, "Garmin: Exchanging API key for access token")
		token, err := client.exchangeAPIKeyForToken()
		if err != nil {
			logHelper(sdk.LogLevelError, fmt.Sprintf("Garmin: Token exchange failed: %v", err))
			return nil, fmt.Errorf("failed to exchange API key for token: %w", err)
		}
		client.accessToken = token
		logHelper(sdk.LogLevelInfo, "Garmin: Token received, waiting 15s for activation")
		// Give token time to activate (Garmin tokens may need 10-30 seconds)
		time.Sleep(15 * time.Second)
		logHelper(sdk.LogLevelInfo, "Garmin: Token activation complete")
	}

	// If we still don't have an access token, fail
	if client.accessToken == "" {
		logHelper(sdk.LogLevelError, "Garmin: No access token available")
		return nil, fmt.Errorf("accessToken or apiKey is required")
	}

	return client, nil
}

// Exchange API key for access token
func (c *garminClient) exchangeAPIKeyForToken() (string, error) {
	// Use the token_exchange connector without auth, manually set ApiKey header
	response, body, err := c.tokenExchangeRequest("/personalapitoken-service/personal-api/personal-api-key/token", []sdk.QueryParam{}, []string{"application/json", "text/plain"})
	if err != nil {
		return "", err
	}

	if response.Status < 200 || response.Status >= 300 {
		return "", fmt.Errorf("token exchange failed (%d): %s", response.Status, string(body))
	}

	var tokenResp tokenExchangeResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		// Try to parse as plain text (some APIs return just the token)
		token := strings.TrimSpace(string(body))
		if token != "" {
			return token, nil
		}
		return "", fmt.Errorf("failed to parse token response: %w", err)
	}

	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("no access token in response")
	}

	return tokenResp.AccessToken, nil
}

// Token exchange request using ApiKey header
// Uses token_exchange connector without auth to manually set ApiKey header
func (c *garminClient) tokenExchangeRequest(path string, query []sdk.QueryParam, contentTypes []string) (sdk.HostResponse, []byte, error) {
	return sdk.HostRequest(sdk.HostRequestSpec{
		Method: "POST",
		Target: sdk.RequestTarget{
			Type:      "connector",
			Connector: "token_exchange",
			Path:      path,
			Query:     query,
		},
		Headers: map[string]string{
			"Authorization": "ApiKey " + c.apiKey,
			// Garmin's token endpoint rejects Accept: application/json with 406.
			// Its documented request omits Accept; */* preserves that behavior
			// while the response is still validated as JSON by the host.
			"Accept": "*/*",
		},
		Expect: sdk.ResponseExpect{
			ContentTypes: contentTypes,
			MaxBytes:     garminDownloadMaxBytes,
		},
	})
}

// Search for activities
func (c *garminClient) searchActivities(page int, limit int, after time.Time) ([]garminActivity, error) {
	query := []sdk.QueryParam{
		{Name: "start", Value: strconv.Itoa((page - 1) * limit)},
		{Name: "limit", Value: strconv.Itoa(limit)},
	}

	// Add date filter if after is specified
	if !after.IsZero() {
		query = append(query, sdk.QueryParam{
			Name:  "startDate",
			Value: after.Format("2006-01-02"),
		})
	}

	response, body, err := c.request("/activitylist-service/activities/search/activities", query, []string{"application/json"})
	if err != nil {
		return nil, err
	}

	if response.Status < 200 || response.Status >= 300 {
		return nil, fmt.Errorf("Garmin API request failed (%d): %s", response.Status, string(body))
	}

	var activities []garminActivity
	trimmedBody := bytes.TrimSpace(body)
	if len(trimmedBody) == 0 {
		return nil, fmt.Errorf("failed to parse activities: empty response")
	}
	if trimmedBody[0] == '[' {
		if err := json.Unmarshal(trimmedBody, &activities); err != nil {
			logActivityResponseShape(trimmedBody)
			return nil, fmt.Errorf("failed to parse activities: %w", err)
		}
	} else {
		var result garminActivitySearch
		if err := json.Unmarshal(trimmedBody, &result); err != nil {
			return nil, fmt.Errorf("failed to parse activities: %w", err)
		}
		activities = result.Activities
	}
	return activities, nil
}

// logActivityResponseShape reports only JSON field names and value types.
// It intentionally omits activity values and identifiers from diagnostics.
func logActivityResponseShape(body []byte) {
	var activities []json.RawMessage
	if err := json.Unmarshal(body, &activities); err != nil {
		logHelper(sdk.LogLevelInfo, "Garmin: activity response is not a JSON array")
		return
	}
	if len(activities) == 0 {
		logHelper(sdk.LogLevelInfo, "Garmin: activity response shape is an empty array")
		return
	}
	logHelper(sdk.LogLevelInfo, "Garmin: activity response shape: "+jsonValueShape(activities[0], 0))
}

func jsonValueShape(raw json.RawMessage, depth int) string {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return "empty"
	}
	switch trimmed[0] {
	case '{':
		if depth >= 1 {
			return "object"
		}
		var object map[string]json.RawMessage
		if err := json.Unmarshal(trimmed, &object); err != nil {
			return "object"
		}
		keys := make([]string, 0, len(object))
		for key := range object {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		fields := make([]string, 0, len(keys))
		for _, key := range keys {
			fields = append(fields, key+":"+jsonValueShape(object[key], depth+1))
		}
		return "object{" + strings.Join(fields, ",") + "}"
	case '[':
		return "array"
	case '"':
		return "string"
	case 't', 'f':
		return "boolean"
	case 'n':
		return "null"
	default:
		return "number"
	}
}

// Get activity details
func (c *garminClient) getActivity(activityID int64) (*garminActivityDetail, error) {
	var detail garminActivityDetail

	path := fmt.Sprintf("/activity-service/activity/%d", activityID)
	response, body, err := c.request(path, nil, []string{"application/json"})
	if err != nil {
		return nil, err
	}

	if response.Status < 200 || response.Status >= 300 {
		return nil, fmt.Errorf("Garmin API request failed (%d): %s", response.Status, string(body))
	}

	if err := json.Unmarshal(body, &detail); err != nil {
		return nil, fmt.Errorf("failed to parse activity detail: %w", err)
	}

	return &detail, nil
}

func (c *garminClient) getActivityTrackPoints(activityID int64) ([]trackPoint, error) {
	path := fmt.Sprintf("/activity-service/activity/%d/details", activityID)
	query := []sdk.QueryParam{{Name: "maxPolylineSize", Value: "10000"}}
	response, body, err := c.request(path, query, []string{"application/json"})
	if err != nil {
		return nil, err
	}
	if response.Status < 200 || response.Status >= 300 {
		return nil, fmt.Errorf("Garmin activity track details failed (%d)", response.Status)
	}
	details, points, err := garmintrack.ParseActivityTrackDetails(body)
	if err != nil {
		if !c.trackShapeLogged {
			c.trackShapeLogged = true
			logActivityTrackShape(body)
		}
		return nil, fmt.Errorf("failed to parse Garmin activity track details: %w", err)
	}
	if !c.trackShapeLogged {
		c.trackShapeLogged = true
		descriptorKeys := make([]string, 0, len(details.MetricDescriptors))
		for _, descriptor := range details.MetricDescriptors {
			descriptorKeys = append(descriptorKeys, descriptor.Key)
		}
		logHelper(sdk.LogLevelInfo, fmt.Sprintf("Garmin: track details polyline_points=%d metric_rows=%d metric_keys=%s", len(details.GeoPolylineDTO.Polyline), len(details.ActivityDetailMetrics), strings.Join(descriptorKeys, ",")))
	}
	return points, nil
}

// Get activity workouts
func (c *garminClient) getActivityWorkouts(activityID int64) ([]byte, error) {
	path := fmt.Sprintf("/activity-service/activity/%d/workouts", activityID)
	response, body, err := c.request(path, nil, []string{"application/json"})
	if err != nil {
		return nil, err
	}

	if response.Status < 200 || response.Status >= 300 {
		return nil, fmt.Errorf("Garmin API request failed (%d): %s", response.Status, string(body))
	}

	return body, nil
}

// Generic request method with retry for token activation
func (c *garminClient) request(path string, query []sdk.QueryParam, contentTypes []string) (sdk.HostResponse, []byte, error) {
	accept := "application/json"
	if len(contentTypes) > 0 {
		accept = contentTypes[0]
	}
	if strings.HasPrefix(path, "/download-service/") {
		// Garmin's activity download endpoint rejects a specific ZIP Accept value.
		accept = "*/*"
	}

	// First attempt
	response, body, err := sdk.HostRequest(sdk.HostRequestSpec{
		Method: "GET",
		Target: sdk.RequestTarget{
			Type:      "connector",
			Connector: "api",
			Path:      path,
			Query:     query,
		},
		Headers: map[string]string{
			sdk.AuthHeaderAuthorization: sdk.AuthSchemeBearer + " " + c.accessToken,
			"Accept":                    accept,
		},
		Expect: sdk.ResponseExpect{
			ContentTypes: contentTypes,
			MaxBytes:     garminDownloadMaxBytes,
		},
	})

	// Check if token is not active - retry after additional delay
	if err == nil && response.Status == 401 {
		// Check if the error is "token is not active"
		bodyStr := string(body)
		if strings.Contains(strings.ToLower(bodyStr), "not active") {
			logHelper(sdk.LogLevelWarn, "Garmin: Token not active, waiting additional 15s and retrying...")
			time.Sleep(15 * time.Second)

			// Retry
			response, body, err = sdk.HostRequest(sdk.HostRequestSpec{
				Method: "GET",
				Target: sdk.RequestTarget{
					Type:      "connector",
					Connector: "api",
					Path:      path,
					Query:     query,
				},
				Headers: map[string]string{
					sdk.AuthHeaderAuthorization: sdk.AuthSchemeBearer + " " + c.accessToken,
					"Accept":                    accept,
				},
				Expect: sdk.ResponseExpect{
					ContentTypes: contentTypes,
					MaxBytes:     garminDownloadMaxBytes,
				},
			})

			if err == nil && response.Status == 200 {
				logHelper(sdk.LogLevelInfo, "Garmin: Request succeeded after retry")
			}
		}
	}

	return response, body, err
}

// Sync activities from Garmin
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
	after := parseAfterDate(input.Options)

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
			startTime, parseErr := parseActivityStartTime(activity.StartTime)
			if parseErr != nil && activity.StartTimeLocal != "" {
				startTime, parseErr = parseActivityStartTime(activity.StartTimeLocal)
			}
			if parseErr != nil {
				invalidStartTime++
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

// Sync routes - Garmin Personal API doesn't have a separate routes endpoint
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

// Parse after date from options
func parseAfterDate(options map[string]any) time.Time {
	if options == nil {
		return time.Time{}
	}

	afterStr, ok := options["after"].(string)
	if !ok || afterStr == "" {
		return time.Time{}
	}

	// Try to parse various date formats
	date, err := time.Parse("2006-01-02", afterStr)
	if err == nil {
		return date
	}

	date, err = time.Parse(time.RFC3339, afterStr)
	if err == nil {
		return date
	}

	date, err = time.Parse("2006-01-02T15:04:05.000Z", afterStr)
	if err == nil {
		return date
	}

	return time.Time{}
}

func parseActivityStartTime(value string) (time.Time, error) {
	for _, layout := range []string{
		time.RFC3339Nano,
		"2006-01-02T15:04:05.999999999-0700",
		"2006-01-02T15:04:05-0700",
		"2006-01-02T15:04:05.999999999",
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05.000",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05.999999999",
		"2006-01-02 15:04:05",
	} {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported Garmin activity start time format")
}

// Get activity detail and map to trail import
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

// Map Garmin activity to trail import
func mapActivityToTrail(activity *garminActivityDetail, gpxData []byte) (trailImport, error) {
	privacy := mapPrivacy(activity.Privacy.TypeKey)
	difficulty := mapDifficulty(activity.Difficulty)

	// Build track from GPX or track points
	var trackData track
	if gpxData != nil && len(gpxData) > 0 {
		trackData = track{
			Format:        "gpx",
			ContentBase64: base64.StdEncoding.EncodeToString(gpxData),
		}
	} else if len(activity.Track) > 0 {
		// Build GPX from track points
		gpxData := buildGPXFromTrackPoints(activity.Track)
		trackData = track{
			Format:        "gpx",
			ContentBase64: base64.StdEncoding.EncodeToString(gpxData),
		}
	}

	// Map activity type
	activityType := mapActivityType(activity.ActivityType.TypeKey)

	// Build metadata
	metadata := map[string]any{
		"distance":         activity.Distance,
		"elevationGain":    activity.ElevationGain,
		"elevationLoss":    activity.ElevationLoss,
		"duration":         activity.Duration,
		"calories":         activity.Calories,
		"averageSpeed":     activity.AverageSpeed,
		"maxSpeed":         activity.MaxSpeed,
		"providerCategory": activity.ActivityType.TypeKey,
	}

	// Garmin's startTimeGMT is often formatted as "2006-01-02 15:04:05".
	// Normalize it to RFC3339 because the host decodes startedAt as time.Time.
	startedAt := ""
	if activity.StartTime != "" {
		if parsed, err := parseActivityStartTime(activity.StartTime); err == nil {
			startedAt = parsed.UTC().Format(time.RFC3339)
		} else {
			logHelper(sdk.LogLevelWarn, fmt.Sprintf("Garmin: unable to parse startTimeGMT %q; omitting startedAt", activity.StartTime))
		}
	}

	return trailImport{
		Source: trailImportSource{
			Provider:   "garmin",
			ExternalID: strconv.FormatInt(activity.ActivityID, 10),
			URL:        fmt.Sprintf("https://connect.garmin.com/modern/activity/%d", activity.ActivityID),
		},
		Kind:         "completed",
		Name:         activity.ActivityName,
		Description:  activity.Description,
		StartedAt:    startedAt,
		ActivityType: activityType,
		Privacy:      &privacy,
		Track:        trackData,
		Waypoints:    []waypoint{}, // Garmin activities don't typically have waypoints
		Photos:       mapGarminPhotos(activity.Photos),
		Metadata:     metadata,
		Difficulty:   difficulty,
	}, nil
}

// Map Garmin privacy to wanderer privacy
func mapPrivacy(garminPrivacy string) string {
	// Garmin privacy values: private, public, friends, etc.
	switch strings.ToLower(garminPrivacy) {
	case "private":
		return "private"
	case "public":
		return "public"
	case "friends":
		return "friends"
	default:
		return "private"
	}
}

// Map Garmin difficulty to wanderer difficulty
func mapDifficulty(garminDifficulty string) string {
	// Garmin difficulty values: easy, moderate, hard, etc.
	switch strings.ToLower(garminDifficulty) {
	case "easy":
		return "easy"
	case "moderate":
		return "moderate"
	case "hard", "difficult":
		return "difficult"
	default:
		return ""
	}
}

// Map Garmin activity type to wanderer activity type
func mapActivityType(garminType string) string {
	// Map common Garmin activity types
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
	case "e_bike_mountain":
		return "e_bike_mountain"
	case "e_bike_road":
		return "e_bike_road"
	case "e_bike_cargo":
		return "e_bike_cargo"
	default:
		return "other"
	}
}

// Map Garmin photos to wanderer photos
func mapGarminPhotos(photos []garminPhoto) []photo {
	if len(photos) == 0 {
		return nil
	}

	result := make([]photo, 0, len(photos))
	for _, p := range photos {
		if p.PhotoURL == "" {
			continue
		}

		// Clean up photo URL
		photoURL := strings.TrimSpace(p.PhotoURL)
		if photoURL == "" {
			continue
		}

		// Ensure URL is properly formatted
		if !strings.HasPrefix(photoURL, "http") {
			photoURL = "https:" + photoURL
		}

		// Parse URL to ensure it's valid
		parsedURL, err := url.Parse(photoURL)
		if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
			continue
		}

		lat := p.Latitude
		lon := p.Longitude

		result = append(result, photo{
			ExternalID:  strconv.FormatInt(p.PhotoID, 10),
			Filename:    fmt.Sprintf("garmin-%d.jpg", p.PhotoID),
			ContentType: "image/jpeg",
			Lat:         &lat,
			Lon:         &lon,
			Source: mediaSource{
				Type: "url",
				URL:  photoURL,
			},
		})
	}

	return result
}

// Build GPX from track points
func buildGPXFromTrackPoints(points []trackPoint) []byte {
	// Build a simple GPX 1.1 document from track points
	// This is a simplified implementation

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
		gpx.WriteString(fmt.Sprintf("      <trkpt lat=\"%f\" lon=\"%f\">\n", point.Latitude, point.Longitude))
		if point.Elevation != 0 {
			gpx.WriteString(fmt.Sprintf("        <ele>%.1f</ele>\n", point.Elevation))
		}
		if point.Timestamp != "" {
			gpx.WriteString(fmt.Sprintf("        <time>%s</time>\n", point.Timestamp))
		}
		gpx.WriteString("      </trkpt>\n")
	}

	gpx.WriteString(`    </trkseg>
  </trk>
</gpx>
`)

	return []byte(gpx.String())
}
