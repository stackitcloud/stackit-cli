package create

import (
	"context"
	"testing"
	"time"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/serviceaccount"
)

var projectIdFlag = globalflags.ProjectIdFlag

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &serviceaccount.APIClient{}
var testProjectId = uuid.NewString()
var testServiceAccountEmail = "my-service-account-1234567@sa.stackit.cloud"
var testNow = time.Now()
var test10DaysFromNow = daysFromNow(testNow, 10)
var testPublicKey = "my-public-key"

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:           testProjectId,
		serviceAccountEmailFlag: testServiceAccountEmail,
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
		ServiceAccountEmail: testServiceAccountEmail,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *serviceaccount.ApiCreateServiceAccountKeyRequest)) serviceaccount.ApiCreateServiceAccountKeyRequest {
	request := testClient.CreateServiceAccountKey(testCtx, testProjectId, testServiceAccountEmail)
	request = request.CreateServiceAccountKeyPayload(serviceaccount.CreateServiceAccountKeyPayload{})
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
			description: "with expiring date",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[expiredInDaysFlag] = "10"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.ExpiresInDays = utils.Ptr(int64(10))
			}),
		},
		{
			description: "with public key",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[publicKeyFlag] = testPublicKey
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.PublicKey = utils.Ptr(testPublicKey)
			}),
		},
		{
			description: "with public key",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[publicKeyFlag] = ""
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.PublicKey = utils.Ptr("")
			}),
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
			description: "service account email missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, serviceAccountEmailFlag)
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			cmd := NewCmd(nil)
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

			model, err := parseInput(nil, cmd)
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
		isValid         bool
		expectedRequest serviceaccount.ApiCreateServiceAccountKeyRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			isValid:         true,
			expectedRequest: fixtureRequest(),
		},
		{
			description: "with expiring date",
			model: fixtureInputModel(func(model *inputModel) {
				model.ExpiresInDays = utils.Ptr(int64(10))
			}),
			isValid: true,
			expectedRequest: fixtureRequest().CreateServiceAccountKeyPayload(
				serviceaccount.CreateServiceAccountKeyPayload{
					ValidUntil: utils.Ptr(test10DaysFromNow),
				}),
		},
		{
			description: "with public key",
			model: fixtureInputModel(func(model *inputModel) {
				model.PublicKey = utils.Ptr(testPublicKey)
			}),
			isValid: true,
			expectedRequest: fixtureRequest().CreateServiceAccountKeyPayload(
				serviceaccount.CreateServiceAccountKeyPayload{
					PublicKey: utils.Ptr(testPublicKey),
				}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient, testNow)

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
