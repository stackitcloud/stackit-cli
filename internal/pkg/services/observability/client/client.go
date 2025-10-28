package client

import (
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	genericclient "github.com/stackitcloud/stackit-cli/internal/pkg/generic-client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-sdk-go/services/observability"

	"github.com/spf13/viper"
)

func ConfigureClient(p *print.Printer, cliVersion string) (*observability.APIClient, error) {
	return genericclient.ConfigureClientGeneric(p, cliVersion, viper.GetString(config.ObservabilityCustomEndpointKey), true, genericclient.CreateApiClient[*observability.APIClient](observability.NewAPIClient))
}
