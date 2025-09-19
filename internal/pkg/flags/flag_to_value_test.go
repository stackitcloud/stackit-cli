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
