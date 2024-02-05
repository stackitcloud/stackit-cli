package serviceaccount

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/service-account/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/service-account/delete"
	getjwks "github.com/stackitcloud/stackit-cli/internal/cmd/service-account/get-jwks"
	"github.com/stackitcloud/stackit-cli/internal/cmd/service-account/key"
	"github.com/stackitcloud/stackit-cli/internal/cmd/service-account/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/service-account/token"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service-account",
		Short: "Provides functionality for service accounts",
		Long:  "Provides functionality for service accounts.",
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
