package artifactory

import (
	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/cmd/registry/artifactory/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/registry/artifactory/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/registry/artifactory/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/registry/artifactory/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/registry/artifactory/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "artifactory",
		Short: "Provides functionality for Container Registry artifactories",
		Long:  "Provides functionality for Container Registry artifactories.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
}
