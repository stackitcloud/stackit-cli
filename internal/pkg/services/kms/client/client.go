package client

import (
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	genericclient "github.com/stackitcloud/stackit-cli/internal/pkg/generic-client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/viper"
	kms "github.com/stackitcloud/stackit-sdk-go/services/kms/v1api"
)

func ConfigureClient(p *print.Printer, cliVersion string) (*kms.APIClient, error) {
	return genericclient.ConfigureClientGeneric(p, cliVersion, viper.GetString(config.KMSCustomEndpointKey), false, kms.NewAPIClient)
}
