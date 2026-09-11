package utils

import (
	"strings"
	"testing"

	telemetryrouter "github.com/stackitcloud/stackit-sdk-go/services/telemetryrouter/v1api"

	"github.com/google/go-cmp/cmp"
)

func TestParseFilterAttribute(t *testing.T) {
	tests := []struct {
		description string
		raw         string
		want        telemetryrouter.ConfigFilterAttributes
		wantErr     bool
		errContains string
	}{
		{
			description: "valid single attribute",
			raw:         "key=http.method;level=resource;matcher=!=;values=GET",
			want: telemetryrouter.ConfigFilterAttributes{
				Key:     "http.method",
				Level:   telemetryrouter.CONFIGFILTERLEVEL_RESOURCE,
				Matcher: telemetryrouter.CONFIGFILTERMATCHER_NOT_EQUAL,
				Values:  []string{"GET"},
			},
		},
		{
			description: "valid equal matcher",
			raw:         "key=http.method;level=resource;matcher==;values=GET,HEAD",
			want: telemetryrouter.ConfigFilterAttributes{
				Key:     "http.method",
				Level:   telemetryrouter.CONFIGFILTERLEVEL_RESOURCE,
				Matcher: telemetryrouter.CONFIGFILTERMATCHER_EQUAL,
				Values:  []string{"GET", "HEAD"},
			},
		},
		{
			description: "whitespace trimming in values",
			raw:         "key=http.method;level=scope;matcher=!=;values= GET , HEAD ,,POST ",
			want: telemetryrouter.ConfigFilterAttributes{
				Key:     "http.method",
				Level:   telemetryrouter.CONFIGFILTERLEVEL_SCOPE,
				Matcher: telemetryrouter.CONFIGFILTERMATCHER_NOT_EQUAL,
				Values:  []string{"GET", "HEAD", "POST"},
			},
		},
		{
			description: "logRecord level",
			raw:         "key=body;level=logRecord;matcher==;values=foo",
			want: telemetryrouter.ConfigFilterAttributes{
				Key:     "body",
				Level:   telemetryrouter.CONFIGFILTERLEVEL_LOG_RECORD,
				Matcher: telemetryrouter.CONFIGFILTERMATCHER_EQUAL,
				Values:  []string{"foo"},
			},
		},
		{
			description: "missing field",
			raw:         "key=http.method;level=resource;values=GET",
			wantErr:     true,
			errContains: "missing field(s): matcher",
		},
		{
			description: "unknown field",
			raw:         "key=http.method;level=resource;matcher==;values=GET;unknown=foo",
			wantErr:     true,
			errContains: "unrecognized field(s): unknown",
		},
		{
			description: "duplicate field",
			raw:         "key=http.method;key=http.status;level=resource;matcher==;values=GET",
			wantErr:     true,
			errContains: "duplicate field(s): key",
		},
		{
			description: "invalid level",
			raw:         "key=http.method;level=bogus;matcher==;values=GET",
			wantErr:     true,
			errContains: `"level" must be one of resource, scope, logRecord`,
		},
		{
			description: "invalid matcher",
			raw:         "key=http.method;level=resource;matcher=~;values=GET",
			wantErr:     true,
			errContains: `"matcher" must be one of =, !=`,
		},
		{
			description: "empty values",
			raw:         "key=http.method;level=resource;matcher==;values=  ,  ",
			wantErr:     true,
			errContains: `"values" must contain at least one non-empty`,
		},
		{
			description: "empty key",
			raw:         "key= ;level=resource;matcher==;values=GET",
			wantErr:     true,
			errContains: `"key" must not be empty`,
		},
		{
			description: "field with no equals sign",
			raw:         "key=http.method;bogusfield;level=resource;matcher==;values=GET",
			wantErr:     true,
			errContains: `must be of the form "name=value"`,
		},
		{
			description: "unknown_default_open_api not accepted as level",
			raw:         "key=http.method;level=unknown_default_open_api;matcher==;values=GET",
			wantErr:     true,
			errContains: `"level" must be one of resource, scope, logRecord`,
		},
		{
			description: "unknown_default_open_api not accepted as matcher",
			raw:         "key=http.method;level=resource;matcher=unknown_default_open_api;values=GET",
			wantErr:     true,
			errContains: `"matcher" must be one of =, !=`,
		},
		{
			description: "empty raw string",
			raw:         "",
			wantErr:     true,
			errContains: "missing field(s): key, level, matcher, values",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			got, err := ParseFilterAttribute(tt.raw)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseFilterAttribute() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if tt.errContains != "" && (err == nil || !strings.Contains(err.Error(), tt.errContains)) {
					t.Fatalf("expected error to contain %q, got %v", tt.errContains, err)
				}
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("ParseFilterAttribute() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestBuildConfigFilter(t *testing.T) {
	tests := []struct {
		description string
		raw         []string
		want        *telemetryrouter.ConfigFilter
		wantErr     bool
		errContains string
	}{
		{
			description: "empty raw slice returns nil, nil",
			raw:         nil,
			want:        nil,
		},
		{
			description: "empty (non-nil) raw slice returns nil, nil",
			raw:         []string{},
			want:        nil,
		},
		{
			description: "single attribute",
			raw: []string{
				"key=http.method;level=resource;matcher==;values=GET,HEAD",
			},
			want: &telemetryrouter.ConfigFilter{
				Attributes: []telemetryrouter.ConfigFilterAttributes{
					{
						Key:     "http.method",
						Level:   telemetryrouter.CONFIGFILTERLEVEL_RESOURCE,
						Matcher: telemetryrouter.CONFIGFILTERMATCHER_EQUAL,
						Values:  []string{"GET", "HEAD"},
					},
				},
			},
		},
		{
			description: "multiple attributes preserve order",
			raw: []string{
				"key=http.method;level=resource;matcher==;values=GET,HEAD",
				"key=http.status;level=logRecord;matcher=!=;values=500",
			},
			want: &telemetryrouter.ConfigFilter{
				Attributes: []telemetryrouter.ConfigFilterAttributes{
					{
						Key:     "http.method",
						Level:   telemetryrouter.CONFIGFILTERLEVEL_RESOURCE,
						Matcher: telemetryrouter.CONFIGFILTERMATCHER_EQUAL,
						Values:  []string{"GET", "HEAD"},
					},
					{
						Key:     "http.status",
						Level:   telemetryrouter.CONFIGFILTERLEVEL_LOG_RECORD,
						Matcher: telemetryrouter.CONFIGFILTERMATCHER_NOT_EQUAL,
						Values:  []string{"500"},
					},
				},
			},
		},
		{
			description: "invalid entry propagates error with context",
			raw: []string{
				"key=http.method;level=resource;matcher==;values=GET",
				"key=http.status;level=bogus;matcher==;values=500",
			},
			wantErr:     true,
			errContains: "parse --filter-attribute value #2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			got, err := BuildConfigFilter(tt.raw)
			if (err != nil) != tt.wantErr {
				t.Fatalf("BuildConfigFilter() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if tt.errContains != "" && (err == nil || !strings.Contains(err.Error(), tt.errContains)) {
					t.Fatalf("expected error to contain %q, got %v", tt.errContains, err)
				}
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("BuildConfigFilter() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

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
