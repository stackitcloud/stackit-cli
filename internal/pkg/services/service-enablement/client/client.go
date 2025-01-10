package client

import (
	"github.com/spf13/viper"
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	sdkConfig "github.com/stackitcloud/stackit-sdk-go/core/config"
	"github.com/stackitcloud/stackit-sdk-go/services/serviceenablement"
)

func ConfigureClient(p *print.Printer) (*serviceenablement.APIClient, error) {
	var err error
	var apiClient *serviceenablement.APIClient
	var cfgOptions []sdkConfig.ConfigurationOption

	authCfgOption, err := auth.AuthenticationConfig(p, auth.AuthorizeUser)
	if err != nil {
		p.Debug(print.ErrorLevel, "configure authentication: %v", err)
		return nil, &errors.AuthError{}
	}
	cfgOptions = append(cfgOptions, authCfgOption)

	customEndpoint := viper.GetString(config.ServiceEnablementCustomEndpointKey)
	if customEndpoint != "" {
		cfgOptions = append(cfgOptions, sdkConfig.WithEndpoint(customEndpoint))
	} else {
		region := viper.GetString(config.RegionKey)
		cfgOptions = append(cfgOptions, sdkConfig.WithRegion(region))
	}

	if p.IsVerbosityDebug() {
		cfgOptions = append(cfgOptions,
			sdkConfig.WithMiddleware(print.RequestResponseCapturer(p, nil)),
		)
	}

	apiClient, err = serviceenablement.NewAPIClient(cfgOptions...)
	if err != nil {
		p.Debug(print.ErrorLevel, "create new API client: %v", err)
		return nil, &errors.AuthError{}
	}

	return apiClient, nil
}
