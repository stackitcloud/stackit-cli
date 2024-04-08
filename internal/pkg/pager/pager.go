package pager

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
)

// Shows the content in the command's stdout using the "less" command
func Display(p *print.Printer, content string) error {
	lessCmd := exec.Command("less", "-F", "-S", "-w")
	lessCmd.Stdin = strings.NewReader(content)
	lessCmd.Stdout = p.Cmd.OutOrStdout()

	err := lessCmd.Run()
	if err != nil {
		return fmt.Errorf("run less command: %w", err)
	}
	return nil
}
