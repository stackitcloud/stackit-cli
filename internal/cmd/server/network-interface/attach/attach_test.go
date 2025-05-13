package attach

import (
	"context"
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

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &iaas.APIClient{}
var testProjectId = uuid.NewString()
var testServerId = uuid.NewString()
var testNicId = uuid.NewString()
var testNetworkId = uuid.NewString()

// contains nic id
func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:          testProjectId,
		serverIdFlag:           testServerId,
		networkInterfaceIdFlag: testNicId,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			Verbosity: globalflags.VerbosityDefault,
			ProjectId: testProjectId,
		},
		ServerId: utils.Ptr(testServerId),
		NicId:    utils.Ptr(testNicId),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequestAttach(mods ...func(request *iaas.ApiAddNicToServerRequest)) iaas.ApiAddNicToServerRequest {
	request := testClient.AddNicToServer(testCtx, testProjectId, testServerId, testNicId)
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixtureRequestCreateAndAttach(mods ...func(request *iaas.ApiAddNetworkToServerRequest)) iaas.ApiAddNetworkToServerRequest {
	request := testClient.AddNetworkToServer(testCtx, testProjectId, testServerId, testNetworkId)
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
			description: "server id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, serverIdFlag)
			}),
			isValid: false,
		},
		{
			description: "server id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[serverIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "server id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[serverIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		// only create
		{
			description: "provided flags invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[createFlag] = "true"
				delete(flagValues, networkInterfaceIdFlag)
			}),
			isValid: false,
		},
		// only network id
		{
			description: "provided flags invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, networkInterfaceIdFlag)
				flagValues[networkIdFlag] = testNetworkId
			}),
			isValid: false,
		},
		// create and nic id
		{
			description: "provided flags invalid 3",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[createFlag] = "true"
			}),
			isValid: false,
		},
		// create and network id (valid)
		{
			description: "provided flags valid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[createFlag] = "true"
				delete(flagValues, networkInterfaceIdFlag)
				flagValues[networkIdFlag] = testNetworkId
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Create = utils.Ptr(true)
				model.NicId = nil
				model.NetworkId = utils.Ptr(testNetworkId)
			}),
		},
		// create, nic id and network id
		{
			description: "provided flags invalid 4",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[createFlag] = "true"
				flagValues[networkIdFlag] = testNetworkId
			}),
			isValid: false,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Create = utils.Ptr(true)
				model.NetworkId = utils.Ptr(testNetworkId)
			}),
		},
		// network id and nic id
		{
			description: "provided flags invalid 5",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[networkIdFlag] = testNetworkId
			}),
			isValid: false,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.NetworkId = utils.Ptr(testNetworkId)
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			p := print.NewPrinter()
			cmd := NewCmd(&params.CmdParams{Printer: p})
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

			err = cmd.ValidateFlagGroups()
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flag groups: %v", err)
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
				t.Fatalf("error parsing input: %v", err)
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

func TestBuildRequestAttach(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest iaas.ApiAddNicToServerRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequestAttach(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequestAttach(testCtx, tt.model, testClient)

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

func TestBuildRequestCreateAndAttach(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest iaas.ApiAddNetworkToServerRequest
	}{
		{
			description: "base",
			model: fixtureInputModel(func(model *inputModel) {
				model.NicId = nil
				model.NetworkId = utils.Ptr(testNetworkId)
				model.Create = utils.Ptr(true)
			}),
			expectedRequest: fixtureRequestCreateAndAttach(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequestCreateAndAttach(testCtx, tt.model, testClient)

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
