package key

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/kms/key/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/kms/key/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/kms/key/importKey"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/kms/key/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/kms/key/restore"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/kms/key/rotate"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "key",
		Short: "Manage KMS Keys",
		Long:  "Provides CRUD functionality for Key operations inside the KMS",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(importKey.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(restore.NewCmd(params))
	cmd.AddCommand(rotate.NewCmd(params))
}
