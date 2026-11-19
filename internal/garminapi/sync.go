package garminapi

import (
	"strconv"
	"time"

	"github.com/open-wanderer/wanderer/plugins/garmin/internal/garminutil"
	"github.com/open-wanderer/wanderer/plugins/sdk"
)

// SummaryResult is the outcome of turning one search page into trail summaries
type SummaryResult struct {
	Items            []sdk.TrailSummary
	FilteredByDate   int
	InvalidStartTime []Activity
}

// BuildSummaries maps one page of activities to trail summaries, applying the
// optional after-date filter. Activities with unparseable start times are
// skipped and reported in InvalidStartTime; activities before after are
// counted in FilteredByDate. At most limit summaries are returned.
func BuildSummaries(activities []Activity, after time.Time, limit int) SummaryResult {
	res := SummaryResult{Items: make([]sdk.TrailSummary, 0, len(activities))}
	for _, activity := range activities {
		if !after.IsZero() {
			startTime, parseErr := garminutil.ParseActivityStartTime(activity.StartTime)
			if parseErr != nil && activity.StartTimeLocal != "" {
				startTime, parseErr = garminutil.ParseActivityStartTime(activity.StartTimeLocal)
			}
			if parseErr != nil {
				res.InvalidStartTime = append(res.InvalidStartTime, activity)
				continue
			}
			if startTime.Before(after) {
				res.FilteredByDate++
				continue
			}
		}

		res.Items = append(res.Items, sdk.TrailSummary{
			Source: sdk.TrailImportSource{
				Provider:   "garmin",
				ExternalID: strconv.FormatInt(activity.ActivityID, 10),
			},
			Kind: "completed",
		})

		if limit > 0 && len(res.Items) >= limit {
			break
		}
	}
	return res
}

// PageResult computes the next-page state for a search page
func PageResult(page, received, limit int) (nextPage int, hasMore bool) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		return page + 1, false
	}
	return page + 1, received >= limit
}
