package client

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-sdk-go/services/foo"
	// (...)
)

func ConfigureClient(cmd *cobra.Command) (*foo.APIClient, error) {
	var err error
	var apiClient foo.APIClient
	var cfgOptions []sdkConfig.ConfigurationOption

	authCfgOption, err := auth.AuthenticationConfig(cmd, auth.AuthorizeUser)
	if err != nil {
		return nil, &errors.AuthError{}
	}
	cfgOptions = append(cfgOptions, authCfgOption, sdkConfig.WithRegion("eu01")) // Configuring region is needed if "foo" is a regional API

	customEndpoint := viper.GetString(config.fooCustomEndpointKey)

	if customEndpoint != "" {
		cfgOptions = append(cfgOptions, sdkConfig.WithEndpoint(customEndpoint))
	}

	apiClient, err = foo.NewAPIClient(cfgOptions...)
	if err != nil {
		return nil, &errors.AuthError{}
	}

	return apiClient, nil
}
