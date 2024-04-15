package client

import (
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/viper"
	sdkConfig "github.com/stackitcloud/stackit-sdk-go/core/config"
	"github.com/stackitcloud/stackit-sdk-go/services/authorization"
)

func ConfigureClient(p *print.Printer) (*authorization.APIClient, error) {
	var err error
	var apiClient *authorization.APIClient
	var cfgOptions []sdkConfig.ConfigurationOption

	authCfgOption, err := auth.AuthenticationConfig(p, auth.AuthorizeUser)
	if err != nil {
		p.Debug(print.ErrorLevel, "auth err: %v", err)
		return nil, &errors.AuthError{}
	}
	cfgOptions = append(cfgOptions, authCfgOption)

	customEndpoint := viper.GetString(config.AuthorizationCustomEndpointKey)

	if customEndpoint != "" {
		cfgOptions = append(cfgOptions, sdkConfig.WithEndpoint(customEndpoint))
	}

	apiClient, err = authorization.NewAPIClient(cfgOptions...)
	if err != nil {
		p.Debug(print.ErrorLevel, "auth err: %v", err)
		return nil, &errors.AuthError{}
	}

	return apiClient, nil
}
