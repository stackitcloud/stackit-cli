package auth

import (
	"github.com/spf13/viper"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
)

func getIDPEndpoint() string {
	customIDPEndpoint := viper.GetString(config.IdentityProviderCustomEndpointKey)
	if customIDPEndpoint != "" {
		return customIDPEndpoint
	}

	return defaultIDPEndpoint
}

func getIDPClientID() string {
	customIDPClientID := viper.GetString(config.IdentityProviderCustomClientIdKey)
	if customIDPClientID != "" {
		return customIDPClientID
	}

	return defaultIDPClientID
}
