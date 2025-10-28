package client

import (
	genericclient "github.com/stackitcloud/stackit-cli/internal/pkg/generic-client"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/viper"
	"github.com/stackitcloud/stackit-sdk-go/services/dns"
)

func ConfigureClient(p *print.Printer, cliVersion string) (*dns.APIClient, error) {
	return genericclient.ConfigureClientGeneric(p, cliVersion, viper.GetString(config.DNSCustomEndpointKey), false, genericclient.CreateApiClient[*dns.APIClient](dns.NewAPIClient))
}
