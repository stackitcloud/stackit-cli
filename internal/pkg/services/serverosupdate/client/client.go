package client

import (
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/viper"
	sdkConfig "github.com/stackitcloud/stackit-sdk-go/core/config"
	"github.com/stackitcloud/stackit-sdk-go/services/serverupdate"
)

func ConfigureClient(p *print.Printer, cliVersion string) (*serverupdate.APIClient, error) {
	authCfgOption, err := auth.AuthenticationConfig(p, auth.AuthorizeUser)
	if err != nil {
		p.Debug(print.ErrorLevel, "configure authentication: %v", err)
		return nil, &errors.AuthError{}
	}
	cfgOptions := []sdkConfig.ConfigurationOption{
		utils.UserAgentConfigOption(cliVersion),
		authCfgOption,
	}

	customEndpoint := viper.GetString(config.ServerOsUpdateCustomEndpointKey)
	if customEndpoint != "" {
		cfgOptions = append(cfgOptions, sdkConfig.WithEndpoint(customEndpoint))
	} else {
		cfgOptions = append(cfgOptions, authCfgOption)
	}

	if p.IsVerbosityDebug() {
		cfgOptions = append(cfgOptions,
			sdkConfig.WithMiddleware(print.RequestResponseCapturer(p, nil)),
		)
	}

	apiClient, err := serverupdate.NewAPIClient(cfgOptions...)
	if err != nil {
		p.Debug(print.ErrorLevel, "create new API client: %v", err)
		return nil, &errors.AuthError{}
	}

	return apiClient, nil
}
