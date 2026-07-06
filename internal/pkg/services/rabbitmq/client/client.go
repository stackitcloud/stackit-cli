package client

import (
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	genericclient "github.com/stackitcloud/stackit-cli/internal/pkg/generic-client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/viper"
	rabbitmq "github.com/stackitcloud/stackit-sdk-go/services/rabbitmq/v2api"
)

func ConfigureClient(p *print.Printer, cliVersion string) (*rabbitmq.APIClient, error) {
	return genericclient.ConfigureClientGeneric(p, cliVersion, viper.GetString(config.RabbitMQCustomEndpointKey), false, genericclient.CreateApiClient[*rabbitmq.APIClient](rabbitmq.NewAPIClient))
}
