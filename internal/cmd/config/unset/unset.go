package unset

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	asyncFlag        = globalflags.AsyncFlag
	outputFormatFlag = globalflags.OutputFormatFlag
	projectIdFlag    = globalflags.ProjectIdFlag
	verbosityFlag    = globalflags.VerbosityFlag

	sessionTimeLimitFlag = "session-time-limit"

	argusCustomEndpointFlag           = "argus-custom-endpoint"
	authorizationCustomEndpointFlag   = "authorization-custom-endpoint"
	dnsCustomEndpointFlag             = "dns-custom-endpoint"
	loadBalancerCustomEndpointFlag    = "load-balancer-custom-endpoint"
	logMeCustomEndpointFlag           = "logme-custom-endpoint"
	mariaDBCustomEndpointFlag         = "mariadb-custom-endpoint"
	mongoDBFlexCustomEndpointFlag     = "mongodbflex-custom-endpoint"
	objectStorageCustomEndpointFlag   = "object-storage-custom-endpoint"
	openSearchCustomEndpointFlag      = "opensearch-custom-endpoint"
	postgresFlexCustomEndpointFlag    = "postgresflex-custom-endpoint"
	rabbitMQCustomEndpointFlag        = "rabbitmq-custom-endpoint"
	redisCustomEndpointFlag           = "redis-custom-endpoint"
	resourceManagerCustomEndpointFlag = "resource-manager-custom-endpoint"
	secretsManagerCustomEndpointFlag  = "secrets-manager-custom-endpoint"
	serviceAccountCustomEndpointFlag  = "service-account-custom-endpoint"
	skeCustomEndpointFlag             = "ske-custom-endpoint"
)

