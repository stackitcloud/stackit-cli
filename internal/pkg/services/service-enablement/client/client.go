package client

import (
	"github.com/spf13/viper"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	genericclient "github.com/stackitcloud/stackit-cli/internal/pkg/generic-client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/stackitcloud/stackit-sdk-go/services/serviceenablement"
)

func ConfigureClient(p *print.Printer, cliVersion string) (*serviceenablement.APIClient, error) {
	return genericclient.ConfigureClientGeneric(p, cliVersion, viper.GetString(config.ServiceEnablementCustomEndpointKey), viper.GetString(config.RegionKey), genericclient.CreateApiClient[*serviceenablement.APIClient](serviceenablement.NewAPIClient))
}
