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
	Issuer                string `json:"issuer"`
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
}

// getIDPWellKnownConfig gets the well-known OpenID configuration and returns it as a JSON
// it uses a default URL for the well-known configuration, unless a custom one is provided in the CLI configuration
// the method also stores the IDP token endpoint in the authentication storage
func getIDPWellKnownConfig(p *print.Printer) (wellKnownConfigResponse *wellKnownConfig, customIDP bool, err error) {
	wellKnownConfigURL := defaultWellKnownConfig
	customIDP = false

	customWellKnownConfig := viper.GetString(config.IdentityProviderCustomWellKnownConfigurationKey)
	if customWellKnownConfig != "" {
		customIDP = true
		wellKnownConfigURL = customWellKnownConfig
		err := utils.ValidateURLDomain(wellKnownConfigURL)
		if err != nil {
			return nil, true, fmt.Errorf("validate custom identity provider well-known configuration: %w", err)
		}
	}

	p.Debug(print.DebugLevel, "get IDP well-known configuration from %s", wellKnownConfigURL)
	req, _ := http.NewRequest("GET", wellKnownConfigURL, nil)
	httpClient := &http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, false, fmt.Errorf("make the request: %w", err)
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
		return nil, false, fmt.Errorf("read response body: %w", err)
	}

	err = json.Unmarshal(body, &wellKnownConfigResponse)
	if err != nil {
		return nil, false, fmt.Errorf("unmarshal response: %w", err)
	}
	if wellKnownConfigResponse == nil {
		return nil, false, fmt.Errorf("nil well-known configuration response")
	}
	if wellKnownConfigResponse.Issuer == "" {
		return nil, false, fmt.Errorf("found no issuer")
	}
	if wellKnownConfigResponse.TokenEndpoint == "" {
		return nil, false, fmt.Errorf("found no token endpoint")
	}
	if wellKnownConfigResponse.AuthorizationEndpoint == "" {
		return nil, false, fmt.Errorf("found no authorization endpoint")
	}
	if wellKnownConfigResponse.TokenEndpoint == "" {
		return nil, false, fmt.Errorf("found no token endpoint")
	}

	err = SetAuthField(IDP_TOKEN_ENDPOINT, wellKnownConfigResponse.TokenEndpoint)
	if err != nil {
		return nil, false, fmt.Errorf("set token endpoint in the authentication storage: %w", err)
	}

	return wellKnownConfigResponse, customIDP, nil
}

func getIDPClientID() (string, error) {
	idpClientID := defaultCLIClientID

	customIDPClientID := viper.GetString(config.IdentityProviderCustomClientIdKey)
	if customIDPClientID != "" {
		idpClientID = customIDPClientID
	}

	return idpClientID, nil
}
