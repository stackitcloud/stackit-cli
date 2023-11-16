package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:               "stackit",
	Short:             "The root command of the STACKIT CLI",
	Long:              "The root command of the STACKIT CLI",
	SilenceUsage:      true,
	DisableAutoGenTag: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Welcome to the STACKIT CLI! Functionality will be added soon.")
		return nil
	},
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
