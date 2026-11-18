//go:build tinygo

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type garminClient struct {
	accessToken      string
	apiKey           string
	trackShapeLogged bool // Track if shape was logged for this client instance
}

type tokenExchangeResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// logHelper sends a log message to the host
func logHelper(level sdk.LogLevel, message string) {
	sdk.Log(level, message)
}

func newGarminClient(auth map[string]any) (*garminClient, error) {
	// Try to get access token first (pre-exchanged)
	accessToken := sdk.StringField(auth, "accessToken")
	apiKey := sdk.StringField(auth, "apiKey")

	client := &garminClient{
		accessToken: accessToken,
		apiKey:      apiKey,
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

// doRequest performs a single authenticated GET against the Garmin connector
func (c *garminClient) doRequest(path string, query []sdk.QueryParam, contentTypes []string, accept string) (sdk.HostResponse, []byte, error) {
	return sdk.HostRequest(sdk.HostRequestSpec{
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
	response, body, err := c.doRequest(path, query, contentTypes, accept)

	// Check if token is not active - retry after additional delay
	if err == nil && response.Status == 401 {
		// Check if the error is "token is not active"
		bodyStr := string(body)
		if strings.Contains(strings.ToLower(bodyStr), "not active") {
			logHelper(sdk.LogLevelWarn, "Garmin: Token not active, waiting additional 15s and retrying...")
			time.Sleep(15 * time.Second)

			// Retry
			response, body, err = c.doRequest(path, query, contentTypes, accept)

			if err == nil && response.Status == 200 {
				logHelper(sdk.LogLevelInfo, "Garmin: Request succeeded after retry")
			}
		}
	}

	return response, body, err
}
