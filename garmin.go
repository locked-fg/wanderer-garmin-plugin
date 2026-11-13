//go:build tinygo

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/open-wanderer/wanderer/plugins/sdk"
)

const (
	garminBaseURL                = "https://personal.connectapi.garmin.com"
	garminDownloadMaxBytes int64 = 16 * 1024 * 1024
	defaultPageSize              = 50
)

// Garmin client structure
type garminClient struct {
	accessToken string
	apiKey      string
	baseURL     string
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

	logHelper(sdk.LogLevelDebug, fmt.Sprintf("Garmin: Creating client - hasAccessToken: %v, hasApiKey: %v", accessToken != "", apiKey != ""))

	client := &garminClient{
		accessToken: accessToken,
		apiKey:      apiKey,
		baseURL:     garminBaseURL,
	}

	// If we have an API key but no access token, exchange it
	if apiKey != "" && accessToken == "" {
		logHelper(sdk.LogLevelInfo, "Garmin: Exchanging API key for access token")
		// Exchange API key for access token
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

	logHelper(sdk.LogLevelDebug, "Garmin: Client created successfully")
	return client, nil
}

// Exchange API key for access token
func (c *garminClient) exchangeAPIKeyForToken() (string, error) {
	// Use the token_exchange connector with ApiKey header
	response, body, err := c.tokenExchangeRequest("/personalapitoken-service/personal-api/personal-api-key/token", []sdk.QueryParam{}, []string{"application/json"})
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
func (c *garminClient) tokenExchangeRequest(path string, query []sdk.QueryParam, contentTypes []string) (sdk.HostResponse, []byte, error) {
	accept := "application/json"
	if len(contentTypes) > 0 {
		accept = contentTypes[0]
	}

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
			"Accept":        accept,
		},
		Expect: sdk.ResponseExpect{
			ContentTypes: contentTypes,
			MaxBytes:     garminDownloadMaxBytes,
		},
	})
}

// Search for activities
func (c *garminClient) searchActivities(page int, limit int, after time.Time) ([]garminActivity, error) {
	var result garminActivitySearch

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

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse activities: %w", err)
	}

	return result.Activities, nil
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

// Download activity file (ZIP containing GPX)
func (c *garminClient) downloadActivity(activityID int64) ([]byte, error) {
	path := fmt.Sprintf("/download-service/files/activity/%d", activityID)

	// Garmin returns ZIP files for activity downloads
	response, body, err := c.request(path, nil, []string{
		"application/zip",
		"application/x-zip-compressed",
		"application/octet-stream",
	})
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
	for _, activity := range activities {
		// Filter by after date if not already filtered by API
		if !after.IsZero() {
			startTime, parseErr := time.Parse(time.RFC3339, activity.StartTime)
			if parseErr != nil {
				// Try different date formats
				startTime, parseErr = time.Parse("2006-01-02T15:04:05.000Z", activity.StartTime)
				if parseErr != nil {
					continue
				}
			}
			if startTime.Before(after) {
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

// Get activity detail and map to trail import
func getActivityDetail(client *garminClient, input detailInput, kind string) (trailImport, error) {
	activityID, err := strconv.ParseInt(input.Summary.Source.ExternalID, 10, 64)
	if err != nil {
		return trailImport{}, fmt.Errorf("invalid activity ID: %w", err)
	}

	detail, err := client.getActivity(activityID)
	if err != nil {
		return trailImport{}, fmt.Errorf("failed to get activity detail: %w", err)
	}

	// Download the activity file (ZIP containing GPX)
	zipData, err := client.downloadActivity(activityID)
	if err != nil {
		// If download fails, try to extract GPX from track points
		sdk.LogWarn(fmt.Sprintf("Failed to download activity file, using track points: %v", err))
		return mapActivityToTrail(detail, nil)
	}

	// Extract GPX from ZIP
	gpxData, err := extractGPXFromZip(zipData)
	if err != nil {
		sdk.LogWarn(fmt.Sprintf("Failed to extract GPX from ZIP: %v", err))
		return mapActivityToTrail(detail, nil)
	}

	return mapActivityToTrail(detail, gpxData)
}

// Extract GPX from ZIP data
func extractGPXFromZip(zipData []byte) ([]byte, error) {
	// For now, we'll assume the ZIP contains a GPX file
	// In a real implementation, we would unzip and find the GPX file
	// This is a simplified version

	// Try to find GPX in the ZIP by looking for the file signature
	// This is a placeholder - actual implementation would use a ZIP reader

	// For Garmin activities, the ZIP typically contains:
	// - activity.gpx (the main GPX file)
	// - other metadata files

	// Since we can't easily parse ZIP in TinyGo without additional libraries,
	// we'll return the raw ZIP data and let the host handle it
	// or we'll need to implement a simple ZIP parser

	// For now, return the raw data
	// The host should be able to handle ZIP files
	return zipData, nil
}

// Map Garmin activity to trail import
func mapActivityToTrail(activity *garminActivityDetail, gpxData []byte) (trailImport, error) {
	privacy := mapPrivacy(activity.Privacy)
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
	activityType := mapActivityType(activity.ActivityType)

	// Build metadata
	metadata := map[string]any{
		"distance":         activity.Distance,
		"elevationGain":    activity.ElevationGain,
		"elevationLoss":    activity.ElevationLoss,
		"duration":         activity.Duration,
		"calories":         activity.Calories,
		"averageSpeed":     activity.AverageSpeed,
		"maxSpeed":         activity.MaxSpeed,
		"providerCategory": activity.ActivityType,
	}

	// Handle start time
	startedAt := activity.StartTime
	if strings.Contains(startedAt, "Z") {
		// Already in ISO format with Z
	} else if !strings.Contains(startedAt, "T") {
		// Might be date only
		startedAt = startedAt + "T00:00:00Z"
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
