package unset

import (
	"fmt"

	"stackit/internal/pkg/args"
	"stackit/internal/pkg/config"
	"stackit/internal/pkg/examples"
	"stackit/internal/pkg/flags"
	"stackit/internal/pkg/globalflags"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	asyncFlag        = globalflags.AsyncFlag
	outputFormatFlag = globalflags.OutputFormatFlag
	projectIdFlag    = globalflags.ProjectIdFlag

	dnsCustomEndpointFlag             = "dns-custom-endpoint"
	membershipCustomEndpointFlag      = "membership-custom-endpoint"
	mongoDBFlexCustomEndpointFlag     = "mongodbflex-custom-endpoint"
	serviceAccountCustomEndpointFlag  = "service-account-custom-endpoint"
	skeCustomEndpointFlag             = "ske-custom-endpoint"
	resourceManagerCustomEndpointFlag = "resource-manager-custom-endpoint"
	openSearchCustomEndpointFlag      = "opensearch-custom-endpoint"
)

type inputModel struct {
	AsyncFlag    bool
	OutputFormat bool
	ProjectId    bool

	DNSCustomEndpoint             bool
	MembershipCustomEndpoint      bool
	MongoDBFlexCustomEndpoint     bool
	ServiceAccountCustomEndpoint  bool
	SKECustomEndpoint             bool
	ResourceManagerCustomEndpoint bool
	OpenSearchCustomEndpoint      bool
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unset",
		Short: "Unset CLI configuration options",
		Long:  "Unset CLI configuration options",
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

			if model.DNSCustomEndpoint {
				viper.Set(config.DNSCustomEndpointKey, "")
			}
			if model.MembershipCustomEndpoint {
				viper.Set(config.MembershipCustomEndpointKey, "")
			}
			if model.MongoDBFlexCustomEndpoint {
				viper.Set(config.MongoDBFlexCustomEndpointKey, "")
			}
			if model.ServiceAccountCustomEndpoint {
				viper.Set(config.ServiceAccountCustomEndpointKey, "")
			}
			if model.SKECustomEndpoint {
				viper.Set(config.SKECustomEndpointKey, "")
			}
			if model.ResourceManagerCustomEndpoint {
				viper.Set(config.ResourceManagerEndpointKey, "")
			}
			if model.OpenSearchCustomEndpoint {
				viper.Set(config.OpenSearchCustomEndpointKey, "")
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

	cmd.Flags().Bool(dnsCustomEndpointFlag, false, "DNS custom endpoint")
	cmd.Flags().Bool(membershipCustomEndpointFlag, false, "Membership custom endpoint")
	cmd.Flags().Bool(mongoDBFlexCustomEndpointFlag, false, "MongoDB Flex custom endpoint")
	cmd.Flags().Bool(serviceAccountCustomEndpointFlag, false, "SKE custom endpoint")
	cmd.Flags().Bool(skeCustomEndpointFlag, false, "SKE custom endpoint")
	cmd.Flags().Bool(resourceManagerCustomEndpointFlag, false, "Resource Manager custom endpoint")
	cmd.Flags().Bool(openSearchCustomEndpointFlag, false, "OpenSearch custom endpoint")
}

func parseInput(cmd *cobra.Command) *inputModel {
	return &inputModel{
		AsyncFlag:    flags.FlagToBoolValue(cmd, asyncFlag),
		OutputFormat: flags.FlagToBoolValue(cmd, outputFormatFlag),
		ProjectId:    flags.FlagToBoolValue(cmd, projectIdFlag),

		DNSCustomEndpoint:             flags.FlagToBoolValue(cmd, dnsCustomEndpointFlag),
		MembershipCustomEndpoint:      flags.FlagToBoolValue(cmd, membershipCustomEndpointFlag),
		MongoDBFlexCustomEndpoint:     flags.FlagToBoolValue(cmd, mongoDBFlexCustomEndpointFlag),
		ServiceAccountCustomEndpoint:  flags.FlagToBoolValue(cmd, serviceAccountCustomEndpointFlag),
		SKECustomEndpoint:             flags.FlagToBoolValue(cmd, skeCustomEndpointFlag),
		ResourceManagerCustomEndpoint: flags.FlagToBoolValue(cmd, resourceManagerCustomEndpointFlag),
		OpenSearchCustomEndpoint:      flags.FlagToBoolValue(cmd, openSearchCustomEndpointFlag),
	}
}
