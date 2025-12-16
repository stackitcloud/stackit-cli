package set

import (
	"fmt"
	"time"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	sessionTimeLimitFlag                             = "session-time-limit"
	identityProviderCustomWellKnownConfigurationFlag = "identity-provider-custom-well-known-configuration"
	identityProviderCustomClientIdFlag               = "identity-provider-custom-client-id"
	allowedUrlDomainFlag                             = "allowed-url-domain"

	authorizationCustomEndpointFlag     = "authorization-custom-endpoint"
	dnsCustomEndpointFlag               = "dns-custom-endpoint"
	loadBalancerCustomEndpointFlag      = "load-balancer-custom-endpoint"
	logMeCustomEndpointFlag             = "logme-custom-endpoint"
	mariaDBCustomEndpointFlag           = "mariadb-custom-endpoint"
	mongoDBFlexCustomEndpointFlag       = "mongodbflex-custom-endpoint"
	objectStorageCustomEndpointFlag     = "object-storage-custom-endpoint"
	observabilityCustomEndpointFlag     = "observability-custom-endpoint"
	openSearchCustomEndpointFlag        = "opensearch-custom-endpoint"
	postgresFlexCustomEndpointFlag      = "postgresflex-custom-endpoint"
	rabbitMQCustomEndpointFlag          = "rabbitmq-custom-endpoint"
	redisCustomEndpointFlag             = "redis-custom-endpoint"
	resourceManagerCustomEndpointFlag   = "resource-manager-custom-endpoint"
	secretsManagerCustomEndpointFlag    = "secrets-manager-custom-endpoint"
	kmsCustomEndpointFlag               = "kms-custom-endpoint"
	serverBackupCustomEndpointFlag      = "serverbackup-custom-endpoint"
	serverOsUpdateCustomEndpointFlag    = "server-osupdate-custom-endpoint"
	runCommandCustomEndpointFlag        = "runcommand-custom-endpoint"
	serviceAccountCustomEndpointFlag    = "service-account-custom-endpoint"
	serviceEnablementCustomEndpointFlag = "service-enablement-custom-endpoint"
	skeCustomEndpointFlag               = "ske-custom-endpoint"
	sqlServerFlexCustomEndpointFlag     = "sqlserverflex-custom-endpoint"
	iaasCustomEndpointFlag              = "iaas-custom-endpoint"
	tokenCustomEndpointFlag             = "token-custom-endpoint"
	intakeCustomEndpointFlag            = "intake-custom-endpoint"
	logsCustomEndpointFlag              = "logs-custom-endpoint"
)

