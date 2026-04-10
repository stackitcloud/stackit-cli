package client

import (
	"github.com/spf13/viper"
	intake "github.com/stackitcloud/stackit-sdk-go/services/intake/v1betaapi"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	genericclient "github.com/stackitcloud/stackit-cli/internal/pkg/generic-client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
)

// ConfigureClient creates and configures a new Intake API client
func ConfigureClient(p *print.Printer, cliVersion string) (*intake.APIClient, error) {
	return genericclient.ConfigureClientGeneric(p, cliVersion, viper.GetString(config.IntakeCustomEndpointKey), true, genericclient.CreateApiClient[*intake.APIClient](intake.NewAPIClient))
}
