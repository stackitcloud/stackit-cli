package list

import (
	"testing"

	"github.com/google/go-cmp/cmp"
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
		globalFlags   map[string]string
		expectedModel *inputModel
		isValid       bool
	}{
		{
			description: "valid project id",
			globalFlags: map[string]string{
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
			globalFlags: map[string]string{},
			expectedModel: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					Verbosity: globalflags.InfoVerbosity,
				},
			},
			isValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			p := print.NewPrinter()
			cmd := NewCmd(&params.CmdParams{Printer: p})
			err := globalflags.Configure(cmd.Flags())
			if err != nil {
				t.Fatal(err)
			}

			for flag, value := range tt.globalFlags {
				if err := cmd.Flags().Set(flag, value); err != nil {
					t.Fatalf("Failed to set global flag %s: %v", flag, err)
				}
			}

			model, err := parseInput(p, cmd)
			if !tt.isValid && err == nil {
				t.Fatalf("parseInput() error = %v, wantErr %v", err, !tt.isValid)
			}

			if tt.isValid {
				if diff := cmp.Diff(model, tt.expectedModel); diff != "" {
					t.Fatalf("Model mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	tests := []struct {
		name           string
		outputFormat   string
		networkList    iaas.PublicNetworkListResponse
		expectedOutput string
		wantErr        bool
	}{
		{
			name:         "JSON output single",
			outputFormat: "json",
			networkList: iaas.PublicNetworkListResponse{
				Items: &[]iaas.PublicNetwork{
					{Cidr: utils.Ptr("192.168.0.0/24")},
				},
			},
			wantErr: false,
		},
		{
			name:         "JSON output multiple",
			outputFormat: "json",
			networkList: iaas.PublicNetworkListResponse{
				Items: &[]iaas.PublicNetwork{
					{Cidr: utils.Ptr("192.168.0.0/24")},
					{Cidr: utils.Ptr("192.167.0.0/24")},
				},
			},
			wantErr: false,
		},
		{
			name:         "YAML output single",
			outputFormat: "yaml",
			networkList: iaas.PublicNetworkListResponse{
				Items: &[]iaas.PublicNetwork{
					{Cidr: utils.Ptr("192.168.0.0/24")},
				},
			},
			wantErr: false,
		},
		{
			name:         "YAML output multiple",
			outputFormat: "yaml",
			networkList: iaas.PublicNetworkListResponse{
				Items: &[]iaas.PublicNetwork{
					{Cidr: utils.Ptr("192.168.0.0/24")},
					{Cidr: utils.Ptr("192.167.0.0/24")},
				},
			},
			wantErr: false,
		},
		{
			name:         "pretty output single",
			outputFormat: "pretty",
			networkList: iaas.PublicNetworkListResponse{
				Items: &[]iaas.PublicNetwork{
					{Cidr: utils.Ptr("192.168.0.0/24")},
				},
			},
			wantErr: false,
		},
		{
			name:         "pretty output multiple",
			outputFormat: "pretty",
			networkList: iaas.PublicNetworkListResponse{
				Items: &[]iaas.PublicNetwork{
					{Cidr: utils.Ptr("192.168.0.0/24")},
					{Cidr: utils.Ptr("192.167.0.0/24")},
				},
			},
			wantErr: false,
		},
		{
			name:         "default output",
			outputFormat: "",
			networkList: iaas.PublicNetworkListResponse{
				Items: &[]iaas.PublicNetwork{
					{Cidr: utils.Ptr("192.168.0.0/24")},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := print.NewPrinter()
			p.Cmd = NewCmd(&params.CmdParams{Printer: p})
			err := outputResult(p, tt.outputFormat, tt.networkList)
			if (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
