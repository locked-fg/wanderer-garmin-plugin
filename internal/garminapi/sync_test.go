package garminapi

import (
	"encoding/json"
	"testing"
	"time"
)

const searchFixture = `{"totalCount":4,"start":0,"limit":50,"activities":[
 {"activityId":1,"activityName":"Old","startTimeGMT":"2025-01-01 10:00:00","activityType":"running","privacy":"private"},
 {"activityId":2,"activityName":"New","startTimeGMT":"2026-06-01T10:00:00Z","activityType":{"typeKey":"hiking"},"privacy":{"typeKey":"public"}},
 {"activityId":3,"activityName":"BadGMT","startTimeGMT":"garbage","startTimeLocal":"2026-06-02 10:00:00"},
 {"activityId":4,"activityName":"BadBoth","startTimeGMT":"garbage","startTimeLocal":"also bad"}
]}`

func loadSearch(t *testing.T) ActivitySearch {
	t.Helper()
	var search ActivitySearch
	if err := json.Unmarshal([]byte(searchFixture), &search); err != nil {
		t.Fatal(err)
	}
	return search
}

func TestActivityDecodingStringOrObject(t *testing.T) {
	search := loadSearch(t)
	if len(search.Activities) != 4 {
		t.Fatalf("got %d activities; want 4", len(search.Activities))
	}
	if got := search.Activities[0].ActivityType.TypeKey; got != "running" {
		t.Fatalf("string typeKey = %q; want running", got)
	}
	if got := search.Activities[0].Privacy.TypeKey; got != "private" {
		t.Fatalf("string privacy = %q; want private", got)
	}
	if got := search.Activities[1].ActivityType.TypeKey; got != "hiking" {
		t.Fatalf("object typeKey = %q; want hiking", got)
	}
	if got := search.Activities[1].Privacy.TypeKey; got != "public" {
		t.Fatalf("object privacy = %q; want public", got)
	}
}

func TestBuildSummariesNoAfterKeepsEverything(t *testing.T) {
	search := loadSearch(t)
	res := BuildSummaries(search.Activities, time.Time{}, 50)
	if len(res.Items) != 4 {
		t.Fatalf("got %d items; want 4 (no filter without after)", len(res.Items))
	}
	if res.FilteredByDate != 0 || len(res.InvalidStartTime) != 0 {
		t.Fatalf("unexpected filtering: filtered=%d invalid=%d", res.FilteredByDate, len(res.InvalidStartTime))
	}
	if res.Items[0].Source.ExternalID != "1" || res.Items[0].Source.Provider != "garmin" || res.Items[0].Kind != "completed" {
		t.Fatalf("unexpected summary: %#v", res.Items[0])
	}
}

func TestBuildSummariesWithAfter(t *testing.T) {
	search := loadSearch(t)
	after := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	res := BuildSummaries(search.Activities, after, 50)

	if len(res.Items) != 2 {
		t.Fatalf("got %d items; want 2 (New + BadGMT via local fallback)", len(res.Items))
	}
	ids := map[string]bool{}
	for _, item := range res.Items {
		ids[item.Source.ExternalID] = true
	}
	if !ids["2"] || !ids["3"] {
		t.Fatalf("kept IDs = %v; want 2 and 3", ids)
	}
	if res.FilteredByDate != 1 {
		t.Fatalf("filteredByDate = %d; want 1 (Old)", res.FilteredByDate)
	}
	if len(res.InvalidStartTime) != 1 || res.InvalidStartTime[0].ActivityID != 4 {
		t.Fatalf("invalid = %#v; want activity 4", res.InvalidStartTime)
	}
}

func TestBuildSummariesLimitBreaksEarly(t *testing.T) {
	search := loadSearch(t)
	res := BuildSummaries(search.Activities, time.Time{}, 2)
	if len(res.Items) != 2 {
		t.Fatalf("got %d items; want 2 (limit)", len(res.Items))
	}
}

func TestPageResult(t *testing.T) {
	tests := []struct {
		page, received, limit int
		wantNext              int
		wantMore              bool
	}{
		{1, 50, 50, 2, true},
		{1, 8, 50, 2, false},
		{3, 50, 50, 4, true},
		{0, 10, 50, 2, false},
		{1, 10, 0, 2, false},
	}
	for _, tt := range tests {
		next, more := PageResult(tt.page, tt.received, tt.limit)
		if next != tt.wantNext || more != tt.wantMore {
			t.Errorf("PageResult(%d,%d,%d) = (%d,%t); want (%d,%t)", tt.page, tt.received, tt.limit, next, more, tt.wantNext, tt.wantMore)
		}
	}
}
