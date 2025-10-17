package list

import (
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"

	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

func TestParseInput(t *testing.T) {
	projectId := uuid.New().String()
	tests := []struct {
		description   string
		argValues     []string
		flagValues    map[string]string
		expectedModel *inputModel
		isValid       bool
	}{
		{
			description: "valid project id",
			flagValues: map[string]string{
				"project-id": projectId,
			},
			expectedModel: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: projectId,
					Verbosity: globalflags.InfoVerbosity,
				},
			},
			isValid: true,
		},
		{
			description: "missing project id does not lead into error",
			flagValues:  map[string]string{},
			expectedModel: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					Verbosity: globalflags.InfoVerbosity,
				},
			},
			isValid: true,
		},
		{
			description: "valid input with limit",
			flagValues: map[string]string{
				"limit": "10",
			},
			expectedModel: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					Verbosity: globalflags.InfoVerbosity,
				},
				Limit: utils.Ptr(int64(10)),
			},
			isValid: true,
		},
		{
			description: "valid input without limit",
			flagValues:  map[string]string{},
			expectedModel: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					Verbosity: globalflags.InfoVerbosity,
				},
			},
			isValid: true,
		},
		{
			description: "invalid limit (zero)",
			flagValues: map[string]string{
				"limit": "0",
			},
			expectedModel: nil,
			isValid:       false,
		},
		{
			description: "invalid limit (negative)",
			flagValues: map[string]string{
				"limit": "-1",
			},
			expectedModel: nil,
			isValid:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}

func TestOutputResult(t *testing.T) {
	tests := []struct {
		name           string
		outputFormat   string
		publicIpRanges []iaas.PublicNetwork
		expectedOutput string
		wantErr        bool
	}{
		{
			name:         "JSON output single",
			outputFormat: "json",
			publicIpRanges: []iaas.PublicNetwork{
				{Cidr: utils.Ptr("192.168.0.0/24")},
			},
			wantErr: false,
		},
		{
			name:         "JSON output multiple",
			outputFormat: "json",
			publicIpRanges: []iaas.PublicNetwork{
				{Cidr: utils.Ptr("192.168.0.0/24")},
				{Cidr: utils.Ptr("192.167.0.0/24")},
			},
			wantErr: false,
		},
		{
			name:         "YAML output single",
			outputFormat: "yaml",
			publicIpRanges: []iaas.PublicNetwork{
				{Cidr: utils.Ptr("192.168.0.0/24")},
			},
			wantErr: false,
		},
		{
			name:         "YAML output multiple",
			outputFormat: "yaml",
			publicIpRanges: []iaas.PublicNetwork{
				{Cidr: utils.Ptr("192.168.0.0/24")},
				{Cidr: utils.Ptr("192.167.0.0/24")},
			},
			wantErr: false,
		},
		{
			name:         "pretty output single",
			outputFormat: "pretty",
			publicIpRanges: []iaas.PublicNetwork{
				{Cidr: utils.Ptr("192.168.0.0/24")},
			},
			wantErr: false,
		},
		{
			name:         "pretty output multiple",
			outputFormat: "pretty",
			publicIpRanges: []iaas.PublicNetwork{
				{Cidr: utils.Ptr("192.168.0.0/24")},
				{Cidr: utils.Ptr("192.167.0.0/24")},
			},
			wantErr: false,
		},
		{
			name:         "default output",
			outputFormat: "",
			publicIpRanges: []iaas.PublicNetwork{
				{Cidr: utils.Ptr("192.168.0.0/24")},
			},
			wantErr: false,
		},
		{
			name:           "empty list",
			outputFormat:   "json",
			publicIpRanges: []iaas.PublicNetwork{},
			wantErr:        false,
		},
		{
			name:         "nil CIDR",
			outputFormat: "pretty",
			publicIpRanges: []iaas.PublicNetwork{
				{Cidr: nil},
				{Cidr: utils.Ptr("192.168.0.0/24")},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := print.NewPrinter()
			p.Cmd = NewCmd(&params.CmdParams{Printer: p})
			err := outputResult(p, tt.outputFormat, tt.publicIpRanges)
			if (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
