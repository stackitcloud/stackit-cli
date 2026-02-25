package kms

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/kms/key"
	"github.com/stackitcloud/stackit-cli/internal/cmd/kms/keyring"
	"github.com/stackitcloud/stackit-cli/internal/cmd/kms/version"
	"github.com/stackitcloud/stackit-cli/internal/cmd/kms/wrappingkey"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kms",
		Short: "Provides functionality for KMS",
		Long:  "Provides functionality for KMS.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(keyring.NewCmd(params))
	cmd.AddCommand(wrappingkey.NewCmd(params))
	cmd.AddCommand(key.NewCmd(params))
	cmd.AddCommand(version.NewCmd(params))
}
