//go:build tinygo

package main

import (
	"encoding/json"

	"github.com/extism/go-pdk"
)

func main() {}

//export list_routes_v1
func listRoutesV1() int32 {
	var input listInput
	if err := pdk.InputJSON(&input); err != nil {
		return fail("invalid_request", "invalid list_routes input: "+err.Error())
	}

	client, err := newGarminClient(input.Auth)
	if err != nil {
		return fail("auth_failed", err.Error())
	}

	output, err := syncRoutes(client, input)
	if err != nil {
		return fail("provider_unavailable", err.Error())
	}

	if err := pdk.OutputJSON(output); err != nil {
		return fail("internal_error", err.Error())
	}

	return 0
}

//export list_activities_v1
func listActivitiesV1() int32 {
	var input listInput
	if err := pdk.InputJSON(&input); err != nil {
		return fail("invalid_request", "invalid list_activities input: "+err.Error())
	}

	client, err := newGarminClient(input.Auth)
	if err != nil {
		return fail("auth_failed", err.Error())
	}

	output, err := syncActivities(client, input)
	if err != nil {
		return fail("provider_unavailable", err.Error())
	}

	if err := pdk.OutputJSON(output); err != nil {
		return fail("internal_error", err.Error())
	}

	return 0
}

//export get_route_detail_v1
func getRouteDetailV1() int32 {
	return getTrailDetail("planned")
}

//export get_activity_detail_v1
func getActivityDetailV1() int32 {
	return getTrailDetail("completed")
}

// Generic function to handle both route and activity detail requests
func getTrailDetail(kind string) int32 {
	var input detailInput
	if err := pdk.InputJSON(&input); err != nil {
		return fail("invalid_request", "invalid detail input: "+err.Error())
	}

	client, err := newGarminClient(input.Auth)
	if err != nil {
		return fail("auth_failed", err.Error())
	}

	// For Garmin, we only support activities for now
	// Routes are not separately available in Personal API
	if kind == "planned" {
		// Garmin Personal API doesn't have a separate routes endpoint
		// For now, we'll return an error or empty result
		// Users should use activities instead
		return fail("not_supported", "Garmin Personal API does not support separate routes, use activities instead")
	}

	item, err := getActivityDetail(client, input, kind)
	if err != nil {
		return fail("provider_unavailable", err.Error())
	}

	output := struct {
		Item  trailImport  `json:"item"`
		Error *pluginError `json:"error,omitempty"`
	}{
		Item: item,
	}

	if err := pdk.OutputJSON(output); err != nil {
		return fail("internal_error", err.Error())
	}

	return 0
}

// Helper function to fail with an error. The host reads the structured error
// from the extism error channel (pdk.SetErrorString), not from the output, so
// the JSON-encoded pluginError must go there.
func fail(code, message string) int32 {
	data, err := json.Marshal(pluginError{Code: code, Message: message})
	if err != nil {
		pdk.SetErrorString(message)
		return 1
	}
	pdk.SetErrorString(string(data))
	return 1
}
