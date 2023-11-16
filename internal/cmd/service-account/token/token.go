package token

import (
	"stackit/internal/cmd/service-account/token/create"
	"stackit/internal/cmd/service-account/token/list"
	"stackit/internal/cmd/service-account/token/revoke"
	"stackit/internal/pkg/args"
	"stackit/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "Provides functionality regarding service account tokens",
		Long:  "Provides functionality regarding service account tokens",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(create.NewCmd())
	cmd.AddCommand(list.NewCmd())
	cmd.AddCommand(revoke.NewCmd())
}
