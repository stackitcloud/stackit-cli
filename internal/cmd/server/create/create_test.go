package create

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
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
var testImageId = uuid.NewString()
var testNetworkId = uuid.NewString()
var testVolumeId = uuid.NewString()

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,

		agentProvisionedFlag:              "false",
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
			Region:    testRegion,
			Verbosity: globalflags.VerbosityDefault,
		},
		AgentProvisioned:              utils.Ptr(false),
		AvailabilityZone:              utils.Ptr("eu01-1"),
		Name:                          utils.Ptr("test-server-name"),
		MachineType:                   utils.Ptr("t1.1"),
		AffinityGroup:                 utils.Ptr("test-affinity-group"),
		BootVolumePerformanceClass:    utils.Ptr("test-perf-class"),
		BootVolumeSize:                utils.Ptr(int64(5)),
		BootVolumeSourceId:            utils.Ptr(testSourceId),
		BootVolumeSourceType:          utils.Ptr("test-source-type"),
		BootVolumeDeleteOnTermination: utils.Ptr(false),
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
	request := testClient.CreateServer(testCtx, testProjectId, testRegion)
	request = request.CreateServerPayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixtureRequiredRequest(mods ...func(request *iaas.ApiCreateServerRequest)) iaas.ApiCreateServerRequest {
	request := testClient.CreateServer(testCtx, testProjectId, testRegion)
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
		MachineType: utils.Ptr("t1.1"),
		Name:        utils.Ptr("test-server-name"),
		Agent: &iaas.ServerAgent{
			Provisioned: utils.Ptr(false),
		},
		AvailabilityZone:    utils.Ptr("eu01-1"),
		AffinityGroup:       utils.Ptr("test-affinity-group"),
		KeypairName:         utils.Ptr("test-keypair-name"),
		SecurityGroups:      utils.Ptr([]string{"test-security-groups"}),
		ServiceAccountMails: utils.Ptr([]string{"test-service-account"}),
		UserData:            utils.Ptr([]byte("test-user-data")),
		Volumes:             utils.Ptr([]string{testVolumeId}),
		BootVolume: &iaas.ServerBootVolume{
			PerformanceClass:    utils.Ptr("test-perf-class"),
			Size:                utils.Ptr(int64(5)),
			DeleteOnTermination: utils.Ptr(false),
			Source: &iaas.BootVolumeSource{
				Id:   utils.Ptr(testSourceId),
				Type: utils.Ptr("test-source-type"),
			},
		},
		Networking: &iaas.CreateServerPayloadAllOfNetworking{
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
				delete(flagValues, affinityGroupFlag)
				delete(flagValues, agentProvisionedFlag)
				delete(flagValues, availabilityZoneFlag)
				delete(flagValues, labelFlag)
				delete(flagValues, bootVolumeSourceIdFlag)
				delete(flagValues, bootVolumeSourceTypeFlag)
				delete(flagValues, bootVolumeSizeFlag)
				delete(flagValues, bootVolumePerformanceClassFlag)
				delete(flagValues, bootVolumeDeleteOnTerminationFlag)
				delete(flagValues, keypairNameFlag)
				delete(flagValues, networkInterfaceIdsFlag)
				delete(flagValues, securityGroupsFlag)
				delete(flagValues, serviceAccountEmailsFlag)
				delete(flagValues, userDataFlag)
				delete(flagValues, volumesFlag)
				flagValues[imageIdFlag] = testImageId
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.AffinityGroup = nil
				model.AgentProvisioned = nil
				model.AvailabilityZone = nil
				model.Labels = nil
				model.BootVolumeSourceId = nil
				model.BootVolumeSourceType = nil
				model.BootVolumeSize = nil
				model.BootVolumePerformanceClass = nil
				model.BootVolumeDeleteOnTermination = nil
				model.KeypairName = nil
				model.NetworkInterfaceIds = nil
				model.SecurityGroups = nil
				model.ServiceAccountMails = nil
				model.UserData = nil
				model.Volumes = nil
				model.ImageId = utils.Ptr(testImageId)
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
				flagValues[imageIdFlag] = testImageId
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.BootVolumeSourceId = nil
				model.BootVolumeSourceType = nil
				model.BootVolumeSize = nil
				model.ImageId = utils.Ptr(testImageId)
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
		{
			description: "valid with agent-provisioned flag missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, agentProvisionedFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.AgentProvisioned = nil
			}),
		},
		{
			description: "agent-provisioned flag properly handled",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[agentProvisionedFlag] = "true"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.AgentProvisioned = utils.Ptr(true)
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
					Region:    testRegion,
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

func TestOutputResult(t *testing.T) {
	type args struct {
		outputFormat string
		projectLabel string
		server       *iaas.Server
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
			name: "empty with iaas server",
			args: args{
				server: &iaas.Server{},
			},
			wantErr: false,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&types.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.outputFormat, tt.args.projectLabel, tt.args.server); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
