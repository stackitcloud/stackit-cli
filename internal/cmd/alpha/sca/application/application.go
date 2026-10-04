package application

import (
	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/cmd/sca/application/create"
	createfrompayload "github.com/stackitcloud/stackit-cli/internal/cmd/sca/application/create-from-payload"
	"github.com/stackitcloud/stackit-cli/internal/cmd/sca/application/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/sca/application/describe"
	generatepayload "github.com/stackitcloud/stackit-cli/internal/cmd/sca/application/generate-payload"
	"github.com/stackitcloud/stackit-cli/internal/cmd/sca/application/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/sca/application/update"
	updatefrompayload "github.com/stackitcloud/stackit-cli/internal/cmd/sca/application/update-from-payload"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "application",
		Short: "Provides functionality for SCA applications",
		Long:  "Provides functionality for STACKIT Container Applications (SCA) cluster.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(createfrompayload.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(generatepayload.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
	cmd.AddCommand(updatefrompayload.NewCmd(params))
}
