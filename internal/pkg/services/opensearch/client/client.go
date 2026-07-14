package client

import (
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	genericclient "github.com/stackitcloud/stackit-cli/internal/pkg/generic-client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/viper"
	opensearch "github.com/stackitcloud/stackit-sdk-go/services/opensearch/v2api"
)

func ConfigureClient(p *print.Printer, cliVersion string) (*opensearch.APIClient, error) {
	return genericclient.ConfigureClientGeneric(p, cliVersion, viper.GetString(config.OpenSearchCustomEndpointKey), false, genericclient.CreateApiClient[*opensearch.APIClient](opensearch.NewAPIClient))
}
