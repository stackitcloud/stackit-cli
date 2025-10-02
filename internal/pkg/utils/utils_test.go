package utils

import (
	"reflect"
	"testing"
	"time"

	sdkConfig "github.com/stackitcloud/stackit-sdk-go/core/config"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"

	"github.com/spf13/viper"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
)

func TestConvertInt64PToFloat64P(t *testing.T) {
	tests := []struct {
		name     string
		input    *int64
		expected *float64
	}{
		{
			name:     "positive",
			input:    Int64Ptr(1),
			expected: Float64Ptr(1.0),
		},
		{
			name:     "negative",
			input:    Int64Ptr(-1),
			expected: Float64Ptr(-1.0),
		},
		{
			name:     "zero",
			input:    Int64Ptr(0),
			expected: Float64Ptr(0.0),
		},
		{
			name:     "nil",
			input:    nil,
			expected: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := ConvertInt64PToFloat64P(tt.input)

			if expected == nil && tt.expected == nil && tt.input == nil {
				return
			}

			if *expected != *tt.expected {
				t.Errorf("ConvertInt64ToFloat64() = %v, want %v", *expected, *tt.expected)
			}
		})
	}
}

func TestValidateURLDomain(t *testing.T) {
	tests := []struct {
		name             string
		allowedUrlDomain string
		isValid          bool
		input            string
	}{
		{
			name:             "STACKIT URL valid",
			allowedUrlDomain: "stackit.cloud",
			input:            "https://example.stackit.cloud",
			isValid:          true,
		},
		{
			name:             "STACKIT URL invalid",
			allowedUrlDomain: "example.com",
			input:            "https://example.stackit.cloud",
			isValid:          false,
		},
		{
			name:             "non-STACKIT URL invalid",
			allowedUrlDomain: "stackit.cloud",
			input:            "https://www.very-suspicious-website.com/",
			isValid:          false,
		},
		{
			name:             "non-STACKIT URL valid",
			allowedUrlDomain: "example.com",
			input:            "https://www.test.example.com/",
			isValid:          true,
		},
		{
			name:             "every URL valid",
			allowedUrlDomain: "",
			input:            "https://www.test.example.com/",
			isValid:          true,
		},
		{
			name:    "invalid URL",
			input:   "",
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			viper.Set(config.AllowedUrlDomainKey, tt.allowedUrlDomain)

			err := ValidateURLDomain(tt.input)
			if tt.isValid && err != nil {
				t.Errorf("expected URL to be valid, got error: %v", err)
			}
			if !tt.isValid && err == nil {
				t.Errorf("expected URL to be invalid, got no error")
			}
		})
	}
}

func TestUserAgentConfigOption(t *testing.T) {
	type args struct {
		providerVersion string
	}
	tests := []struct {
		name string
		args args
		want sdkConfig.ConfigurationOption
	}{
		{
			name: "TestUserAgentConfigOption",
			args: args{
				providerVersion: "1.0.0",
			},
			want: sdkConfig.WithUserAgent("stackit-cli/1.0.0"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientConfigActual := sdkConfig.Configuration{}
			err := tt.want(&clientConfigActual)
			if err != nil {
				t.Errorf("error configuring client: %v", err)
			}

			clientConfigExpected := sdkConfig.Configuration{}
			err = UserAgentConfigOption(tt.args.providerVersion)(&clientConfigExpected)
			if err != nil {
				t.Errorf("error configuring client: %v", err)
			}

			if !reflect.DeepEqual(clientConfigActual, clientConfigExpected) {
				t.Errorf("UserAgentConfigOption() = %v, want %v", clientConfigActual, clientConfigExpected)
			}
		})
	}
}

func TestConvertStringMapToInterfaceMap(t *testing.T) {
	tests := []struct {
		name     string
		input    *map[string]string
		expected *map[string]interface{}
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
		{
			name:     "empty map",
			input:    &map[string]string{},
			expected: nil,
		},
		{
			name: "single key-value pair",
			input: &map[string]string{
				"key1": "value1",
			},
			expected: &map[string]interface{}{
				"key1": "value1",
			},
		},
		{
			name: "multiple key-value pairs",
			input: &map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
			expected: &map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
		},
		{
			name: "special characters in values",
			input: &map[string]string{
				"key1": "value with spaces",
				"key2": "value,with,commas",
				"key3": "value\nwith\nnewlines",
			},
			expected: &map[string]interface{}{
				"key1": "value with spaces",
				"key2": "value,with,commas",
				"key3": "value\nwith\nnewlines",
			},
		},
		{
			name: "empty values",
			input: &map[string]string{
				"key1": "",
				"key2": "value2",
			},
			expected: &map[string]interface{}{
				"key1": "",
				"key2": "value2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertStringMapToInterfaceMap(tt.input)

			// Check if both are nil
			if result == nil && tt.expected == nil {
				return
			}

			// Check if one is nil and other isn't
			if (result == nil && tt.expected != nil) || (result != nil && tt.expected == nil) {
				t.Errorf("ConvertStringMapToInterfaceMap() = %v, want %v", result, tt.expected)
				return
			}

			// Compare maps
			if len(*result) != len(*tt.expected) {
				t.Errorf("ConvertStringMapToInterfaceMap() map length = %d, want %d", len(*result), len(*tt.expected))
				return
			}

			for k, v := range *result {
				expectedVal, ok := (*tt.expected)[k]
				if !ok {
					t.Errorf("ConvertStringMapToInterfaceMap() unexpected key %s in result", k)
					continue
				}
				if v != expectedVal {
					t.Errorf("ConvertStringMapToInterfaceMap() value for key %s = %v, want %v", k, v, expectedVal)
				}
			}
		})
	}
}

