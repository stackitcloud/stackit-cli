package flags

import (
	"errors"
	"io"
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
)

type forbiddenSecretReader struct{ t *testing.T }

func (r forbiddenSecretReader) Read([]byte) (int, error) {
	r.t.Fatal("unexpected stdin read")
	return 0, io.EOF
}

type forbiddenSecretFS struct{ t *testing.T }

func (f forbiddenSecretFS) Open(string) (fs.File, error) {
	f.t.Fatal("unexpected file read")
	return nil, fs.ErrNotExist
}

// These test commands exercise pflag token parsing, Cobra argument validation,
// and resolution in RunE.
func newSecretTestCommand(t *testing.T, params *testparams.TestParams, autoPrompt bool, got **string, applied *bool) *cobra.Command {
	t.Helper()
	cmd := &cobra.Command{
		Args:          cobra.ExactArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var err error
			if autoPrompt {
				*got, err = ResolveSecretFlagWithAutoPrompt(params.Printer, cmd, "password")
			} else {
				*got, err = ResolveSecretFlag(params.Printer, cmd, "password")
			}
			if err != nil {
				return err
			}
			*applied = true
			return nil
		},
	}
	flag := SecretFlag("password", params.CmdParams)
	cmd.Flags().Var(flag, "password", flag.Usage())
	cmd.Flags().Bool("other", false, "another option")
	cmd.SetOut(params.Out)
	cmd.SetErr(params.Err)
	return cmd
}

