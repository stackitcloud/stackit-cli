package client

import (
	"stackit/internal/pkg/auth"
	"stackit/internal/pkg/config"
	"stackit/internal/pkg/errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	sdkConfig "github.com/stackitcloud/stackit-sdk-go/core/config"
	"github.com/stackitcloud/stackit-sdk-go/services/postgresflex"
)

func ConfigureClient(cmd *cobra.Command) (*postgresflex.APIClient, error) {
	var err error
	var apiClient *postgresflex.APIClient
	var cfgOptions []sdkConfig.ConfigurationOption

	authCfgOption, err := auth.AuthenticationConfig(cmd, auth.AuthorizeUser)
	if err != nil {
		return nil, &errors.AuthError{}
	}
	cfgOptions = append(cfgOptions, authCfgOption, sdkConfig.WithRegion("eu01"))

	customEndpoint := viper.GetString(config.PostgresFlexCustomEndpointKey)

	if customEndpoint != "" {
		cfgOptions = append(cfgOptions, sdkConfig.WithEndpoint(customEndpoint))
	}

	apiClient, err = postgresflex.NewAPIClient(cfgOptions...)
	if err != nil {
		return nil, &errors.AuthError{}
	}

	return apiClient, nil
}
