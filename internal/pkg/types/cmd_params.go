package types

import (
	"io/fs"

	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
)

type CmdParams struct {
	Printer    *print.Printer
	CliVersion string
	Date       string
	Fs         fs.FS
	Args       []string
}
