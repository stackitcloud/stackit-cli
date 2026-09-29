package create

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	sca "github.com/stackitcloud/stackit-sdk-go/services/sca/v1alphaapi"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	testRegion = "eu01"
)

var testClient = &sca.APIClient{DefaultAPI: &sca.DefaultAPIService{}}
var testCtx = context.Background()

var (
	testProjectID     = uuid.NewString()
	testEnvironmentID = uuid.NewString()
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectID,
		globalflags.RegionFlag:    testRegion,

		environmentIDFlag:      testEnvironmentID,
		nameFlag:               "test-application-name",
		imageFlag:              "test-image",
		externalPortFlag:       "8888",
		cpuFlag:                "2000",
		memoryFlag:             "2048",
		instancesFlag:          "2",
		publicFlag:             "true",
		scalingTypeFlag.Name(): scalingTypeManual,
		// minInstancesFlag:  "",
		// maxInstancesFlag:  "",
		// scaleToZeroFlag:   "",
		envVarsFlag:  "ENV1=val1,ENV2=val2",
		commandsFlag: "/bin/sh,-c",
		argsFlag:     "echo 'test'",
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
		EnvironmentID:         testEnvironmentID,
		Name:                  "test-application-name",
		Image:                 "test-image",
		ContainerExternalPort: 8888,
		ScalingType:           scalingTypeManual,
		Public:                true,
		CPU:                   2000,
		Memory:                2048,
		Instances:             2,
		EnvironmentVars: map[string]string{
			"ENV1": "val1",
			"ENV2": "val2",
		},
		Commands: []string{"/bin/sh", "-c"},
		Args:     []string{"echo 'test'"},
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *sca.ApiCreateApplicationRequest)) sca.ApiCreateApplicationRequest {
	request := testClient.DefaultAPI.CreateApplication(testCtx, testProjectID, testEnvironmentID)
	request = request.CreateApplicationPayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayload(mods ...func(payload *sca.CreateApplicationPayload)) sca.CreateApplicationPayload {
	payload := sca.CreateApplicationPayload{
		DisplayName: "test-application-name",
		Network: sca.Network{
			PublicIngress: true,
			Port:          utils.Ptr(int32(8888)),
		},
		Scaling: sca.Scaling{
			Type:          sca.SCALINGTYPE_SCALING_TYPE_MANUAL,
			ManualScaling: &sca.ManualScaling{},
		},
		Containers: []sca.Container{{
			Name:    "container-1",
			Image:   "test-image",
			Cpu:     utils.Ptr(int32(2000)),
			Memory:  utils.Ptr(int32(2048)),
			Command: []string{"/bin/sh", "-c"},
			Args:    []string{"echo 'test'"},
			EnvironmentVariables: []sca.EnvVar{
				{Key: "ENV1", Value: "val1", Origin: sca.ENVVARTYPE_ENV_FROM_SOURCE_TYPE_MANUAL.Ptr()},
				{Key: "ENV2", Value: "val2", Origin: sca.ENVVARTYPE_ENV_FROM_SOURCE_TYPE_MANUAL.Ptr()},
			},
		}},
	}
	for _, mod := range mods {
		mod(&payload)
	}
	return payload
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
			desc: "required only",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, environmentIDFlag)
				delete(flagValues, externalPortFlag)
				delete(flagValues, cpuFlag)
				delete(flagValues, memoryFlag)
				delete(flagValues, instancesFlag)
				delete(flagValues, publicFlag)
				delete(flagValues, scalingTypeFlag.String())
				delete(flagValues, envVarsFlag)
				delete(flagValues, commandsFlag)
				delete(flagValues, argsFlag)
			}),
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.CPU = 1000
				model.Memory = 1024
				model.ScalingType = scalingTypeManual
				model.Instances = 1
				model.ContainerExternalPort = 8080
				model.Public = true
				model.EnvironmentID = testProjectID
				model.EnvironmentVars = nil
				model.Commands = nil
				model.Args = nil
			}),
			isValid: true,
		},
		{
			desc: "missing name",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, nameFlag)
			}),
			isValid: false,
		},
		{
			desc: "missing image",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, imageFlag)
			}),
			isValid: false,
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
			desc: "invalid port",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[externalPortFlag] = "80"
			}),
			isValid: false,
		},
		{
			desc: "invalid scaling type",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[scalingTypeFlag.Name()] = "invalid-scaling"
			}),
			isValid: false,
		},
		{
			desc: "invalid instances",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[instancesFlag] = "-2"
			}),
			isValid: false,
		},
		{
			desc: "invalid cpu",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[cpuFlag] = "1500"
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
