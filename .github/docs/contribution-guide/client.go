package client

import (
	"github.com/spf13/viper"
	genericclient "github.com/stackitcloud/stackit-cli/internal/pkg/generic-client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-sdk-go/services/foo"
	// (...)
)

func ConfigureClient(p *print.Printer, cliVersion string) (*foo.APIClient, error) {
	return genericclient.ConfigureClientGeneric(p, cliVersion, viper.GetString(config.fooCustomEndpointKey), false, genericclient.CreateApiClient[*foo.APIClient](foo.NewAPIClient))
}
