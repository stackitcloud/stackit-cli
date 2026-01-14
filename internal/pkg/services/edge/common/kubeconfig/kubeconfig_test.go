// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 STACKIT GmbH & Co. KG

package kubeconfig

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	testUtils "github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	testErrorMessage = "test error message"
	errStringErrTest = errors.New(testErrorMessage)
)

const (
	kubeconfig_1_yaml = `
apiVersion: v1
clusters:
- cluster:
    server: https://server-1.com
  name: cluster-1
contexts:
- context:
    cluster: cluster-1
    user: user-1
  name: context-1
current-context: context-1
kind: Config
preferences: {}
users:
- name: user-1
  user: {}
`
	kubeconfig_2_yaml = `
apiVersion: v1
clusters:
- cluster:
    server: https://server-2.com
  name: cluster-2
contexts:
- context:
    cluster: cluster-2
    user: user-2
  name: context-2
current-context: context-2
kind: Config
users:
- name: user-2
  user: {}
`
	overwriteKubeconfigTarget = `
apiVersion: v1
clusters:
- cluster:
    server: https://server-1.com
  name: cluster-1
contexts:
- context:
    cluster: cluster-1
    user: user-1
  name: context-1
current-context: context-1
kind: Config
users:
- name: user-1
  user:
    token: old-token
`
	overwriteKubeconfigSource = `
apiVersion: v1
clusters:
- cluster:
    server: https://server-1-new.com
  name: cluster-1
contexts:
- context:
    cluster: cluster-1
    user: user-1
  name: context-1
current-context: context-1
kind: Config
users:
- name: user-1
  user:
    token: new-token
`
)

func TestValidateExpiration(t *testing.T) {
	type args struct {
		expiration *uint64
	}
	tests := []struct {
		name string
		args *args
		want error
	}{
		// Valid cases
		{
			name: "nil expiration",
			args: &args{
				expiration: nil,
			},
		},
		{
			name: "valid expiration - minimum value",
			args: &args{
				expiration: utils.Ptr(uint64(expirationSecondsMin)),
			},
		},
		{
			name: "valid expiration - maximum value",
			args: &args{
				expiration: utils.Ptr(uint64(expirationSecondsMax)),
			},
		},
		{
			name: "valid expiration - default value",
			args: &args{
				expiration: utils.Ptr(uint64(ExpirationSecondsDefault)),
			},
		},
		{
			name: "valid expiration - middle value",
			args: &args{
				expiration: utils.Ptr(uint64(86400)), // 1 day
			},
		},

		// Error cases - below minimum
		{
			name: "expiration too small - below minimum",
			args: &args{
				expiration: utils.Ptr(uint64(expirationSecondsMin - 1)),
			},
			want: fmt.Errorf("%s is too small (minimum is %d seconds)", ExpirationFlag, expirationSecondsMin),
		},
		{
			name: "expiration too small - zero",
			args: &args{
				expiration: utils.Ptr(uint64(0)),
			},
			want: fmt.Errorf("%s is too small (minimum is %d seconds)", ExpirationFlag, expirationSecondsMin),
		},

		// Error cases - above maximum
		{
			name: "expiration too large - above maximum",
			args: &args{
				expiration: utils.Ptr(uint64(expirationSecondsMax + 1)),
			},
			want: fmt.Errorf("%s is too large (maximum is %d seconds)", ExpirationFlag, expirationSecondsMax),
		},
		{
			name: "expiration too large - way above maximum",
			args: &args{
				expiration: utils.Ptr(uint64(9999999999999999999)),
			},
			want: fmt.Errorf("%s is too large (maximum is %d seconds)", ExpirationFlag, expirationSecondsMax),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateExpiration(tt.args.expiration)
			testUtils.AssertError(t, err, tt.want)
		})
	}
}

func TestErrors(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name    string
		args    *args
		wantErr error
	}{
		// EmptyKubeconfigError
		{
			name: "EmptyKubeconfigError",
			args: &args{
				err: &EmptyKubeconfigError{},
			},
			wantErr: &EmptyKubeconfigError{},
		},

		// LoadKubeconfigError
		{
			name: "LoadKubeconfigError",
			args: &args{
				err: &LoadKubeconfigError{Err: errStringErrTest},
			},
			wantErr: errStringErrTest,
		},

		// WriteKubeconfigError
		{
			name: "WriteKubeconfigError",
			args: &args{
				err: &WriteKubeconfigError{Err: errStringErrTest},
			},
			wantErr: errStringErrTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testUtils.AssertError(t, tt.args.err, tt.wantErr)
		})
	}
}

