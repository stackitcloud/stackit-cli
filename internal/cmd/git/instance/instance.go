package git

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/git/instance/create"
	delete2 "github.com/stackitcloud/stackit-cli/internal/cmd/git/instance/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/git/instance/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/git/instance/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instance",
		Short: "Provides functionality for STACKIT Git",
		Long:  "Provides functionality for STACKIT Git.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(
		list.NewCmd(params),
		describe.NewCmd(params),
		create.NewCmd(params),
		delete2.NewCmd(params),
	)
}