func TestSecretFlagCommandParsing(t *testing.T) {
	t.Parallel()
	const fileBytes = " \x00\r\n@other-file\n-\n "
	for _, tt := range []struct {
		name        string
		args        []string
		want        *string
		stdin       bool
		wantErr     string
		wantNoRun   bool
		forbidFiles bool
	}{
		{
			name: "omitted",
			args: []string{"RESOURCE"},
		},
		{
			name: "spaced file",
			args: []string{"RESOURCE", "--password", "@secret.txt"},
			want: new("file-secret\n"),
		},
		{
			name: "attached file",
			args: []string{"RESOURCE", "--password=@secret.txt"},
			want: new("file-secret\n"),
		},
		{
			name: "file before resource",
			args: []string{"--password", "@secret.txt", "RESOURCE"},
			want: new("file-secret\n"),
		},
		{
			name: "path with spaces",
			args: []string{"RESOURCE", "--password", "@secret file.txt"},
			want: new("spaced-path-secret"),
		},
		{
			name: "double-quoted path",
			args: []string{"RESOURCE", "--password", `@"secret file.txt"`},
			want: new("spaced-path-secret"),
		},
		{
			name: "single-quoted path",
			args: []string{"RESOURCE", "--password", "@'secret file.txt'"},
			want: new("spaced-path-secret"),
		},
		{
			name:  "spaced stdin",
			args:  []string{"RESOURCE", "--password", "-"},
			stdin: true,
			want:  new("stdin-secret"),
		},
		{
			name:  "attached stdin",
			args:  []string{"RESOURCE", "--password=-"},
			stdin: true,
			want:  new("stdin-secret"),
		},
		{
			name:  "stdin before resource",
			args:  []string{"--password", "-", "RESOURCE"},
			stdin: true,
			want:  new("stdin-secret"),
		},
		{
			name:  "other flag after source",
			args:  []string{"RESOURCE", "--password", "-", "--other"},
			stdin: true,
			want:  new("stdin-secret"),
		},
		{
			name:      "bare flag",
			args:      []string{"RESOURCE", "--password"},
			wantErr:   "flag needs an argument",
			wantNoRun: true,
		},
		{
			name:      "consumed resource",
			args:      []string{"--password", "RESOURCE"},
			wantErr:   "accepts 1 arg(s)",
			wantNoRun: true,
		},
		{
			name:      "consumed adjacent flag",
			args:      []string{"RESOURCE", "--password", "--other"},
			wantErr:   "invalid source",
			wantNoRun: true,
		},
		{
			name:  "terminator after source",
			args:  []string{"--password=-", "--", "RESOURCE"},
			stdin: true,
			want:  new("stdin-secret"),
		},
		{
			name:      "consumed terminator",
			args:      []string{"--password", "--", "RESOURCE"},
			wantErr:   "invalid source",
			wantNoRun: true,
		},
		{
			name:      "flag after terminator is positional",
			args:      []string{"RESOURCE", "--", "--password=-"},
			wantErr:   "accepts 1 arg(s)",
			wantNoRun: true,
		},
		{
			name:      "explicit empty",
			args:      []string{"RESOURCE", "--password="},
			wantErr:   "invalid source",
			wantNoRun: true,
		},
		{
			name:      "empty file path",
			args:      []string{"RESOURCE", "--password=@"},
			wantErr:   "invalid source",
			wantNoRun: true,
		},
		{
			name:      "quoted empty file path",
			args:      []string{"RESOURCE", `--password=@""`},
			wantErr:   "invalid source",
			wantNoRun: true,
		},
		{
			name:      "literal",
			args:      []string{"RESOURCE", "--password", "literal-secret"},
			wantErr:   "literal secrets are not supported",
			wantNoRun: true,
		},
		{
			name:      "attached literal",
			args:      []string{"RESOURCE", "--password=literal-secret"},
			wantErr:   "literal secrets are not supported",
			wantNoRun: true,
		},
		{
			name: "empty file is explicit",
			args: []string{"RESOURCE", "--password=@empty.txt"},
			want: new(""),
		},
		{
			name: "file containing stdin marker",
			args: []string{"RESOURCE", "--password=@marker.txt"},
			want: new("-"),
		},
		{
			name: "file bytes preserved",
			args: []string{"RESOURCE", "--password=@bytes.txt"},
			want: new(fileBytes),
		},
		{
			name:      "missing file",
			args:      []string{"RESOURCE", "--password=@missing.txt"},
			wantErr:   "reading secret --password",
			wantNoRun: true,
		},
		{
			name: "last source wins",
			args: []string{"RESOURCE", "--password=-", "--password=@secret.txt"},
			want: new("file-secret\n"),
		},
		{
			name:  "overwritten missing file",
			args:  []string{"RESOURCE", "--password=@missing.txt", "--password=-"},
			stdin: true,
			want:  new("stdin-secret"),
		},
		{
			name:  "overwritten literal",
			args:  []string{"RESOURCE", "--password=literal-secret", "--password=-"},
			stdin: true,
			want:  new("stdin-secret"),
		},
		{
			name:      "invalid final source",
			args:      []string{"RESOURCE", "--password=-", "--password=literal-secret"},
			wantErr:   "invalid source",
			wantNoRun: true,
		},
		{
			name:        "invalid arguments with stdin",
			args:        []string{"RESOURCE", "EXTRA", "--password=-"},
			wantErr:     "accepts 1 arg(s)",
			wantNoRun:   true,
			forbidFiles: true,
		},
		{
			name:        "invalid arguments with file",
			args:        []string{"RESOURCE", "EXTRA", "--password=@secret.txt"},
			wantErr:     "accepts 1 arg(s)",
			wantNoRun:   true,
			forbidFiles: true,
		},
		{
			name:        "help does not read stdin",
			args:        []string{"--password=-", "--help"},
			wantNoRun:   true,
			forbidFiles: true,
		},
		{
			name:        "help does not read file",
			args:        []string{"--password=@missing.txt", "--help"},
			wantNoRun:   true,
			forbidFiles: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			params := testparams.NewTestParams()
			params.Fs = fstest.MapFS{
				"secret.txt":      {Data: []byte("file-secret\n")},
				"secret file.txt": {Data: []byte("spaced-path-secret")},
				"empty.txt":       {Data: []byte{}},
				"marker.txt":      {Data: []byte("-")},
				"bytes.txt":       {Data: []byte(fileBytes)},
			}
			if tt.forbidFiles {
				params.Fs = forbiddenSecretFS{t}
			}
			if tt.stdin {
				params.Printer.AssumeYes = true
				params.In.WriteString("stdin-secret\n")
			} else {
				params.Printer.StdIn = forbiddenSecretReader{t}
			}
			var got *string
			var applied bool
			cmd := newSecretTestCommand(t, params, false, &got, &applied)
			cmd.SetArgs(tt.args)
			err := cmd.Execute()
			if tt.wantErr == "" && err != nil || tt.wantErr != "" && (err == nil || !strings.Contains(err.Error(), tt.wantErr)) {
				t.Fatalf("got error %v, want %q", err, tt.wantErr)
			}
			if applied == tt.wantNoRun {
				t.Fatalf("mutation reached: %v, want %v", applied, !tt.wantNoRun)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("unexpected secret (-want +got):\n%s", diff)
			}
		})
	}
}

