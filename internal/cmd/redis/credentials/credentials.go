package credentials

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/redis/credentials/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/redis/credentials/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/redis/credentials/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/redis/credentials/list"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

// Deprecated: Will be removed after 2027-08-31.
func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:        "credentials",
		Short:      "Provides functionality for Redis credentials",
		Long:       "Provides functionality for Redis credentials.",
		Deprecated: "\nCommand `stackit redis credentials` is deprecated and will be removed after 2027-08-31. Please use `stackit valkey credentials` instead.",
		Args:       args.NoArgs,
		Run:        utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

// Deprecated: Will be removed after 2027-08-31.
func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
}
