package list

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	testRegion = "eu01"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &iaas.APIClient{}
var testProjectId = uuid.NewString()

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		limitFlag:                 "10",
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			Verbosity: globalflags.VerbosityDefault,
			ProjectId: testProjectId,
			Region:    testRegion,
		},
		Limit:    utils.Ptr(int64(10)),
		MinVCPUs: nil,
		MinRAM:   nil,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *iaas.ApiListMachineTypesRequest)) iaas.ApiListMachineTypesRequest {
	request := testClient.ListMachineTypes(testCtx, testProjectId, testRegion)
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description   string
		argValues     []string
		flagValues    map[string]string
		isValid       bool
		expectedModel *inputModel
	}{
		{
			description:   "base",
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "no values",
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "filter by resources valid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[minVcpuFlag] = "16"
				flagValues[minRamFlag] = "32"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.MinVCPUs = utils.Ptr(int64(16))
				model.MinRAM = utils.Ptr(int64(32))
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}

func TestFilterMachineTypes(t *testing.T) {
	items := []iaas.MachineType{
		{
			Name:        utils.Ptr("c3i.16"),
			Vcpus:       utils.Ptr(int64(16)),
			Ram:         utils.Ptr(int64(32768)),
			Description: utils.Ptr("Intel Emerald Rapids 8580 CPU instance"),
			ExtraSpecs: &map[string]interface{}{
				"hw:cpu_sockets":   "1",
				"hw:mem_page_size": "large",
				"aggregate":        "intel-gen3-oc-cpu-optimized",
			},
		},
		{
			Name:        utils.Ptr("c3i.28"),
			Vcpus:       utils.Ptr(int64(28)),
			Ram:         utils.Ptr(int64(60416)),
			Description: utils.Ptr("Intel Emerald Rapids 8580 CPU instance"),
			ExtraSpecs: &map[string]interface{}{
				"cpu":              "intel-emerald-rapids-8580-dual-socket",
				"overcommit":       "4",
				"hw:mem_page_size": "large",
			},
		},
		{
			Name:        utils.Ptr("g2i.1"),
			Vcpus:       utils.Ptr(int64(1)),
			Ram:         utils.Ptr(int64(4096)),
			Description: utils.Ptr("Intel Ice Lake 4316 CPU instance"),
		},
	}

	tests := []struct {
		description string
		items       *[]iaas.MachineType
		model       *inputModel
		expectedLen int
	}{
		{
			description: "base filters",
			items:       &items,
			model:       fixtureInputModel(),
			expectedLen: 3,
		},
		{
			description: "nil items slice",
			items:       nil,
			model:       fixtureInputModel(),
			expectedLen: 0,
		},
		{
			description: "filter min-vcpu 20",
			items:       &items,
			model:       fixtureInputModel(func(m *inputModel) { m.MinVCPUs = utils.Ptr(int64(20)) }),
			expectedLen: 1, // c3i.28 only
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			result := filterMachineTypes(tt.items, tt.model)
			if len(result) != tt.expectedLen {
				t.Errorf("expected %d items, got %d", tt.expectedLen, len(result))
			}
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest iaas.ApiListMachineTypesRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	type args struct {
		outputFormat string
		machineTypes iaas.MachineTypeListResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty response",
			args:    args{},
			wantErr: false,
		},
		{
			name: "response with extra specs",
			args: args{
				outputFormat: "table",
				machineTypes: iaas.MachineTypeListResponse{
					Items: &[]iaas.MachineType{
						{
							Name:  utils.Ptr("c3i.16"),
							Vcpus: utils.Ptr(int64(16)),
							Ram:   utils.Ptr(int64(32768)),
							ExtraSpecs: &map[string]interface{}{
								"aggregate":  "intel-gen3",
								"overcommit": 4,
							},
							Description: utils.Ptr("Intel Emerald Rapids 8580 CPU instance"),
						},
					},
				},
			},
			wantErr: false,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&types.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.outputFormat, tt.args.machineTypes); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