type inputModel struct {
	SessionTimeLimit *string
	// If true, projectId has been set
	ProjectIdSet bool
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Sets CLI configuration options",
		Long: fmt.Sprintf("%s\n%s\n%s\n%s",
			"Sets CLI configuration options.",
			"All of the configuration options can be set using an environment variable, which takes precedence over what is configured using this command.",
			`The environment variable is the name of the flag, with underscores ("_") instead of dashes ("-") and the "STACKIT" prefix.`,
			"Example: to set the project ID you can set the environment variable STACKIT_PROJECT_ID.",
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Set a project ID in your active configuration. This project ID will be used by every command (unless overridden by the "STACKIT_PROJECT_ID" environment variable)`,
				"$ stackit config set --project-id xxx"),
			examples.NewExample(
				`Set the session time limit to 1 hour`,
				"$ stackit config set --session-time-limit 1h"),
			examples.NewExample(
				`Set the DNS custom endpoint. This endpoint will be used on all calls to the DNS API (unless overridden by the "STACKIT_DNS_CUSTOM_ENDPOINT" environment variable)`,
				"$ stackit config set --dns-custom-endpoint https://dns.stackit.cloud"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			if model.SessionTimeLimit != nil {
				params.Printer.Warn("Authenticate again to apply changes to session time limit\n")
				viper.Set(config.SessionTimeLimitKey, *model.SessionTimeLimit)
			}

			// If project ID was set, remove the value for project name stored in config
			if model.ProjectIdSet {
				viper.Set(config.ProjectNameKey, "")
			}

			err = config.Write()
			if err != nil {
				return fmt.Errorf("write config to file: %w", err)
			}
			return nil
		},
	}
	configureFlags(cmd)
	cmd.SetHelpTemplate( // Slightly modified help template to merge the command level Flags and Global Flags in a single section
		`Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}

Available Commands:{{range $cmds}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{else}}{{range $group := .Groups}}

{{.Title}}{{range $cmds}}{{if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if not .AllChildCommandsHaveGroup}}

Additional Commands:{{range $cmds}}{{if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.Flags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(sessionTimeLimitFlag, "", "Maximum time before authentication is required again. After this time, you will be prompted to login again to execute commands that require authentication. Can't be larger than 24h. Requires authentication after being set to take effect. Examples: 3h, 5h30m40s (BETA: currently values greater than 2h have no effect)")
	cmd.Flags().String(identityProviderCustomWellKnownConfigurationFlag, "", "Identity Provider well-known OpenID configuration URL, used for user authentication")
	cmd.Flags().String(identityProviderCustomClientIdFlag, "", "Identity Provider client ID, used for user authentication")
	cmd.Flags().String(allowedUrlDomainFlag, "", `Domain name, used for the verification of the URLs that are given in the custom identity provider endpoint and "STACKIT curl" command`)
	cmd.Flags().String(observabilityCustomEndpointFlag, "", "Observability API base URL, used in calls to this API")
	cmd.Flags().String(authorizationCustomEndpointFlag, "", "Authorization API base URL, used in calls to this API")
	cmd.Flags().String(dnsCustomEndpointFlag, "", "DNS API base URL, used in calls to this API")
	cmd.Flags().String(loadBalancerCustomEndpointFlag, "", "Load Balancer API base URL, used in calls to this API")
	cmd.Flags().String(logMeCustomEndpointFlag, "", "LogMe API base URL, used in calls to this API")
	cmd.Flags().String(mariaDBCustomEndpointFlag, "", "MariaDB API base URL, used in calls to this API")
	cmd.Flags().String(mongoDBFlexCustomEndpointFlag, "", "MongoDB Flex API base URL, used in calls to this API")
	cmd.Flags().String(objectStorageCustomEndpointFlag, "", "Object Storage API base URL, used in calls to this API")
	cmd.Flags().String(openSearchCustomEndpointFlag, "", "OpenSearch API base URL, used in calls to this API")
	cmd.Flags().String(postgresFlexCustomEndpointFlag, "", "PostgreSQL Flex API base URL, used in calls to this API")
	cmd.Flags().String(rabbitMQCustomEndpointFlag, "", "RabbitMQ API base URL, used in calls to this API")
	cmd.Flags().String(redisCustomEndpointFlag, "", "Redis API base URL, used in calls to this API")
	cmd.Flags().String(resourceManagerCustomEndpointFlag, "", "Resource Manager API base URL, used in calls to this API")
	cmd.Flags().String(secretsManagerCustomEndpointFlag, "", "Secrets Manager API base URL, used in calls to this API")
	cmd.Flags().String(kmsCustomEndpointFlag, "", "KMS API base URL, used in calls to this API")
	cmd.Flags().String(serviceAccountCustomEndpointFlag, "", "Service Account API base URL, used in calls to this API")
	cmd.Flags().String(serviceEnablementCustomEndpointFlag, "", "Service Enablement API base URL, used in calls to this API")
	cmd.Flags().String(serverBackupCustomEndpointFlag, "", "Server Backup API base URL, used in calls to this API")
	cmd.Flags().String(serverOsUpdateCustomEndpointFlag, "", "Server Update Management API base URL, used in calls to this API")
	cmd.Flags().String(runCommandCustomEndpointFlag, "", "Run Command API base URL, used in calls to this API")
	cmd.Flags().String(skeCustomEndpointFlag, "", "SKE API base URL, used in calls to this API")
	cmd.Flags().String(sqlServerFlexCustomEndpointFlag, "", "SQLServer Flex API base URL, used in calls to this API")
	cmd.Flags().String(iaasCustomEndpointFlag, "", "IaaS API base URL, used in calls to this API")
	cmd.Flags().String(tokenCustomEndpointFlag, "", "Custom token endpoint of the Service Account API, which is used to request access tokens when the service account authentication is activated. Not relevant for user authentication.")
	cmd.Flags().String(intakeCustomEndpointFlag, "", "Intake API base URL, used in calls to this API")
	cmd.Flags().String(logsCustomEndpointFlag, "", "Logs API base URL, used in calls to this API")

	err := viper.BindPFlag(config.SessionTimeLimitKey, cmd.Flags().Lookup(sessionTimeLimitFlag))
	cobra.CheckErr(err)
	err = viper.BindPFlag(config.IdentityProviderCustomWellKnownConfigurationKey, cmd.Flags().Lookup(identityProviderCustomWellKnownConfigurationFlag))
	cobra.CheckErr(err)
	err = viper.BindPFlag(config.IdentityProviderCustomClientIdKey, cmd.Flags().Lookup(identityProviderCustomClientIdFlag))
	cobra.CheckErr(err)
	err = viper.BindPFlag(config.AllowedUrlDomainKey, cmd.Flags().Lookup(allowedUrlDomainFlag))
	cobra.CheckErr(err)

	err = viper.BindPFlag(config.ObservabilityCustomEndpointKey, cmd.Flags().Lookup(observabilityCustomEndpointFlag))
	cobra.CheckErr(err)
	err = viper.BindPFlag(config.AuthorizationCustomEndpointKey, cmd.Flags().Lookup(authorizationCustomEndpointFlag))
	cobra.CheckErr(err)
	err = viper.BindPFlag(config.DNSCustomEndpointKey, cmd.Flags().Lookup(dnsCustomEndpointFlag))
	cobra.CheckErr(err)
	err = viper.BindPFlag(config.LoadBalancerCustomEndpointKey, cmd.Flags().Lookup(loadBalancerCustomEndpointFlag))
	cobra.CheckErr(err)
	err = viper.BindPFlag(config.LogMeCustomEndpointKey, cmd.Flags().Lookup(logMeCustomEndpointFlag))
	cobra.CheckErr(err)
	err = viper.BindPFlag(config.MariaDBCustomEndpointKey, cmd.Flags().Lookup(mariaDBCustomEndpointFlag))
	cobra.CheckErr(err)
	err = viper.BindPFlag(config.MongoDBFlexCustomEndpointKey, cmd.Flags().Lookup(mongoDBFlexCustomEndpointFlag))
	cobra.CheckErr(err)
	err = viper.BindPFlag(config.ObjectStorageCustomEndpointKey, cmd.Flags().Lookup(objectStorageCustomEndpointFlag))
	cobra.CheckErr(err)
	err = viper.BindPFlag(config.OpenSearchCustomEndpointKey, cmd.Flags().Lookup(openSearchCustomEndpointFlag))
	cobra.CheckErr(err)
	err = viper.BindPFlag(config.PostgresFlexCustomEndpointKey, cmd.Flags().Lookup(postgresFlexCustomEndpointFlag))
	cobra.CheckErr(err)
	err = viper.BindPFlag(config.RabbitMQCustomEndpointKey, cmd.Flags().Lookup(rabbitMQCustomEndpointFlag))
	cobra.CheckErr(err)
	err = viper.BindPFlag(config.RedisCustomEndpointKey, cmd.Flags().Lookup(redisCustomEndpointFlag))
	cobra.CheckErr(err)
	err = viper.BindPFlag(config.ResourceManagerEndpointKey, cmd.Flags().Lookup(resourceManagerCustomEndpointFlag))
	cobra.CheckErr(err)
	err = viper.BindPFlag(config.SecretsManagerCustomEndpointKey, cmd.Flags().Lookup(secretsManagerCustomEndpointFlag))
	cobra.CheckErr(err)
	err = viper.BindPFlag(config.KMSCustomEndpointKey, cmd.Flags().Lookup(kmsCustomEndpointFlag))
	cobra.CheckErr(err)
	err = viper.BindPFlag(config.ServerBackupCustomEndpointKey, cmd.Flags().Lookup(serverBackupCustomEndpointFlag))
	cobra.CheckErr(err)
	err = viper.BindPFlag(config.ServerOsUpdateCustomEndpointKey, cmd.Flags().Lookup(serverOsUpdateCustomEndpointFlag))
	cobra.CheckErr(err)
	err = viper.BindPFlag(config.RunCommandCustomEndpointKey, cmd.Flags().Lookup(runCommandCustomEndpointFlag))
	cobra.CheckErr(err)
	err = viper.BindPFlag(config.ServiceAccountCustomEndpointKey, cmd.Flags().Lookup(serviceAccountCustomEndpointFlag))
	cobra.CheckErr(err)
	err = viper.BindPFlag(config.ServiceEnablementCustomEndpointKey, cmd.Flags().Lookup(serviceEnablementCustomEndpointFlag))
	cobra.CheckErr(err)
	err = viper.BindPFlag(config.SKECustomEndpointKey, cmd.Flags().Lookup(skeCustomEndpointFlag))
	cobra.CheckErr(err)
	err = viper.BindPFlag(config.SQLServerFlexCustomEndpointKey, cmd.Flags().Lookup(sqlServerFlexCustomEndpointFlag))
	cobra.CheckErr(err)
	err = viper.BindPFlag(config.IaaSCustomEndpointKey, cmd.Flags().Lookup(iaasCustomEndpointFlag))
	cobra.CheckErr(err)
	err = viper.BindPFlag(config.TokenCustomEndpointKey, cmd.Flags().Lookup(tokenCustomEndpointFlag))
	cobra.CheckErr(err)
	err = viper.BindPFlag(config.IntakeCustomEndpointKey, cmd.Flags().Lookup(intakeCustomEndpointFlag))
	cobra.CheckErr(err)
	err = viper.BindPFlag(config.LogsCustomEndpointKey, cmd.Flags().Lookup(logsCustomEndpointFlag))
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	sessionTimeLimit, err := parseSessionTimeLimit(p, cmd)
	if err != nil {
		return nil, &errors.FlagValidationError{
			Flag:    sessionTimeLimitFlag,
			Details: err.Error(),
		}
	}

	// values.FlagToStringPointer pulls the projectId from passed flags
	// globalflags.Parse uses the flags, and fallsback to config file
	// To check if projectId was passed, we use the first rather than the second
	projectIdFromFlag := flags.FlagToStringPointer(p, cmd, globalflags.ProjectIdFlag)
	projectIdSet := projectIdFromFlag != nil

	allowedUrlDomainFromFlag := flags.FlagToStringPointer(p, cmd, allowedUrlDomainFlag)
	allowedUrlDomainFlagValue := flags.FlagToStringValue(p, cmd, allowedUrlDomainFlag)
	if allowedUrlDomainFromFlag != nil && allowedUrlDomainFlagValue == "" {
		p.Warn("The allowed URL domain is set to empty. All URLs will be accepted regardless of their domain.\n")
	}

	model := inputModel{
		SessionTimeLimit: sessionTimeLimit,
		ProjectIdSet:     projectIdSet,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func parseSessionTimeLimit(p *print.Printer, cmd *cobra.Command) (*string, error) {
	sessionTimeLimit := flags.FlagToStringPointer(p, cmd, sessionTimeLimitFlag)
	if sessionTimeLimit == nil {
		return nil, nil
	}

	// time.ParseDuration doesn't recognize unit "d", for simplicity we allow the value "1d"
	if *sessionTimeLimit == "1d" {
		*sessionTimeLimit = "24h"
	}

	duration, err := time.ParseDuration(*sessionTimeLimit)
	if err != nil {
		return nil, fmt.Errorf("parse value \"%s\": %w", *sessionTimeLimit, err)
	}
	if duration <= 0 {
		return nil, fmt.Errorf("value must be positive")
	}
	if duration > time.Duration(24)*time.Hour {
		return nil, fmt.Errorf("value can't be larger than 24h")
	}

	return sessionTimeLimit, nil
}
