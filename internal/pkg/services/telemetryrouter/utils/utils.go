package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	telemetryrouter "github.com/stackitcloud/stackit-sdk-go/services/telemetryrouter/v1api"
)

var (
	ErrResponseNil = errors.New("response is nil")
)

func GetInstanceName(ctx context.Context, apiClient telemetryrouter.DefaultAPI, projectId, regionId, instanceId string) (string, error) {
	resp, err := apiClient.GetTelemetryRouter(ctx, projectId, regionId, instanceId).Execute()
	if err != nil {
		return "", fmt.Errorf("get TelemetryRouter instance: %w", err)
	} else if resp == nil {
		return "", ErrResponseNil
	}
	return resp.DisplayName, nil
}

func GetDestinationName(ctx context.Context, apiClient telemetryrouter.DefaultAPI, projectId, regionId, instanceId, destinationId string) (string, error) {
	resp, err := apiClient.GetDestination(ctx, projectId, regionId, instanceId, destinationId).Execute()
	if err != nil {
		return "", fmt.Errorf("get TelemetryRouter destination: %w", err)
	} else if resp == nil {
		return "", ErrResponseNil
	}
	return resp.DisplayName, nil
}

func GetAccessTokenName(ctx context.Context, apiClient telemetryrouter.DefaultAPI, projectId, regionId, instanceId, accessTokenId string) (string, error) {
	resp, err := apiClient.GetAccessToken(ctx, projectId, regionId, instanceId, accessTokenId).Execute()
	if err != nil {
		return "", fmt.Errorf("get TelemetryRouter access token: %w", err)
	} else if resp == nil {
		return "", ErrResponseNil
	}
	return resp.DisplayName, nil
}

// DefaultCreateDestinationPayloadOpenTelemetry is a default payload for creating a TelemetryRouter destination
// that routes to an OpenTelemetry (OTLP) endpoint, secured with basic auth.
var DefaultCreateDestinationPayloadOpenTelemetry = telemetryrouter.CreateDestinationPayload{
	DisplayName: "default-name",
	Config: telemetryrouter.DestinationConfig{
		ConfigType: telemetryrouter.DESTINATIONCONFIGTYPE_OPEN_TELEMETRY,
		OpenTelemetry: &telemetryrouter.DestinationConfigOpenTelemetry{
			Uri: "https://otel-collector.example.com:4317",
			BasicAuth: &telemetryrouter.DestinationConfigOpenTelemetryBasicAuth{
				Username: "username",
				Password: "password",
			},
		},
	},
}

// DefaultCreateDestinationPayloadS3 is a default payload for creating a TelemetryRouter destination
// that routes to an S3-compatible bucket, secured with an access key.
var DefaultCreateDestinationPayloadS3 = telemetryrouter.CreateDestinationPayload{
	DisplayName: "default-name",
	Config: telemetryrouter.DestinationConfig{
		ConfigType: telemetryrouter.DESTINATIONCONFIGTYPE_S3,
		S3: &telemetryrouter.DestinationConfigS3{
			Bucket:   "my-bucket",
			Endpoint: "https://object.storage.eu01.onstackit.cloud",
			AccessKey: &telemetryrouter.DestinationConfigS3AccessKey{ // nolint:gosec // false positive
				Id:     "access-key-id",
				Secret: "access-key-secret",
			},
		},
	},
}

// DefaultConfigFilter is a default example filter configuration for a TelemetryRouter instance,
// used as a starting point for the --filter flag of "instance create"/"instance update" when no
// existing instance is used as a seed.
var DefaultConfigFilter = telemetryrouter.ConfigFilter{
	Attributes: []telemetryrouter.ConfigFilterAttributes{
		{
			Key:     "http.method",
			Level:   telemetryrouter.CONFIGFILTERLEVEL_RESOURCE,
			Matcher: telemetryrouter.CONFIGFILTERMATCHER_EQUAL,
			Values:  []string{"GET", "HEAD"},
		},
	},
}

// MapToUpdateDestinationPayload maps a DestinationResponse (as returned by the API) to an
// UpdateDestinationPayload, so it can be used as a seed payload for the "update" command.
func MapToUpdateDestinationPayload(resp *telemetryrouter.DestinationResponse) (*telemetryrouter.UpdateDestinationPayload, error) {
	if resp == nil {
		return nil, fmt.Errorf("no TelemetryRouter destination provided")
	}

	payload := &telemetryrouter.UpdateDestinationPayload{
		Config:      &resp.Config,
		Description: resp.Description,
		DisplayName: &resp.DisplayName,
	}

	return payload, nil
}

const (
	filterAttributeFieldKey     = "key"
	filterAttributeFieldLevel   = "level"
	filterAttributeFieldMatcher = "matcher"
	filterAttributeFieldValues  = "values"
)

