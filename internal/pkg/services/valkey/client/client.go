package client

import (
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	genericclient "github.com/stackitcloud/stackit-cli/internal/pkg/generic-client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/viper"
	valkey "github.com/stackitcloud/stackit-sdk-go/services/valkey/v2api"
)

func ConfigureClient(p *print.Printer, cliVersion string) (*valkey.APIClient, error) {
	return genericclient.ConfigureClientGeneric(p, cliVersion, viper.GetString(config.ValkeyCustomEndpointKey), false, genericclient.CreateApiClient[*valkey.APIClient](valkey.NewAPIClient))
}
