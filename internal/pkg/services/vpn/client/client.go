package client

import (
	vpn "github.com/stackitcloud/stackit-sdk-go/services/vpn/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	genericclient "github.com/stackitcloud/stackit-cli/internal/pkg/generic-client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/viper"
)

func ConfigureClient(p *print.Printer, cliVersion string) (*vpn.APIClient, error) {
	return genericclient.ConfigureClientGeneric(p, cliVersion, viper.GetString(config.VpnCustomEndpointKey), false, genericclient.CreateApiClient[*vpn.APIClient](vpn.NewAPIClient))
}
