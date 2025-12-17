package grafana

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/observability/grafana/describe"
	publicreadaccess "github.com/stackitcloud/stackit-cli/internal/cmd/observability/grafana/public-read-access"
	singlesignon "github.com/stackitcloud/stackit-cli/internal/cmd/observability/grafana/single-sign-on"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grafana",
		Short: "Provides functionality for the Grafana configuration of Observability instances",
		Long:  "Provides functionality for the Grafana configuration of Observability instances.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(publicreadaccess.NewCmd(params))
	cmd.AddCommand(singlesignon.NewCmd(params))
}
