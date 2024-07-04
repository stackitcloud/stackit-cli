package auth

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
)

func TestGetIDPEndpoint(t *testing.T) {
	tests := []struct {
		name              string
		idpCustomEndpoint string
		expected          string
	}{
		{
			name:              "custom endpoint specified",
			idpCustomEndpoint: "https://custom.endpoint",
			expected:          "https://custom.endpoint",
		},
		{
			name:              "custom endpoint not specified",
			idpCustomEndpoint: "",
			expected:          defaultIDPEndpoint,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()

			if tt.idpCustomEndpoint != "" {
				viper.Set(config.IdentityProviderCustomEndpointKey, tt.idpCustomEndpoint)
			}

			got := getIDPEndpoint()

			if got != tt.expected {
				t.Fatalf("expected idp endpoint %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestGetIDPClientID(t *testing.T) {
	tests := []struct {
		name              string
		idpCustomClientID string
		expected          string
	}{
		{
			name:              "custom client ID specified",
			idpCustomClientID: "custom-client-id",
			expected:          "custom-client-id",
		},
		{
			name:              "custom client ID not specified",
			idpCustomClientID: "",
			expected:          defaultIDPClientID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()

			if tt.idpCustomClientID != "" {
				viper.Set(config.IdentityProviderCustomClientIdKey, tt.idpCustomClientID)
			}

			got := getIDPClientID()

			if got != tt.expected {
				t.Fatalf("expected idp client ID %q, got %q", tt.expected, got)
			}
		})
	}
}
