package client

import (
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/viper"
	sdkConfig "github.com/stackitcloud/stackit-sdk-go/core/config"
	"github.com/stackitcloud/stackit-sdk-go/services/logme"
)

func ConfigureClient(p *print.Printer) (*logme.APIClient, error) {
	var err error
	var apiClient *logme.APIClient
	var cfgOptions []sdkConfig.ConfigurationOption

	authCfgOption, err := auth.AuthenticationConfig(p, auth.AuthorizeUser)
	if err != nil {
		return nil, &errors.AuthError{}
	}
	cfgOptions = append(cfgOptions, authCfgOption, sdkConfig.WithRegion("eu01"))

	customEndpoint := viper.GetString(config.LogMeCustomEndpointKey)

	if customEndpoint != "" {
		cfgOptions = append(cfgOptions, sdkConfig.WithEndpoint(customEndpoint))
	}

	apiClient, err = logme.NewAPIClient(cfgOptions...)
	if err != nil {
		return nil, &errors.AuthError{}
	}

	return apiClient, nil
}
