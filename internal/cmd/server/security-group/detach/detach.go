package detach

import (
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detach",
		Short: "Detach a security group from a server",
		Long:  "Detach a security group from a server.",
		Run: func(cmd *cobra.Command, args []string) {
			params.Printer.Info("Detaching security group from server...")
		},
	}
	return cmd
}
