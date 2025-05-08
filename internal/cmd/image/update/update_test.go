package update

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
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

var (
	testCtx       = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient    = &iaas.APIClient{}
	testProjectId = uuid.NewString()

	testImageId                      = []string{uuid.NewString()}
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
		Id:              testImageId[0],
		Name:            &testName,
		DiskFormat:      &testDiskFormat,
		Labels:          utils.Ptr(parseLabels(testLabels)),
		Config: &imageConfig{
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
			Uefi:                   &testUefi,
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

func fixtureCreatePayload(mods ...func(payload *iaas.UpdateImagePayload)) (payload iaas.UpdateImagePayload) {
	payload = iaas.UpdateImagePayload{
		Config: &iaas.ImageConfig{
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

func fixtureRequest(mods ...func(*iaas.ApiUpdateImageRequest)) iaas.ApiUpdateImageRequest {
	request := testClient.UpdateImage(testCtx, testProjectId, testImageId[0])

	request = request.UpdateImagePayload(fixtureCreatePayload())

	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description   string
		flagValues    map[string]string
		args          []string
		isValid       bool
		expectedModel *inputModel
	}{
		{
			description:   "base",
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			args:          testImageId,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "no values but valid image id",
			flagValues: map[string]string{
				projectIdFlag: testProjectId,
			},
			args:    testImageId,
			isValid: false,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Labels = nil
				model.Name = nil
			}),
		},
		{
			description: "project id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, projectIdFlag)
			}),
			args:    testImageId,
			isValid: false,
		},
		{
			description: "project id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = ""
			}),
			args:    testImageId,
			isValid: false,
		},
		{
			description: "project id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = "invalid-uuid"
			}),
			args:    testImageId,
			isValid: false,
		},
		{
			description: "no name passed",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, nameFlag)
			}),
			args: testImageId,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Name = nil
			}),
			isValid: true,
		},
		{
			description: "no labels",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, labelsFlag)
			}),
			args: testImageId,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Labels = nil
			}),
			isValid: true,
		},
		{
			description: "single label",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[labelsFlag] = "foo=bar"
			}),
			args:    testImageId,
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Labels = &map[string]string{
					"foo": "bar",
				}
			}),
		},
		{
			description: "no image id passed",
			flagValues:  fixtureFlagValues(),
			args:        nil,
			isValid:     false,
		},
		{
			description: "invalid image id passed",
			flagValues:  fixtureFlagValues(),
			args:        []string{"foobar"},
			isValid:     false,
		},
		{
			description: "multiple image ids passed",
			flagValues:  fixtureFlagValues(),
			args:        []string{uuid.NewString(), uuid.NewString()},
			isValid:     false,
		},
		{
			description: "only rescue bus is invalid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, rescueDeviceFlag)
			}),
			args:    []string{testImageId[0]},
			isValid: false,
		},
		{
			description: "only rescue device is invalid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, rescueBusFlag)
			}),
			args:    []string{testImageId[0]},
			isValid: false,
		},
		{
			description: "no rescue device and no bus is valid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, rescueBusFlag)
				delete(flagValues, rescueDeviceFlag)
			}),
			isValid: true,
			args:    []string{testImageId[0]},
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
				if err := cmd.Flags().Set(flag, value); err != nil {
					if !tt.isValid {
						return
					}
					t.Fatalf("setting flag --%s=%s: %v", flag, value, err)
				}
			}

			if err := cmd.ValidateRequiredFlags(); err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flags: %v", err)
			}

			if err := cmd.ValidateFlagGroups(); err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flag groups: %v", err)
			}

			if err := cmd.ValidateArgs(tt.args); err != nil {
				if !tt.isValid {
					return
				}
			}

			model, err := parseInput(p, cmd, tt.args)
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
		expectedRequest iaas.ApiUpdateImageRequest
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
			expectedRequest: fixtureRequest(func(request *iaas.ApiUpdateImageRequest) {
				*request = request.UpdateImagePayload(fixtureCreatePayload(func(payload *iaas.UpdateImagePayload) {
					payload.Labels = nil
				}))
			}),
		},
		{
			description: "change name",
			model: fixtureInputModel(func(model *inputModel) {
				model.Name = utils.Ptr("something else")
			}),
			expectedRequest: fixtureRequest(func(request *iaas.ApiUpdateImageRequest) {
				*request = request.UpdateImagePayload(fixtureCreatePayload(func(payload *iaas.UpdateImagePayload) {
					payload.Name = utils.Ptr("something else")
				}))
			}),
		},
		{
			description: "change cdrom",
			model: fixtureInputModel(func(model *inputModel) {
				model.Config.CdromBus = utils.Ptr("something else")
			}),
			expectedRequest: fixtureRequest(func(request *iaas.ApiUpdateImageRequest) {
				*request = request.UpdateImagePayload(fixtureCreatePayload(func(payload *iaas.UpdateImagePayload) {
					payload.Config.CdromBus.Set(utils.Ptr("something else"))
				}))
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)
			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest, iaas.NullableString{}),
				cmpopts.EquateComparable(testCtx),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}
