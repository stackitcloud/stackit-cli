package client

import (
	"github.com/spf13/viper"
	scf "github.com/stackitcloud/stackit-sdk-go/services/scf/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	genericclient "github.com/stackitcloud/stackit-cli/internal/pkg/generic-client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
)

// ConfigureClient creates and configures a new SCF API client
func ConfigureClient(p *print.Printer, cliVersion string) (*scf.APIClient, error) {
	return genericclient.ConfigureClientGeneric(p, cliVersion, viper.GetString(config.ScfCustomEndpointKey), false, genericclient.CreateApiClient[*scf.APIClient](scf.NewAPIClient))
}
