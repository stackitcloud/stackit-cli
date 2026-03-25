package client

import (
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	genericclient "github.com/stackitcloud/stackit-cli/internal/pkg/generic-client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/viper"
	objectstorage "github.com/stackitcloud/stackit-sdk-go/services/objectstorage/v2api"
)

func ConfigureClient(p *print.Printer, cliVersion string) (*objectstorage.APIClient, error) {
	return genericclient.ConfigureClientGeneric(p, cliVersion, viper.GetString(config.ObjectStorageCustomEndpointKey), false, genericclient.CreateApiClient[*objectstorage.APIClient](objectstorage.NewAPIClient))
}
