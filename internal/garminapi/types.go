package garminapi

import (
	"github.com/open-wanderer/wanderer/plugins/garmin/internal/garmintrack"
	"github.com/open-wanderer/wanderer/plugins/garmin/internal/garminutil"
)

// ActivitySearch is the Garmin activity search response
type ActivitySearch struct {
	TotalCount int        `json:"totalCount"`
	Start      int        `json:"start"`
	Limit      int        `json:"limit"`
	Activities []Activity `json:"activities"`
}

// ActivityType may arrive as a plain string or as {"typeKey": "..."}
type ActivityType struct {
	TypeKey string `json:"typeKey"`
}

func (t *ActivityType) UnmarshalJSON(data []byte) error {
	key, err := garminutil.UnmarshalTypeKey(data)
	if err != nil {
		return err
	}
	t.TypeKey = key
	return nil
}

// Privacy may arrive as a plain string or as {"typeKey": "..."}
type Privacy struct {
	TypeKey string `json:"typeKey"`
}

func (p *Privacy) UnmarshalJSON(data []byte) error {
	key, err := garminutil.UnmarshalTypeKey(data)
	if err != nil {
		return err
	}
	p.TypeKey = key
	return nil
}

// Activity is a Garmin activity search summary
type Activity struct {
	ActivityID     int64        `json:"activityId"`
	ActivityName   string       `json:"activityName"`
	StartTime      string       `json:"startTimeGMT"`
	StartTimeLocal string       `json:"startTimeLocal"`
	Duration       float64      `json:"duration"`
	Distance       float64      `json:"distance"`
	ElevationGain  float64      `json:"elevationGain"`
	ElevationLoss  float64      `json:"elevationLoss"`
	AverageSpeed   float64      `json:"averageSpeed"`
	MaxSpeed       float64      `json:"maxSpeed"`
	Calories       float64      `json:"calories"`
	ActivityType   ActivityType `json:"activityType"`
	StartLat       float64      `json:"startLatitude"`
	StartLon       float64      `json:"startLongitude"`
	EndLat         float64      `json:"endLatitude"`
	EndLon         float64      `json:"endLongitude"`
	HasTrack       bool         `json:"hasTrack"`
	HasPhotos      bool         `json:"hasPhotos"`
	Privacy        Privacy      `json:"privacy"`
}

// ActivityDetail is the Garmin activity detail response
type ActivityDetail struct {
	ActivityID    int64               `json:"activityId"`
	ActivityName  string              `json:"activityName"`
	Description   string              `json:"description"`
	StartTime     string              `json:"startTimeGMT"`
	Duration      float64             `json:"duration"`
	Distance      float64             `json:"distance"`
	ElevationGain float64             `json:"elevationGain"`
	ElevationLoss float64             `json:"elevationLoss"`
	AverageSpeed  float64             `json:"averageSpeed"`
	MaxSpeed      float64             `json:"maxSpeed"`
	Calories      float64             `json:"calories"`
	ActivityType  ActivityType        `json:"activityType"`
	StartLat      float64             `json:"startLatitude"`
	StartLon      float64             `json:"startLongitude"`
	EndLat        float64             `json:"endLatitude"`
	EndLon        float64             `json:"endLongitude"`
	Privacy       Privacy             `json:"privacy"`
	Difficulty    string              `json:"difficultyRating"`
	Location      string              `json:"locationName"`
	Photos        []Photo             `json:"photos"`
	Track         []garmintrack.Point `json:"trackPoints"`
}

// Photo is a Garmin activity photo reference
type Photo struct {
	PhotoID   int64   `json:"photoId"`
	PhotoURL  string  `json:"photoUrl"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp string  `json:"timestamp"`
}
