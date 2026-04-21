package flags

import (
	"fmt"
	"io"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

type testFile struct {
	path, content string
}

func TestSecretFlag(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		value      string
		want       *string
		file       *testFile
		stdin      string
		wantErr    bool
		wantStdErr string
	}{
		{
			name:  "no value: prompts",
			value: "",
			want:  utils.Ptr("from stdin"),
			stdin: "from stdin",
		},
		{
			name:       "a value: prints deprecation",
			value:      "a value",
			want:       utils.Ptr("a value"),
			wantStdErr: "Warning: Passing a secret value on the command line is insecure and deprecated. This usage will stop working October 2026.\n",
		},
		{
			name:  "from an existing file",
			value: "@some-file.txt",
			want:  utils.Ptr("from file"),
			file: &testFile{
				path:    "some-file.txt",
				content: "from file",
			},
		},
		{
			name:    "from a non-existing file",
			value:   "@some-file-with-typo.txt",
			wantErr: true,
			file: &testFile{
				path:    "some-file.txt",
				content: "from file",
			},
		},
		{
			name:  "from an existing double-quoted file",
			value: `@"some-file.txt"`,
			want:  utils.Ptr("from file"),
			file: &testFile{
				path:    "some-file.txt",
				content: "from file",
			},
		},
		{
			name:  "from an existing single-quoted file",
			value: "@'some-file.txt'",
			want:  utils.Ptr("from file"),
			file: &testFile{
				path:    "some-file.txt",
				content: "from file",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			params := testparams.NewTestParams()
			if tt.file != nil {
				params.Fs = fstest.MapFS{
					tt.file.path: &fstest.MapFile{
						Data: []byte(tt.file.content),
					},
				}
			}
			flag := SecretFlag("test", params.CmdParams)
			cmd := cobra.Command{}
			cmd.Flags().Var(flag, "test", flag.Usage())
			if tt.stdin != "" {
				params.In.WriteString(tt.stdin)
				params.In.WriteString("\n")
			}

			if tt.value != "" { // emulate pflag only calling set when flag is specified on the command line
				err := cmd.Flags().Set("test", tt.value)
				if err != nil && !tt.wantErr {
					t.Fatalf("unexpected error: %v", err)
				}
				if err == nil && tt.wantErr {
					t.Fatalf("expected error, got none")
				}
			}

			got := SecretFlagToStringPointer(params.Printer, &cmd, "test")

			if got != tt.want && *got != *tt.want {
				t.Fatalf("unexpected value: got %q, want %q", *got, *tt.want)
			}
			if tt.wantStdErr != "" {
				message, err := params.Err.ReadString('\n')
				if err != nil && err != io.EOF {
					t.Fatalf("reading stderr: %v", err)
				}
				if message != tt.wantStdErr {
					t.Fatalf("unexpected stderr: got %q, want %q", message, tt.wantStdErr)
				}
			}
		})
	}
}

func TestSecretFlag_Usage(t *testing.T) {
	t.Parallel()
	tests := []struct{
		in string
		want string
	} {
		{
			in: "password",
			want: "Password",
		},
		{
			in: "Password",
			want: "Password",
		},
		{
			in: "",
			want: "",
		},
		{
			in: "secret-key",
			want: "Secret-Key",
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%q -> %q", tt.in, tt.want), func(t *testing.T) {
			t.Parallel()
			params := testparams.NewTestParams()
			flag := SecretFlag(tt.in, params.CmdParams)
			got := flag.Usage()
			if !strings.HasPrefix(got, tt.want) {
				t.Fatalf("unexpected usage: got %q, want %q", got, tt.want)
			}
		})
	}
}
