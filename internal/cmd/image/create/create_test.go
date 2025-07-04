package create

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

var projectIdFlag = globalflags.ProjectIdFlag

type testCtxKey struct{}

var (
	testCtx       = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient    = &iaas.APIClient{}
	testProjectId = uuid.NewString()

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

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag: testProjectId,

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
		GlobalFlagModel: &globalflags.GlobalFlagModel{ProjectId: testProjectId, Verbosity: globalflags.VerbosityDefault},
		Name:            testName,
		DiskFormat:      testDiskFormat,
		LocalFilePath:   testLocalImagePath,
		Labels:          utils.Ptr(parseLabels(testLabels)),
		Config: &imageConfig{
			Architecture:           &testArchitecture,
			BootMenu:               &testBootmenu,
			CdromBus:               &testCdRomBus,
			DiskBus:                &testDiskBus,
			NicModel:               &testNicModel,
			OperatingSystem:        &testOperatingSystem,
			OperatingSystemDistro:  &testOperatingSystemDistro,
			OperatingSystemVersion: &testOperatingSystemVersion,
			RescueBus:              &testRescueBus,
			RescueDevice:           &testRescueDevice,
			SecureBoot:             &testSecureBoot,
			Uefi:                   testUefi,
			VideoModel:             &testVideoModel,
			VirtioScsi:             &testVirtioScsi,
		},
		MinDiskSize: &testDiskSize,
		MinRam:      &testRamSize,
		Protected:   &testProtected,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureCreatePayload(mods ...func(payload *iaas.CreateImagePayload)) (payload iaas.CreateImagePayload) {
	payload = iaas.CreateImagePayload{
		Config: &iaas.ImageConfig{
			Architecture:           &testArchitecture,
			BootMenu:               &testBootmenu,
			CdromBus:               iaas.NewNullableString(&testCdRomBus),
			DiskBus:                iaas.NewNullableString(&testDiskBus),
			NicModel:               iaas.NewNullableString(&testNicModel),
			OperatingSystem:        &testOperatingSystem,
			OperatingSystemDistro:  iaas.NewNullableString(&testOperatingSystemDistro),
			OperatingSystemVersion: iaas.NewNullableString(&testOperatingSystemVersion),
			RescueBus:              iaas.NewNullableString(&testRescueBus),
			RescueDevice:           iaas.NewNullableString(&testRescueDevice),
			SecureBoot:             &testSecureBoot,
			Uefi:                   &testUefi,
			VideoModel:             iaas.NewNullableString(&testVideoModel),
			VirtioScsi:             &testVirtioScsi,
		},
		DiskFormat: &testDiskFormat,
		Labels: &map[string]interface{}{
			"foo": "FOO",
			"bar": "BAR",
			"baz": "BAZ",
		},
		MinDiskSize: &testDiskSize,
		MinRam:      &testRamSize,
		Name:        &testName,
		Protected:   &testProtected,
	}
	for _, mod := range mods {
		mod(&payload)
	}
	return payload
}

func fixtureRequest(mods ...func(request *iaas.ApiCreateImageRequest)) iaas.ApiCreateImageRequest {
	request := testClient.CreateImage(testCtx, testProjectId)

	request = request.CreateImagePayload(fixtureCreatePayload())

	for _, mod := range mods {
		mod(&request)
	}
	return request
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
			p := print.NewPrinter()
			cmd := NewCmd(&params.CmdParams{Printer: p})
			if err := globalflags.Configure(cmd.Flags()); err != nil {
				t.Errorf("cannot configure global flags: %v", err)
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

			if err := cmd.ValidateFlagGroups(); err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flag groups: %v", err)
			}

			if err := cmd.ValidateRequiredFlags(); err != nil {
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
	p.Cmd = NewCmd(&params.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.model, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
