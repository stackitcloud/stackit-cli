package client

import (
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	genericclient "github.com/stackitcloud/stackit-cli/internal/pkg/generic-client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/viper"
	albwaf "github.com/stackitcloud/stackit-sdk-go/services/albwaf/v1api"
)

func ConfigureClient(p *print.Printer, cliVersion string) (*albwaf.APIClient, error) {
	return genericclient.ConfigureClientGeneric(p, cliVersion, viper.GetString(config.AlbWafCustomEndpointKey), false, genericclient.CreateApiClient[*albwaf.APIClient](albwaf.NewAPIClient))
}
