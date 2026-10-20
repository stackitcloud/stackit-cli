package create

import (
	"context"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	iaas "github.com/stackitcloud/stackit-sdk-go/services/iaas/v2api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	testRegion = "eu01"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &iaas.DefaultAPIServiceMock{}

var testProjectId = uuid.NewString()
var testSourceId = uuid.NewString()
var testKeyId = uuid.NewString()
var testKeyringId = uuid.NewString()

func fixtureEncryptionFlagValues(mods ...func(map[string]string)) map[string]string {
	values := fixtureFlagValues()
	values[encryptionKekKeyIdFlag] = testKeyId
	values[encryptionKekKeyringIdFlag] = testKeyringId
	values[encryptionKekKeyVersionFlag] = "1"
	values[encryptionServiceAccountFlag] = "volume-kms@example.com"
	for _, mod := range mods {
		mod(values)
	}
	return values
}

func fixtureEncryptedInputModel(mods ...func(*inputModel)) *inputModel {
	model := fixtureInputModel()
	model.EncryptionKekKeyId = new(testKeyId)
	model.EncryptionKekKeyringId = new(testKeyringId)
	model.EncryptionKekKeyVersion = new(int64(1))
	model.EncryptionServiceAccount = new("volume-kms@example.com")
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureEncryptionParameters(mods ...func(*iaas.VolumeEncryptionParameter)) *iaas.VolumeEncryptionParameter {
	parameters := &iaas.VolumeEncryptionParameter{
		KekKeyId:       testKeyId,
		KekKeyringId:   testKeyringId,
		KekKeyVersion:  1,
		ServiceAccount: "volume-kms@example.com",
	}
	for _, mod := range mods {
		mod(parameters)
	}
	return parameters
}

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
		AvailabilityZone: "eu01-1",
		Name:             utils.Ptr("example-volume-name"),
		Description:      utils.Ptr("example-volume-description"),
		PerformanceClass: utils.Ptr("example-perf-class"),
		Size:             utils.Ptr(int64(5)),
		SourceId:         &testSourceId,
		SourceType:       utils.Ptr("example-source-type"),
		Labels: map[string]any{
			"key": "value",
		},
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
		AvailabilityZone: "eu01-1",
	})
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayload(mods ...func(payload *iaas.CreateVolumePayload)) iaas.CreateVolumePayload {
	payload := iaas.CreateVolumePayload{
		AvailabilityZone: "eu01-1",
		Name:             utils.Ptr("example-volume-name"),
		Description:      utils.Ptr("example-volume-description"),
		PerformanceClass: utils.Ptr("example-perf-class"),
		Size:             utils.Ptr(int64(5)),
		Labels: map[string]any{
			"key": "value",
		},
		Source: &iaas.VolumeSource{
			Id:   testSourceId,
			Type: "example-source-type",
		},
	}
	for _, mod := range mods {
		mod(&payload)
	}
	return payload
}

func TestParseInput(t *testing.T) {
	const secretFile = "secret.txt"
	const secret = "c2VjcmV0"
	tests := []struct {
		description   string
		argValues     []string
		stdin         string
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
			description:   "encryption with generated secret",
			flagValues:    fixtureEncryptionFlagValues(),
			isValid:       true,
			expectedModel: fixtureEncryptedInputModel(),
		},
		{
			description: "encryption with predefined secret from file",
			flagValues: fixtureEncryptionFlagValues(func(values map[string]string) {
				values[encryptionKeyPayloadFlag] = "@" + secretFile
			}),
			isValid: true,
			expectedModel: fixtureEncryptedInputModel(func(model *inputModel) {
				model.EncryptionKeyPayload = new(secret)
			}),
		},
		{
			description: "missing key id",
			flagValues: fixtureEncryptionFlagValues(func(values map[string]string) {
				delete(values, encryptionKekKeyIdFlag)
			}),
		},
		{
			description: "missing keyring id",
			flagValues: fixtureEncryptionFlagValues(func(values map[string]string) {
				delete(values, encryptionKekKeyringIdFlag)
			}),
		},
		{
			description: "missing key version",
			flagValues: fixtureEncryptionFlagValues(func(values map[string]string) {
				delete(values, encryptionKekKeyVersionFlag)
			}),
		},
		{
			description: "missing service account",
			flagValues: fixtureEncryptionFlagValues(func(values map[string]string) {
				delete(values, encryptionServiceAccountFlag)
			}),
		},
		{
			description: "invalid key id",
			flagValues: fixtureEncryptionFlagValues(func(values map[string]string) {
				values[encryptionKekKeyIdFlag] = "invalid"
			}),
		},
		{
			description: "invalid keyring id",
			flagValues: fixtureEncryptionFlagValues(func(values map[string]string) {
				values[encryptionKekKeyringIdFlag] = "invalid"
			}),
		},
		{
			description: "payload without encryption parameters",
			flagValues: fixtureFlagValues(func(values map[string]string) {
				values[encryptionKeyPayloadFlag] = "@secret.txt"
			}),
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
				model.PerformanceClass = new("example-perf-class")
				model.Size = new(int64(5))
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
				model.SourceId = &testSourceId
				model.SourceType = new("example-source-type")
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, func(params *types.CmdParams) *cobra.Command {
				params.Fs = fstest.MapFS{
					secretFile: {Data: []byte(secret)},
				}
				params.Printer.StdIn = strings.NewReader(tt.stdin)
				return NewCmd(params)
			}, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
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
			description: "encryption with generated secret",
			model:       fixtureEncryptedInputModel(),
			expectedRequest: fixtureRequest(func(request *iaas.ApiCreateVolumeRequest) {
				*request = request.CreateVolumePayload(fixturePayload(func(payload *iaas.CreateVolumePayload) {
					payload.EncryptionParameters = fixtureEncryptionParameters()
				}))
			}),
		},
		{
			description: "encryption with predefined secret",
			model: fixtureEncryptedInputModel(func(model *inputModel) {
				model.EncryptionKeyPayload = new("c2VjcmV0")
			}),
			expectedRequest: fixtureRequest(func(request *iaas.ApiCreateVolumeRequest) {
				*request = request.CreateVolumePayload(fixturePayload(func(payload *iaas.CreateVolumePayload) {
					payload.EncryptionParameters = fixtureEncryptionParameters(func(parameters *iaas.VolumeEncryptionParameter) {
						parameters.KeyPayload = new("c2VjcmV0")
					})
				}))
			}),
		},
		{
			description: "only availability zone in payload",
			model: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Verbosity: globalflags.VerbosityDefault,
					Region:    testRegion,
				},
				AvailabilityZone: "eu01-1",
			},
			expectedRequest: fixtureRequiredRequest(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx, testClient),
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
	params := testparams.NewTestParams()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.model, tt.args.projectLabel, tt.args.volume); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
