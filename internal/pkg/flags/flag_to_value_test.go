package flags

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func TestFlagToStringToStringPointer(t *testing.T) {
	const flagName = "labels"

	tests := []struct {
		name      string
		flagValue *string
		want      *map[string]string
	}{
		{
			name:      "flag unset",
			flagValue: nil,
			want:      nil,
		},
		{
			name:      "flag set with single value",
			flagValue: utils.Ptr("foo=bar"),
			want: &map[string]string{
				"foo": "bar",
			},
		},
		{
			name:      "flag set with multiple values",
			flagValue: utils.Ptr("foo=bar,label1=value1,label2=value2"),
			want: &map[string]string{
				"foo":    "bar",
				"label1": "value1",
				"label2": "value2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := print.NewPrinter()
			// create a new, simple test command with a string-to-string flag
			cmd := func() *cobra.Command {
				cmd := &cobra.Command{
					Use:   "greet",
					Short: "A simple greeting command",
					Long:  "A simple greeting command",
					Run: func(_ *cobra.Command, _ []string) {
						fmt.Println("Hello world")
					},
				}
				cmd.Flags().StringToString(flagName, nil, "Labels are key-value string pairs.")
				return cmd
			}()

			// set the flag value if a value use given, else consider the flag unset
			if tt.flagValue != nil {
				err := cmd.Flags().Set(flagName, *tt.flagValue)
				if err != nil {
					t.Error(err)
				}
			}

			if got := FlagToStringToStringPointer(p, cmd, flagName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FlagToStringToStringPointer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlagToStringArrayValue(t *testing.T) {
	const flagName = "geofencing"
	tests := []struct {
		name       string
		flagValues []string
		want       []string
	}{
		{
			name:       "flag unset",
			flagValues: nil,
			want:       nil,
		},
		{
			name: "single flag value",
			flagValues: []string{
				"https://foo.example.com DE,CH",
			},
			want: []string{
				"https://foo.example.com DE,CH",
			},
		},
		{
			name: "multiple flag value",
			flagValues: []string{
				"https://foo.example.com DE,CH",
				"https://bar.example.com AT",
			},
			want: []string{
				"https://foo.example.com DE,CH",
				"https://bar.example.com AT",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := print.NewPrinter()
			cmd := func() *cobra.Command {
				cmd := &cobra.Command{
					Use:   "greet",
					Short: "A simple greeting command",
					Long:  "A simple greeting command",
					Run: func(_ *cobra.Command, _ []string) {
						fmt.Println("Hello world")
					},
				}
				cmd.Flags().StringArray(flagName, []string{}, "url to multiple region codes, repeatable")
				return cmd
			}()
			// set the flag value if a value use given, else consider the flag unset
			if tt.flagValues != nil {
				for _, val := range tt.flagValues {
					err := cmd.Flags().Set(flagName, val)
					if err != nil {
						t.Error(err)
					}
				}
			}

			if got := FlagToStringArrayValue(p, cmd, flagName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FlagToStringArrayValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlagToInt32Pointer(t *testing.T) {
	const flagName = "limit"
	tests := []struct {
		name      string
		flagValue *string
		want      *int32
	}{
		{
			name:      "flag unset",
			flagValue: nil,
			want:      nil,
		},
		{
			name:      "flag value",
			flagValue: utils.Ptr("42"),
			want:      utils.Ptr(int32(42)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := print.NewPrinter()
			cmd := func() *cobra.Command {
				cmd := &cobra.Command{
					Use:   "greet",
					Short: "A simple greeting command",
					Long:  "A simple greeting command",
					Run: func(_ *cobra.Command, _ []string) {
						fmt.Println("Hello world")
					},
				}
				cmd.Flags().Int32(flagName, 0, "limit")
				return cmd
			}()
			// set the flag value if a value use given, else consider the flag unset
			if tt.flagValue != nil {
				err := cmd.Flags().Set(flagName, *tt.flagValue)
				if err != nil {
					t.Error(err)
				}
			}

			if got := FlagToInt32Pointer(p, cmd, flagName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FlagToInt32Pointer() = %v, want %v", got, tt.want)
			}
		})
	}
}
