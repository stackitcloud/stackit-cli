package create

import (
	"context"
	"strconv"
	"strings"
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
	testRegion                       = "eu01"
	testLocalImagePath               = "/does/not/exist"
	testDiskFormat                   = "raw"
	testDiskSize               int64 = 16 * 1024 * 1024 * 1024
	testRamSize                int64 = 8 * 1024 * 1024 * 1024
	testName                         = "test-image"
	testProtected                    = true
	testCdRomBus                     = "test-cdrom"
	testDiskBus                      = "test-diskbus"
	testNicModel                     = "test-nic"
	testOperatingSystem              = "test-os"
	testOperatingSystemDistro        = "test-distro"
	testOperatingSystemVersion       = "test-distro-version"
	testRescueBus                    = "test-rescue-bus"
	testRescueDevice                 = "test-rescue-device"
	testArchitecture                 = "arm64"
	testBootmenu                     = true
	testSecureBoot                   = true
	testUefi                         = true
	testVideoModel                   = "test-video-model"
	testVirtioScsi                   = true
	testLabels                       = "foo=FOO,bar=BAR,baz=BAZ"
)

type testCtxKey struct{}

var (
	testCtx       = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient    = &iaas.APIClient{}
	testProjectId = uuid.NewString()
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,

		nameFlag:                   testName,
		diskFormatFlag:             testDiskFormat,
		localFilePathFlag:          testLocalImagePath,
		architectureFlag:           testArchitecture,
		bootMenuFlag:               strconv.FormatBool(testBootmenu),
		cdromBusFlag:               testCdRomBus,
		diskBusFlag:                testDiskBus,
		nicModelFlag:               testNicModel,
		operatingSystemFlag:        testOperatingSystem,
		operatingSystemDistroFlag:  testOperatingSystemDistro,
		operatingSystemVersionFlag: testOperatingSystemVersion,
		rescueBusFlag:              testRescueBus,
		rescueDeviceFlag:           testRescueDevice,
		secureBootFlag:             strconv.FormatBool(testSecureBoot),
		uefiFlag:                   strconv.FormatBool(testUefi),
		videoModelFlag:             testVideoModel,
		virtioScsiFlag:             strconv.FormatBool(testVirtioScsi),
		labelsFlag:                 testLabels,
		minDiskSizeFlag:            strconv.Itoa(int(testDiskSize)),
		minRamFlag:                 strconv.Itoa(int(testRamSize)),
		protectedFlag:              strconv.FormatBool(testProtected),
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func parseLabels(labelstring string) map[string]string {
	labels := map[string]string{}
	for _, part := range strings.Split(labelstring, ",") {
		v := strings.Split(part, "=")
		labels[v[0]] = v[1]
	}

	return labels
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Verbosity: globalflags.VerbosityDefault,
			Region:    testRegion,
		},
		Name:          testName,
		DiskFormat:    testDiskFormat,
		LocalFilePath: testLocalImagePath,
		Labels:        utils.Ptr(parseLabels(testLabels)),
		Config: &imageConfig{
			Architecture:           utils.Ptr(testArchitecture),
			BootMenu:               utils.Ptr(testBootmenu),
			CdromBus:               utils.Ptr(testCdRomBus),
			DiskBus:                utils.Ptr(testDiskBus),
			NicModel:               utils.Ptr(testNicModel),
			OperatingSystem:        utils.Ptr(testOperatingSystem),
			OperatingSystemDistro:  utils.Ptr(testOperatingSystemDistro),
			OperatingSystemVersion: utils.Ptr(testOperatingSystemVersion),
			RescueBus:              utils.Ptr(testRescueBus),
			RescueDevice:           utils.Ptr(testRescueDevice),
			SecureBoot:             utils.Ptr(testSecureBoot),
			Uefi:                   testUefi,
			VideoModel:             utils.Ptr(testVideoModel),
			VirtioScsi:             utils.Ptr(testVirtioScsi),
		},
		MinDiskSize: utils.Ptr(testDiskSize),
		MinRam:      utils.Ptr(testRamSize),
		Protected:   utils.Ptr(testProtected),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureCreatePayload(mods ...func(payload *iaas.CreateImagePayload)) (payload iaas.CreateImagePayload) {
	payload = iaas.CreateImagePayload{
		Config: &iaas.ImageConfig{
			Architecture:           utils.Ptr(testArchitecture),
			BootMenu:               utils.Ptr(testBootmenu),
			CdromBus:               iaas.NewNullableString(utils.Ptr(testCdRomBus)),
			DiskBus:                iaas.NewNullableString(utils.Ptr(testDiskBus)),
			NicModel:               iaas.NewNullableString(utils.Ptr(testNicModel)),
			OperatingSystem:        utils.Ptr(testOperatingSystem),
			OperatingSystemDistro:  iaas.NewNullableString(utils.Ptr(testOperatingSystemDistro)),
			OperatingSystemVersion: iaas.NewNullableString(utils.Ptr(testOperatingSystemVersion)),
			RescueBus:              iaas.NewNullableString(utils.Ptr(testRescueBus)),
			RescueDevice:           iaas.NewNullableString(utils.Ptr(testRescueDevice)),
			SecureBoot:             utils.Ptr(testSecureBoot),
			Uefi:                   utils.Ptr(testUefi),
			VideoModel:             iaas.NewNullableString(utils.Ptr(testVideoModel)),
			VirtioScsi:             utils.Ptr(testVirtioScsi),
		},
		DiskFormat: utils.Ptr(testDiskFormat),
		Labels: &map[string]interface{}{
			"foo": "FOO",
			"bar": "BAR",
			"baz": "BAZ",
		},
		MinDiskSize: utils.Ptr(testDiskSize),
		MinRam:      utils.Ptr(testRamSize),
		Name:        utils.Ptr(testName),
		Protected:   utils.Ptr(testProtected),
	}
	for _, mod := range mods {
		mod(&payload)
	}
	return payload
}

func fixtureRequest(mods ...func(request *iaas.ApiCreateImageRequest)) iaas.ApiCreateImageRequest {
	request := testClient.CreateImage(testCtx, testProjectId, testRegion)

	request = request.CreateImagePayload(fixtureCreatePayload())

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
			description: "name missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, nameFlag)
			}),
			isValid: false,
		},
		{
			description: "no labels",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, labelsFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Labels = nil
			}),
		},
		{
			description: "single label",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[labelsFlag] = "foo=bar"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Labels = &map[string]string{
					"foo": "bar",
				}
			}),
		},
		{
			description: "only rescue bus is invalid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, rescueDeviceFlag)
			}),
			isValid: false,
		},
		{
			description: "only rescue device is invalid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, rescueBusFlag)
			}),
			isValid: false,
		},
		{
			description: "uefi flag is set to false",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[uefiFlag] = strconv.FormatBool(false)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Config.Uefi = false
			}),
		},
		{
			description: "no rescue device and no bus is valid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, rescueBusFlag)
				delete(flagValues, rescueDeviceFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Config.RescueBus = nil
				model.Config.RescueDevice = nil
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
		expectedRequest iaas.ApiCreateImageRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "no labels",
			model: fixtureInputModel(func(model *inputModel) {
				model.Labels = nil
			}),
			expectedRequest: fixtureRequest(func(request *iaas.ApiCreateImageRequest) {
				*request = (*request).CreateImagePayload(fixtureCreatePayload(func(payload *iaas.CreateImagePayload) {
					payload.Labels = nil
				}))
			}),
		},
		{
			description: "cd rom bus",
			model: fixtureInputModel(func(model *inputModel) {
				model.Config.CdromBus = utils.Ptr("foobar")
			}),
			expectedRequest: fixtureRequest(func(request *iaas.ApiCreateImageRequest) {
				*request = (*request).CreateImagePayload(fixtureCreatePayload(func(payload *iaas.CreateImagePayload) {
					payload.Config.CdromBus = iaas.NewNullableString(utils.Ptr("foobar"))
				}))
			}),
		},
		{
			description: "uefi flag",
			model: fixtureInputModel(func(model *inputModel) {
				model.Config.Uefi = false
			}),
			expectedRequest: fixtureRequest(func(request *iaas.ApiCreateImageRequest) {
				*request = (*request).CreateImagePayload(fixtureCreatePayload(func(payload *iaas.CreateImagePayload) {
					payload.Config.Uefi = utils.Ptr(false)
				}))
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)
			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx),
				cmp.AllowUnexported(iaas.NullableString{}),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	type args struct {
		model *inputModel
		resp  *iaas.ImageCreateResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil",
			args: args{
				model: nil,
				resp:  nil,
			},
			wantErr: true,
		},
		{
			name: "empty input",
			args: args{
				model: &inputModel{},
				resp:  &iaas.ImageCreateResponse{},
			},
			wantErr: false,
		},
		{
			name: "output json",
			args: args{
				model: &inputModel{
					GlobalFlagModel: &globalflags.GlobalFlagModel{
						OutputFormat: print.JSONOutputFormat,
					},
				},
				resp: nil,
			},
			wantErr: false,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&types.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.model, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
