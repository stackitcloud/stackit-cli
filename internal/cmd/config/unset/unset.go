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

	sessionTimeLimitFlag               = "session-time-limit"
	identityProviderCustomEndpointFlag = "identity-provider-custom-endpoint"
	identityProviderCustomClientIdFlag = "identity-provider-custom-client-id"
	allowedUrlDomainFlag               = "allowed-url-domain"

	argusCustomEndpointFlag             = "argus-custom-endpoint"
	authorizationCustomEndpointFlag     = "authorization-custom-endpoint"
	dnsCustomEndpointFlag               = "dns-custom-endpoint"
	loadBalancerCustomEndpointFlag      = "load-balancer-custom-endpoint"
	logMeCustomEndpointFlag             = "logme-custom-endpoint"
	mariaDBCustomEndpointFlag           = "mariadb-custom-endpoint"
	mongoDBFlexCustomEndpointFlag       = "mongodbflex-custom-endpoint"
	objectStorageCustomEndpointFlag     = "object-storage-custom-endpoint"
	openSearchCustomEndpointFlag        = "opensearch-custom-endpoint"
	postgresFlexCustomEndpointFlag      = "postgresflex-custom-endpoint"
	rabbitMQCustomEndpointFlag          = "rabbitmq-custom-endpoint"
	redisCustomEndpointFlag             = "redis-custom-endpoint"
	resourceManagerCustomEndpointFlag   = "resource-manager-custom-endpoint"
	secretsManagerCustomEndpointFlag    = "secrets-manager-custom-endpoint"
	serviceAccountCustomEndpointFlag    = "service-account-custom-endpoint"
	serviceEnablementCustomEndpointFlag = "service-enablement-custom-endpoint"
	serverBackupCustomEndpointFlag      = "serverbackup-custom-endpoint"
	runCommandCustomEndpointFlag        = "runcommand-custom-endpoint"
	skeCustomEndpointFlag               = "ske-custom-endpoint"
	sqlServerFlexCustomEndpointFlag     = "sqlserverflex-custom-endpoint"
	iaasCustomEndpointFlag              = "iaas-custom-endpoint"
	tokenCustomEndpointFlag             = "token-custom-endpoint"
	jwksCustomEndpointFlag              = "jwks-custom-endpoint"
)

