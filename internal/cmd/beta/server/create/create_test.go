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
var testImageId = uuid.NewString()
var testNetworkId = uuid.NewString()
var testVolumeId = uuid.NewString()

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:                     testProjectId,
		availabilityZoneFlag:              "eu01-1",
		nameFlag:                          "test-server-name",
		machineTypeFlag:                   "t1.1",
		affinityGroupFlag:                 "test-affinity-group",
		labelFlag:                         "key=value",
		bootVolumePerformanceClassFlag:    "test-perf-class",
		bootVolumeSizeFlag:                "5",
		bootVolumeSourceIdFlag:            testSourceId,
		bootVolumeSourceTypeFlag:          "test-source-type",
		bootVolumeDeleteOnTerminationFlag: "false",
		imageIdFlag:                       testImageId,
		keypairNameFlag:                   "test-keypair-name",
		networkIdFlag:                     testNetworkId,
		securityGroupsFlag:                "test-security-groups",
		serviceAccountEmailsFlag:          "test-service-account",
		userDataFlag:                      "test-user-data",
		volumesFlag:                       testVolumeId,
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
		AvailabilityZone:              utils.Ptr("eu01-1"),
		Name:                          utils.Ptr("test-server-name"),
		MachineType:                   utils.Ptr("t1.1"),
		AffinityGroup:                 utils.Ptr("test-affinity-group"),
		BootVolumePerformanceClass:    utils.Ptr("test-perf-class"),
		BootVolumeSize:                utils.Ptr(int64(5)),
		BootVolumeSourceId:            utils.Ptr(testSourceId),
		BootVolumeSourceType:          utils.Ptr("test-source-type"),
		BootVolumeDeleteOnTermination: utils.Ptr(false),
		ImageId:                       utils.Ptr(testImageId),
		KeypairName:                   utils.Ptr("test-keypair-name"),
		NetworkId:                     utils.Ptr(testNetworkId),
		SecurityGroups:                utils.Ptr([]string{"test-security-groups"}),
		ServiceAccountMails:           utils.Ptr([]string{"test-service-account"}),
		UserData:                      utils.Ptr("test-user-data"),
		Volumes:                       utils.Ptr([]string{testVolumeId}),
		Labels: utils.Ptr(map[string]string{
			"key": "value",
		}),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *iaas.ApiCreateServerRequest)) iaas.ApiCreateServerRequest {
	request := testClient.CreateServer(testCtx, testProjectId)
	request = request.CreateServerPayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixtureRequiredRequest(mods ...func(request *iaas.ApiCreateServerRequest)) iaas.ApiCreateServerRequest {
	request := testClient.CreateServer(testCtx, testProjectId)
	request = request.CreateServerPayload(iaas.CreateServerPayload{
		MachineType: utils.Ptr("t1.1"),
		Name:        utils.Ptr("test-server-name"),
	})
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayload(mods ...func(payload *iaas.CreateServerPayload)) iaas.CreateServerPayload {
	payload := iaas.CreateServerPayload{
		Labels: utils.Ptr(map[string]interface{}{
			"key": "value",
		}),
		MachineType:         utils.Ptr("t1.1"),
		Name:                utils.Ptr("test-server-name"),
		AvailabilityZone:    utils.Ptr("eu01-1"),
		AffinityGroup:       utils.Ptr("test-affinity-group"),
		ImageId:             utils.Ptr(testImageId),
		KeypairName:         utils.Ptr("test-keypair-name"),
		SecurityGroups:      utils.Ptr([]string{"test-security-groups"}),
		ServiceAccountMails: utils.Ptr([]string{"test-service-account"}),
		UserData:            utils.Ptr("test-user-data"),
		Volumes:             utils.Ptr([]string{testVolumeId}),
		BootVolume: &iaas.CreateServerPayloadBootVolume{
			PerformanceClass:    utils.Ptr("test-perf-class"),
			Size:                utils.Ptr(int64(5)),
			DeleteOnTermination: utils.Ptr(false),
			Source: &iaas.BootVolumeSource{
				Id:   utils.Ptr(testSourceId),
				Type: utils.Ptr("test-source-type"),
			},
		},
		Networking: &iaas.CreateServerPayloadNetworking{
			CreateServerNetworking: &iaas.CreateServerNetworking{
				NetworkId: utils.Ptr(testNetworkId),
			},
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
				delete(flagValues, affinityGroupFlag)
				delete(flagValues, availabilityZoneFlag)
				delete(flagValues, labelFlag)
				delete(flagValues, bootVolumeSourceIdFlag)
				delete(flagValues, bootVolumeSourceTypeFlag)
				delete(flagValues, bootVolumeSizeFlag)
				delete(flagValues, bootVolumePerformanceClassFlag)
				delete(flagValues, bootVolumeDeleteOnTerminationFlag)
				delete(flagValues, keypairNameFlag)
				delete(flagValues, networkIdFlag)
				delete(flagValues, networkInterfaceIdsFlag)
				delete(flagValues, securityGroupsFlag)
				delete(flagValues, serviceAccountEmailsFlag)
				delete(flagValues, userDataFlag)
				delete(flagValues, volumesFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.AffinityGroup = nil
				model.AvailabilityZone = nil
				model.Labels = nil
				model.BootVolumeSourceId = nil
				model.BootVolumeSourceType = nil
				model.BootVolumeSize = nil
				model.BootVolumePerformanceClass = nil
				model.BootVolumeDeleteOnTermination = nil
				model.KeypairName = nil
				model.NetworkId = nil
				model.NetworkInterfaceIds = nil
				model.SecurityGroups = nil
				model.ServiceAccountMails = nil
				model.UserData = nil
				model.Volumes = nil
			}),
		},
		{
			description: "machine type missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, machineTypeFlag)
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
			description: "use network id",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[networkIdFlag] = testNetworkId
				flagValues[nameFlag] = "test-server-name"
				flagValues[machineTypeFlag] = "t1.1"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.NetworkId = utils.Ptr(testNetworkId)
				model.Name = utils.Ptr("test-server-name")
				model.MachineType = utils.Ptr("t1.1")
			}),
		},
		{
			description: "use boot volume source id and type",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[bootVolumeSourceIdFlag] = testImageId
				flagValues[bootVolumeSourceTypeFlag] = "image"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.BootVolumeSourceId = utils.Ptr(testImageId)
				model.BootVolumeSourceType = utils.Ptr("image")
			}),
		},
		{
			description: "invalid without image-id, boot-volume-source-id and type",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, bootVolumeSourceIdFlag)
				delete(flagValues, bootVolumeSourceTypeFlag)
				delete(flagValues, imageIdFlag)
			}),
			isValid: false,
		},
		{
			description: "invalid with boot-volume-source-id and without type",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, bootVolumeSourceTypeFlag)
			}),
			isValid: false,
		},
		{
			description: "invalid with boot-volume-source-type is image and without size",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, bootVolumeSizeFlag)
				flagValues[bootVolumeSourceIdFlag] = testImageId
				flagValues[bootVolumeSourceTypeFlag] = "image"
			}),
			isValid: false,
		},
		{
			description: "valid with image-id",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, bootVolumeSourceIdFlag)
				delete(flagValues, bootVolumeSourceTypeFlag)
				delete(flagValues, bootVolumeSizeFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.BootVolumeSourceId = nil
				model.BootVolumeSourceType = nil
				model.BootVolumeSize = nil
			}),
		},
		{
			description: "valid with boot-volume-source-id and type volume",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, imageIdFlag)
				delete(flagValues, bootVolumeSizeFlag)

			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.ImageId = nil
				model.BootVolumeSize = nil
			}),
		},
		{
			description: "valid with boot-volume-source-id, type volume and size",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, imageIdFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.ImageId = nil
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
		expectedRequest iaas.ApiCreateServerRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "only name and machine type in payload",
			model: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Verbosity: globalflags.VerbosityDefault,
				},
				MachineType: utils.Ptr("t1.1"),
				Name:        utils.Ptr("test-server-name"),
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
