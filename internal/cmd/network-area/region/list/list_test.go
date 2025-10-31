package list

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &iaas.APIClient{}

var (
	testAreaId = uuid.NewString()
	testOrgId  = uuid.NewString()
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		networkAreaIdFlag:  testAreaId,
		organizationIdFlag: testOrgId,
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
		},
		OrganizationId: testOrgId,
		NetworkAreaId:  testAreaId,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *iaas.ApiListNetworkAreaRegionsRequest)) iaas.ApiListNetworkAreaRegionsRequest {
	request := testClient.ListNetworkAreaRegions(testCtx, testOrgId, testAreaId)
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
			description: "org id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, organizationIdFlag)
			}),
			isValid: false,
		},
		{
			description: "org id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[organizationIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "org id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[organizationIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "area id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, networkAreaIdFlag)
			}),
			isValid: false,
		},
		{
			description: "area id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[networkAreaIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "area id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[networkAreaIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest iaas.ApiListNetworkAreaRegionsRequest
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

func Test_outputResult(t *testing.T) {
	type args struct {
		outputFormat     string
		networkAreaLabel string
		regionalArea     iaas.RegionalAreaListResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty",
			args:    args{},
			wantErr: false,
		},
		{
			name: "set empty response",
			args: args{
				regionalArea: iaas.RegionalAreaListResponse{},
			},
			wantErr: false,
		},
		{
			name: "set nil for regions map in response",
			args: args{
				regionalArea: iaas.RegionalAreaListResponse{
					Regions: nil,
				},
			},
			wantErr: false,
		},
		{
			name: "set empty map for regions map in response",
			args: args{
				regionalArea: iaas.RegionalAreaListResponse{
					Regions: utils.Ptr(map[string]iaas.RegionalArea{}),
				},
			},
			wantErr: false,
		},
		{
			name: "set empty region in response",
			args: args{
				regionalArea: iaas.RegionalAreaListResponse{
					Regions: utils.Ptr(map[string]iaas.RegionalArea{
						"eu01": {},
					}),
				},
			},
			wantErr: false,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&params.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.outputFormat, tt.args.networkAreaLabel, tt.args.regionalArea); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
