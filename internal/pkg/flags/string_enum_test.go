package flags

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestStringEnumFlag_Set(t *testing.T) {
	tests := []struct {
		name       string
		options    []string
		ignoreCase bool
		setValue   string
		want       string
		wantErr    bool
	}{
		{
			name:     "valid value",
			options:  []string{"a", "b", "c"},
			setValue: "a",
			want:     "a",
		},
		{
			name:     "invalid value",
			options:  []string{"a", "b", "c"},
			setValue: "d",
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
			want:       "a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := []StringEnumFlagOption[string]{}
			if tt.ignoreCase {
				opts = append(opts, StringEnumIgnoreCase[string]())
			}
			f := StringEnumFlag("test", tt.options, "docs", opts...)

			err := f.Set(tt.setValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				got := f.Get()
				if got != tt.want {
					t.Errorf("Set() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestStringEnumFlag_DefaultValue(t *testing.T) {
	f := StringEnumFlag("test", []string{"a", "b"}, "docs", StringEnumDefaultValue("a"))

	got := f.Get()
	if got != "a" {
		t.Errorf("Expected default value a, got %v", got)
	}

	// Setting a value should override the default
	err := f.Set("b")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	got = f.Get()
	if got != "b" {
		t.Errorf("Expected value b after Set, got %v", got)
	}
}

func TestStringEnumFlag_Usage(t *testing.T) {
	f := StringEnumFlag("test", []string{"a", "b"}, "docs")
	usage := f.Usage()
	if usage != "docs (possible values: [a, b])" {
		t.Errorf("Expected usage 'docs (possible values: [a, b])', got %q", usage)
	}
}

func TestStringEnumFlag_UnknownDefaultOpenAPI(t *testing.T) {
	f := StringEnumFlag("test", []string{"a", "unknown_default_open_api", "b"}, "docs")
	usage := f.Usage()
	if usage != "docs (possible values: [a, b])" {
		t.Errorf("Expected unknown_default_open_api to be filtered out, got %q", usage)
	}
}

func TestStringEnumFlag_Register(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	f := StringEnumFlag("my-flag", []string{"a", "b"}, "docs")
	f.Register(cmd.Flags())

	flag := cmd.Flags().Lookup("my-flag")
	if flag == nil {
		t.Fatalf("Expected flag 'my-flag' to be registered")
	}
	if flag.Usage != "docs (possible values: [a, b])" {
		t.Errorf("Expected flag usage to be set correctly")
	}
}

func TestStringEnumFlag_Ptr(t *testing.T) {
	f := StringEnumFlag("test", []string{"a", "b"}, "docs")
	if f.Ptr() != nil {
		t.Errorf("Expected Ptr() to be nil initially, got %v", *f.Ptr())
	}

	err := f.Set("a")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	ptr := f.Ptr()
	if ptr == nil {
		t.Errorf("Expected Ptr() to not be nil after Set")
	} else if *ptr != "a" {
		t.Errorf("Expected Ptr() to point to 'a', got %v", *ptr)
	}

	fWithDefault := StringEnumFlag("test_default", []string{"a", "b"}, "docs", StringEnumDefaultValue("b"))
	ptrDefault := fWithDefault.Ptr()
	if ptrDefault == nil {
		t.Errorf("Expected Ptr() to not be nil with default value")
	} else if *ptrDefault != "b" {
		t.Errorf("Expected Ptr() to point to 'b' with default value, got %v", *ptrDefault)
	}
}
