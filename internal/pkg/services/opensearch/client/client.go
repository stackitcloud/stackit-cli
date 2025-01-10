package client

import (
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/viper"
	sdkConfig "github.com/stackitcloud/stackit-sdk-go/core/config"
	"github.com/stackitcloud/stackit-sdk-go/services/opensearch"
)

func ConfigureClient(p *print.Printer) (*opensearch.APIClient, error) {
	var err error
	var apiClient *opensearch.APIClient
	var cfgOptions []sdkConfig.ConfigurationOption

	authCfgOption, err := auth.AuthenticationConfig(p, auth.AuthorizeUser)
	if err != nil {
		p.Debug(print.ErrorLevel, "configure authentication: %v", err)
		return nil, &errors.AuthError{}
	}
	region := viper.GetString(config.RegionKey)
	cfgOptions = append(cfgOptions, authCfgOption, sdkConfig.WithRegion(region))

	customEndpoint := viper.GetString(config.OpenSearchCustomEndpointKey)

	if customEndpoint != "" {
		cfgOptions = append(cfgOptions, sdkConfig.WithEndpoint(customEndpoint))
	}

	if p.IsVerbosityDebug() {
		cfgOptions = append(cfgOptions,
			sdkConfig.WithMiddleware(print.RequestResponseCapturer(p, nil)),
		)
	}

	apiClient, err = opensearch.NewAPIClient(cfgOptions...)
	if err != nil {
		p.Debug(print.ErrorLevel, "create new API client: %v", err)
		return nil, &errors.AuthError{}
	}

	return apiClient, nil
}
