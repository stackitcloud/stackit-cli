package update

import (
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "update security groups",
		Long:  "update security groups",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(`example 1`, `foo bar baz`),
			examples.NewExample(`example 2`, `foo bar baz`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeUpdate(cmd, p, args)
		},
	}
	cmd.Flags().String("dummy", "foo", "fooify")
	return cmd
}

func executeUpdate(cmd *cobra.Command, p *print.Printer, args []string) error {
	p.Info("executing update command")
	return nil
}
