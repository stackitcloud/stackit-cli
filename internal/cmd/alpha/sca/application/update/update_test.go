package update

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

		environmentIDFlag:      testEnvironmentID,
		imageFlag:              "test-image",
		publicFlag:             "true",
		externalPortFlag:       "8888",
		cpuFlag:                "2000",
		memoryFlag:             "2048",
		instancesFlag:          "2",
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
		ApplicationID:         testApplicationID,
		EnvironmentID:         testEnvironmentID,
		Image:                 "test-image",
		ContainerExternalPort: new(int32(8888)),
		ScalingType:           new(scalingTypeManual),
		Public:                new(true),
		CPU:                   new(int32(2000)),
		Memory:                new(int32(2048)),
		Instances:             new(int32(2)),
		EnvironmentVariables: map[string]string{
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

func fixtureCurrentApplication(mods ...func(application *sca.Application)) *sca.Application {
	application := &sca.Application{
		DisplayName:   "test-name",
		EnvironmentId: &testEnvironmentID,
		Containers: []sca.Container{{
			Image:  "container-1",
			Memory: sca.PtrInt32(1024),
			Cpu:    sca.PtrInt32(1000),
		}},
		Scaling: sca.Scaling{
			Type: sca.SCALINGTYPE_SCALING_TYPE_MANUAL,
			ManualScaling: &sca.ManualScaling{
				Instances: 1,
			},
		},
		Network: sca.Network{
			PublicIngress: true,
			Port:          sca.PtrInt32(8080),
		},
	}
	for _, mod := range mods {
		mod(application)
	}
	return application
}

func fixtureRequest(mods ...func(request *sca.ApiUpdateApplicationRequest)) sca.ApiUpdateApplicationRequest {
	request := testClient.DefaultAPI.UpdateApplication(testCtx, testProjectID, testEnvironmentID, testApplicationID)
	request = request.UpdateApplicationPayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayload(mods ...func(payload *sca.UpdateApplicationPayload)) sca.UpdateApplicationPayload {
	payload := sca.UpdateApplicationPayload{
		Network: &sca.Network{
			PublicIngress: true,
			Port:          new(int32(8888)),
		},
		Scaling: &sca.Scaling{
			Type: sca.SCALINGTYPE_SCALING_TYPE_MANUAL,
			ManualScaling: &sca.ManualScaling{
				Instances: 2,
			},
		},
		Containers: []sca.Container{{
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
		argValues     []string
		flagValues    map[string]string
		expectedModel *inputModel
		isValid       bool
	}{
		{
			desc:          "base",
			argValues:     fixtureArgsValues(),
			flagValues:    fixtureFlagValues(),
			expectedModel: fixtureInputModel(),
			isValid:       true,
		},
		{
			desc:       "no arg values",
			argValues:  nil,
			flagValues: fixtureFlagValues(),
			isValid:    false,
		},
		{
			desc:       "no flag values",
			argValues:  fixtureArgsValues(),
			flagValues: nil,
			isValid:    false,
		},
		{
			desc:      "project id missing",
			argValues: fixtureArgsValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			desc:      "project id invalid",
			argValues: fixtureArgsValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			desc:      "invalid port",
			argValues: fixtureArgsValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[externalPortFlag] = "80"
			}),
			isValid: false,
		},
		{
			desc:      "invalid scaling type",
			argValues: fixtureArgsValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[scalingTypeFlag.Name()] = "invalid-scaling"
			}),
			isValid: false,
		},
		{
			desc:      "invalid instances",
			argValues: fixtureArgsValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[instancesFlag] = "-2"
			}),
			isValid: false,
		},
		{
			desc:      "invalid cpu",
			argValues: fixtureArgsValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[cpuFlag] = "1500"
			}),
			isValid: false,
		},
		{
			desc:      "invalid min instances",
			argValues: fixtureArgsValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[minInstancesFlag] = "-2"
			}),
			isValid: false,
		},
		{
			desc:      "invalid max instances",
			argValues: fixtureArgsValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[scalingTypeFlag.Name()] = scalingTypeAuto
				delete(flagValues, instancesFlag)
				flagValues[maxInstancesFlag] = "12"
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		current         *sca.Application
		expectedRequest sca.ApiUpdateApplicationRequest
		invalidRequest  bool
	}{
		{
			description:     "base",
			current:         fixtureCurrentApplication(),
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "update only image",
			current:     fixtureCurrentApplication(),
			model: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectID,
				},
				ApplicationID: testApplicationID,
				EnvironmentID: testEnvironmentID,
				Image:         "new-image",
			},
			expectedRequest: fixtureRequest(func(request *sca.ApiUpdateApplicationRequest) {
				*request = request.UpdateApplicationPayload(fixturePayload(func(payload *sca.UpdateApplicationPayload) {
					app := fixtureCurrentApplication()
					baseContainers := app.Containers
					baseContainers[0].Image = "new-image"
					payload.Containers = baseContainers
					payload.Network = nil
					payload.Scaling = nil
				}))

			}),
		},
		{
			description:    "missing required flag if replacing scaling type to auto",
			invalidRequest: true,
			current:        fixtureCurrentApplication(),
			model: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectID,
				},
				ApplicationID: testApplicationID,
				EnvironmentID: testEnvironmentID,
				ScalingType:   new(scalingTypeAuto),
			},
		},
		{
			description:    "missing required flag if replacing scaling type to manual",
			invalidRequest: true,
			current: fixtureCurrentApplication(func(application *sca.Application) {
				application.Scaling.Type = sca.SCALINGTYPE_SCALING_TYPE_AUTO
			}),
			model: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectID,
				},
				ApplicationID: testApplicationID,
				EnvironmentID: testEnvironmentID,
				ScalingType:   new(scalingTypeManual),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request, err := buildRequest(testCtx, tt.model, testClient, tt.current)
			if err != nil {
				if tt.invalidRequest {
					return
				}
				t.Fatalf("error building request: %v", err)
			}

			diff := cmp.Diff(*request, tt.expectedRequest,
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
