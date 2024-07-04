package auth

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func getIDPEndpoint() (string, error) {
	idpEndpoint := defaultIDPEndpoint

	customIDPEndpoint := viper.GetString(config.IdentityProviderCustomEndpointKey)
	if customIDPEndpoint != "" {
		idpEndpoint = customIDPEndpoint
	}
	err := utils.ValidateSTACKITURL(idpEndpoint)
	if err != nil {
		return "", fmt.Errorf("validate custom identity provider endpoint: %w", err)
	}

	return idpEndpoint, nil
}

func getIDPClientID() string {
	customIDPClientID := viper.GetString(config.IdentityProviderCustomClientIdKey)
	if customIDPClientID != "" {
		return customIDPClientID
	}

	return defaultIDPClientID
}
