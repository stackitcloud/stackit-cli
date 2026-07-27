package client

import (
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	genericclient "github.com/stackitcloud/stackit-cli/internal/pkg/generic-client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/viper"
	edge "github.com/stackitcloud/stackit-sdk-go/services/edge/v1beta1api"
)

// ConfigureClient configures and returns a new API client for the Edge service.
func ConfigureClient(p *print.Printer, cliVersion string) (*edge.APIClient, error) {
	return genericclient.ConfigureClientGeneric(p, cliVersion, viper.GetString(config.EdgeCustomEndpointKey), false, genericclient.CreateApiClient[*edge.APIClient](edge.NewAPIClient))
}
