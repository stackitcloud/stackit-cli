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
		containerNameFlag:      "my-container",
		imageFlag:              "test-image",
		externalPortFlag:       "8888",
		cpuFlag:                "2000",
		memoryFlag:             "2048",
		instancesFlag:          "2",
		publicFlag:             "true",
		scalingTypeFlag.Name(): scalingTypeManual,
		envVarsFlag:            "ENV1=val1,ENV2=val2",
		commandsFlag:           "/bin/sh,-c",
		argsFlag:               "echo 'test'",
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
		ContainerName:         "my-container",
		Image:                 "test-image",
		ContainerExternalPort: 8888,
		ScalingType:           scalingTypeManual,
		Public:                true,
		CPU:                   2000,
		Memory:                2048,
		Instances:             2,
		MinInstances:          1,
		MaxInstances:          1,
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
			Port:          new(int32(8888)),
		},
		Scaling: sca.Scaling{
			Type: sca.SCALINGTYPE_SCALING_TYPE_MANUAL,
			ManualScaling: &sca.ManualScaling{
				Instances: 2,
			},
		},
		Containers: []sca.Container{{
			Name:    "container-1",
			Image:   "test-image",
			Cpu:     new(int32(2000)),
			Memory:  new(int32(2048)),
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
		{
			desc: "invalid min instances",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[minInstancesFlag] = "-2"
			}),
			isValid: false,
		},
		{
			desc: "valid autoscaling config",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[scalingTypeFlag.Name()] = scalingTypeAuto
				delete(flagValues, instancesFlag)
				flagValues[minInstancesFlag] = "2"
				flagValues[maxInstancesFlag] = "3"
				flagValues[rpsFlag] = "5"
			}),
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.ScalingType = scalingTypeAuto
				model.Instances = defaultInstances
				model.MinInstances = 2
				model.MaxInstances = 3
				model.RPS = 5
			}),
			isValid: true,
		},
		{
			desc: "set min instances as max instances if not defined",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[scalingTypeFlag.Name()] = scalingTypeAuto
				delete(flagValues, instancesFlag)
				delete(flagValues, maxInstancesFlag)
				flagValues[minInstancesFlag] = "2"
				flagValues[rpsFlag] = "5"
			}),
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.ScalingType = scalingTypeAuto
				model.Instances = defaultInstances
				model.MinInstances = 2
				model.MaxInstances = 2
				model.RPS = 5
			}),
			isValid: true,
		},
		{
			desc: "invalid max instances",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[scalingTypeFlag.Name()] = scalingTypeAuto
				delete(flagValues, instancesFlag)
				flagValues[minInstancesFlag] = "2"
				flagValues[maxInstancesFlag] = "1"
				flagValues[rpsFlag] = "5"
			}),
			isValid: false,
		},
		{
			desc: "missing rps and concurrency if autoscaling is enabled",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[scalingTypeFlag.Name()] = scalingTypeAuto
				delete(flagValues, instancesFlag)
				flagValues[minInstancesFlag] = "2"
				flagValues[maxInstancesFlag] = "5"
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
