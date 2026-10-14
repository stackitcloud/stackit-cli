package patch

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

var testContext = context.WithValue(context.Background(), testContextKey{}, "token-patch")
var testClient = &modelexperiments.APIClient{DefaultAPI: modelexperiments.DefaultAPIServiceMock{}}
var testProjectID = uuid.NewString()
var testInstanceID = uuid.NewString()
var testTokenID = uuid.NewString()

func fixtureFlagValues(mods ...func(map[string]string)) map[string]string {
	values := map[string]string{projectIDFlag: testProjectID, instanceIDFlag: testInstanceID, globalflags.RegionFlag: "eu01"}
	for _, mod := range mods {
		mod(values)
	}
	return values
}
func fixtureInputModel(mods ...func(*inputModel)) *inputModel {
	model := &inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{ProjectId: testProjectID, Region: "eu01", Verbosity: globalflags.VerbosityDefault}, TokenID: testTokenID, InstanceID: testInstanceID}
	for _, mod := range mods {
		mod(model)
	}
	return model
}
func fixtureRequest(payload modelexperiments.PartialUpdateInstanceTokenPayload) modelexperiments.ApiPartialUpdateInstanceTokenRequest {
	return testClient.DefaultAPI.PartialUpdateInstanceToken(testContext, testProjectID, "eu01", testTokenID, testInstanceID).PartialUpdateInstanceTokenPayload(payload)
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
			values[nameFlag] = "updated"
			values[descriptionFlag] = "updated token"
			values[labelFlag] = "env=prod"
		}), valid: true, want: fixtureInputModel(func(model *inputModel) {
			model.Name = utils.Ptr("updated")
			model.Description = utils.Ptr("updated token")
			model.Labels = &map[string]string{"env": "prod"}
		})},
		{name: "project missing", flags: fixtureFlagValues(func(values map[string]string) { delete(values, projectIDFlag) }), valid: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.want, []string{testTokenID}, tt.flags, tt.valid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	model := fixtureInputModel(func(model *inputModel) {
		model.Name = utils.Ptr("updated")
		model.Description = utils.Ptr("updated token")
		model.Labels = &map[string]string{"env": "prod"}
	})
	got := buildRequest(testContext, model, testClient)
	labels := map[string]*string{"env": utils.Ptr("prod")}
	want := fixtureRequest(modelexperiments.PartialUpdateInstanceTokenPayload{Name: utils.Ptr("updated"), Description: utils.Ptr("updated token"), Labels: &labels})
	if diff := cmp.Diff(got, want, cmp.AllowUnexported(want), cmpopts.EquateComparable(testContext)); diff != "" {
		t.Fatalf("request mismatch (-got +want):\n%s", diff)
	}
}

func TestOutputResult(t *testing.T) {
	params := testparams.NewTestParams()
	if err := outputResult(params.Printer, "", nil); err == nil {
		t.Fatal("expected nil response error")
	}
	if err := outputResult(params.Printer, "", &modelexperiments.PartialUpdateInstanceTokenResponse{}); err == nil {
		t.Fatal("expected nil token error")
	}
}
