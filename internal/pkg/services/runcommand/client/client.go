package client

import (
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	genericclient "github.com/stackitcloud/stackit-cli/internal/pkg/generic-client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/viper"
	runcommand "github.com/stackitcloud/stackit-sdk-go/services/runcommand/v2api"
)

func ConfigureClient(p *print.Printer, cliVersion string) (*runcommand.APIClient, error) {
	return genericclient.ConfigureClientGeneric(p, cliVersion, viper.GetString(config.RunCommandCustomEndpointKey), false, genericclient.CreateApiClient[*runcommand.APIClient](runcommand.NewAPIClient))
}