type inputModel struct {
	Async        bool
	OutputFormat bool
	ProjectId    bool
	Verbosity    bool

	SessionTimeLimit               bool
	IdentityProviderCustomEndpoint bool
	IdentityProviderCustomClientID bool
	AllowedUrlDomain               bool

	ArgusCustomEndpoint             bool
	AuthorizationCustomEndpoint     bool
	DNSCustomEndpoint               bool
	LoadBalancerCustomEndpoint      bool
	LogMeCustomEndpoint             bool
	MariaDBCustomEndpoint           bool
	MongoDBFlexCustomEndpoint       bool
	ObjectStorageCustomEndpoint     bool
	OpenSearchCustomEndpoint        bool
	PostgresFlexCustomEndpoint      bool
	RabbitMQCustomEndpoint          bool
	RedisCustomEndpoint             bool
	ResourceManagerCustomEndpoint   bool
	SecretsManagerCustomEndpoint    bool
	ServerBackupCustomEndpoint      bool
	RunCommandCustomEndpoint        bool
	ServiceAccountCustomEndpoint    bool
	ServiceEnablementCustomEndpoint bool
	SKECustomEndpoint               bool
	SQLServerFlexCustomEndpoint     bool
	IaaSCustomEndpoint              bool
	TokenCustomEndpoint             bool
	JwksCustomEndpoint              bool
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
			if model.Verbosity {
				viper.Set(config.VerbosityKey, globalflags.VerbosityDefault)
			}

			if model.SessionTimeLimit {
				viper.Set(config.SessionTimeLimitKey, config.SessionTimeLimitDefault)
			}
			if model.IdentityProviderCustomEndpoint {
				viper.Set(config.IdentityProviderCustomEndpointKey, "")
			}
			if model.IdentityProviderCustomClientID {
				viper.Set(config.IdentityProviderCustomClientIdKey, "")
			}
			if model.AllowedUrlDomain {
				viper.Set(config.AllowedUrlDomainKey, "")
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
			if model.ServiceEnablementCustomEndpoint {
				viper.Set(config.ServiceEnablementCustomEndpointKey, "")
			}
			if model.ServerBackupCustomEndpoint {
				viper.Set(config.ServerBackupCustomEndpointKey, "")
			}
			if model.RunCommandCustomEndpoint {
				viper.Set(config.RunCommandCustomEndpointKey, "")
			}
			if model.SKECustomEndpoint {
				viper.Set(config.SKECustomEndpointKey, "")
			}
			if model.SQLServerFlexCustomEndpoint {
				viper.Set(config.SQLServerFlexCustomEndpointKey, "")
			}
			if model.IaaSCustomEndpoint {
				viper.Set(config.IaaSCustomEndpointKey, "")
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
	cmd.Flags().Bool(verbosityFlag, false, "Verbosity of the CLI")

	cmd.Flags().Bool(sessionTimeLimitFlag, false, fmt.Sprintf("Maximum time before authentication is required again. If unset, defaults to %s", config.SessionTimeLimitDefault))
	cmd.Flags().Bool(identityProviderCustomEndpointFlag, false, "Identity Provider base URL. If unset, uses the default base URL")
	cmd.Flags().Bool(identityProviderCustomClientIdFlag, false, "Identity Provider client ID, used for user authentication")
	cmd.Flags().Bool(allowedUrlDomainFlag, false, "Domain name, used for the verification of the URLs that are given in the IDP endpoint and curl commands")

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
	cmd.Flags().Bool(serviceAccountCustomEndpointFlag, false, "Service Account API base URL. If unset, uses the default base URL")
	cmd.Flags().Bool(serviceEnablementCustomEndpointFlag, false, "Service Enablement API base URL. If unset, uses the default base URL")
	cmd.Flags().Bool(serverBackupCustomEndpointFlag, false, "Server Backup base URL. If unset, uses the default base URL")
	cmd.Flags().Bool(runCommandCustomEndpointFlag, false, "Server Command base URL. If unset, uses the default base URL")
	cmd.Flags().Bool(skeCustomEndpointFlag, false, "SKE API base URL. If unset, uses the default base URL")
	cmd.Flags().Bool(sqlServerFlexCustomEndpointFlag, false, "SQLServer Flex API base URL. If unset, uses the default base URL")
	cmd.Flags().Bool(iaasCustomEndpointFlag, false, "IaaS API base URL. If unset, uses the default base URL")
	cmd.Flags().Bool(tokenCustomEndpointFlag, false, "Custom endpoint for the token API, which is used to request access tokens when the service-account authentication is activated")
	cmd.Flags().Bool(jwksCustomEndpointFlag, false, "Custom endpoint for the jwks API, which is used to get the json web key sets (jwks) to validate tokens when the service-account authentication is activated")
}

func parseInput(p *print.Printer, cmd *cobra.Command) *inputModel {
	model := inputModel{
		Async:        flags.FlagToBoolValue(p, cmd, asyncFlag),
		OutputFormat: flags.FlagToBoolValue(p, cmd, outputFormatFlag),
		ProjectId:    flags.FlagToBoolValue(p, cmd, projectIdFlag),
		Verbosity:    flags.FlagToBoolValue(p, cmd, verbosityFlag),

		SessionTimeLimit:               flags.FlagToBoolValue(p, cmd, sessionTimeLimitFlag),
		IdentityProviderCustomEndpoint: flags.FlagToBoolValue(p, cmd, identityProviderCustomEndpointFlag),
		IdentityProviderCustomClientID: flags.FlagToBoolValue(p, cmd, identityProviderCustomClientIdFlag),
		AllowedUrlDomain:               flags.FlagToBoolValue(p, cmd, allowedUrlDomainFlag),

		ArgusCustomEndpoint:             flags.FlagToBoolValue(p, cmd, argusCustomEndpointFlag),
		AuthorizationCustomEndpoint:     flags.FlagToBoolValue(p, cmd, authorizationCustomEndpointFlag),
		DNSCustomEndpoint:               flags.FlagToBoolValue(p, cmd, dnsCustomEndpointFlag),
		LoadBalancerCustomEndpoint:      flags.FlagToBoolValue(p, cmd, loadBalancerCustomEndpointFlag),
		LogMeCustomEndpoint:             flags.FlagToBoolValue(p, cmd, logMeCustomEndpointFlag),
		MariaDBCustomEndpoint:           flags.FlagToBoolValue(p, cmd, mariaDBCustomEndpointFlag),
		MongoDBFlexCustomEndpoint:       flags.FlagToBoolValue(p, cmd, mongoDBFlexCustomEndpointFlag),
		ObjectStorageCustomEndpoint:     flags.FlagToBoolValue(p, cmd, objectStorageCustomEndpointFlag),
		OpenSearchCustomEndpoint:        flags.FlagToBoolValue(p, cmd, openSearchCustomEndpointFlag),
		PostgresFlexCustomEndpoint:      flags.FlagToBoolValue(p, cmd, postgresFlexCustomEndpointFlag),
		RabbitMQCustomEndpoint:          flags.FlagToBoolValue(p, cmd, rabbitMQCustomEndpointFlag),
		RedisCustomEndpoint:             flags.FlagToBoolValue(p, cmd, redisCustomEndpointFlag),
		ResourceManagerCustomEndpoint:   flags.FlagToBoolValue(p, cmd, resourceManagerCustomEndpointFlag),
		SecretsManagerCustomEndpoint:    flags.FlagToBoolValue(p, cmd, secretsManagerCustomEndpointFlag),
		ServiceAccountCustomEndpoint:    flags.FlagToBoolValue(p, cmd, serviceAccountCustomEndpointFlag),
		ServiceEnablementCustomEndpoint: flags.FlagToBoolValue(p, cmd, serviceEnablementCustomEndpointFlag),
		ServerBackupCustomEndpoint:      flags.FlagToBoolValue(p, cmd, serverBackupCustomEndpointFlag),
		RunCommandCustomEndpoint:        flags.FlagToBoolValue(p, cmd, runCommandCustomEndpointFlag),
		SKECustomEndpoint:               flags.FlagToBoolValue(p, cmd, skeCustomEndpointFlag),
		SQLServerFlexCustomEndpoint:     flags.FlagToBoolValue(p, cmd, sqlServerFlexCustomEndpointFlag),
		IaaSCustomEndpoint:              flags.FlagToBoolValue(p, cmd, iaasCustomEndpointFlag),
		TokenCustomEndpoint:             flags.FlagToBoolValue(p, cmd, tokenCustomEndpointFlag),
		JwksCustomEndpoint:              flags.FlagToBoolValue(p, cmd, jwksCustomEndpointFlag),
	}

	if p.IsVerbosityDebug() {
		modelStr, err := print.BuildDebugStrFromInputModel(model)
		if err != nil {
			p.Debug(print.ErrorLevel, "convert model to string for debugging: %v", err)
		} else {
			p.Debug(print.DebugLevel, "parsed input values: %s", modelStr)
		}
	}

	return &model
}
