package login

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to the STACKIT CLI",
		Long:  "Login to the STACKIT CLI",
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

			cmd.Println("Successfully logged into STACKIT CLI.")
			return nil
		},
	}
	return cmd
}