// ParseFilterAttribute parses a single --filter-attribute flag value in the form
//
//	"key=<attr-key>;level=<resource|scope|logRecord>;matcher=<=|!=>;values=<v1>,<v2>,..."
//
// into a telemetryrouter.ConfigFilterAttributes.
func ParseFilterAttribute(raw string) (telemetryrouter.ConfigFilterAttributes, error) {
	fields := strings.Split(raw, ";")

	found := make(map[string]string, len(fields))
	var unknownFields []string
	var duplicateFields []string

	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}

		// Split on the first "=" only, in case a value itself contains "=".
		parts := strings.SplitN(field, "=", 2)
		if len(parts) != 2 {
			return telemetryrouter.ConfigFilterAttributes{}, fmt.Errorf(
				"invalid --filter-attribute value %q: field %q must be of the form \"name=value\"", raw, field)
		}
		name := strings.TrimSpace(parts[0])
		value := parts[1]

		switch name {
		case filterAttributeFieldKey, filterAttributeFieldLevel, filterAttributeFieldMatcher, filterAttributeFieldValues:
			if _, exists := found[name]; exists {
				duplicateFields = append(duplicateFields, name)
				continue
			}
			found[name] = value
		default:
			unknownFields = append(unknownFields, name)
		}
	}

	var missingFields []string
	for _, required := range []string{filterAttributeFieldKey, filterAttributeFieldLevel, filterAttributeFieldMatcher, filterAttributeFieldValues} {
		if _, ok := found[required]; !ok {
			missingFields = append(missingFields, required)
		}
	}

	if len(missingFields) > 0 || len(unknownFields) > 0 || len(duplicateFields) > 0 {
		var problems []string
		if len(missingFields) > 0 {
			problems = append(problems, fmt.Sprintf("missing field(s): %s", strings.Join(missingFields, ", ")))
		}
		if len(unknownFields) > 0 {
			problems = append(problems, fmt.Sprintf("unrecognized field(s): %s", strings.Join(unknownFields, ", ")))
		}
		if len(duplicateFields) > 0 {
			problems = append(problems, fmt.Sprintf("duplicate field(s): %s", strings.Join(duplicateFields, ", ")))
		}
		return telemetryrouter.ConfigFilterAttributes{}, fmt.Errorf(
			`invalid --filter-attribute value %q: %s (expected exactly the fields "key", "level", "matcher" and "values", each specified once)`,
			raw, strings.Join(problems, "; "))
	}

	key := strings.TrimSpace(found[filterAttributeFieldKey])
	if key == "" {
		return telemetryrouter.ConfigFilterAttributes{}, fmt.Errorf(
			`invalid --filter-attribute value %q: "key" must not be empty`, raw)
	}

	level := strings.TrimSpace(found[filterAttributeFieldLevel])
	if !isValidConfigFilterLevel(level) {
		return telemetryrouter.ConfigFilterAttributes{}, fmt.Errorf(
			`invalid --filter-attribute value %q: "level" must be one of resource, scope, logRecord (got %q)`, raw, level)
	}

	matcher := strings.TrimSpace(found[filterAttributeFieldMatcher])
	if !isValidConfigFilterMatcher(matcher) {
		return telemetryrouter.ConfigFilterAttributes{}, fmt.Errorf(
			`invalid --filter-attribute value %q: "matcher" must be one of =, != (got %q)`, raw, matcher)
	}

	var values []string
	for _, v := range strings.Split(found[filterAttributeFieldValues], ",") {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		values = append(values, v)
	}
	if len(values) == 0 {
		return telemetryrouter.ConfigFilterAttributes{}, fmt.Errorf(
			`invalid --filter-attribute value %q: "values" must contain at least one non-empty, comma-separated value`, raw)
	}

	return telemetryrouter.ConfigFilterAttributes{
		Key:     key,
		Level:   telemetryrouter.ConfigFilterLevel(level),
		Matcher: telemetryrouter.ConfigFilterMatcher(matcher),
		Values:  values,
	}, nil
}

// isValidConfigFilterLevel reports whether level is one of the SDK's allowed
// ConfigFilterLevel enum values, excluding the "unknown_default_open_api" sentinel.
func isValidConfigFilterLevel(level string) bool {
	for _, allowed := range telemetryrouter.AllowedConfigFilterLevelEnumValues {
		if allowed == telemetryrouter.CONFIGFILTERLEVEL_UNKNOWN_DEFAULT_OPEN_API {
			continue
		}
		if string(allowed) == level {
			return true
		}
	}
	return false
}

// isValidConfigFilterMatcher reports whether matcher is one of the SDK's allowed
// ConfigFilterMatcher enum values, excluding the "unknown_default_open_api" sentinel.
func isValidConfigFilterMatcher(matcher string) bool {
	for _, allowed := range telemetryrouter.AllowedConfigFilterMatcherEnumValues {
		if allowed == telemetryrouter.CONFIGFILTERMATCHER_UNKNOWN_DEFAULT_OPEN_API {
			continue
		}
		if string(allowed) == matcher {
			return true
		}
	}
	return false
}

// BuildConfigFilter maps a slice of raw --filter-attribute flag values (as collected by
// pflag.StringArray) into a *ConfigFilter. Returns nil (not an error) if raw is empty,
// so the Filter field stays unset for API calls where no filtering was requested.
func BuildConfigFilter(raw []string) (*telemetryrouter.ConfigFilter, error) {
	if len(raw) == 0 {
		return nil, nil
	}

	attrs := make([]telemetryrouter.ConfigFilterAttributes, 0, len(raw))
	for i, r := range raw {
		attr, err := ParseFilterAttribute(r)
		if err != nil {
			return nil, fmt.Errorf("parse --filter-attribute value #%d (%q): %w", i+1, r, err)
		}
		attrs = append(attrs, attr)
	}

	return &telemetryrouter.ConfigFilter{Attributes: attrs}, nil
}

// ParseFilterJSON parses the raw JSON value of a --filter flag (a JSON object of the form
// {"attributes": [{"key": ..., "level": ..., "matcher": ..., "values": [...]}, ...]}, matching
// the API's ConfigFilter schema) into a *ConfigFilter. This is the practical alternative to
// repeating --filter-attribute when many filter attributes need to be set at once.
func ParseFilterJSON(raw string) (*telemetryrouter.ConfigFilter, error) {
	var filter telemetryrouter.ConfigFilter
	if err := json.Unmarshal([]byte(raw), &filter); err != nil {
		return nil, fmt.Errorf("decode filter: %w", err)
	}
	return &filter, nil
}
