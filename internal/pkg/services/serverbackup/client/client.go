package client

import (
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	genericclient "github.com/stackitcloud/stackit-cli/internal/pkg/generic-client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/viper"
	serverbackup "github.com/stackitcloud/stackit-sdk-go/services/serverbackup/v2api"
)

func ConfigureClient(p *print.Printer, cliVersion string) (*serverbackup.APIClient, error) {
	return genericclient.ConfigureClientGeneric(p, cliVersion, viper.GetString(config.ServerBackupCustomEndpointKey), false, genericclient.CreateApiClient[*serverbackup.APIClient](serverbackup.NewAPIClient))
}
