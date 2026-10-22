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

type secretSource uint8

const (
	secretUnset secretSource = iota
	secretFile
	secretStdin
	secretInvalid
)

// Fixed deadline, six months after the introduction of the compatibility helper.
const autoPromptRemovalDate = "2027-04"

type secretFlag struct {
	fs     fs.FS
	name   string
	source secretSource
	path   string
}

// SecretFlag accepts @file or - (stdin).
// ResolveSecretFlag must be called after argument validation: parsing neither validates nor reads secret sources.
func SecretFlag(name string, params *types.CmdParams) *secretFlag {
	return &secretFlag{fs: params.Fs, name: name}
}

var _ pflag.Value = &secretFlag{}

// String describes the source, never the resolved secret or an invalid literal.
func (f *secretFlag) String() string {
	switch f.source {
	case secretFile:
		return "@" + f.path
	case secretStdin:
		return "-"
	case secretInvalid:
		return "<invalid>"
	default:
		return ""
	}
}

func (f *secretFlag) Set(value string) error {
	f.path = ""
	switch {
	case value == "-":
		f.source = secretStdin
	case strings.HasPrefix(value, "@"):
		f.path = strings.Trim(value[1:], `"'`)
		f.source = secretFile
		if f.path == "" {
			f.source = secretInvalid
		}
	default:
		f.source = secretInvalid
	}
	// pflag includes the supplied token when wrapping a Set error. Defer errors
	// until resolution, and do not retain rejected literals. The final occurrence
	// wins, without performing I/O for overwritten sources.
	return nil
}

func (f *secretFlag) Type() string {
	return "secret"
}

func (f *secretFlag) Usage() string {
	name := cases.Title(language.AmericanEnglish).String(f.name)
	return fmt.Sprintf("%s. Use @<file> to read from a file (example: @./secret.txt), or - to read from stdin (hidden input on a terminal). Literal values are not supported.", name)
}

// ResolveSecretFlag returns nil for an omitted flag, without reading input.
// Files are read verbatim, including empty files. Stdin supplies a single line;
// LF/CRLF is removed, spaces are preserved, and EOF after data is accepted.
// Call this after command-specific validation and propagate errors before making
// API requests. Do not use GetString: the flag stores a source, not a secret.
func ResolveSecretFlag(p *print.Printer, cmd *cobra.Command, name string) (*string, error) {
	flag := cmd.Flag(name)
	if flag == nil {
		return nil, fmt.Errorf("secret flag --%s not found", name)
	}
	value, ok := flag.Value.(*secretFlag)
	if !ok {
		return nil, fmt.Errorf("flag --%s is not a secret flag", name)
	}
	if !flag.Changed {
		return nil, nil
	}

	var secret string
	switch value.source {
	case secretFile:
		data, err := fs.ReadFile(value.fs, value.path)
		if err != nil {
			return nil, fmt.Errorf("reading secret --%s: %w", name, err)
		}
		secret = string(data)
	case secretStdin:
		input, err := p.PromptForPassword(fmt.Sprintf("enter %s: ", name))
		if err != nil {
			return nil, fmt.Errorf("reading secret --%s from stdin: %w", name, err)
		}
		secret = input
	default:
		return nil, fmt.Errorf("invalid source for --%s: use @<file> or - for stdin; literal secrets are not supported", name)
	}
	return &secret, nil
}

// ResolveSecretFlagWithAutoPrompt preserves implicit prompting for legacy callers.
// Only omission triggers the warning and prompt, never an explicitly empty file.
//
// Deprecated: Use ResolveSecretFlag and require explicit - for stdin. Do not use
// this helper in new commands. Automatic prompting will be removed 2027-04.
func ResolveSecretFlagWithAutoPrompt(p *print.Printer, cmd *cobra.Command, name string) (*string, error) {
	value, err := ResolveSecretFlag(p, cmd, name)
	if err != nil || value != nil {
		return value, err
	}
	p.Warn("Automatically reading --%s from stdin is deprecated and will stop working %s. Use --%s - instead.\n", name, autoPromptRemovalDate, name)
	input, err := p.PromptForPassword(fmt.Sprintf("enter %s: ", name))
	if err != nil {
		return nil, fmt.Errorf("reading secret --%s from stdin: %w", name, err)
	}
	return &input, nil
}
