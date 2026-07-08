package client

import (
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	genericclient "github.com/stackitcloud/stackit-cli/internal/pkg/generic-client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/viper"
	redis "github.com/stackitcloud/stackit-sdk-go/services/redis/v2api"
)

func ConfigureClient(p *print.Printer, cliVersion string) (*redis.APIClient, error) {
	return genericclient.ConfigureClientGeneric(p, cliVersion, viper.GetString(config.RedisCustomEndpointKey), false, genericclient.CreateApiClient[*redis.APIClient](redis.NewAPIClient))
}
