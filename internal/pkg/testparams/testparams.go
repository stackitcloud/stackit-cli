package testparams

import (
	"bytes"

	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

type TestParams struct {
	*types.CmdParams
	In, Out, Err *bytes.Buffer
}

func NewTestParams() *TestParams {
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	err := &bytes.Buffer{}
	return &TestParams{
		&types.CmdParams{
			Printer: print.NewPrinter(
				in, out, err,
			),
		},
		in,
		out,
		err,
	}
}
