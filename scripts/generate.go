package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/cmd"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra/doc"
)

const (
	DocsFolder = "docs"
)

func main() {
	repoRoot, err := getGitRepoRoot()
	if err != nil {
		log.Fatalf("Error determining Git repository root: %v", err)
	}
	docsDir := filepath.Join(repoRoot, DocsFolder)
	err = os.RemoveAll(docsDir)
	if err != nil {
		log.Fatalf("Error removing old documentation directory: %v", err)
	}
	err = os.Mkdir(docsDir, 0o750)
	if err != nil {
		log.Fatalf("Error creating new documentation directory: %v", err)
	}

	filePrepender := func(_ string) string {
		return ""
	}
	linkHandler := func(filename string) string {
		return fmt.Sprintf("./%s", filename)
	}
	printer := print.NewPrinter(
		os.Stdin,
		os.Stdout,
		os.Stderr,
	)
	params := &types.CmdParams{
		Printer: printer,
		Args:    os.Args,
	}
	err = doc.GenMarkdownTreeCustom(cmd.NewRootCmd(params), docsDir, filePrepender, linkHandler)
	if err != nil {
		log.Fatalf("Error generating documentation: %v", err)
	}
}

func getGitRepoRoot() (string, error) {
	output, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
