package keyring

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/kms/keyring/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/kms/keyring/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/kms/keyring/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/kms/keyring/list"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keyring",
		Short: "Manage KMS key rings",
		Long:  "Provides functionality for key ring operations inside the KMS",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
}
