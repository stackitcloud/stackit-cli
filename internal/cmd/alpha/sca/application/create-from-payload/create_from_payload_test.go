package createfrompayload

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	sca "github.com/stackitcloud/stackit-sdk-go/services/sca/v1alphaapi"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
)

const (
	testRegion = "eu01"
)

var testClient = &sca.APIClient{DefaultAPI: &sca.DefaultAPIService{}}
var testCtx = context.Background()

var (
	testProjectID     = uuid.NewString()
	testEnvironmentID = uuid.NewString()

	testPayload = &sca.CreateApplicationPayload{
		DisplayName:          "test-application",
		AdditionalProperties: map[string]any{},
		Containers: []sca.Container{{
			Name:                 "test-container",
			Image:                "test-image",
			Cpu:                  sca.PtrInt32(2000),
			Memory:               sca.PtrInt32(2048),
			AdditionalProperties: map[string]any{},
		}},
		Scaling: sca.Scaling{
			Type: sca.SCALINGTYPE_SCALING_TYPE_MANUAL,
			ManualScaling: &sca.ManualScaling{
				Instances:            2,
				AdditionalProperties: map[string]any{},
			},
			AdditionalProperties: map[string]any{},
		},
		Network: sca.Network{
			PublicIngress:        true,
			Port:                 sca.PtrInt32(8888),
			AdditionalProperties: map[string]any{},
		},
	}
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectID,
		globalflags.RegionFlag:    testRegion,

		environmentIDFlag: testEnvironmentID,
		nameFlag:          "test-application",
		payloadFlag: `{
			"displayName": "",
			"containers": [{
				"name": "test-container",
				"image": "test-image",
				"cpu": 2000,
				"memory": 2048
			}],
			"scaling": {
				"type": "SCALING_TYPE_MANUAL",
				"manualScaling": {
					"instances": 2
				}
			},
			"network": {
				"publicIngress": true,
				"port": 8888
			}
		}`,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectID,
			Region:    testRegion,
			Verbosity: globalflags.VerbosityDefault,
		},
		EnvironmentID: testEnvironmentID,
		Payload:       testPayload,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *sca.ApiCreateApplicationRequest)) sca.ApiCreateApplicationRequest {
	request := testClient.DefaultAPI.CreateApplication(testCtx, testProjectID, testEnvironmentID)
	request = request.CreateApplicationPayload(*testPayload)
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		desc          string
		flagValues    map[string]string
		expectedModel *inputModel
		isValid       bool
	}{
		{
			desc:          "base",
			flagValues:    fixtureFlagValues(),
			expectedModel: fixtureInputModel(),
			isValid:       true,
		},
		{
			desc:       "no values",
			flagValues: nil,
			isValid:    false,
		},
		{
			desc: "project id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			desc: "project id invalid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			desc: "invalid json",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[payloadFlag] = "not json"
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, nil, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest sca.ApiCreateApplicationRequest
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
				cmpopts.EquateComparable(testCtx, sca.DefaultAPIService{}),
				cmpopts.SortSlices(func(a, b sca.EnvVar) bool {
					return a.Key < b.Key
				}),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}
