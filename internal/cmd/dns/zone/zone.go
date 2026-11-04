package zone

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns/zone/clone"
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns/zone/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns/zone/delete"
	deletemovecode "github.com/stackitcloud/stackit-cli/internal/cmd/dns/zone/delete-move-code"
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns/zone/describe"
	generatemovecode "github.com/stackitcloud/stackit-cli/internal/cmd/dns/zone/generate-move-code"
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns/zone/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns/zone/move"
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns/zone/update"
	validatemovecode "github.com/stackitcloud/stackit-cli/internal/cmd/dns/zone/validate-move-code"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zone",
		Short: "Provides functionality for DNS zones",
		Long:  "Provides functionality for DNS zones.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(clone.NewCmd(params))
	cmd.AddCommand(generatemovecode.NewCmd(params))
	cmd.AddCommand(move.NewCmd(params))
	cmd.AddCommand(validatemovecode.NewCmd(params))
	cmd.AddCommand(deletemovecode.NewCmd(params))
}
