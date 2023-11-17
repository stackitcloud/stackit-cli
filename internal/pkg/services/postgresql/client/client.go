package client

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	sdkConfig "github.com/stackitcloud/stackit-sdk-go/core/config"
	"github.com/stackitcloud/stackit-sdk-go/services/postgresql"
)

func ConfigureClient(cmd *cobra.Command) (*postgresql.APIClient, error) {
	var err error
	var apiClient *postgresql.APIClient
	var cfgOptions []sdkConfig.ConfigurationOption

	authCfgOption, err := auth.AuthenticationConfig(cmd, auth.AuthorizeUser)
	if err != nil {
		return nil, fmt.Errorf("configure authentication: %w", err)
	}
	cfgOptions = append(cfgOptions, authCfgOption, sdkConfig.WithRegion("eu01"))

	customEndpoint := viper.GetString(config.PostgreSQLCustomEndpointKey)

	if customEndpoint != "" {
		cfgOptions = append(cfgOptions, sdkConfig.WithEndpoint(customEndpoint))
	}

	apiClient, err = postgresql.NewAPIClient(cfgOptions...)
	if err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}

	return apiClient, nil
}
