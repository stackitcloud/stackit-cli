package list

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/authorization"
)

var projectIdFlag = globalflags.ProjectIdFlag

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &authorization.APIClient{}
var testProjectId = uuid.NewString()

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag: testProjectId,
		limitFlag:     "10",
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
		Limit:  utils.Ptr(int64(10)),
		SortBy: "subject",
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *authorization.ApiListMembersRequest)) authorization.ApiListMembersRequest {
	request := testClient.ListMembers(testCtx, projectResourceType, testProjectId)
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
			description: "with subject",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[subjectFlag] = "someone@domain.com"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(
				func(model *inputModel) {
					model.Subject = utils.Ptr("someone@domain.com")
				},
			),
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
			description: "limit invalid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[limitFlag] = "invalid"
			}),
			isValid: false,
		},
		{
			description: "limit invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[limitFlag] = "0"
			}),
			isValid: false,
		},
		{
			description: "sort by role",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[sortByFlag] = "role"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.SortBy = "role"
			}),
		},
		{
			description: "sort by invalid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[sortByFlag] = "invalid"
			}),
			isValid: false,
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
		expectedRequest authorization.ApiListMembersRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "with subject",
			model: fixtureInputModel(func(model *inputModel) {
				model.Subject = utils.Ptr("someone@domain.com")
			}),
			expectedRequest: fixtureRequest().Subject("someone@domain.com"),
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

func Test_outputResult(t *testing.T) {
	type args struct {
		model   inputModel
		members []authorization.Member
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty", args{model: inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{}}}, false},
		{"base", args{inputModel{
			GlobalFlagModel: &globalflags.GlobalFlagModel{},
			Subject:         utils.Ptr("subject"),
			Limit:           nil,
			SortBy:          "",
		}, nil}, false},
		{"complete", args{inputModel{
			GlobalFlagModel: &globalflags.GlobalFlagModel{},
			Subject:         utils.Ptr("subject"),
			Limit:           nil,
			SortBy:          "",
		},
			[]authorization.Member{
				{Role: utils.Ptr("role1"), Subject: utils.Ptr("subject1")},
				{Role: utils.Ptr("role2"), Subject: utils.Ptr("subject2")},
				{Role: utils.Ptr("role3"), Subject: utils.Ptr("subject3")},
			}},
			false},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&params.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.model, tt.args.members); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
