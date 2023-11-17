package client

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	sdkConfig "github.com/stackitcloud/stackit-sdk-go/core/config"
	"github.com/stackitcloud/stackit-sdk-go/services/dns"
)

func ConfigureClient(cmd *cobra.Command) (*dns.APIClient, error) {
	var err error
	var apiClient *dns.APIClient
	var cfgOptions []sdkConfig.ConfigurationOption

	authCfgOption, err := auth.AuthenticationConfig(cmd, auth.AuthorizeUser)
	if err != nil {
		return nil, fmt.Errorf("configure authentication: %w", err)
	}
	cfgOptions = append(cfgOptions, authCfgOption)

	customEndpoint := viper.GetString(config.DNSCustomEndpointKey)

	if customEndpoint != "" {
		cfgOptions = append(cfgOptions, sdkConfig.WithEndpoint(customEndpoint))
	}

	apiClient, err = dns.NewAPIClient(cfgOptions...)
	if err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}

	return apiClient, nil
}
