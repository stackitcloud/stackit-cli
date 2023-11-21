package login

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "login",
	Short:   "Login to the provider",
	Long:    "Login to the provider",
	Example: `$ stackit auth login`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := auth.AuthorizeUser()
		if err != nil {
			return fmt.Errorf("authorization failed: %w", err)
		}

		cmd.Println("Successfully logged into STACKIT CLI.")
		return nil
	},
}

func init() {
}
