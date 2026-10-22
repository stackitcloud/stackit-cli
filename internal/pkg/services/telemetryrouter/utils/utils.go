package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	telemetryrouter "github.com/stackitcloud/stackit-sdk-go/services/telemetryrouter/v1api"
)

var (
	ErrResponseNil = errors.New("response is nil")
)

func GetInstanceName(ctx context.Context, apiClient telemetryrouter.DefaultAPI, projectId, regionId, instanceId string) (string, error) {
	resp, err := apiClient.GetTelemetryRouter(ctx, projectId, regionId, instanceId).Execute()
	if err != nil {
		return "", fmt.Errorf("get Telemetry Router instance: %w", err)
	} else if resp == nil {
		return "", ErrResponseNil
	}
	return resp.DisplayName, nil
}

func GetDestinationName(ctx context.Context, apiClient telemetryrouter.DefaultAPI, projectId, regionId, instanceId, destinationId string) (string, error) {
	resp, err := apiClient.GetDestination(ctx, projectId, regionId, instanceId, destinationId).Execute()
	if err != nil {
		return "", fmt.Errorf("get Telemetry Router destination: %w", err)
	} else if resp == nil {
		return "", ErrResponseNil
	}
	return resp.DisplayName, nil
}

func GetAccessTokenName(ctx context.Context, apiClient telemetryrouter.DefaultAPI, projectId, regionId, instanceId, accessTokenId string) (string, error) {
	resp, err := apiClient.GetAccessToken(ctx, projectId, regionId, instanceId, accessTokenId).Execute()
	if err != nil {
		return "", fmt.Errorf("get Telemetry Router access token: %w", err)
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
		return nil, fmt.Errorf("no Telemetry Router destination provided")
	}

	payload := &telemetryrouter.UpdateDestinationPayload{
		Config:      &resp.Config,
		Description: resp.Description,
		DisplayName: &resp.DisplayName,
	}

	return payload, nil
}

// ParseFilterJSON parses the raw JSON value of a --filter flag (a JSON object of the form
// {"attributes": [{"key": ..., "level": ..., "matcher": ..., "values": [...]}, ...]}, matching
// the API's ConfigFilter schema) into a *ConfigFilter.
func ParseFilterJSON(raw string) (*telemetryrouter.ConfigFilter, error) {
	var filter telemetryrouter.ConfigFilter
	if err := json.Unmarshal([]byte(raw), &filter); err != nil {
		return nil, fmt.Errorf("decode filter: %w", err)
	}
	return &filter, nil
}
