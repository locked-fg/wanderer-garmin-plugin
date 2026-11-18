//go:build tinygo

package main

import (
	"encoding/json"
	"errors"

	"github.com/extism/go-pdk"
	"github.com/open-wanderer/wanderer/plugins/garmin/internal/garminutil"
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
		var configErr *garminutil.ConfigError
		if errors.As(err, &configErr) {
			return fail("invalid_request", err.Error())
		}
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

type pluginError struct {
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
}

func fail(code, message string) int32 {
	data, err := json.Marshal(pluginError{Code: code, Message: message})
	if err != nil {
		pdk.SetErrorString(message)
		return 1
	}
	pdk.SetErrorString(string(data))
	return 1
}
