package cmd

import (
	"os"

	"github.com/stackitcloud/stackit-cli/internal/cmd/auth"
	"github.com/stackitcloud/stackit-cli/internal/cmd/config"
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresql"
	configPkg "github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	projectIdFlag = "project-id"
)

var RootCmd = &cobra.Command{
	Use:               "stackit",
	Short:             "The root command of the STACKIT CLI",
	Long:              "The root command of the STACKIT CLI",
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Set up configuration files
	configPkg.InitConfig()

	// Add all direct child commands
	RootCmd.AddCommand(auth.Cmd)
	RootCmd.AddCommand(config.Cmd)
	RootCmd.AddCommand(dns.Cmd)
	RootCmd.AddCommand(postgresql.Cmd)

	configureFlags(RootCmd)
}

func configureFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().Var(flags.UUIDFlag(), projectIdFlag, "Project ID")

	err := viper.BindPFlag(configPkg.ProjectIdKey, cmd.PersistentFlags().Lookup(projectIdFlag))
	cobra.CheckErr(err)
}
