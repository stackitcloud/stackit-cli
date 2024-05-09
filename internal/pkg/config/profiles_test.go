package config

import "testing"

func TestValidateProfile(t *testing.T) {
	tests := []struct {
		description string
		profile     string
		isValid     bool
	}{
		{
			description: "valid profile with letters",
			profile:     "myprofile",
			isValid:     true,
		},
		{
			description: "valid profile with uppercase letters",
			profile:     "myProfile",
			isValid:     true,
		},
		{
			description: "valid with letters and hyphen",
			profile:     "my-profile",
			isValid:     true,
		},
		{
			description: "valid with letters, numbers, and hyphen",
			profile:     "my-profile-123",
			isValid:     true,
		},
		{
			description: "valid empty",
			profile:     "",
			isValid:     true,
		},
		{
			description: "invalid with special characters",
			profile:     "my_profile",
			isValid:     false,
		},
		{
			description: "invalid with spaces",
			profile:     "my profile",
			isValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := validateProfile(tt.profile)
			if tt.isValid && err != nil {
				t.Errorf("expected profile to be valid but got error: %v", err)
			}
			if !tt.isValid && err == nil {
				t.Errorf("expected profile to be invalid but got no error")
			}
		})
	}
}
