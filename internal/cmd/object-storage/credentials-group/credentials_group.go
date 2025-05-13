package credentialsgroup

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/object-storage/credentials-group/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/object-storage/credentials-group/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/object-storage/credentials-group/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "credentials-group",
		Short: "Provides functionality for Object Storage credentials group",
		Long:  "Provides functionality for Object Storage credentials group.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
}