func TestConvertToBase64PatchedServer(t *testing.T) {
	now := time.Now()
	userData := []byte("test")
	emptyUserData := []byte("")

	tests := []struct {
		name     string
		input    *iaas.Server
		expected *Base64PatchedServer
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
		{
			name: "server with user data",
			input: &iaas.Server{
				Id:               Ptr("server-123"),
				Name:             Ptr("test-server"),
				Status:           Ptr("ACTIVE"),
				AvailabilityZone: Ptr("eu01-1"),
				MachineType:      Ptr("t1.1"),
				UserData:         &userData,
				CreatedAt:        &now,
				PowerStatus:      Ptr("RUNNING"),
				AffinityGroup:    Ptr("group-1"),
				ImageId:          Ptr("image-123"),
				KeypairName:      Ptr("keypair-1"),
			},
			expected: &Base64PatchedServer{
				Id:               Ptr("server-123"),
				Name:             Ptr("test-server"),
				Status:           Ptr("ACTIVE"),
				AvailabilityZone: Ptr("eu01-1"),
				MachineType:      Ptr("t1.1"),
				UserData:         Ptr(Base64Bytes(userData)),
				CreatedAt:        &now,
				PowerStatus:      Ptr("RUNNING"),
				AffinityGroup:    Ptr("group-1"),
				ImageId:          Ptr("image-123"),
				KeypairName:      Ptr("keypair-1"),
			},
		},
		{
			name: "server with empty user data",
			input: &iaas.Server{
				Id:               Ptr("server-456"),
				Name:             Ptr("test-server-2"),
				Status:           Ptr("STOPPED"),
				AvailabilityZone: Ptr("eu01-2"),
				MachineType:      Ptr("t1.2"),
				UserData:         &emptyUserData,
			},
			expected: &Base64PatchedServer{
				Id:               Ptr("server-456"),
				Name:             Ptr("test-server-2"),
				Status:           Ptr("STOPPED"),
				AvailabilityZone: Ptr("eu01-2"),
				MachineType:      Ptr("t1.2"),
				UserData:         Ptr(Base64Bytes(emptyUserData)),
			},
		},
		{
			name: "server without user data",
			input: &iaas.Server{
				Id:               Ptr("server-789"),
				Name:             Ptr("test-server-3"),
				Status:           Ptr("CREATING"),
				AvailabilityZone: Ptr("eu01-3"),
				MachineType:      Ptr("t1.3"),
				UserData:         nil,
			},
			expected: &Base64PatchedServer{
				Id:               Ptr("server-789"),
				Name:             Ptr("test-server-3"),
				Status:           Ptr("CREATING"),
				AvailabilityZone: Ptr("eu01-3"),
				MachineType:      Ptr("t1.3"),
				UserData:         nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertToBase64PatchedServer(tt.input)

			if result == nil && tt.expected == nil {
				return
			}

			if (result == nil && tt.expected != nil) || (result != nil && tt.expected == nil) {
				t.Errorf("ConvertToBase64PatchedServer() = %v, want %v", result, tt.expected)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ConvertToBase64PatchedServer() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestConvertToBase64PatchedServers(t *testing.T) {
	now := time.Now()
	userData1 := []byte("test1")
	userData2 := []byte("test2")
	emptyUserData := []byte("")

	tests := []struct {
		name     string
		input    []iaas.Server
		expected []Base64PatchedServer
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
		{
			name:     "empty slice",
			input:    []iaas.Server{},
			expected: []Base64PatchedServer{},
		},
		{
			name: "single server with user data",
			input: []iaas.Server{
				{
					Id:               Ptr("server-1"),
					Name:             Ptr("test-server-1"),
					Status:           Ptr("ACTIVE"),
					MachineType:      Ptr("t1.1"),
					AvailabilityZone: Ptr("eu01-1"),
					UserData:         &userData1,
					CreatedAt:        &now,
				},
			},
			expected: []Base64PatchedServer{
				{
					Id:               Ptr("server-1"),
					Name:             Ptr("test-server-1"),
					Status:           Ptr("ACTIVE"),
					MachineType:      Ptr("t1.1"),
					AvailabilityZone: Ptr("eu01-1"),
					UserData:         Ptr(Base64Bytes(userData1)),
					CreatedAt:        &now,
				},
			},
		},
		{
			name: "multiple servers mixed",
			input: []iaas.Server{
				{
					Id:               Ptr("server-1"),
					Name:             Ptr("test-server-1"),
					Status:           Ptr("ACTIVE"),
					MachineType:      Ptr("t1.1"),
					AvailabilityZone: Ptr("eu01-1"),
					UserData:         &userData1,
					CreatedAt:        &now,
				},
				{
					Id:               Ptr("server-2"),
					Name:             Ptr("test-server-2"),
					Status:           Ptr("STOPPED"),
					MachineType:      Ptr("t1.2"),
					AvailabilityZone: Ptr("eu01-2"),
					UserData:         &userData2,
				},
				{
					Id:               Ptr("server-3"),
					Name:             Ptr("test-server-3"),
					Status:           Ptr("CREATING"),
					MachineType:      Ptr("t1.3"),
					AvailabilityZone: Ptr("eu01-3"),
					UserData:         &emptyUserData,
				},
				{
					Id:               Ptr("server-4"),
					Name:             Ptr("test-server-4"),
					Status:           Ptr("ERROR"),
					MachineType:      Ptr("t1.4"),
					AvailabilityZone: Ptr("eu01-4"),
					UserData:         nil,
				},
			},
			expected: []Base64PatchedServer{
				{
					Id:               Ptr("server-1"),
					Name:             Ptr("test-server-1"),
					Status:           Ptr("ACTIVE"),
					MachineType:      Ptr("t1.1"),
					AvailabilityZone: Ptr("eu01-1"),
					UserData:         Ptr(Base64Bytes(userData1)),
					CreatedAt:        &now,
				},
				{
					Id:               Ptr("server-2"),
					Name:             Ptr("test-server-2"),
					Status:           Ptr("STOPPED"),
					MachineType:      Ptr("t1.2"),
					AvailabilityZone: Ptr("eu01-2"),
					UserData:         Ptr(Base64Bytes(userData2)),
				},
				{
					Id:               Ptr("server-3"),
					Name:             Ptr("test-server-3"),
					Status:           Ptr("CREATING"),
					MachineType:      Ptr("t1.3"),
					AvailabilityZone: Ptr("eu01-3"),
					UserData:         Ptr(Base64Bytes(emptyUserData)),
				},
				{
					Id:               Ptr("server-4"),
					Name:             Ptr("test-server-4"),
					Status:           Ptr("ERROR"),
					MachineType:      Ptr("t1.4"),
					AvailabilityZone: Ptr("eu01-4"),
					UserData:         nil,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertToBase64PatchedServers(tt.input)

			if result == nil && tt.expected == nil {
				return
			}

			if (result == nil && tt.expected != nil) || (result != nil && tt.expected == nil) {
				t.Errorf("ConvertToBase64PatchedServers() = %v, want %v", result, tt.expected)
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("ConvertToBase64PatchedServers() length = %d, want %d", len(result), len(tt.expected))
				return
			}

			for i, server := range result {
				if !reflect.DeepEqual(server, tt.expected[i]) {
					t.Errorf("ConvertToBase64PatchedServers() [%d] = %v, want %v", i, server, tt.expected[i])
				}
			}
		})
	}
}

func TestBase64Bytes_MarshalYAML(t *testing.T) {
	tests := []struct {
		name     string
		input    Base64Bytes
		expected interface{}
	}{
		{
			name:     "empty bytes",
			input:    Base64Bytes{},
			expected: "",
		},
		{
			name:     "nil bytes",
			input:    Base64Bytes(nil),
			expected: "",
		},
		{
			name:     "simple text",
			input:    Base64Bytes("test"),
			expected: "dGVzdA==",
		},
		{
			name:     "special characters",
			input:    Base64Bytes("test@#$%"),
			expected: "dGVzdEAjJCU=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.input.MarshalYAML()
			if err != nil {
				t.Errorf("MarshalYAML() error = %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("MarshalYAML() = %v, want %v", result, tt.expected)
			}
		})
	}
}
func TestGetSliceFromPointer(t *testing.T) {
	tests := []struct {
		name     string
		input    *[]string
		expected []string
	}{
		{
			name:     "nil pointer",
			input:    nil,
			expected: []string{},
		},
		{
			name: "pointer to nil slice",
			input: func() *[]string {
				var s []string
				return &s
			}(),
			expected: []string{},
		},
		{
			name:     "empty slice",
			input:    &[]string{},
			expected: []string{},
		},
		{
			name:     "populated slice",
			input:    &[]string{"item1", "item2"},
			expected: []string{"item1", "item2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetSliceFromPointer(tt.input)

			if result == nil {
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GetSliceFromPointer() = %v, want %v", result, tt.expected)
			}
		})
	}
}
