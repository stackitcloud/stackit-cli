package client

import (
	"github.com/spf13/viper"
	genericclient "github.com/stackitcloud/stackit-cli/internal/pkg/generic-client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-sdk-go/services/foo"
	// (...)
)

// TODO: region parameter will be removed when every service implemented the region adjustment
func ConfigureClient(p *print.Printer, cliVersion string) (*foo.APIClient, error) {
	return genericclient.ConfigureClientGeneric(p, cliVersion, viper.GetString(config.fooCustomEndpointKey), "", genericclient.CreateApiClient[*foo.APIClient](foo.NewAPIClient))
}