type inputModel struct {
	Async            bool
	OutputFormat     bool
	ProjectId        bool
	SessionTimeLimit bool
	Verbosity        bool

	ArgusCustomEndpoint           bool
	AuthorizationCustomEndpoint   bool
	DNSCustomEndpoint             bool
	LoadBalancerCustomEndpoint    bool
	LogMeCustomEndpoint           bool
	MariaDBCustomEndpoint         bool
	MongoDBFlexCustomEndpoint     bool
	ObjectStorageCustomEndpoint   bool
	OpenSearchCustomEndpoint      bool
	PostgresFlexCustomEndpoint    bool
	RabbitMQCustomEndpoint        bool
	RedisCustomEndpoint           bool
	ResourceManagerCustomEndpoint bool
	SecretsManagerCustomEndpoint  bool
	ServiceAccountCustomEndpoint  bool
	SKECustomEndpoint             bool
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unset",
		Short: "Unsets CLI configuration options",
		Long:  "Unsets CLI configuration options, undoing past usages of the `stackit config set` command.",
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
			model := parseInput(p, cmd)

			if model.Async {
				viper.Set(config.AsyncKey, config.AsyncDefault)
			}
			if model.OutputFormat {
				viper.Set(config.OutputFormatKey, "")
			}
			if model.ProjectId {
				viper.Set(config.ProjectIdKey, "")
			}
			if model.SessionTimeLimit {
				viper.Set(config.SessionTimeLimitKey, config.SessionTimeLimitDefault)
			}
			if model.Verbosity {
				viper.Set(config.VerbosityKey, globalflags.VerbosityDefault)
			}

			if model.ArgusCustomEndpoint {
				viper.Set(config.ArgusCustomEndpointKey, "")
			}
			if model.AuthorizationCustomEndpoint {
				viper.Set(config.AuthorizationCustomEndpointKey, "")
			}
			if model.DNSCustomEndpoint {
				viper.Set(config.DNSCustomEndpointKey, "")
			}
			if model.LoadBalancerCustomEndpoint {
				viper.Set(config.LoadBalancerCustomEndpointKey, "")
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
			if model.ObjectStorageCustomEndpoint {
				viper.Set(config.ObjectStorageCustomEndpointKey, "")
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
			if model.SecretsManagerCustomEndpoint {
				viper.Set(config.SecretsManagerCustomEndpointKey, "")
			}
			if model.ServiceAccountCustomEndpoint {
				viper.Set(config.ServiceAccountCustomEndpointKey, "")
			}
			if model.SKECustomEndpoint {
				viper.Set(config.SKECustomEndpointKey, "")
			}

			err := config.Write()
			if err != nil {
				return fmt.Errorf("write config to file: %w", err)
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
	cmd.Flags().Bool(sessionTimeLimitFlag, false, fmt.Sprintf("Maximum time before authentication is required again. If unset, defaults to %s", config.SessionTimeLimitDefault))
	cmd.Flags().Bool(verbosityFlag, false, "Verbosity of the CLI")

	cmd.Flags().Bool(argusCustomEndpointFlag, false, "Argus API base URL. If unset, uses the default base URL")
	cmd.Flags().Bool(authorizationCustomEndpointFlag, false, "Authorization API base URL. If unset, uses the default base URL")
	cmd.Flags().Bool(dnsCustomEndpointFlag, false, "DNS API base URL. If unset, uses the default base URL")
	cmd.Flags().Bool(loadBalancerCustomEndpointFlag, false, "Load Balancer API base URL. If unset, uses the default base URL")
	cmd.Flags().Bool(logMeCustomEndpointFlag, false, "LogMe API base URL. If unset, uses the default base URL")
	cmd.Flags().Bool(mariaDBCustomEndpointFlag, false, "MariaDB API base URL. If unset, uses the default base URL")
	cmd.Flags().Bool(mongoDBFlexCustomEndpointFlag, false, "MongoDB Flex API base URL. If unset, uses the default base URL")
	cmd.Flags().Bool(objectStorageCustomEndpointFlag, false, "Object Storage API base URL. If unset, uses the default base URL")
	cmd.Flags().Bool(openSearchCustomEndpointFlag, false, "OpenSearch API base URL. If unset, uses the default base URL")
	cmd.Flags().Bool(postgresFlexCustomEndpointFlag, false, "PostgreSQL Flex API base URL. If unset, uses the default base URL")
	cmd.Flags().Bool(rabbitMQCustomEndpointFlag, false, "RabbitMQ API base URL. If unset, uses the default base URL")
	cmd.Flags().Bool(redisCustomEndpointFlag, false, "Redis API base URL. If unset, uses the default base URL")
	cmd.Flags().Bool(resourceManagerCustomEndpointFlag, false, "Resource Manager API base URL. If unset, uses the default base URL")
	cmd.Flags().Bool(secretsManagerCustomEndpointFlag, false, "Secrets Manager API base URL. If unset, uses the default base URL")
	cmd.Flags().Bool(serviceAccountCustomEndpointFlag, false, "SKE API base URL. If unset, uses the default base URL")
	cmd.Flags().Bool(skeCustomEndpointFlag, false, "SKE API base URL. If unset, uses the default base URL")
}

func parseInput(p *print.Printer, cmd *cobra.Command) *inputModel {
	return &inputModel{
		Async:            flags.FlagToBoolValue(p, cmd, asyncFlag),
		OutputFormat:     flags.FlagToBoolValue(p, cmd, outputFormatFlag),
		ProjectId:        flags.FlagToBoolValue(p, cmd, projectIdFlag),
		SessionTimeLimit: flags.FlagToBoolValue(p, cmd, sessionTimeLimitFlag),
		Verbosity:        flags.FlagToBoolValue(p, cmd, verbosityFlag),

		ArgusCustomEndpoint:           flags.FlagToBoolValue(p, cmd, argusCustomEndpointFlag),
		AuthorizationCustomEndpoint:   flags.FlagToBoolValue(p, cmd, authorizationCustomEndpointFlag),
		DNSCustomEndpoint:             flags.FlagToBoolValue(p, cmd, dnsCustomEndpointFlag),
		LoadBalancerCustomEndpoint:    flags.FlagToBoolValue(p, cmd, loadBalancerCustomEndpointFlag),
		LogMeCustomEndpoint:           flags.FlagToBoolValue(p, cmd, logMeCustomEndpointFlag),
		MariaDBCustomEndpoint:         flags.FlagToBoolValue(p, cmd, mariaDBCustomEndpointFlag),
		MongoDBFlexCustomEndpoint:     flags.FlagToBoolValue(p, cmd, mongoDBFlexCustomEndpointFlag),
		ObjectStorageCustomEndpoint:   flags.FlagToBoolValue(p, cmd, objectStorageCustomEndpointFlag),
		OpenSearchCustomEndpoint:      flags.FlagToBoolValue(p, cmd, openSearchCustomEndpointFlag),
		PostgresFlexCustomEndpoint:    flags.FlagToBoolValue(p, cmd, postgresFlexCustomEndpointFlag),
		RabbitMQCustomEndpoint:        flags.FlagToBoolValue(p, cmd, rabbitMQCustomEndpointFlag),
		RedisCustomEndpoint:           flags.FlagToBoolValue(p, cmd, redisCustomEndpointFlag),
		ResourceManagerCustomEndpoint: flags.FlagToBoolValue(p, cmd, resourceManagerCustomEndpointFlag),
		SecretsManagerCustomEndpoint:  flags.FlagToBoolValue(p, cmd, secretsManagerCustomEndpointFlag),
		ServiceAccountCustomEndpoint:  flags.FlagToBoolValue(p, cmd, serviceAccountCustomEndpointFlag),
		SKECustomEndpoint:             flags.FlagToBoolValue(p, cmd, skeCustomEndpointFlag),
	}
}
