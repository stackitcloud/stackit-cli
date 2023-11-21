package auth

import (
	activateserviceaccount "github.com/stackitcloud/stackit-cli/internal/cmd/auth/activate-service-account"
	"github.com/stackitcloud/stackit-cli/internal/cmd/auth/login"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "auth",
		Short:   "Provides authentication functionality",
		Long:    "Provides authentication functionality",
		Example: `$ stackit auth login`,
	}
	addChilds(cmd)
	return cmd
}

func addChilds(cmd *cobra.Command) {
	cmd.AddCommand(login.NewCmd())
	cmd.AddCommand(activateserviceaccount.NewCmd())
}
