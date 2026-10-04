package updatefrompayload

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
	testApplicationID = uuid.NewString()

	testPayload = &sca.UpdateApplicationPayload{
		Containers: []sca.Container{{
			Name:                 "test-container",
			Image:                "test-image",
			Cpu:                  sca.PtrInt32(2000),
			Memory:               sca.PtrInt32(2048),
			AdditionalProperties: map[string]any{},
		}},
		Scaling: &sca.Scaling{
			Type: sca.SCALINGTYPE_SCALING_TYPE_MANUAL,
			ManualScaling: &sca.ManualScaling{
				Instances:            2,
				AdditionalProperties: map[string]any{},
			},
			AdditionalProperties: map[string]any{},
		},
		Network: &sca.Network{
			PublicIngress:        true,
			Port:                 sca.PtrInt32(8888),
			AdditionalProperties: map[string]any{},
		},
	}
)

func fixtureArgsValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testApplicationID,
	}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectID,
		globalflags.RegionFlag:    testRegion,
		environmentIDFlag:         testEnvironmentID,
		payloadFlag: `{
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
		ApplicationID: testApplicationID,
		EnvironmentID: testEnvironmentID,
		Payload:       testPayload,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *sca.ApiUpdateApplicationRequest)) sca.ApiUpdateApplicationRequest {
	request := testClient.DefaultAPI.UpdateApplication(testCtx, testProjectID, testEnvironmentID, testApplicationID)
	request = request.UpdateApplicationPayload(*testPayload)
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		desc          string
		argsValues    []string
		flagValues    map[string]string
		expectedModel *inputModel
		isValid       bool
	}{
		{
			desc:          "base",
			argsValues:    fixtureArgsValues(),
			flagValues:    fixtureFlagValues(),
			expectedModel: fixtureInputModel(),
			isValid:       true,
		},
		{
			desc:       "no values",
			argsValues: nil,
			flagValues: nil,
			isValid:    false,
		},
		{
			desc:       "application id missing",
			argsValues: []string{},
			flagValues: fixtureFlagValues(),
			isValid:    false,
		},
		{
			desc:       "project id missing",
			argsValues: fixtureArgsValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			desc:       "project id invalid",
			argsValues: fixtureArgsValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			desc:       "invalid json",
			argsValues: fixtureArgsValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[payloadFlag] = "not json"
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argsValues, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest sca.ApiUpdateApplicationRequest
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
