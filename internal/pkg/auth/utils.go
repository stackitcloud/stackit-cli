package auth

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

type wellKnownConfig struct {
	Issuer                string `json:"issuer"`
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
}

func getIDPWellKnownConfigURL() (wellKnownConfigURL string, err error) {
	wellKnownConfigURL = defaultWellKnownConfig

	customWellKnownConfig := viper.GetString(config.IdentityProviderCustomWellKnownConfigurationKey)
	if customWellKnownConfig != "" {
		wellKnownConfigURL = customWellKnownConfig
		err := utils.ValidateURLDomain(wellKnownConfigURL)
		if err != nil {
			return "", fmt.Errorf("validate custom identity provider well-known configuration: %w", err)
		}
	}

	return wellKnownConfigURL, nil
}

func getIDPClientID() (string, error) {
	idpClientID := defaultCLIClientID

	customIDPClientID := viper.GetString(config.IdentityProviderCustomClientIdKey)
	if customIDPClientID != "" {
		idpClientID = customIDPClientID
	}

	return idpClientID, nil
}
