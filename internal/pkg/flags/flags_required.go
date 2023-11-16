package flags

import "github.com/spf13/cobra"

// Marks all given flags as required, causing the command to report an error if invoked without them.
func MarkFlagsRequired(cmd *cobra.Command, flags ...string) error {
	for _, flag := range flags {
		err := cmd.MarkFlagRequired(flag)
		if err != nil {
			return err
		}
	}
	return nil
}
