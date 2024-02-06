package set

import (
	"fmt"
	"time"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	sessionTimeLimitFlag = "session-time-limit"

	dnsCustomEndpointFlag             = "dns-custom-endpoint"
	membershipCustomEndpointFlag      = "membership-custom-endpoint"
	mongoDBFlexCustomEndpointFlag     = "mongodbflex-custom-endpoint"
	openSearchCustomEndpointFlag      = "opensearch-custom-endpoint"
	postgresFlexCustomEndpointFlag    = "postgresflex-custom-endpoint"
	resourceManagerCustomEndpointFlag = "resource-manager-custom-endpoint"
	serviceAccountCustomEndpointFlag  = "service-account-custom-endpoint"
	skeCustomEndpointFlag             = "ske-custom-endpoint"
)

type inputModel struct {
	SessionTimeLimit *string
	// If true, projectId has been set
	ProjectIdSet bool
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Sets CLI configuration options",
		Long: fmt.Sprintf("%s\n%s\n%s\n%s",
			"Sets CLI configuration options.",
			"All of the configuration options can be set using an environment variable, which takes precedence over what is configured.",
			`The environment variable is the name of the flag, with underscores ("_") instead of dashes ("-") and the "STACKIT" prefix.`,
			"Example: to set the project ID you can set the environment variable STACKIT_PROJECT_ID.",
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Set a project ID in your active configuration. This project ID will be used by every command, as long as it's not overridden by the "STACKIT_PROJECT_ID" environment variable or the command flag`,
				"$ stackit config set --project-id xxx"),
			examples.NewExample(
				`Set the session time limit to 1 hour. After this time you will be prompted to login again to be able to execute commands that need authentication`,
				"$ stackit config set --session-time-limit 1h"),
			examples.NewExample(
				`Set the DNS custom endpoint. This endpoint will be used on all calls to the DNS API, unless overridden by the "STACKIT_DNS_CUSTOM_ENDPOINT" environment variable`,
				"$ stackit config set --dns-custom-endpoint https://dns.stackit.cloud"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			model, err := parseInput(cmd)
			if err != nil {
				return err
			}

			if model.SessionTimeLimit != nil {
				cmd.Println("Authenticate again to apply changes to session time limit")
				viper.Set(config.SessionTimeLimitKey, *model.SessionTimeLimit)
			}

			// If project ID was set, remove the value for project name stored in config
			if model.ProjectIdSet {
				viper.Set(config.ProjectNameKey, "")
			}

			err = viper.WriteConfig()
			if err != nil {
				return fmt.Errorf("write new config to file: %w", err)
			}
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(sessionTimeLimitFlag, "", "Maximum time before authentication is required again. Can't be larger than 24h. Examples: 3h, 5h30m40s (BETA: currently values greater than 2h have no effect)")

	cmd.Flags().String(dnsCustomEndpointFlag, "", "DNS custom endpoint")
	cmd.Flags().String(membershipCustomEndpointFlag, "", "Membership custom endpoint")
	cmd.Flags().String(mongoDBFlexCustomEndpointFlag, "", "MongoDB Flex custom endpoint")
	cmd.Flags().String(openSearchCustomEndpointFlag, "", "OpenSearch custom endpoint")
	cmd.Flags().String(postgresFlexCustomEndpointFlag, "", "PostgreSQL Flex custom endpoint")
	cmd.Flags().String(resourceManagerCustomEndpointFlag, "", "Resource manager custom endpoint")
	cmd.Flags().String(serviceAccountCustomEndpointFlag, "", "Service Account custom endpoint")
	cmd.Flags().String(skeCustomEndpointFlag, "", "SKE custom endpoint")

	err := viper.BindPFlag(config.DNSCustomEndpointKey, cmd.Flags().Lookup(dnsCustomEndpointFlag))
	cobra.CheckErr(err)
	err = viper.BindPFlag(config.MembershipCustomEndpointKey, cmd.Flags().Lookup(membershipCustomEndpointFlag))
	cobra.CheckErr(err)
	err = viper.BindPFlag(config.MongoDBFlexCustomEndpointKey, cmd.Flags().Lookup(mongoDBFlexCustomEndpointFlag))
	cobra.CheckErr(err)
	err = viper.BindPFlag(config.OpenSearchCustomEndpointKey, cmd.Flags().Lookup(openSearchCustomEndpointFlag))
	cobra.CheckErr(err)
	err = viper.BindPFlag(config.PostgresFlexCustomEndpointKey, cmd.Flags().Lookup(postgresFlexCustomEndpointFlag))
	cobra.CheckErr(err)
	err = viper.BindPFlag(config.ResourceManagerEndpointKey, cmd.Flags().Lookup(skeCustomEndpointFlag))
	cobra.CheckErr(err)
	err = viper.BindPFlag(config.ServiceAccountCustomEndpointKey, cmd.Flags().Lookup(serviceAccountCustomEndpointFlag))
	cobra.CheckErr(err)
	err = viper.BindPFlag(config.SKECustomEndpointKey, cmd.Flags().Lookup(skeCustomEndpointFlag))
	cobra.CheckErr(err)
}

func parseInput(cmd *cobra.Command) (*inputModel, error) {
	sessionTimeLimit, err := parseSessionTimeLimit(cmd)
	if err != nil {
		return nil, &errors.FlagValidationError{
			Flag:    sessionTimeLimitFlag,
			Details: err.Error(),
		}
	}

	// values.FlagToStringPointer pulls the projectId from passed flags
	// globalflags.Parse uses the flags, and fallsback to config file
	// To check if projectId was passed, we use the first rather than the second
	projectIdFromFlag := flags.FlagToStringPointer(cmd, globalflags.ProjectIdFlag)
	projectIdSet := false
	if projectIdFromFlag != nil {
		projectIdSet = true
	}

	return &inputModel{
		SessionTimeLimit: sessionTimeLimit,
		ProjectIdSet:     projectIdSet,
	}, nil
}

func parseSessionTimeLimit(cmd *cobra.Command) (*string, error) {
	sessionTimeLimit := flags.FlagToStringPointer(cmd, sessionTimeLimitFlag)
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