// Already have comprehensive tests for WriteKubeconfig

func TestWriteOptions(t *testing.T) {
	confirmFn := func(_ string) error { return nil }

	type args struct {
		modify func(WriteOptions) WriteOptions
		check  func(*testing.T, WriteOptions)
	}
	tests := []struct {
		name string
		args *args
	}{
		// Default options
		{
			name: "NewWriteOptions creates default options",
			args: &args{
				modify: func(o WriteOptions) WriteOptions { return o },
				check: func(t *testing.T, opts WriteOptions) {
					if opts.Overwrite {
						t.Error("expected Overwrite to be false by default")
					}
					if opts.SwitchContext {
						t.Error("expected SwitchContext to be false by default")
					}
					if opts.ConfirmFn != nil {
						t.Error("expected ConfirmFn to be nil by default")
					}
				},
			},
		},

		// Individual option tests
		{
			name: "WithOverwrite sets overwrite flag",
			args: &args{
				modify: func(o WriteOptions) WriteOptions { return o.WithOverwrite(true) },
				check: func(t *testing.T, opts WriteOptions) {
					if !opts.Overwrite {
						t.Error("expected Overwrite to be true")
					}
				},
			},
		},
		{
			name: "WithSwitchContext sets switch context flag",
			args: &args{
				modify: func(o WriteOptions) WriteOptions { return o.WithSwitchContext(true) },
				check: func(t *testing.T, opts WriteOptions) {
					if !opts.SwitchContext {
						t.Error("expected SwitchContext to be true")
					}
				},
			},
		},
		{
			name: "WithConfirmation sets confirmation callback",
			args: &args{
				modify: func(o WriteOptions) WriteOptions { return o.WithConfirmation(confirmFn) },
				check: func(t *testing.T, opts WriteOptions) {
					if opts.ConfirmFn == nil {
						t.Error("expected ConfirmFn to be set")
					}
				},
			},
		},

		// Chained options
		{
			name: "options are chainable",
			args: &args{
				modify: func(o WriteOptions) WriteOptions {
					return o.WithOverwrite(true).
						WithSwitchContext(true).
						WithConfirmation(confirmFn)
				},
				check: func(t *testing.T, opts WriteOptions) {
					if !opts.Overwrite {
						t.Error("expected Overwrite to be true")
					}
					if !opts.SwitchContext {
						t.Error("expected SwitchContext to be true")
					}
					if opts.ConfirmFn == nil {
						t.Error("expected ConfirmFn to be set")
					}
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := tt.args.modify(NewWriteOptions())
			tt.args.check(t, opts)
		})
	}
}

func TestGetDefaultKubeconfigPath(t *testing.T) {
	type args struct {
		kubeconfigEnv *string // nil means unset
	}
	tests := []struct {
		name string
		args *args
		want string
	}{
		// KUBECONFIG not set
		{
			name: "returns a non-empty path when KUBECONFIG is not set",
			args: &args{kubeconfigEnv: nil},
			want: "",
		},

		// Single path
		{
			name: "returns path from KUBECONFIG if set",
			args: &args{kubeconfigEnv: utils.Ptr("/test/kubeconfig_1_yaml")},
			want: "/test/kubeconfig_1_yaml",
		},

		// Multiple paths
		{
			name: "returns first path from KUBECONFIG if multiple are set",
			args: &args{kubeconfigEnv: utils.Ptr("/test/kubeconfig_1_yaml" + string(os.PathListSeparator) + "/test/kubeconfig_2_yaml")},
			want: "/test/kubeconfig_1_yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original env and restore after test
			oldKubeconfig := os.Getenv("KUBECONFIG")
			defer func() {
				if err := os.Setenv("KUBECONFIG", oldKubeconfig); err != nil {
					t.Logf("failed to restore KUBECONFIG: %v", err)
				}
			}()

			// Setup test environment
			if tt.args.kubeconfigEnv == nil {
				if err := os.Unsetenv("KUBECONFIG"); err != nil {
					t.Fatalf("failed to unset KUBECONFIG: %v", err)
				}
			} else {
				if err := os.Setenv("KUBECONFIG", *tt.args.kubeconfigEnv); err != nil {
					t.Fatalf("failed to set KUBECONFIG: %v", err)
				}
			}

			// Run test
			got := getDefaultKubeconfigPath()

			// If want is empty only make sure the returned path is not empty
			// In that case we don't care about what path is default, only that one is.
			want := filepath.Clean(tt.want)
			if want == filepath.Clean("") {
				if filepath.Clean(got) != "" {
					return
				}
			}

			// Verify results
			testUtils.AssertValue(t, filepath.Clean(got), want)
		})
	}
}

