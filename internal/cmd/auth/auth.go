package auth

import (
	activateserviceaccount "github.com/stackitcloud/stackit-cli/internal/cmd/auth/activate-service-account"
	"github.com/stackitcloud/stackit-cli/internal/cmd/auth/login"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "auth",
	Short:   "Provides authentication functionality",
	Long:    "Provides authentication functionality",
	Example: `$ stackit auth login`,
}

func init() {
	// Add all direct child commands
	Cmd.AddCommand(login.Cmd)
	Cmd.AddCommand(activateserviceaccount.Cmd)
}
