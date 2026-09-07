package create

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"

	modelexperiments "github.com/stackitcloud/stackit-sdk-go/services/modelexperiments/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

var projectIDFlag = globalflags.ProjectIdFlag

type testContextKey struct{}

var testContext = context.WithValue(context.Background(), testContextKey{}, "token-create")
var testClient = &modelexperiments.APIClient{DefaultAPI: modelexperiments.DefaultAPIServiceMock{}}
var testProjectID = uuid.NewString()
var testInstanceID = uuid.NewString()

func fixtureFlagValues(mods ...func(map[string]string)) map[string]string {
	values := map[string]string{projectIDFlag: testProjectID, instanceIDFlag: testInstanceID, regionFlag: "eu01", nameFlag: "example"}
	for _, mod := range mods {
		mod(values)
	}
	return values
}

func fixtureInputModel(mods ...func(*inputModel)) *inputModel {
	model := &inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{ProjectId: testProjectID, Region: "eu01", Verbosity: globalflags.VerbosityDefault}, InstanceID: testInstanceID, Region: "eu01", Name: "example"}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(*modelexperiments.ApiCreateInstanceTokenRequest)) modelexperiments.ApiCreateInstanceTokenRequest {
	request := testClient.DefaultAPI.CreateInstanceToken(testContext, testProjectID, "eu01", testInstanceID).CreateInstanceTokenPayload(modelexperiments.CreateInstanceTokenPayload{Name: "example"})
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		name  string
		flags map[string]string
		valid bool
		want  *inputModel
	}{
		{name: "base", flags: fixtureFlagValues(), valid: true, want: fixtureInputModel()},
		{name: "optional fields", flags: fixtureFlagValues(func(values map[string]string) {
			values[descriptionFlag] = "service token"
			values[labelFlag] = "env=prod"
			values[ttlDurationFlag] = "5h"
		}), valid: true, want: fixtureInputModel(func(model *inputModel) {
			model.Description = utils.Ptr("service token")
			model.Labels = &map[string]string{"env": "prod"}
			model.TTLDuration = utils.Ptr("5h")
		})},
		{name: "name missing", flags: fixtureFlagValues(func(values map[string]string) { delete(values, nameFlag) }), valid: false},
		{name: "project missing", flags: fixtureFlagValues(func(values map[string]string) { delete(values, projectIDFlag) }), valid: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { testutils.TestParseInput(t, NewCmd, parseInput, tt.want, nil, tt.flags, tt.valid) })
	}
}

func TestBuildRequest(t *testing.T) {
	model := fixtureInputModel(func(model *inputModel) {
		model.Description = utils.Ptr("service token")
		model.Labels = &map[string]string{"env": "prod"}
		model.TTLDuration = utils.Ptr("5h")
	})
	want := fixtureRequest(func(request *modelexperiments.ApiCreateInstanceTokenRequest) {
		*request = request.CreateInstanceTokenPayload(modelexperiments.CreateInstanceTokenPayload{Name: "example", Description: utils.Ptr("service token"), Labels: &map[string]string{"env": "prod"}, TtlDuration: utils.Ptr("5h")})
	})
	got := buildRequest(testContext, model, testClient)
	if diff := cmp.Diff(got, want, cmp.AllowUnexported(want), cmpopts.EquateComparable(testContext)); diff != "" {
		t.Fatalf("request mismatch (-got +want):\n%s", diff)
	}
}

func TestOutputResult(t *testing.T) {
	params := testparams.NewTestParams()
	if err := outputResult(params.Printer, "", nil); err == nil {
		t.Fatal("expected nil response error")
	}
	if err := outputResult(params.Printer, "", &modelexperiments.CreateInstanceTokenResponse{}); err == nil {
		t.Fatal("expected nil token error")
	}
}
