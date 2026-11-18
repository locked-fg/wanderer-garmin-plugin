package garminmap

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/open-wanderer/wanderer/plugins/garmin/internal/garminapi"
	"github.com/open-wanderer/wanderer/plugins/garmin/internal/garmintrack"
	"github.com/open-wanderer/wanderer/plugins/garmin/internal/garminutil"
	"github.com/open-wanderer/wanderer/plugins/sdk"
)

// MapActivityToTrail converts a Garmin activity detail into a Wanderer trail
// import. Non-fatal issues (e.g. an unparseable startTimeGMT) are returned as
// warnings so the caller can log them in the host.
func MapActivityToTrail(activity *garminapi.ActivityDetail, gpxData []byte) (sdk.TrailImport, []string, error) {
	var warnings []string

	privacy := MapPrivacy(activity.Privacy.TypeKey)
	difficulty := MapDifficulty(activity.Difficulty)

	var trackData sdk.Track
	if len(gpxData) > 0 {
		trackData = sdk.Track{
			Format:        "gpx",
			ContentBase64: base64.StdEncoding.EncodeToString(gpxData),
		}
	} else if len(activity.Track) > 0 {
		trackData = sdk.Track{
			Format:        "gpx",
			ContentBase64: base64.StdEncoding.EncodeToString(garmintrack.BuildGPX(activity.Track)),
		}
	}

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
			warnings = append(warnings, fmt.Sprintf("unable to parse startTimeGMT %q; omitting startedAt", activity.StartTime))
		}
	}

	return sdk.TrailImport{
		Source: sdk.TrailImportSource{
			Provider:   "garmin",
			ExternalID: strconv.FormatInt(activity.ActivityID, 10),
			URL:        fmt.Sprintf("https://connect.garmin.com/modern/activity/%d", activity.ActivityID),
		},
		Kind:         "completed",
		Name:         activity.ActivityName,
		Description:  activity.Description,
		StartedAt:    startedAt,
		ActivityType: MapActivityType(activity.ActivityType.TypeKey),
		Privacy:      &privacy,
		Track:        trackData,
		Waypoints:    []sdk.Waypoint{},
		Photos:       MapPhotos(activity.Photos),
		Metadata:     metadata,
		Difficulty:   difficulty,
	}, warnings, nil
}

// MapPrivacy maps a Garmin privacy value to a Wanderer privacy value
func MapPrivacy(garminPrivacy string) string {
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

// MapDifficulty maps a Garmin difficulty rating to a Wanderer difficulty
func MapDifficulty(garminDifficulty string) string {
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

// MapActivityType maps a Garmin activity type key to a Wanderer activity type
func MapActivityType(garminType string) string {
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

// MapPhotos converts Garmin photo references to Wanderer photos, skipping
// entries without a usable URL
func MapPhotos(photos []garminapi.Photo) []sdk.Photo {
	if len(photos) == 0 {
		return nil
	}

	result := make([]sdk.Photo, 0, len(photos))
	for _, p := range photos {
		photoURL := strings.TrimSpace(p.PhotoURL)
		if photoURL == "" {
			continue
		}
		if !strings.HasPrefix(photoURL, "http") {
			photoURL = "https:" + photoURL
		}
		parsedURL, err := url.Parse(photoURL)
		if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
			continue
		}

		lat := p.Latitude
		lon := p.Longitude

		result = append(result, sdk.Photo{
			ExternalID:  strconv.FormatInt(p.PhotoID, 10),
			Filename:    fmt.Sprintf("garmin-%d.jpg", p.PhotoID),
			ContentType: "image/jpeg",
			Lat:         &lat,
			Lon:         &lon,
			Source: sdk.MediaSource{
				Type: "url",
				URL:  photoURL,
			},
		})
	}

	return result
}
