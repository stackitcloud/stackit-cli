package get

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
)

var projectIDFlag = globalflags.ProjectIdFlag

type testContextKey struct{}

var testContext = context.WithValue(context.Background(), testContextKey{}, "token-get")
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
func fixtureInputModel() *inputModel {
	return &inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{ProjectId: testProjectID, Region: "eu01", Verbosity: globalflags.VerbosityDefault}, TokenID: testTokenID, InstanceID: testInstanceID}
}
func fixtureRequest() modelexperiments.ApiGetInstanceTokenRequest {
	return testClient.DefaultAPI.GetInstanceToken(testContext, testProjectID, "eu01", testTokenID, testInstanceID)
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		name  string
		flags map[string]string
		valid bool
	}{{name: "base", flags: fixtureFlagValues(), valid: true}, {name: "project missing", flags: fixtureFlagValues(func(values map[string]string) { delete(values, projectIDFlag) }), valid: false}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var want *inputModel
			if tt.valid {
				want = fixtureInputModel()
			}
			testutils.TestParseInput(t, NewCmd, parseInput, want, []string{testTokenID}, tt.flags, tt.valid)
		})
	}
}
func TestBuildRequest(t *testing.T) {
	got := buildRequest(testContext, fixtureInputModel(), testClient)
	want := fixtureRequest()
	if diff := cmp.Diff(got, want, cmp.AllowUnexported(want), cmpopts.EquateComparable(testContext)); diff != "" {
		t.Fatalf("request mismatch (-got +want):\n%s", diff)
	}
}
func TestOutputResult(t *testing.T) {
	params := testparams.NewTestParams()
	if err := outputResult(params.Printer, "", nil); err == nil {
		t.Fatal("expected nil response error")
	}
	if err := outputResult(params.Printer, "", &modelexperiments.GetInstanceTokenResponse{}); err == nil {
		t.Fatal("expected nil token error")
	}
}
