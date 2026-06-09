package flags

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestStringEnumSliceFlag_Set(t *testing.T) {
	tests := []struct {
		name       string
		options    []string
		ignoreCase bool
		setValue   string
		want       []string
		wantErr    bool
	}{
		{
			name:     "valid value",
			options:  []string{"a", "b", "c"},
			setValue: "a",
			want:     []string{"a"},
			wantErr:  false,
		},
		{
			name:     "multiple valid values",
			options:  []string{"a", "b", "c"},
			setValue: "a,b",
			want:     []string{"a", "b"},
			wantErr:  false,
		},
		{
			name:     "multiple valid values with spaces",
			options:  []string{"a", "b", "c"},
			setValue: "a, b ,c",
			want:     []string{"a", "b", "c"},
			wantErr:  false,
		},
		{
			name:     "invalid value",
			options:  []string{"a", "b", "c"},
			setValue: "d",
			wantErr:  true,
		},
		{
			name:     "partially invalid value",
			options:  []string{"a", "b", "c"},
			setValue: "a,d",
			wantErr:  true,
		},
		{
			name:     "empty value",
			options:  []string{"a", "b", "c"},
			setValue: "",
			wantErr:  true,
		},
		{
			name:       "case sensitive mismatch",
			options:    []string{"A", "B"},
			setValue:   "a",
			ignoreCase: false,
			wantErr:    true,
		},
		{
			name:       "case insensitive match",
			options:    []string{"A", "B"},
			setValue:   "a",
			ignoreCase: true,
			want:       []string{"a"},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := []StringEnumSliceFlagOption[string]{}
			if tt.ignoreCase {
				opts = append(opts, IgnoreCase[string]())
			}
			f := StringEnumSliceFlag("test", tt.options, "docs", opts...)

			err := f.Set(tt.setValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				got := f.Get()
				if len(got) != len(tt.want) {
					t.Errorf("Set() got = %v, want %v", got, tt.want)
					return
				}
				for i := range got {
					if got[i] != tt.want[i] {
						t.Errorf("Set() got = %v, want %v", got, tt.want)
						break
					}
				}
			}
		})
	}
}

func TestStringEnumSliceFlag_DefaultValues(t *testing.T) {
	f := StringEnumSliceFlag("test", []string{"a", "b"}, "docs", DefaultValues("a"))

	got := f.Get()
	if len(got) != 1 || got[0] != "a" {
		t.Errorf("Expected default value [a], got %v", got)
	}

	// Setting a value should override the default
	err := f.Set("b")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	got = f.Get()
	if len(got) != 1 || got[0] != "b" {
		t.Errorf("Expected value [b] after Set, got %v", got)
	}

	// Setting another value should append
	err = f.Set("a")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	got = f.Get()
	if len(got) != 2 || got[0] != "b" || got[1] != "a" {
		t.Errorf("Expected value [b, a] after second Set, got %v", got)
	}
}

func TestStringEnumSliceFlag_Usage(t *testing.T) {
	f := StringEnumSliceFlag("test", []string{"a", "b"}, "docs")
	usage := f.Usage()
	if usage != "docs (possible values: [a, b])" {
		t.Errorf("Expected usage 'docs (possible values: [a, b])', got %q", usage)
	}
}

func TestStringEnumSliceFlag_UnknownDefaultOpenAPI(t *testing.T) {
	f := StringEnumSliceFlag("test", []string{"a", "unknown_default_open_api", "b"}, "docs")
	usage := f.Usage()
	if usage != "docs (possible values: [a, b])" {
		t.Errorf("Expected unknown_default_open_api to be filtered out, got %q", usage)
	}
}

func TestStringEnumSliceFlag_Register(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	f := StringEnumSliceFlag("my-flag", []string{"a", "b"}, "docs")
	f.Register(cmd)

	flag := cmd.Flags().Lookup("my-flag")
	if flag == nil {
		t.Errorf("Expected flag 'my-flag' to be registered")
	}
	if flag.Usage != "docs (possible values: [a, b])" {
		t.Errorf("Expected flag usage to be set correctly")
	}
}
