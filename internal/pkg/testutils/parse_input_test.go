// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 STACKIT GmbH & Co. KG

package testutils

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

type parseInputTestModel struct {
	Value       string
	Args        []string
	RepeatValue []string
	hidden      string
}

func newTestCmdFactory(flagSetup func(*cobra.Command)) func(*types.CmdParams) *cobra.Command {
	return func(*types.CmdParams) *cobra.Command {
		cmd := &cobra.Command{Use: "test"}
		if flagSetup != nil {
			flagSetup(cmd)
		}
		return cmd
	}
}

func TestRunParseInputCase(t *testing.T) {
	sentinel := errors.New("parse failed")
	tests := []struct {
		name            string
		flagSetup       func(*cobra.Command)
		flags           map[string]string
		repeatFlags     map[string][]string
		args            []string
		cmpOpts         []ParseInputCaseOption
		wantModel       *parseInputTestModel
		wantErr         any
		parseFunc       func(*print.Printer, *cobra.Command, []string) (*parseInputTestModel, error)
		expectParseCall bool
	}{
		{
			name: "success",
			flagSetup: func(cmd *cobra.Command) {
				cmd.Flags().String("name", "", "")
			},
			flags:     map[string]string{"name": "edge"},
			cmpOpts:   []ParseInputCaseOption{WithParseInputCmpOptions(WithAllowUnexported(parseInputTestModel{}))},
			wantModel: &parseInputTestModel{Value: "edge", hidden: "protected"},
			parseFunc: func(_ *print.Printer, cmd *cobra.Command, _ []string) (*parseInputTestModel, error) {
				val, _ := cmd.Flags().GetString("name")
				return &parseInputTestModel{Value: val, hidden: "protected"}, nil
			},
			expectParseCall: true,
		},
		{
			name: "flag set failure",
			flagSetup: func(cmd *cobra.Command) {
				cmd.Flags().Int("count", 0, "")
			},
			flags:   map[string]string{"count": "invalid"},
			wantErr: "invalid syntax",
			parseFunc: func(_ *print.Printer, _ *cobra.Command, _ []string) (*parseInputTestModel, error) {
				return &parseInputTestModel{}, nil
			},
			expectParseCall: false,
		},
		{
			name: "flag group validation",
			flagSetup: func(cmd *cobra.Command) {
				cmd.Flags().String("first", "", "")
				cmd.Flags().String("second", "", "")
				cmd.MarkFlagsRequiredTogether("first", "second")
			},
			flags:   map[string]string{"first": "only"},
			wantErr: "must all be set",
			parseFunc: func(_ *print.Printer, _ *cobra.Command, _ []string) (*parseInputTestModel, error) {
				return &parseInputTestModel{}, nil
			},
			expectParseCall: false,
		},
		{
			name: "parse func error",
			flagSetup: func(cmd *cobra.Command) {
				cmd.Flags().Bool("ok", false, "")
			},
			flags:   map[string]string{"ok": "true"},
			wantErr: sentinel,
			parseFunc: func(_ *print.Printer, _ *cobra.Command, _ []string) (*parseInputTestModel, error) {
				return nil, sentinel
			},
			expectParseCall: true,
		},
		{
			name: "args success",
			flagSetup: func(cmd *cobra.Command) {
				cmd.Args = cobra.ExactArgs(1)
			},
			args:      []string{"arg1"},
			cmpOpts:   []ParseInputCaseOption{WithParseInputCmpOptions(WithAllowUnexported(parseInputTestModel{}))},
			wantModel: &parseInputTestModel{Args: []string{"arg1"}},
			parseFunc: func(_ *print.Printer, _ *cobra.Command, args []string) (*parseInputTestModel, error) {
				return &parseInputTestModel{Args: args}, nil
			},
			expectParseCall: true,
		},
		{
			name: "args validation failure",
			flagSetup: func(cmd *cobra.Command) {
				cmd.Args = cobra.NoArgs
			},
			args:    []string{"arg1"},
			wantErr: "unknown command",
			parseFunc: func(_ *print.Printer, _ *cobra.Command, _ []string) (*parseInputTestModel, error) {
				return &parseInputTestModel{}, nil
			},
			expectParseCall: false,
		},
		{
			name: "repeat flags success",
			flagSetup: func(cmd *cobra.Command) {
				cmd.Flags().StringSlice("tags", []string{}, "")
			},
			repeatFlags: map[string][]string{"tags": {"tag1", "tag2"}},
			cmpOpts:     []ParseInputCaseOption{WithParseInputCmpOptions(WithAllowUnexported(parseInputTestModel{}))},
			wantModel:   &parseInputTestModel{RepeatValue: []string{"tag1", "tag2"}},
			parseFunc: func(_ *print.Printer, cmd *cobra.Command, _ []string) (*parseInputTestModel, error) {
				val, _ := cmd.Flags().GetStringSlice("tags")
				return &parseInputTestModel{RepeatValue: val}, nil
			},
			expectParseCall: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmdFactory := newTestCmdFactory(tt.flagSetup)
			var parseCalled bool
			parseFn := tt.parseFunc
			if parseFn == nil {
				parseFn = func(*print.Printer, *cobra.Command, []string) (*parseInputTestModel, error) {
					return &parseInputTestModel{}, nil
				}
			}

			RunParseInputCase(t, ParseInputTestCase[*parseInputTestModel]{
				Name:        tt.name,
				Flags:       tt.flags,
				RepeatFlags: tt.repeatFlags,
				Args:        tt.args,
				WantModel:   tt.wantModel,
				WantErr:     tt.wantErr,
				CmdFactory:  cmdFactory,
				ParseInputFunc: func(pr *print.Printer, cmd *cobra.Command, args []string) (*parseInputTestModel, error) {
					parseCalled = true
					return parseFn(pr, cmd, args)
				},
			}, tt.cmpOpts...)

			if parseCalled != tt.expectParseCall {
				t.Fatalf("parseCalled = %v, expect %v", parseCalled, tt.expectParseCall)
			}
		})
	}
}
