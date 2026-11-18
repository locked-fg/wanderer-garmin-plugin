package garminutil

import "testing"

func TestUnmarshalTypeKey(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    string
		wantErr bool
	}{
		{name: "plain string", data: `"running"`, want: "running"},
		{name: "object", data: `{"typeKey": "hiking"}`, want: "hiking"},
		{name: "object with extra fields", data: `{"typeKey": "cycling", "other": 1}`, want: "cycling"},
		{name: "empty string", data: `""`, want: ""},
		{name: "invalid json", data: `{`, wantErr: true},
		{name: "wrong type", data: `42`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalTypeKey([]byte(tt.data))
			if (err != nil) != tt.wantErr {
				t.Fatalf("UnmarshalTypeKey(%s) error = %v, wantErr %v", tt.data, err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Fatalf("UnmarshalTypeKey(%s) = %q, want %q", tt.data, got, tt.want)
			}
		})
	}
}
