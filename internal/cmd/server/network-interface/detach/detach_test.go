package detach

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
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

func fixtureRequestDetach(mods ...func(request *iaas.ApiRemoveNicFromServerRequest)) iaas.ApiRemoveNicFromServerRequest {
	request := testClient.RemoveNicFromServer(testCtx, testProjectId, testServerId, testNicId)
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixtureRequestDetachAndDelete(mods ...func(request *iaas.ApiRemoveNetworkFromServerRequest)) iaas.ApiRemoveNetworkFromServerRequest {
	request := testClient.RemoveNetworkFromServer(testCtx, testProjectId, testServerId, testNetworkId)
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
		// only delete
		{
			description: "provided flags invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[deleteFlag] = "true"
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
		// delete and nic id
		{
			description: "provided flags invalid 3",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[deleteFlag] = "true"
			}),
			isValid: false,
		},
		// delete and network id (valid)
		{
			description: "provided flags valid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[deleteFlag] = "true"
				delete(flagValues, networkInterfaceIdFlag)
				flagValues[networkIdFlag] = testNetworkId
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Delete = utils.Ptr(true)
				model.NicId = nil
				model.NetworkId = utils.Ptr(testNetworkId)
			}),
		},
		// delete, nic id and network id
		{
			description: "provided flags invalid 4",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[deleteFlag] = "true"
				flagValues[networkIdFlag] = testNetworkId
			}),
			isValid: false,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Delete = utils.Ptr(true)
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
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildRequestDetach(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest iaas.ApiRemoveNicFromServerRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequestDetach(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequestDetach(testCtx, tt.model, testClient)

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

func TestBuildRequestDetachAndDelete(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest iaas.ApiRemoveNetworkFromServerRequest
	}{
		{
			description: "base",
			model: fixtureInputModel(func(model *inputModel) {
				model.NicId = nil
				model.NetworkId = utils.Ptr(testNetworkId)
				model.Delete = utils.Ptr(true)
			}),
			expectedRequest: fixtureRequestDetachAndDelete(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequestDetachAndDelete(testCtx, tt.model, testClient)

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
