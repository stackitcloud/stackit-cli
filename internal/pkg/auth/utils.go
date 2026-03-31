package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/viper"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

type wellKnownConfig struct {
	Issuer                      string   `json:"issuer"`
	AuthorizationEndpoint       string   `json:"authorization_endpoint"`
	TokenEndpoint               string   `json:"token_endpoint"`
	DeviceAuthorizationEndpoint string   `json:"device_authorization_endpoint"`
	GrantTypesSupported         []string `json:"grant_types_supported"`
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

func retrieveIDPWellKnownConfig(p *print.Printer) (*wellKnownConfig, error) {
	idpWellKnownConfigURL, err := getIDPWellKnownConfigURL()
	if err != nil {
		return nil, fmt.Errorf("get IDP well-known configuration: %w", err)
	}
	if idpWellKnownConfigURL != defaultWellKnownConfig {
		p.Warn("You are using a custom identity provider well-known configuration (%s) for authentication.\n", idpWellKnownConfigURL)
		err := p.PromptForEnter("Press Enter to proceed with the login...")
		if err != nil {
			return nil, err
		}
	}

	p.Debug(print.DebugLevel, "get IDP well-known configuration from %s", idpWellKnownConfigURL)
	httpClient := &http.Client{}
	idpWellKnownConfig, err := parseWellKnownConfiguration(httpClient, idpWellKnownConfigURL)
	if err != nil {
		return nil, fmt.Errorf("parse IDP well-known configuration: %w", err)
	}
	return idpWellKnownConfig, nil
}

// parseWellKnownConfiguration gets the well-known OpenID configuration from the provided URL and returns it as a JSON
// the method also stores the IDP token endpoint in the authentication storage
func parseWellKnownConfiguration(httpClient apiClient, wellKnownConfigURL string) (wellKnownConfig *wellKnownConfig, err error) {
	req, _ := http.NewRequest("GET", wellKnownConfigURL, http.NoBody)
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("make the request: %w", err)
	}

	// Process the response
	defer func() {
		closeErr := res.Body.Close()
		if closeErr != nil {
			err = fmt.Errorf("close response body: %w", closeErr)
		}
	}()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	err = json.Unmarshal(body, &wellKnownConfig)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	if wellKnownConfig == nil {
		return nil, fmt.Errorf("nil well-known configuration response")
	}
	if wellKnownConfig.Issuer == "" {
		return nil, fmt.Errorf("found no issuer")
	}
	if wellKnownConfig.AuthorizationEndpoint == "" {
		return nil, fmt.Errorf("found no authorization endpoint")
	}
	if wellKnownConfig.TokenEndpoint == "" {
		return nil, fmt.Errorf("found no token endpoint")
	}

	err = SetAuthField(IDP_TOKEN_ENDPOINT, wellKnownConfig.TokenEndpoint)
	if err != nil {
		return nil, fmt.Errorf("set token endpoint in the authentication storage: %w", err)
	}
	return wellKnownConfig, err
}
