package auth

import (
	activateserviceaccount "stackit/internal/cmd/auth/activate-service-account"
	"stackit/internal/cmd/auth/login"
	"stackit/internal/pkg/args"
	"stackit/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Provides authentication functionality",
		Long:  "Provides authentication functionality",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(login.NewCmd())
	cmd.AddCommand(activateserviceaccount.NewCmd())
}
