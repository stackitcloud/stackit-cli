package utils

import (
	"strings"
	"testing"

	telemetryrouter "github.com/stackitcloud/stackit-sdk-go/services/telemetryrouter/v1api"

	"github.com/google/go-cmp/cmp"
)

func TestParseFilterJSON(t *testing.T) {
	tests := []struct {
		description string
		raw         string
		want        *telemetryrouter.ConfigFilter
		wantErr     bool
		errContains string
	}{
		{
			// json.Unmarshal on the generated SDK types always yields a non-nil (possibly
			// empty) AdditionalProperties map, unlike hand-constructed struct literals.
			description: "single attribute",
			raw:         `{"attributes": [{"key": "http.method", "level": "resource", "matcher": "=", "values": ["GET", "HEAD"]}]}`,
			want: &telemetryrouter.ConfigFilter{
				Attributes: []telemetryrouter.ConfigFilterAttributes{
					{
						Key:                  "http.method",
						Level:                telemetryrouter.CONFIGFILTERLEVEL_RESOURCE,
						Matcher:              telemetryrouter.CONFIGFILTERMATCHER_EQUAL,
						Values:               []string{"GET", "HEAD"},
						AdditionalProperties: map[string]interface{}{},
					},
				},
				AdditionalProperties: map[string]interface{}{},
			},
		},
		{
			description: "many attributes",
			raw: `{"attributes": [` +
				`{"key": "http.method", "level": "resource", "matcher": "=", "values": ["GET", "HEAD"]},` +
				`{"key": "http.status", "level": "logRecord", "matcher": "!=", "values": ["500"]}` +
				`]}`,
			want: &telemetryrouter.ConfigFilter{
				Attributes: []telemetryrouter.ConfigFilterAttributes{
					{
						Key:                  "http.method",
						Level:                telemetryrouter.CONFIGFILTERLEVEL_RESOURCE,
						Matcher:              telemetryrouter.CONFIGFILTERMATCHER_EQUAL,
						Values:               []string{"GET", "HEAD"},
						AdditionalProperties: map[string]interface{}{},
					},
					{
						Key:                  "http.status",
						Level:                telemetryrouter.CONFIGFILTERLEVEL_LOG_RECORD,
						Matcher:              telemetryrouter.CONFIGFILTERMATCHER_NOT_EQUAL,
						Values:               []string{"500"},
						AdditionalProperties: map[string]interface{}{},
					},
				},
				AdditionalProperties: map[string]interface{}{},
			},
		},
		{
			description: "empty attributes list",
			raw:         `{"attributes": []}`,
			want: &telemetryrouter.ConfigFilter{
				Attributes:           []telemetryrouter.ConfigFilterAttributes{},
				AdditionalProperties: map[string]interface{}{},
			},
		},
		{
			description: "invalid json",
			raw:         `{"attributes": [`,
			wantErr:     true,
			errContains: "decode filter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			got, err := ParseFilterJSON(tt.raw)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseFilterJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if tt.errContains != "" && (err == nil || !strings.Contains(err.Error(), tt.errContains)) {
					t.Fatalf("expected error to contain %q, got %v", tt.errContains, err)
				}
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("ParseFilterJSON() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestMapToUpdateDestinationPayload(t *testing.T) {
	config := telemetryrouter.DestinationConfig{
		ConfigType: telemetryrouter.DESTINATIONCONFIGTYPE_OPEN_TELEMETRY,
		OpenTelemetry: &telemetryrouter.DestinationConfigOpenTelemetry{
			Uri: "https://otel-collector.example.com:4317",
		},
	}
	description := "my description"
	displayName := "my-destination"

	tests := []struct {
		description string
		resp        *telemetryrouter.DestinationResponse
		want        *telemetryrouter.UpdateDestinationPayload
		wantErr     bool
	}{
		{
			description: "maps config, display name and description",
			resp: &telemetryrouter.DestinationResponse{
				Config:      config,
				DisplayName: displayName,
				Description: &description,
			},
			want: &telemetryrouter.UpdateDestinationPayload{
				Config:      &config,
				DisplayName: &displayName,
				Description: &description,
			},
		},
		{
			description: "maps destination without description",
			resp: &telemetryrouter.DestinationResponse{
				Config:      config,
				DisplayName: displayName,
			},
			want: &telemetryrouter.UpdateDestinationPayload{
				Config:      &config,
				DisplayName: &displayName,
			},
		},
		{
			description: "nil response returns error",
			resp:        nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			got, err := MapToUpdateDestinationPayload(tt.resp)
			if (err != nil) != tt.wantErr {
				t.Fatalf("MapToUpdateDestinationPayload() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("MapToUpdateDestinationPayload() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
