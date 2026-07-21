package create

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	iaas "github.com/stackitcloud/stackit-sdk-go/services/iaas/v2api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
)

const (
	testRegion                    = "eu01"
	testName                      = "example-network-area-name"
	testTransferNetwork           = "100.0.0.0/24"
	testDefaultPrefixLength int64 = 25
	testMaxPrefixLength     int64 = 26
	testMinPrefixLength     int64 = 24
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &iaas.APIClient{DefaultAPI: &iaas.DefaultAPIService{}}
var testOrgId = uuid.NewString()

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.RegionFlag: testRegion,

		nameFlag:           testName,
		organizationIdFlag: testOrgId,
		labelFlag:          "key=value",
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
			Region:    testRegion,
		},
		Name:           "example-network-area-name",
		OrganizationId: testOrgId,
		Labels: map[string]any{
			"key": "value",
		},
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *iaas.ApiCreateNetworkAreaRequest)) iaas.ApiCreateNetworkAreaRequest {
	request := testClient.DefaultAPI.CreateNetworkArea(testCtx, testOrgId)
	request = request.CreateNetworkAreaPayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayload(mods ...func(payload *iaas.CreateNetworkAreaPayload)) iaas.CreateNetworkAreaPayload {
	payload := iaas.CreateNetworkAreaPayload{
		Name: "example-network-area-name",
		Labels: map[string]any{
			"key": "value",
		},
	}
	for _, mod := range mods {
		mod(&payload)
	}
	return payload
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description   string
		argValues     []string
		flagValues    map[string]string
		aclValues     []string
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
			description: "name missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, nameFlag)
			}),
			isValid: false,
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
			description: "labels missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, labelFlag)
			}),
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Labels = nil
			}),
			isValid: true,
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
		expectedRequest iaas.ApiCreateNetworkAreaRequest
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
				cmpopts.EquateComparable(testCtx, iaas.DefaultAPIService{}),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func Test_outputResult(t *testing.T) {
	type args struct {
		outputFormat string
		orgLabel     string
		responses    *NetworkAreaResponses
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty",
			args:    args{},
			wantErr: true,
		},
		{
			name: "set empty response",
			args: args{
				responses: &NetworkAreaResponses{},
			},
			wantErr: false,
		},
		{
			name: "set empty network area",
			args: args{
				responses: &NetworkAreaResponses{
					NetworkArea: iaas.NetworkArea{},
				},
			},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.outputFormat, tt.args.orgLabel, tt.args.responses); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
