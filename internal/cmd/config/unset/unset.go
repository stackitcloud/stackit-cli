package unset

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	asyncFlag        = globalflags.AsyncFlag
	outputFormatFlag = globalflags.OutputFormatFlag
	projectIdFlag    = globalflags.ProjectIdFlag

	authorizationCustomEndpointFlag   = "authorization-custom-endpoint"
	dnsCustomEndpointFlag             = "dns-custom-endpoint"
	logMeCustomEndpointFlag           = "logme-custom-endpoint"
	mariaDBCustomEndpointFlag         = "mariadb-custom-endpoint"
	mongoDBFlexCustomEndpointFlag     = "mongodbflex-custom-endpoint"
	openSearchCustomEndpointFlag      = "opensearch-custom-endpoint"
	postgresFlexCustomEndpointFlag    = "postgresflex-custom-endpoint"
	rabbitMQCustomEndpointFlag        = "rabbitmq-custom-endpoint"
	redisCustomEndpointFlag           = "redis-custom-endpoint"
	resourceManagerCustomEndpointFlag = "resource-manager-custom-endpoint"
	serviceAccountCustomEndpointFlag  = "service-account-custom-endpoint"
	skeCustomEndpointFlag             = "ske-custom-endpoint"
)

type inputModel struct {
	AsyncFlag    bool
	OutputFormat bool
	ProjectId    bool

	AuthorizationCustomEndpoint   bool
	DNSCustomEndpoint             bool
	LogMeCustomEndpoint           bool
	MariaDBCustomEndpoint         bool
	MongoDBFlexCustomEndpoint     bool
	OpenSearchCustomEndpoint      bool
	PostgresFlexCustomEndpoint    bool
	RabbitMQCustomEndpoint        bool
	RedisCustomEndpoint           bool
	ResourceManagerCustomEndpoint bool
	ServiceAccountCustomEndpoint  bool
	SKECustomEndpoint             bool
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unset",
		Short: "Unsets CLI configuration options",
		Long:  "Unsets CLI configuration options.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Unset the project ID stored in your configuration`,
				"$ stackit config unset --project-id"),
			examples.NewExample(
				`Unset the session time limit stored in your configuration`,
				"$ stackit config unset --session-time-limit"),
			examples.NewExample(
				`Unset the DNS custom endpoint stored in your configuration`,
				"$ stackit config unset --dns-custom-endpoint"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			model := parseInput(cmd)

			if model.AsyncFlag {
				viper.Set(config.AsyncKey, "")
			}
			if model.OutputFormat {
				viper.Set(config.OutputFormatKey, "")
			}
			if model.ProjectId {
				viper.Set(config.ProjectIdKey, "")
			}

			if model.AuthorizationCustomEndpoint {
				viper.Set(config.AuthorizationCustomEndpointKey, "")
			}
			if model.DNSCustomEndpoint {
				viper.Set(config.DNSCustomEndpointKey, "")
			}
			if model.LogMeCustomEndpoint {
				viper.Set(config.LogMeCustomEndpointKey, "")
			}
			if model.MariaDBCustomEndpoint {
				viper.Set(config.MariaDBCustomEndpointKey, "")
			}
			if model.MongoDBFlexCustomEndpoint {
				viper.Set(config.MongoDBFlexCustomEndpointKey, "")
			}
			if model.OpenSearchCustomEndpoint {
				viper.Set(config.OpenSearchCustomEndpointKey, "")
			}
			if model.PostgresFlexCustomEndpoint {
				viper.Set(config.PostgresFlexCustomEndpointKey, "")
			}
			if model.RabbitMQCustomEndpoint {
				viper.Set(config.RabbitMQCustomEndpointKey, "")
			}
			if model.RedisCustomEndpoint {
				viper.Set(config.RedisCustomEndpointKey, "")
			}
			if model.ResourceManagerCustomEndpoint {
				viper.Set(config.ResourceManagerEndpointKey, "")
			}
			if model.ServiceAccountCustomEndpoint {
				viper.Set(config.ServiceAccountCustomEndpointKey, "")
			}
			if model.SKECustomEndpoint {
				viper.Set(config.SKECustomEndpointKey, "")
			}

			err := viper.WriteConfig()
			if err != nil {
				return fmt.Errorf("write updated config to file: %w", err)
			}
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Bool(asyncFlag, false, "Configuration option to run commands asynchronously")
	cmd.Flags().Bool(projectIdFlag, false, "Project ID")
	cmd.Flags().Bool(outputFormatFlag, false, "Output format")

	cmd.Flags().Bool(authorizationCustomEndpointFlag, false, "Authorization custom endpoint")
	cmd.Flags().Bool(dnsCustomEndpointFlag, false, "DNS custom endpoint")
	cmd.Flags().Bool(logMeCustomEndpointFlag, false, "LogMe custom endpoint")
	cmd.Flags().Bool(mariaDBCustomEndpointFlag, false, "MariaDB custom endpoint")
	cmd.Flags().Bool(mongoDBFlexCustomEndpointFlag, false, "MongoDB Flex custom endpoint")
	cmd.Flags().Bool(openSearchCustomEndpointFlag, false, "OpenSearch custom endpoint")
	cmd.Flags().Bool(postgresFlexCustomEndpointFlag, false, "PostgreSQL Flex custom endpoint")
	cmd.Flags().Bool(rabbitMQCustomEndpointFlag, false, "RabbitMQ custom endpoint")
	cmd.Flags().Bool(redisCustomEndpointFlag, false, "Redis custom endpoint")
	cmd.Flags().Bool(resourceManagerCustomEndpointFlag, false, "Resource Manager custom endpoint")
	cmd.Flags().Bool(serviceAccountCustomEndpointFlag, false, "SKE custom endpoint")
	cmd.Flags().Bool(skeCustomEndpointFlag, false, "SKE custom endpoint")
}

