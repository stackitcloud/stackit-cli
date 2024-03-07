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
	"github.com/stackitcloud/stackit-sdk-go/services/objectstorage"
)

var projectIdFlag = globalflags.ProjectIdFlag

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &objectstorage.APIClient{}
var testProjectId = uuid.NewString()
var testCredentialsGroupId = uuid.NewString()
var testExpirationDate = "2024-01-01T00:00:00Z"

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:          testProjectId,
		credentialsGroupIdFlag: testCredentialsGroupId,
		expireDateFlag:         testExpirationDate,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	testExpirationDate, err := time.Parse(expirationTimeFormat, testExpirationDate)
	if err != nil {
		return &inputModel{}
	}

	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
		},
		ExpireDate:         utils.Ptr(testExpirationDate),
		CredentialsGroupId: testCredentialsGroupId,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixturePayload(mods ...func(payload *objectstorage.CreateAccessKeyPayload)) objectstorage.CreateAccessKeyPayload {
	testExpirationDate, err := time.Parse(expirationTimeFormat, testExpirationDate)
	if err != nil {
		return objectstorage.CreateAccessKeyPayload{}
	}
	payload := objectstorage.CreateAccessKeyPayload{
		Expires: utils.Ptr(testExpirationDate),
	}
	for _, mod := range mods {
		mod(&payload)
	}
	return payload
}

func fixtureRequest(mods ...func(request *objectstorage.ApiCreateAccessKeyRequest)) objectstorage.ApiCreateAccessKeyRequest {
	request := testClient.CreateAccessKey(testCtx, testProjectId)
	request = request.CreateAccessKeyPayload(fixturePayload())
	request = request.CredentialsGroup(testCredentialsGroupId)
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
			description: "credentials group id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, credentialsGroupIdFlag)
			}),
			isValid: false,
		},
		{
			description: "credentials group id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[credentialsGroupIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "credentials group id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[credentialsGroupIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "expiration date is missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, expireDateFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.ExpireDate = nil
			}),
		},
		{
			description: "expiration date is empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[expireDateFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "expiration date is invalid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[expireDateFlag] = "test"
			}),
			isValid: false,
		},
		{
			description: "expiration date is invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[expireDateFlag] = "11:00 12/12/2024"
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			cmd := NewCmd()
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

			model, err := parseInput(cmd)
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
		expectedRequest objectstorage.ApiCreateAccessKeyRequest
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
