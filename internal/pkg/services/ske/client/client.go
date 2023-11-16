package client

import (
	"stackit/internal/pkg/auth"
	"stackit/internal/pkg/config"
	"stackit/internal/pkg/errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	sdkConfig "github.com/stackitcloud/stackit-sdk-go/core/config"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
)

func ConfigureClient(cmd *cobra.Command) (*ske.APIClient, error) {
	var err error
	var apiClient *ske.APIClient
	var cfgOptions []sdkConfig.ConfigurationOption

	authCfgOption, err := auth.AuthenticationConfig(cmd, auth.AuthorizeUser)
	if err != nil {
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
		return nil, &errors.AuthError{}
	}

	return apiClient, nil
}
