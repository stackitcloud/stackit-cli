package client

import (
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	genericclient "github.com/stackitcloud/stackit-cli/internal/pkg/generic-client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/viper"
	"github.com/stackitcloud/stackit-sdk-go/services/alb"
)

func ConfigureClient(p *print.Printer, cliVersion string) (*alb.APIClient, error) {
	return genericclient.ConfigureClientGeneric(p, cliVersion, viper.GetString(config.AlbCustomEndpoint), true, genericclient.CreateApiClient[*alb.APIClient](alb.NewAPIClient))
}
