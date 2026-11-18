//go:build tinygo

package main

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/open-wanderer/wanderer/plugins/garmin/internal/garmintrack"
	"github.com/open-wanderer/wanderer/plugins/garmin/internal/garminutil"
	"github.com/open-wanderer/wanderer/plugins/sdk"
)

// Map Garmin activity to trail import
func mapActivityToTrail(activity *garminActivityDetail, gpxData []byte) (trailImport, error) {
	privacy := mapPrivacy(activity.Privacy.TypeKey)
	difficulty := mapDifficulty(activity.Difficulty)

	var trackData track
	if gpxData != nil && len(gpxData) > 0 {
		trackData = track{
			Format:        "gpx",
			ContentBase64: base64.StdEncoding.EncodeToString(gpxData),
		}
	} else if len(activity.Track) > 0 {
		gpxData := garmintrack.BuildGPX(activity.Track)
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
		if parsed, err := garminutil.ParseActivityStartTime(activity.StartTime); err == nil {
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
