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
		allowedUrlDomain  string
		isValid           bool
		expected          string
	}{
		{
			name:              "custom endpoint specified",
			idpCustomEndpoint: "https://example.stackit.cloud",
			allowedUrlDomain:  "stackit.cloud",
			isValid:           true,
			expected:          "https://example.stackit.cloud",
		},
		{
			name:              "custom endpoint outside STACKIT",
			idpCustomEndpoint: "https://www.very-suspicious-website.com/",
			allowedUrlDomain:  "stackit.cloud",
			isValid:           false,
		},
		{
			name:              "non-STACKIT custom endpoint invalid",
			idpCustomEndpoint: "https://www.very-suspicious-website.com/",
			allowedUrlDomain:  "stackit.cloud",
			isValid:           false,
		},
		{
			name:              "non-STACKIT custom endpoint valid",
			idpCustomEndpoint: "https://www.test.example.com/",
			allowedUrlDomain:  "example.com",
			isValid:           true,
			expected:          "https://www.test.example.com/",
		},
		{
			name:              "every URL valid",
			idpCustomEndpoint: "https://www.test.example.com/",
			allowedUrlDomain:  "",
			isValid:           true,
			expected:          "https://www.test.example.com/",
		},
		{
			name:              "custom endpoint not specified",
			idpCustomEndpoint: "",
			allowedUrlDomain:  "",
			isValid:           true,
			expected:          defaultIDPEndpoint,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			viper.Set(config.IdentityProviderCustomEndpointKey, tt.idpCustomEndpoint)
			viper.Set(config.AllowedUrlDomainKey, tt.allowedUrlDomain)

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

func TestGetIDPClientID(t *testing.T) {
	tests := []struct {
		name              string
		idpCustomClientID string
		isValid           bool
		expected          string
	}{
		{
			name:              "custom client ID specified",
			idpCustomClientID: "custom-client-id",
			isValid:           true,
			expected:          "custom-client-id",
		},
		{
			name:              "custom client ID not specified",
			idpCustomClientID: "",
			isValid:           true,
			expected:          defaultCLIClientID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			viper.Set(config.IdentityProviderCustomClientIdKey, tt.idpCustomClientID)

			got, err := getIDPClientID()

			if tt.isValid && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if !tt.isValid && err == nil {
				t.Fatalf("expected error, got none")
			}

			if got != tt.expected {
				t.Fatalf("expected idp client ID %q, got %q", tt.expected, got)
			}
		})
	}
}
