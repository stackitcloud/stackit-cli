package unset

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	projectIdFlag                = globalflags.ProjectIdFlag
	outputFormatFlag             = globalflags.OutputFormatFlag
	dnsCustomEndpointFlag        = "dns-custom-endpoint"
	postgreSQLCustomEndpointFlag = "postgresql-custom-endpoint"
	skeCustomEndpointFlag        = "ske-custom-endpoint"
)

type flagModel struct {
	ProjectId                bool
	OutputFormat             bool
	DNSCustomEndpoint        bool
	PostgreSQLCustomEndpoint bool
	SKECustomEndpoint        bool
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "unset",
		Short:   "Unset CLI configuration options",
		Long:    "Unset CLI configuration options",
		Example: `$ stackit config unset --project-id`,
		RunE: func(cmd *cobra.Command, args []string) error {
			model := parseFlags(cmd)

			if model.ProjectId {
				viper.Set(config.ProjectIdKey, "")
			}
			if model.OutputFormat {
				viper.Set(config.OutputFormatKey, "default")
			}
			if model.DNSCustomEndpoint {
				viper.Set(config.DNSCustomEndpointKey, "")
			}
			if model.PostgreSQLCustomEndpoint {
				viper.Set(config.PostgreSQLCustomEndpointKey, "")
			}
			if model.PostgreSQLCustomEndpoint {
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
	cmd.Flags().Bool(projectIdFlag.FlagName(), false, "Project ID")
	cmd.Flags().Bool(outputFormatFlag.FlagName(), false, "Output format")
	cmd.Flags().Bool(dnsCustomEndpointFlag, false, "DNS custom endpoint")
	cmd.Flags().Bool(postgreSQLCustomEndpointFlag, false, "PostgreSQL custom endpoint")
	cmd.Flags().Bool(skeCustomEndpointFlag, false, "SKE custom endpoint")
}

func parseFlags(cmd *cobra.Command) *flagModel {
	return &flagModel{
		ProjectId:                utils.FlagToBoolValue(cmd, projectIdFlag.FlagName()),
		OutputFormat:             utils.FlagToBoolValue(cmd, outputFormatFlag.FlagName()),
		DNSCustomEndpoint:        utils.FlagToBoolValue(cmd, dnsCustomEndpointFlag),
		PostgreSQLCustomEndpoint: utils.FlagToBoolValue(cmd, postgreSQLCustomEndpointFlag),
		SKECustomEndpoint:        utils.FlagToBoolValue(cmd, skeCustomEndpointFlag),
	}
}
