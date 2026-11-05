package client

import (
	"github.com/spf13/viper"
	registry "github.com/stackitcloud/stackit-sdk-go/services/registry/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	genericclient "github.com/stackitcloud/stackit-cli/internal/pkg/generic-client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
)

func ConfigureClient(p *print.Printer, cliVersion string) (*registry.APIClient, error) {
	return genericclient.ConfigureClientGeneric(p, cliVersion, viper.GetString(config.RegistryCustomEndpointKey), false, genericclient.CreateApiClient[*registry.APIClient](registry.NewAPIClient))
}
