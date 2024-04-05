package login

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Logs in to the STACKIT CLI",
		Long:  "Logs in to the STACKIT CLI using a user account.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Login to the STACKIT CLI. This command will open a browser window where you can login to your STACKIT account`,
				"$ stackit auth login"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := auth.AuthorizeUser()
			if err != nil {
				return fmt.Errorf("authorization failed: %w", err)
			}

			p.Info("Successfully logged into STACKIT CLI.")
			return nil
		},
	}
	return cmd
}
