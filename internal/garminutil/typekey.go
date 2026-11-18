package garminutil

import "encoding/json"

// UnmarshalTypeKey decodes a Garmin type field that may arrive either as a
// plain string ("running") or as an object ({"typeKey": "running"}).
func UnmarshalTypeKey(data []byte) (string, error) {
	var key string
	if err := json.Unmarshal(data, &key); err == nil {
		return key, nil
	}

	var value struct {
		TypeKey string `json:"typeKey"`
	}
	if err := json.Unmarshal(data, &value); err != nil {
		return "", err
	}
	return value.TypeKey, nil
}
