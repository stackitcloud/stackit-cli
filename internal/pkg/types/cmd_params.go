package types

import (
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
)

type CmdParams struct {
	Printer    *print.Printer
	CliVersion string
}
