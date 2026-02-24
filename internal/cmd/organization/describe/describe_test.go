package describe

import (
	"context"
	"testing"
	"time"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/resourcemanager"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &resourcemanager.APIClient{}

var (
	testOrganizationId = uuid.NewString()
)

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testOrganizationId,
	}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			Verbosity: globalflags.VerbosityDefault,
		},
		OrganizationId: testOrganizationId,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *resourcemanager.ApiGetOrganizationRequest)) resourcemanager.ApiGetOrganizationRequest {
	request := testClient.GetOrganization(testCtx, testOrganizationId)
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
			argValues:     fixtureArgValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "uuid as example for an organization id",
			argValues:   []string{"12345678-90ab-cdef-1234-1234567890ab"},
			isValid:     true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.OrganizationId = "12345678-90ab-cdef-1234-1234567890ab"
			}),
		},
		{
			description: "non uuid string as example for a container id",
			argValues:   []string{"foo-bar-organization"},
			isValid:     true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.OrganizationId = "foo-bar-organization"
			}),
		},
		{
			description: "no args",
			argValues:   []string{},
			isValid:     false,
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
		expectedRequest resourcemanager.ApiGetOrganizationRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
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
		organization *resourcemanager.OrganizationResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty",
			args:    args{},
			wantErr: false,
		},
		{
			name: "nil pointer as organization",
			args: args{
				organization: nil,
			},
			wantErr: false,
		},
		{
			name: "empty organization",
			args: args{
				organization: utils.Ptr(resourcemanager.OrganizationResponse{}),
			},
			wantErr: false,
		},
		{
			name: "full response",
			args: args{
				organization: utils.Ptr(resourcemanager.OrganizationResponse{
					OrganizationId: utils.Ptr(uuid.NewString()),
					Name:           utils.Ptr("foo bar"),
					LifecycleState: utils.Ptr(resourcemanager.LIFECYCLESTATE_ACTIVE),
					ContainerId:    utils.Ptr("foo-bar-organization"),
					CreationTime:   utils.Ptr(time.Now()),
					UpdateTime:     utils.Ptr(time.Now()),
					Labels: utils.Ptr(map[string]string{
						"foo": "true",
						"bar": "false",
					}),
				}),
			},
			wantErr: false,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&types.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.outputFormat, tt.args.organization); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
