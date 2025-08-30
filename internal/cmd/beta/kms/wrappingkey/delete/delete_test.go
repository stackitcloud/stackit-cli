package delete

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-sdk-go/services/kms"
)

const (
	testRegion = "eu01"
)

type testCtxKey struct{}

var (
	testCtx               = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient            = &kms.APIClient{}
	testProjectId         = uuid.NewString()
	testKeyRingId         = uuid.NewString()
	testWrappingKeyRingId = uuid.NewString()
)

// Flags
func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		keyRingIdFlag:             testKeyRingId,
		wrappingKeyIdFlag:         testWrappingKeyRingId,
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
		WrappingKey: testWrappingKeyRingId,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

// Request
func fixtureRequest(mods ...func(request *kms.ApiDeleteWrappingKeyRequest)) kms.ApiDeleteWrappingKeyRequest {
	request := testClient.DeleteWrappingKey(testCtx, testProjectId, testRegion, testKeyRingId, testWrappingKeyRingId)
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
			description: "wrapping key id missing (required)",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, wrappingKeyIdFlag)
			}),
			isValid: false,
		},
		{
			description: "wrapping key id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[wrappingKeyIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "wrapping key id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[wrappingKeyIdFlag] = "invalid-uuid"
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

			p := print.NewPrinter()
			model, err := parseInput(p, cmd)
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
		expectedRequest kms.ApiDeleteWrappingKeyRequest
	}{
		{
			description:     "base case",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)

			diff := cmp.Diff(tt.expectedRequest, request,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}
