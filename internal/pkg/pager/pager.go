package pager

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// Shows the content in the command's stdout using the "less" command
func Display(cmd *cobra.Command, content string) error {
	lessCmd := exec.Command("less", "-F", "-S", "-w")
	lessCmd.Stdin = strings.NewReader(content)
	lessCmd.Stdout = cmd.OutOrStdout()

	err := lessCmd.Run()
	if err != nil {
		return fmt.Errorf("run less command: %w", err)
	}
	return nil
}
