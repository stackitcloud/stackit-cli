package client

import (
	"github.com/spf13/viper"
	automation "github.com/stackitcloud/stackit-sdk-go/services/automation/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	genericclient "github.com/stackitcloud/stackit-cli/internal/pkg/generic-client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
)

func ConfigureClient(p *print.Printer, cliVersion string) (*automation.APIClient, error) {
	return genericclient.ConfigureClientGeneric(p, cliVersion, viper.GetString(config.AutomationCustomEndpointKey), false, automation.NewAPIClient)
}
