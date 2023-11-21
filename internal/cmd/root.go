package cmd

import (
	"os"

	"github.com/stackitcloud/stackit-cli/internal/cmd/auth"
	"github.com/stackitcloud/stackit-cli/internal/cmd/config"
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresql"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske"
	configPkg "github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	projectIdFlag = "project-id"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "stackit",
		Short:             "The root command of the STACKIT CLI",
		Long:              "The root command of the STACKIT CLI",
		SilenceUsage:      true,
		DisableAutoGenTag: true,
	}
	configureFlags(cmd)
	addSubcommands(cmd)
	return cmd
}

func Execute() {
	err := NewRootCmd().Execute()
	if err != nil {
		os.Exit(1)
	}
}

func configureFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().Var(flags.UUIDFlag(), projectIdFlag, "Project ID")

	err := viper.BindPFlag(configPkg.ProjectIdKey, cmd.PersistentFlags().Lookup(projectIdFlag))
	cobra.CheckErr(err)
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(auth.NewCmd())
	cmd.AddCommand(config.NewCmd())
	cmd.AddCommand(dns.NewCmd())
	cmd.AddCommand(postgresql.NewCmd())
	cmd.AddCommand(ske.NewCmd())
}
