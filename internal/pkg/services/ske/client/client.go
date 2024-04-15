package client

import (
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/viper"
	sdkConfig "github.com/stackitcloud/stackit-sdk-go/core/config"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
)

func ConfigureClient(p *print.Printer) (*ske.APIClient, error) {
	var err error
	var apiClient *ske.APIClient
	var cfgOptions []sdkConfig.ConfigurationOption

	authCfgOption, err := auth.AuthenticationConfig(p, auth.AuthorizeUser)
	if err != nil {
		p.Debug(print.ErrorLevel, "auth err: %v", err)
		return nil, &errors.AuthError{}
	}
	cfgOptions = append(cfgOptions, authCfgOption)

	customEndpoint := viper.GetString(config.SKECustomEndpointKey)
	if customEndpoint != "" {
		cfgOptions = append(cfgOptions, sdkConfig.WithEndpoint(customEndpoint))
	} else {
		cfgOptions = append(cfgOptions, authCfgOption, sdkConfig.WithRegion("eu01"))
	}

	apiClient, err = ske.NewAPIClient(cfgOptions...)
	if err != nil {
		p.Debug(print.ErrorLevel, "auth err: %v", err)
		return nil, &errors.AuthError{}
	}

	return apiClient, nil
}
