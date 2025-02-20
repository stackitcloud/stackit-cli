package create

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

var projectIdFlag = globalflags.ProjectIdFlag

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &iaas.APIClient{}

var testProjectId = uuid.NewString()
var testSourceId = uuid.NewString()

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:        testProjectId,
		availabilityZoneFlag: "eu01-1",
		nameFlag:             "example-volume-name",
		descriptionFlag:      "example-volume-description",
		labelFlag:            "key=value",
		performanceClassFlag: "example-perf-class",
		sizeFlag:             "5",
		sourceIdFlag:         testSourceId,
		sourceTypeFlag:       "example-source-type",
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Verbosity: globalflags.VerbosityDefault,
		},
		AvailabilityZone: utils.Ptr("eu01-1"),
		Name:             utils.Ptr("example-volume-name"),
		Description:      utils.Ptr("example-volume-description"),
		PerformanceClass: utils.Ptr("example-perf-class"),
		Size:             utils.Ptr(int64(5)),
		SourceId:         utils.Ptr(testSourceId),
		SourceType:       utils.Ptr("example-source-type"),
		Labels: utils.Ptr(map[string]string{
			"key": "value",
		}),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *iaas.ApiCreateVolumeRequest)) iaas.ApiCreateVolumeRequest {
	request := testClient.CreateVolume(testCtx, testProjectId)
	request = request.CreateVolumePayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixtureRequiredRequest(mods ...func(request *iaas.ApiCreateVolumeRequest)) iaas.ApiCreateVolumeRequest {
	request := testClient.CreateVolume(testCtx, testProjectId)
	request = request.CreateVolumePayload(iaas.CreateVolumePayload{
		AvailabilityZone: utils.Ptr("eu01-1"),
	})
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayload(mods ...func(payload *iaas.CreateVolumePayload)) iaas.CreateVolumePayload {
	payload := iaas.CreateVolumePayload{
		AvailabilityZone: utils.Ptr("eu01-1"),
		Name:             utils.Ptr("example-volume-name"),
		Description:      utils.Ptr("example-volume-description"),
		PerformanceClass: utils.Ptr("example-perf-class"),
		Size:             utils.Ptr(int64(5)),
		Labels: utils.Ptr(map[string]interface{}{
			"key": "value",
		}),
		Source: &iaas.VolumeSource{
			Id:   utils.Ptr(testSourceId),
			Type: utils.Ptr("example-source-type"),
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
			description: "required only",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, nameFlag)
				delete(flagValues, descriptionFlag)
				delete(flagValues, labelFlag)
				delete(flagValues, performanceClassFlag)
				delete(flagValues, sizeFlag)
				delete(flagValues, sourceIdFlag)
				delete(flagValues, sourceTypeFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Name = nil
				model.Description = nil
				model.Labels = nil
				model.PerformanceClass = nil
				model.Size = nil
				model.SourceType = nil
				model.SourceId = nil
			}),
		},
		{
			description: "availability zone missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, availabilityZoneFlag)
			}),
			isValid: false,
		},
		{
			description: "no values",
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "project id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, projectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "use performance class and size",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[performanceClassFlag] = "example-perf-class"
				flagValues[sizeFlag] = "5"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.PerformanceClass = utils.Ptr("example-perf-class")
				model.Size = utils.Ptr(int64(5))
			}),
		},
		{
			description: "use source id and type",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[sourceIdFlag] = testSourceId
				flagValues[sourceTypeFlag] = "example-source-type"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.SourceId = utils.Ptr(testSourceId)
				model.SourceType = utils.Ptr("example-source-type")
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			p := print.NewPrinter()
			cmd := NewCmd(p)
			err := globalflags.Configure(cmd.Flags())
			if err != nil {
				t.Fatalf("configure global flags: %v", err)
			}

			for flag, value := range tt.flagValues {
				err := cmd.Flags().Set(flag, value)
				if err != nil {
					if !tt.isValid {
						return
					}
					t.Fatalf("setting flag --%s=%s: %v", flag, value, err)
				}
			}

			err = cmd.ValidateRequiredFlags()
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flags: %v", err)
			}

			model, err := parseInput(p, cmd)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error parsing flags: %v", err)
			}

			if !tt.isValid {
				t.Fatalf("did not fail on invalid input")
			}
			diff := cmp.Diff(model, tt.expectedModel)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest iaas.ApiCreateVolumeRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "only availability zone in payload",
			model: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Verbosity: globalflags.VerbosityDefault,
				},
				AvailabilityZone: utils.Ptr("eu01-1"),
			},
			expectedRequest: fixtureRequiredRequest(),
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
