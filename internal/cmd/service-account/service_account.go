package serviceaccount

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
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

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service-account",
		Short: "Provides functionality for service accounts",
		Long:  "Provides functionality for service accounts.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(getjwks.NewCmd(params))

	cmd.AddCommand(key.NewCmd(params))
	cmd.AddCommand(token.NewCmd(params))
}
