package client

import (
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/viper"
	sdkConfig "github.com/stackitcloud/stackit-sdk-go/core/config"
	"github.com/stackitcloud/stackit-sdk-go/services/opensearch"
)

func ConfigureClient(p *print.Printer, cliVersion string) (*opensearch.APIClient, error) {
	authCfgOption, err := auth.AuthenticationConfig(p, auth.AuthorizeUser)
	if err != nil {
		p.Debug(print.ErrorLevel, "configure authentication: %v", err)
		return nil, &errors.AuthError{}
	}
	region := viper.GetString(config.RegionKey)
	cfgOptions := []sdkConfig.ConfigurationOption{
		utils.UserAgentConfigOption(cliVersion),
		sdkConfig.WithRegion(region),
		authCfgOption,
	}

	customEndpoint := viper.GetString(config.OpenSearchCustomEndpointKey)

	if customEndpoint != "" {
		cfgOptions = append(cfgOptions, sdkConfig.WithEndpoint(customEndpoint))
	}

	if p.IsVerbosityDebug() {
		cfgOptions = append(cfgOptions,
			sdkConfig.WithMiddleware(print.RequestResponseCapturer(p, nil)),
		)
	}

	apiClient, err := opensearch.NewAPIClient(cfgOptions...)
	if err != nil {
		p.Debug(print.ErrorLevel, "create new API client: %v", err)
		return nil, &errors.AuthError{}
	}

	return apiClient, nil
}