func TestGetKubeconfigPath(t *testing.T) {
	type args struct {
		path      *string
		checkPath func(t *testing.T, path string)
	}
	tests := []struct {
		name    string
		args    *args
		wantErr error
	}{
		{
			name: "uses default path when nil provided",
			args: &args{
				path: nil,
				checkPath: func(t *testing.T, path string) {
					if path == "" {
						t.Error("expected non-empty path")
					}
				},
			},
		},
		{
			name: "validates and returns absolute path when valid path provided",
			args: &args{
				path: utils.Ptr("/tmp/kubeconfig"),
				checkPath: func(t *testing.T, path string) {
					if !filepath.IsAbs(path) {
						t.Error("expected absolute path")
					}
				},
			},
		},
		{
			name: "returns error for invalid path",
			args: &args{
				path: utils.Ptr("."),
			},
			wantErr: &InvalidKubeconfigPathError{Path: "."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := getKubeconfigPath(tt.args.path)
			if !testUtils.AssertError(t, err, tt.wantErr) {
				return
			}
			if tt.args.checkPath != nil {
				tt.args.checkPath(t, path)
			}
		})
	}
}

func TestIsValidFilePath(t *testing.T) {
	type args struct {
		path *string
	}
	tests := []struct {
		name string
		args *args

		want bool
	}{
		{
			name: "valid path",
			args: &args{
				path: utils.Ptr("/test/kubeconfig"),
			},
			want: true,
		},
		{
			name: "nil path",
			args: &args{
				path: nil,
			},
			want: false,
		},
		{
			name: "empty path",
			args: &args{
				path: utils.Ptr(""),
			},
			want: false,
		},
		{
			name: "single dot",
			args: &args{
				path: utils.Ptr("."),
			},
			want: false,
		},
		{
			name: "single slash",
			args: &args{
				path: utils.Ptr("/"),
			},
			want: false,
		},
		{
			name: "relative path with parent directory",
			args: &args{
				path: utils.Ptr("../kubeconfig"),
			},
			want: true,
		},
		{
			name: "path with spaces",
			args: &args{
				path: utils.Ptr("/test/kube config"),
			},
			want: true,
		},
		{
			name: "complex but valid path",
			args: &args{
				path: utils.Ptr("/test/kube-config.d/cluster1/config"),
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidFilePath(tt.args.path); got != tt.want {
				t.Errorf("isValidFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriteKubeconfig(t *testing.T) {
	testPath := filepath.Join(t.TempDir(), "config")
	defaultTempFile := filepath.Join(t.TempDir(), "default-kubeconfig")

	type args struct {
		path      *string
		content   string
		options   WriteOptions
		setupEnv  func()
		checkFile func(t *testing.T, path string)
	}
	tests := []struct {
		name     string
		args     *args
		wantPath *string
		wantErr  any
	}{
		{
			name: "writes new file with default options",
			args: &args{
				path:    &testPath,
				content: kubeconfig_1_yaml,
				options: NewWriteOptions(),
				checkFile: func(t *testing.T, path string) {
					if !isExistingFile(&path) {
						t.Error("file was not created")
					}
				},
			},
			wantPath: &testPath,
		},
		{
			name: "handles invalid file path",
			args: &args{
				path:    utils.Ptr("."),
				content: kubeconfig_1_yaml,
				options: NewWriteOptions(),
			},
			wantErr: &InvalidKubeconfigPathError{Path: "."},
		},
		{
			name: "handles empty kubeconfig",
			args: &args{
				path:    &testPath,
				content: "",
				options: NewWriteOptions(),
			},
			wantErr: &EmptyKubeconfigError{},
		},
		{
			name: "uses default path when nil provided",
			args: &args{
				path:    nil,
				content: kubeconfig_1_yaml,
				options: NewWriteOptions(),
				setupEnv: func() {
					t.Setenv("KUBECONFIG", defaultTempFile)
				},
			},
			wantPath: &defaultTempFile,
		},
		{
			name: "overwrites existing file when option is set",
			args: &args{
				path:    &testPath,
				content: kubeconfig_2_yaml,
				options: NewWriteOptions().WithOverwrite(true),
				setupEnv: func() {
					// Pre-write first file
					if _, err := WriteKubeconfig(&testPath, kubeconfig_1_yaml, NewWriteOptions()); err != nil {
						t.Fatalf("failed to setup test: %v", err)
					}
				},
				checkFile: func(t *testing.T, path string) {
					content, err := os.ReadFile(path)
					if err != nil {
						t.Fatalf("failed to read kubeconfig: %v", err)
					}
					if !strings.Contains(string(content), "server-2.com") {
						t.Error("file was not overwritten")
					}
				},
			},
			wantPath: &testPath,
		},
		{
			name: "respects user confirmation - confirmed",
			args: &args{
				path:    &testPath,
				content: kubeconfig_1_yaml,
				options: NewWriteOptions().WithConfirmation(func(_ string) error {
					return nil
				}),
			},
			wantPath: &testPath,
		},
		{
			name: "respects user confirmation - denied",
			args: &args{
				path:    &testPath,
				content: kubeconfig_1_yaml,
				options: NewWriteOptions().WithConfirmation(func(_ string) error {
					return errStringErrTest
				}),
			},
			wantErr: errStringErrTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.setupEnv != nil {
				tt.args.setupEnv()
			}

			got, gotErr := WriteKubeconfig(tt.args.path, tt.args.content, tt.args.options)
			if !testUtils.AssertError(t, gotErr, tt.wantErr) {
				return
			}

			testUtils.AssertValue(t, got, tt.wantPath)

			if tt.args.checkFile != nil {
				tt.args.checkFile(t, *got)
			}
		})
	}
}

func TestMergeKubeconfig(t *testing.T) {
	type args struct {
		path      *string
		content   string
		switchCtx bool
		setupEnv  func()
	}
	tests := []struct {
		name    string
		args    args
		verify  func(t *testing.T, path string)
		wantErr error
	}{
		{
			name: "merges configs with conflicting names",
			args: args{
				path:      utils.Ptr(filepath.Join(t.TempDir(), "kubeconfig")),
				content:   overwriteKubeconfigSource,
				switchCtx: true,
				setupEnv: func() {
					// Pre-write first file
					if _, err := WriteKubeconfig(utils.Ptr(filepath.Join(t.TempDir(), "kubeconfig")), overwriteKubeconfigTarget, NewWriteOptions()); err != nil {
						t.Fatalf("failed to setup test: %v", err)
					}
				},
			},
			verify: func(t *testing.T, path string) {
				config, err := clientcmd.LoadFromFile(path)
				if err != nil {
					t.Fatalf("failed to load merged config: %v", err)
				}

				cluster := config.Clusters["cluster-1"]
				if cluster.Server != "https://server-1-new.com" {
					t.Errorf("expected server to be 'https://server-1-new.com', got '%s'", cluster.Server)
				}

				user := config.AuthInfos["user-1"]
				if user.Token != "new-token" {
					t.Errorf("expected token to be 'new-token', got '%s'", user.Token)
				}
			},
		},
		{
			name: "handles nil file path",
			args: args{
				path:      nil,
				content:   kubeconfig_1_yaml,
				switchCtx: false,
			},
			wantErr: fmt.Errorf("no kubeconfig file provided to be merged"),
		},
		{
			name: "handles invalid config",
			args: args{
				path:      utils.Ptr(filepath.Join(t.TempDir(), "kubeconfig")),
				content:   "invalid yaml",
				switchCtx: false,
			},
			wantErr: &LoadKubeconfigError{},
		},
		{
			name: "handles empty config",
			args: args{
				path:      utils.Ptr(filepath.Join(t.TempDir(), "kubeconfig")),
				content:   "",
				switchCtx: false,
			},
			wantErr: &EmptyKubeconfigError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.setupEnv != nil {
				tt.args.setupEnv()
			}

			err := mergeKubeconfig(tt.args.path, tt.args.content, tt.args.switchCtx)
			if !testUtils.AssertError(t, err, tt.wantErr) {
				return
			}

			if tt.verify != nil {
				if tt.args.path == nil {
					t.Fatalf("expected path to be set")
				}
				tt.verify(t, *tt.args.path)
			}
		})
	}
}
