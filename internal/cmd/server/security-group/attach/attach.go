package attach

import (
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attach",
		Short: "Attach a security group to a server",
		Long:  "Attach a security group to a server.",
		Run: func(cmd *cobra.Command, args []string) {
			params.Printer.Info("Attaching security group to server...")
		},
	}
	return cmd
}
