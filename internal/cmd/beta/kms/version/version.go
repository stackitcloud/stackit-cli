package version

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/kms/version/destroy"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/kms/version/disable"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/kms/version/enable"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/kms/version/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/kms/version/restore"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Manage KMS Key versions",
		Long:  "Provides CRUD functionality for Key Version operations inside the KMS",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(destroy.NewCmd(params))
	cmd.AddCommand(disable.NewCmd(params))
	cmd.AddCommand(enable.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(restore.NewCmd(params))
}
