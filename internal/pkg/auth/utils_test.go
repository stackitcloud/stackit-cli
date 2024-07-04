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
		isValid           bool
		expected          string
	}{
		{
			name:              "custom endpoint specified",
			idpCustomEndpoint: "https://example.stackit.cloud",
			isValid:           true,
			expected:          "https://example.stackit.cloud",
		},
		{
			name:              "custom endpoint outside STACKIT",
			idpCustomEndpoint: "https://www.very-suspicious-website.com/",
			isValid:           false,
		},
		{
			name:              "custom endpoint not specified",
			idpCustomEndpoint: "",
			isValid:           true,
			expected:          defaultIDPEndpoint,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			viper.Set(config.IdentityProviderCustomEndpointKey, tt.idpCustomEndpoint)

			got, err := getIDPEndpoint()

			if tt.isValid && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if !tt.isValid && err == nil {
				t.Fatalf("expected error, got none")
			}

			if got != tt.expected {
				t.Fatalf("expected idp endpoint %q, got %q", tt.expected, got)
			}
		})
	}
}
