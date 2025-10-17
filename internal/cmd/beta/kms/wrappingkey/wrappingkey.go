package wrappingkey

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/kms/wrappingkey/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/kms/wrappingkey/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/kms/wrappingkey/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wrapping-key",
		Short: "Manage KMS wrapping keys",
		Long:  "Provides functionality for wrapping key operations inside the KMS",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(create.NewCmd(params))
}
