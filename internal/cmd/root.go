package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/auth"
	"github.com/stackitcloud/stackit-cli/internal/cmd/config"
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresql"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "stackit",
		Short:             "The root command of the STACKIT CLI",
		Long:              "The root command of the STACKIT CLI",
		SilenceUsage:      true,
		DisableAutoGenTag: true,
	}
	err := globalflags.ConfigureFlags(cmd.PersistentFlags())
	cobra.CheckErr(err)
	addSubcommands(cmd)
	return cmd
}

func Execute() {
	err := NewRootCmd().Execute()
	if err != nil {
		os.Exit(1)
	}
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(auth.NewCmd())
	cmd.AddCommand(config.NewCmd())
	cmd.AddCommand(dns.NewCmd())
	cmd.AddCommand(postgresql.NewCmd())
	cmd.AddCommand(ske.NewCmd())
}
