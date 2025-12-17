package profile

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/cmd/config/profile/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/config/profile/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/config/profile/export"
	importProfile "github.com/stackitcloud/stackit-cli/internal/cmd/config/profile/import"
	"github.com/stackitcloud/stackit-cli/internal/cmd/config/profile/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/config/profile/set"
	"github.com/stackitcloud/stackit-cli/internal/cmd/config/profile/unset"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage the CLI configuration profiles",
		Long: fmt.Sprintf("%s\n%s\n%s\n%s",
			"Manage the CLI configuration profiles.",
			`The profile to be used can be managed via the "STACKIT_CLI_PROFILE" environment variable or using the "stackit config profile set PROFILE" and "stackit config profile unset" commands.`,
			"The environment variable takes precedence over what is set via the commands.",
			"When no profile is set, the default profile is used.",
		),
		Args: args.NoArgs,
		Run:  utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(set.NewCmd(params))
	cmd.AddCommand(unset.NewCmd(params))
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(importProfile.NewCmd(params))
	cmd.AddCommand(export.NewCmd(params))
}