func parseInput(cmd *cobra.Command) *inputModel {
	return &inputModel{
		AsyncFlag:    flags.FlagToBoolValue(cmd, asyncFlag),
		OutputFormat: flags.FlagToBoolValue(cmd, outputFormatFlag),
		ProjectId:    flags.FlagToBoolValue(cmd, projectIdFlag),

		AuthorizationCustomEndpoint:   flags.FlagToBoolValue(cmd, authorizationCustomEndpointFlag),
		DNSCustomEndpoint:             flags.FlagToBoolValue(cmd, dnsCustomEndpointFlag),
		LogMeCustomEndpoint:           flags.FlagToBoolValue(cmd, logMeCustomEndpointFlag),
		MariaDBCustomEndpoint:         flags.FlagToBoolValue(cmd, mariaDBCustomEndpointFlag),
		MongoDBFlexCustomEndpoint:     flags.FlagToBoolValue(cmd, mongoDBFlexCustomEndpointFlag),
		OpenSearchCustomEndpoint:      flags.FlagToBoolValue(cmd, openSearchCustomEndpointFlag),
		PostgresFlexCustomEndpoint:    flags.FlagToBoolValue(cmd, postgresFlexCustomEndpointFlag),
		RabbitMQCustomEndpoint:        flags.FlagToBoolValue(cmd, rabbitMQCustomEndpointFlag),
		RedisCustomEndpoint:           flags.FlagToBoolValue(cmd, redisCustomEndpointFlag),
		ResourceManagerCustomEndpoint: flags.FlagToBoolValue(cmd, resourceManagerCustomEndpointFlag),
		ServiceAccountCustomEndpoint:  flags.FlagToBoolValue(cmd, serviceAccountCustomEndpointFlag),
		SKECustomEndpoint:             flags.FlagToBoolValue(cmd, skeCustomEndpointFlag),
	}
}
