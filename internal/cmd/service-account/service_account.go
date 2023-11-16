package serviceaccount

import (
	"stackit/internal/cmd/service-account/create"
	"stackit/internal/cmd/service-account/delete"
	getjwks "stackit/internal/cmd/service-account/get-jwks"
	"stackit/internal/cmd/service-account/key"
	"stackit/internal/cmd/service-account/list"
	"stackit/internal/cmd/service-account/token"
	"stackit/internal/pkg/args"
	"stackit/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service-account",
		Short: "Provides functionality for service accounts",
		Long:  "Provides functionality for service accounts",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(list.NewCmd())
	cmd.AddCommand(create.NewCmd())
	cmd.AddCommand(delete.NewCmd())
	cmd.AddCommand(getjwks.NewCmd())

	cmd.AddCommand(key.NewCmd())
	cmd.AddCommand(token.NewCmd())
}
