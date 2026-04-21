package flags

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

type secretFlag struct {
	printer *print.Printer
	fs      fs.FS
	value   string
	name    string
}

func SecretFlag(name string, params *types.CmdParams) *secretFlag {
	f := &secretFlag{
		printer: params.Printer,
		fs:      params.Fs,
		name:    name,
	}
	return f
}

var _ pflag.Value = &secretFlag{}

func (f *secretFlag) String() string {
	return f.value
}

func (f *secretFlag) Set(value string) error {
	if strings.HasPrefix(value, "@") {
		path := strings.Trim(value[1:], `"'`)
		bytes, err := fs.ReadFile(f.fs, path)
		if err != nil {
			return fmt.Errorf("reading secret %s: %w", f.name, err)
		}
		f.value = string(bytes)
		return nil
	}
	f.printer.Warn("Passing a secret value on the command line is insecure and deprecated. This usage will stop working October 2026.\n")
	f.value = value
	return nil
}

func (f *secretFlag) Type() string {
	return "string"
}

func (f *secretFlag) Usage() string {
	name := cases.Title(language.AmericanEnglish).String(f.name)
	return fmt.Sprintf("%s. Can be a string (deprecated) or a file path, if prefixed with '@' (example: @./secret.txt). Will be read from stdin when empty.", name)
}

func SecretFlagToStringPointer(p *print.Printer, cmd *cobra.Command, flag string) *string {
	value, err := cmd.Flags().GetString(flag)
	if err != nil {
		p.Debug(print.ErrorLevel, "convert secret flag to string pointer: %v", err)
		return nil
	}
	if value == "" {
		input, err := p.PromptForPassword(fmt.Sprintf("enter %s: ", flag))
		if err != nil {
			p.Debug(print.ErrorLevel, "convert secret flag %q to string pointer: %v", flag, err)
			return nil
		}
		return &input
	}
	if cmd.Flag(flag).Changed {
		return &value
	}
	return nil
}