func TestSecretFlagAutoPrompt(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name    string
		args    []string
		stdin   *string
		want    *string
		wantErr bool
		readErr error
	}{
		{
			name:  "omitted",
			args:  []string{"RESOURCE"},
			stdin: new("stdin-secret\n"),
			want:  new("stdin-secret"),
		},
		{
			name:  "explicit stdin",
			args:  []string{"RESOURCE", "--password=-"},
			stdin: new("stdin-secret\n"),
			want:  new("stdin-secret"),
		},
		{
			name:    "omitted with failed input",
			args:    []string{"RESOURCE"},
			stdin:   new(""),
			wantErr: true,
			readErr: io.EOF,
		},
		{
			name: "empty file",
			args: []string{"RESOURCE", "--password=@empty.txt"},
			want: new(""),
		},
		{
			name:    "literal",
			args:    []string{"RESOURCE", "--password=literal-secret"},
			wantErr: true,
		},
		{
			name:    "missing file",
			args:    []string{"RESOURCE", "--password=@missing.txt"},
			wantErr: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			params := testparams.NewTestParams()
			params.Fs = fstest.MapFS{"empty.txt": {}}
			if tt.stdin != nil {
				params.In.WriteString(*tt.stdin)
			} else {
				params.Printer.StdIn = forbiddenSecretReader{t}
			}
			var got *string
			var applied bool
			cmd := newSecretTestCommand(t, params, true, &got, &applied)
			cmd.SetArgs(tt.args)
			err := cmd.Execute()
			if (err != nil) != tt.wantErr || applied == tt.wantErr {
				t.Fatalf("got applied=%v, err=%v", applied, err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("unexpected secret (-want +got):\n%s", diff)
			}
			if tt.readErr != nil && !errors.Is(err, tt.readErr) {
				t.Fatalf("got error %v, want %v", err, tt.readErr)
			}
		})
	}
}

type failedSecretReader struct{ err error }

func (r failedSecretReader) Read([]byte) (int, error) { return 0, r.err }

func TestSecretFlagStdin(t *testing.T) {
	t.Parallel()
	readErr := errors.New("input failed")
	for _, tt := range []struct {
		name    string
		input   io.Reader
		want    *string
		wantErr error
	}{
		{
			name:  "CRLF",
			input: strings.NewReader("secret\r\n"),
			want:  new("secret"),
		},
		{
			name:  "EOF after data",
			input: strings.NewReader("secret"),
			want:  new("secret"),
		},
		{
			name:  "spaces",
			input: strings.NewReader("  secret  \n"),
			want:  new("  secret  "),
		},
		{
			name:  "bare CR preserved",
			input: strings.NewReader("secret\r"),
			want:  new("secret\r"),
		},
		{
			name:  "empty line is explicit",
			input: strings.NewReader("\n"),
			want:  new(""),
		},
		{
			name:    "empty input",
			input:   strings.NewReader(""),
			wantErr: io.EOF,
		},
		{
			name:    "read failure",
			input:   failedSecretReader{err: readErr},
			wantErr: readErr,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			params := testparams.NewTestParams()
			params.Printer.StdIn = tt.input
			params.Printer.AssumeYes = true
			var got *string
			var applied bool
			cmd := newSecretTestCommand(t, params, false, &got, &applied)
			cmd.SetArgs([]string{"RESOURCE", "--password=-"})
			err := cmd.Execute()
			if !errors.Is(err, tt.wantErr) || applied != (tt.wantErr == nil) {
				t.Fatalf("got applied=%v, err=%v; want err=%v", applied, err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("unexpected secret (-want +got):\n%s", diff)
			}
		})
	}
}

func TestSecretFlagMultipleReadsNoReadAhead(t *testing.T) {
	t.Parallel()
	params := testparams.NewTestParams()
	params.In.WriteString("first-secret\nsecond-secret\ny\n")
	secrets := []struct {
		name string
		want string
	}{
		{
			name: "first",
			want: "first-secret",
		},
		{
			name: "second",
			want: "second-secret",
		},
	}
	cmd := &cobra.Command{
		Use:  "secret-test",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			for _, secret := range secrets {
				got, err := ResolveSecretFlag(params.Printer, cmd, secret.name)
				if err != nil {
					t.Fatalf("resolving %s: %v", secret.name, err)
				}
				if diff := cmp.Diff(new(secret.want), got); diff != "" {
					t.Fatalf("unexpected value for --%s (-want +got):\n%s", secret.name, diff)
				}
			}
			// consumes the `y` in StdIn
			return params.Printer.PromptForConfirmation("Continue?")
		},
	}
	for _, secret := range secrets {
		flag := SecretFlag(secret.name, params.CmdParams)
		cmd.Flags().Var(flag, secret.name, flag.Usage())
	}
	cmd.SetArgs([]string{"--first=-", "--second=-"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("secret reads consumed confirmation input: %v", err)
	}
}
