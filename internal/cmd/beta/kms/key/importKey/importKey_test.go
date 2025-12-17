package importKey

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-sdk-go/services/kms"
)

const (
	testRegion = "eu01"
)

type testCtxKey struct{}

var (
	testCtx           = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient        = &kms.APIClient{}
	testProjectId     = uuid.NewString()
	testKeyRingId     = uuid.NewString()
	testKeyId         = uuid.NewString()
	testWrappingKeyId = uuid.NewString()
	testWrappedKey    = "SnVzdCBzYXlpbmcgaGV5Oyk="
)

// Args
func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testKeyId,
	}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

// Flags
func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		keyRingIdFlag:             testKeyRingId,
		wrappedKeyFlag:            testWrappedKey,
		wrappingKeyIdFlag:         testWrappingKeyId,
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
		KeyRingId:     testKeyRingId,
		KeyId:         testKeyId,
		WrappedKey:    &testWrappedKey,
		WrappingKeyId: &testWrappingKeyId,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

// Request
func fixtureRequest(mods ...func(request *kms.ApiImportKeyRequest)) kms.ApiImportKeyRequest {
	request := testClient.ImportKey(testCtx, testProjectId, testRegion, testKeyRingId, testKeyId)
	request = request.ImportKeyPayload(kms.ImportKeyPayload{
		WrappedKey:    &testWrappedKey,
		WrappingKeyId: &testWrappingKeyId,
	})

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
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "no args (keyId)",
			argValues:   []string{},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "no values provided",
			argValues:   fixtureArgValues(),
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "project id missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "key ring id missing (required)",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, keyRingIdFlag)
			}),
			isValid: false,
		},
		{
			description: "key ring id invalid 1",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[keyRingIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "key ring id invalid 2",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[keyRingIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "key id invalid 1",
			argValues:   []string{""},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "key id invalid 2",
			argValues:   []string{"invalid-key"},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "wrapping key id missing (required)",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, wrappingKeyIdFlag)
			}),
			isValid: false,
		},
		{
			description: "wrapping key id invalid 1",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[wrappingKeyIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "wrapping key id invalid 2",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[wrappingKeyIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "wrapped key missing (required)",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, wrappedKeyFlag)
			}),
			isValid: false,
		},
		{
			description: "wrapped key invalid 1",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[wrappedKeyFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "wrapped key invalid 2 - not base64",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[wrappedKeyFlag] = "Not Base 64"
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			p := print.NewPrinter()
			cmd := NewCmd(&types.CmdParams{Printer: p})
			err := globalflags.Configure(cmd.Flags())
			if err != nil {
				t.Fatalf("configure global flags: %v", err)
			}

			err = cmd.ValidateArgs(tt.argValues)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating args: %v", err)
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

			model, err := parseInput(p, cmd, tt.argValues)
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
		expectedRequest kms.ApiImportKeyRequest
	}{
		{
			description:     "base case",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request, err := buildRequest(testCtx, tt.model, testClient)
			if err != nil {
				t.Fatalf("error building request: %v", err)
			}

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

func TestOutputResult(t *testing.T) {
	tests := []struct {
		description  string
		version      *kms.Version
		outputFormat string
		keyRingName  string
		keyName      string
		wantErr      bool
	}{
		{
			description: "nil response",
			version:     nil,
			wantErr:     true,
		},
		{
			description: "default output",
			version:     &kms.Version{},
			keyRingName: "my-key-ring",
			keyName:     "my-key",
			wantErr:     false,
		},
		{
			description:  "json output",
			version:      &kms.Version{},
			outputFormat: print.JSONOutputFormat,
			keyRingName:  "my-key-ring",
			keyName:      "my-key",
			wantErr:      false,
		},
		{
			description:  "yaml output",
			version:      &kms.Version{},
			outputFormat: print.YAMLOutputFormat,
			keyRingName:  "my-key-ring",
			keyName:      "my-key",
			wantErr:      false,
		},
	}

	p := print.NewPrinter()
	p.Cmd = NewCmd(&types.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := outputResult(p, tt.outputFormat, tt.keyRingName, tt.keyName, tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
