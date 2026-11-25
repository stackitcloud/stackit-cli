package bucket

import (
	"testing"
)

func TestDefaultEndpoint(t *testing.T) {
	region := "eu01"
	expected := "https://object.storage.eu01.onstackit.cloud"
	if result := defaultEndpoint(region); result != expected {
		t.Fatalf("defaultEndpoint(%q) = %q, want %q", region, result, expected)
	}
}

func TestClientConfigResolvedEndpoint(t *testing.T) {
	tests := []struct {
		description string
		cfg         ClientConfig
		expected    string
		isValid     bool
	}{
		{
			description: "explicit endpoint wins",
			cfg: ClientConfig{
				Region:   "eu01",
				Endpoint: "https://custom.endpoint.example",
			},
			expected: "https://custom.endpoint.example",
			isValid:  true,
		},
		{
			description: "derived from region",
			cfg: ClientConfig{
				Region: "eu01",
			},
			expected: "https://object.storage.eu01.onstackit.cloud",
			isValid:  true,
		},
		{
			description: "no endpoint and no region",
			cfg:         ClientConfig{},
			isValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			result, err := tt.cfg.ResolvedEndpoint()
			if (err != nil) != !tt.isValid {
				t.Fatalf("ResolvedEndpoint() error = %v, wantErr %v", err, !tt.isValid)
			}
			if tt.isValid && result != tt.expected {
				t.Fatalf("ResolvedEndpoint() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestResolveFlagValue(t *testing.T) {
	// A set flag value always wins
	if result := resolveFlagValue("flag-value", "object_storage_access_key_id"); result != "flag-value" {
		t.Fatalf("resolveFlagValue() = %q, want %q", result, "flag-value")
	}
	// An empty flag value falls back to the config (empty in tests)
	if result := resolveFlagValue("", "object_storage_access_key_id"); result != "" {
		t.Fatalf("resolveFlagValue() = %q, want empty", result)
	}
}

func TestClientConfigFromReference(t *testing.T) {
	tests := []struct {
		description string
		flags       ClientFlags
		region      string
		ref         Reference
		expected    ClientConfig
	}{
		{
			description: "plain reference, no endpoint from ref",
			flags:       ClientFlags{},
			region:      "eu01",
			ref:         Reference{Bucket: "my-bucket"},
			expected:    ClientConfig{Region: "eu01"},
		},
		{
			description: "https reference provides endpoint and region",
			flags:       ClientFlags{},
			region:      "eu01",
			ref: Reference{
				Bucket:   "my-bucket",
				Endpoint: "https://object.storage.eu02.onstackit.cloud",
				Region:   "eu02",
			},
			expected: ClientConfig{
				Region:   "eu02",
				Endpoint: "https://object.storage.eu02.onstackit.cloud",
			},
		},
		{
			description: "explicit endpoint flag wins over reference",
			flags: ClientFlags{
				Endpoint: "https://custom.endpoint.example",
			},
			region: "eu01",
			ref: Reference{
				Bucket:   "my-bucket",
				Endpoint: "https://object.storage.eu02.onstackit.cloud",
				Region:   "eu02",
			},
			expected: ClientConfig{
				Region:   "eu02",
				Endpoint: "https://custom.endpoint.example",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			result := ClientConfigFromReference(tt.flags, tt.region, tt.ref)
			if result.Region != tt.expected.Region {
				t.Fatalf("Region = %q, want %q", result.Region, tt.expected.Region)
			}
			if result.Endpoint != tt.expected.Endpoint {
				t.Fatalf("Endpoint = %q, want %q", result.Endpoint, tt.expected.Endpoint)
			}
		})
	}
}
