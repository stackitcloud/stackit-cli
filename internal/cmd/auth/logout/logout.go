package logout

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
		Use:   "logout",
		Short: "Logs the user account out of the STACKIT CLI",
		Long:  "Logs the user account out of the STACKIT CLI.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Log out of the STACKIT CLI.`,
				"$ stackit auth logout"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := auth.LogoutUser()
			if err != nil {
				return fmt.Errorf("authorization failed: %w", err)
			}

			p.Info("Successfully logged out of the STACKIT CLI.\n")
			return nil
		},
	}
	return cmd
}
