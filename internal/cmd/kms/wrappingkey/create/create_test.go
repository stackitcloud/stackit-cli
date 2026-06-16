package create

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	kms "github.com/stackitcloud/stackit-sdk-go/services/kms/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	testRegion      = "eu01"
	testAlgorithm   = kms.WRAPPINGALGORITHM_RSA_2048_OAEP_SHA256
	testDisplayName = "my-key"
	testPurpose     = kms.WRAPPINGPURPOSE_WRAP_ASYMMETRIC_KEY
	testDescription = "my key description"
	testProtection  = kms.PROTECTION_SOFTWARE
)

type testCtxKey struct{}

var (
	testCtx       = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient    = &kms.APIClient{DefaultAPI: &kms.DefaultAPIService{}}
	testProjectId = uuid.NewString()
	testKeyRingId = uuid.NewString()
)

// Flags
func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		keyRingIdFlag:             testKeyRingId,
		algorithmFlag.Name():      string(testAlgorithm),
		displayNameFlag:           testDisplayName,
		purposeFlag.Name():        string(testPurpose),
		descriptionFlag:           testDescription,
		protectionFlag.Name():     string(testProtection),
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

// Input Model
func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Region:    testRegion,
			Verbosity: globalflags.VerbosityDefault,
		},
		KeyRingId:   testKeyRingId,
		Algorithm:   testAlgorithm,
		Name:        utils.Ptr(testDisplayName),
		Purpose:     testPurpose,
		Description: utils.Ptr(testDescription),
		Protection:  testProtection,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

// Request
func fixtureRequest(mods ...func(request *kms.ApiCreateWrappingKeyRequest)) kms.ApiCreateWrappingKeyRequest {
	request := testClient.DefaultAPI.CreateWrappingKey(testCtx, testProjectId, testRegion, testKeyRingId)
	request = request.CreateWrappingKeyPayload(kms.CreateWrappingKeyPayload{
		Algorithm:   testAlgorithm,
		DisplayName: testDisplayName,
		Purpose:     testPurpose,
		Description: utils.Ptr(testDescription),
		Protection:  testProtection,
	})

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
			description: "optional flags omitted",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, descriptionFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Description = nil
			}),
		},
		{
			description: "no values provided",
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
			description: "key ring id missing (required)",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, keyRingIdFlag)
			}),
			isValid: false,
		},
		{
			description: "key ring id invalid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[keyRingIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "algorithm missing (required)",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, algorithmFlag.Name())
			}),
			isValid: false,
		},
		{
			description: "name missing (required)",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, displayNameFlag)
			}),
			isValid: false,
		},
		{
			description: "purpose missing (required)",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, purposeFlag.Name())
			}),
			isValid: false,
		},
		{
			description: "protection missing (required)",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, protectionFlag.Name())
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			cmd := &cobra.Command{}
			err := globalflags.Configure(cmd.Flags())
			if err != nil {
				t.Fatalf("configure global flags: %v", err)
			}

			configureFlags(cmd)

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

			params := testparams.NewTestParams()
			model, err := parseInput(params.Printer, cmd)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error parsing flags: %v", err)
			}

			if !tt.isValid {
				t.Fatalf("did not fail on invalid input")
			}
			diff := cmp.Diff(tt.expectedModel, model)
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
		expectedRequest kms.ApiCreateWrappingKeyRequest
	}{
		{
			description:     "base case",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "no optional values",
			model: fixtureInputModel(func(model *inputModel) {
				model.Description = nil
			}),
			expectedRequest: fixtureRequest().CreateWrappingKeyPayload(kms.CreateWrappingKeyPayload{
				Algorithm:   testAlgorithm,
				DisplayName: testDisplayName,
				Purpose:     testPurpose,
				Protection:  testProtection,
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request, err := buildRequest(testCtx, tt.model, testClient.DefaultAPI)
			if err != nil {
				t.Fatalf("error building request: %v", err)
			}

			diff := cmp.Diff(tt.expectedRequest, request,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx, kms.DefaultAPIService{}),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	tests := []struct {
		description string
		model       *inputModel
		wrappingKey *kms.WrappingKey
		wantErr     bool
	}{
		{
			description: "nil response",
			wrappingKey: nil,
			wantErr:     true,
		},
		{
			description: "default output",
			model:       &inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{}},
			wrappingKey: &kms.WrappingKey{},
			wantErr:     false,
		},
		{
			description: "json output",
			model:       &inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{OutputFormat: print.JSONOutputFormat}},
			wrappingKey: &kms.WrappingKey{},
			wantErr:     false,
		},
		{
			description: "yaml output",
			model:       &inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{OutputFormat: print.YAMLOutputFormat}},
			wrappingKey: &kms.WrappingKey{},
			wantErr:     false,
		},
	}

	params := testparams.NewTestParams()

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := outputResult(params.Printer, tt.model, tt.wrappingKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
