package genericclient

import (
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	sdkConfig "github.com/stackitcloud/stackit-sdk-go/core/config"
)

type CreateApiClient[T any] func(opts ...sdkConfig.ConfigurationOption) (T, error)

// ConfigureClientGeneric contains the generic code which needs to be executed in order to configure the api client.
// TODO: region parameter will be removed when every API implemented the new region
func ConfigureClientGeneric[T any](p *print.Printer, cliVersion, customEndpoint, region string, createApiClient CreateApiClient[T]) (T, error) {
	// return value if an error happens
	var zero T
	authCfgOption, err := auth.AuthenticationConfig(p, auth.AuthorizeUser)
	if err != nil {
		p.Debug(print.ErrorLevel, "configure authentication: %v", err)
		return zero, &errors.AuthError{}
	}
	cfgOptions := []sdkConfig.ConfigurationOption{
		utils.UserAgentConfigOption(cliVersion),
		authCfgOption,
	}

	if customEndpoint != "" {
		cfgOptions = append(cfgOptions, sdkConfig.WithEndpoint(customEndpoint))
	}

	// TODO: this will be removed when every API implemented the new region
	if region != "" {
		cfgOptions = append(cfgOptions, authCfgOption, sdkConfig.WithRegion(region))
	}

	if p.IsVerbosityDebug() {
		cfgOptions = append(cfgOptions,
			sdkConfig.WithMiddleware(print.RequestResponseCapturer(p, nil)),
		)
	}

	apiClient, err := createApiClient(cfgOptions...)
	if err != nil {
		p.Debug(print.ErrorLevel, "create new API client: %v", err)
		return zero, &errors.AuthError{}
	}

	return apiClient, nil
}
