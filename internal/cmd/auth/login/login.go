package login

import (
	"fmt"
	"os"

	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "login",
	Short:   "Login to the provider",
	Long:    "Login to the provider",
	Example: `$ stackit auth login`,
	Run: func(cmd *cobra.Command, args []string) {
		err := auth.AuthorizeUser()
		if err != nil {
			fmt.Printf("Authorization failed: %s\n", err)
			os.Exit(1)
		}

		fmt.Println("Successfully logged into STACKIT CLI.")
	},
}

func init() {
}
