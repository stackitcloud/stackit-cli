package create

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	testRegion = "eu01"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &iaas.APIClient{}

var testProjectId = uuid.NewString()
var testSourceId = uuid.NewString()

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,

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
			Region:    testRegion,
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
	request := testClient.CreateVolume(testCtx, testProjectId, testRegion)
	request = request.CreateVolumePayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixtureRequiredRequest(mods ...func(request *iaas.ApiCreateVolumeRequest)) iaas.ApiCreateVolumeRequest {
	request := testClient.CreateVolume(testCtx, testProjectId, testRegion)
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
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = "invalid-uuid"
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
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
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
					Region:    testRegion,
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

func TestOutputResult(t *testing.T) {
	type args struct {
		model        *inputModel
		projectLabel string
		volume       *iaas.Volume
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
			name: "volume as argument",
			args: args{
				model:  fixtureInputModel(),
				volume: &iaas.Volume{},
			},
			wantErr: false,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&types.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.model, tt.args.projectLabel, tt.args.volume); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
