package describe

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/kms"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &kms.APIClient{}
var testProjectId = uuid.NewString()
var testKeyRingID = uuid.NewString()
var testWrappingKeyID = uuid.NewString()
var testTime = time.Time{}

const testRegion = "eu01"

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		flagKeyRingID:             testKeyRingID,
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
			Region:    testRegion,
			Verbosity: globalflags.VerbosityDefault,
		},
		KeyRingID:     testKeyRingID,
		WrappingKeyID: testWrappingKeyID,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
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
			argValues:     []string{testWrappingKeyID},
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "no values",
			argValues:   []string{},
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "no arg values",
			argValues:   []string{},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "no flag values",
			argValues:   []string{testWrappingKeyID},
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "invalid key ring id",
			argValues:   []string{testWrappingKeyID},
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[flagKeyRingID] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "missing project id",
			argValues:   []string{testWrappingKeyID},
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "invalid project id",
			argValues:   []string{testWrappingKeyID},
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = "invalid-uuid"
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
	got := buildRequest(testCtx, fixtureInputModel(), testClient)
	want := testClient.GetWrappingKey(testCtx, testProjectId, testRegion, testKeyRingID, testWrappingKeyID)
	diff := cmp.Diff(got, want,
		cmp.AllowUnexported(want),
		cmpopts.EquateComparable(testCtx),
	)
	if diff != "" {
		t.Fatalf("buildRequest() mismatch (-want +got):\n%s", diff)
	}
}
func TestOutputResult(t *testing.T) {
	tests := []struct {
		description string
		outputFmt   string
		keyRing     *kms.WrappingKey
		wantErr     bool
		expected    string
	}{
		{
			description: "empty",
			outputFmt:   "table",
			wantErr:     true,
		},
		{
			description: "table format",
			outputFmt:   "table",
			keyRing: &kms.WrappingKey{
				Id:          utils.Ptr(testWrappingKeyID),
				DisplayName: utils.Ptr("Test Key Ring"),
				CreatedAt:   utils.Ptr(testTime),
				Description: utils.Ptr("This is a test key ring."),
				State:       utils.Ptr(kms.WRAPPINGKEYSTATE_ACTIVE),
				AccessScope: utils.Ptr(kms.ACCESSSCOPE_PUBLIC),
				Algorithm:   utils.Ptr(kms.WRAPPINGALGORITHM__2048_OAEP_SHA256),
				ExpiresAt:   utils.Ptr(testTime),
				KeyRingId:   utils.Ptr(testKeyRingID),
				Protection:  utils.Ptr(kms.PROTECTION_SOFTWARE),
				PublicKey:   utils.Ptr("-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQ...\n-----END PUBLIC KEY-----"),
				Purpose:     utils.Ptr(kms.WRAPPINGPURPOSE_ASYMMETRIC_KEY),
			},
			expected: fmt.Sprintf(`
 ID           │ %-46s
──────────────┼───────────────────────────────────────────────
 DISPLAY NAME │ Test Key Ring                                 
──────────────┼───────────────────────────────────────────────
 CREATED AT   │ %-46s
──────────────┼───────────────────────────────────────────────
 STATE        │ active                                        
──────────────┼───────────────────────────────────────────────
 DESCRIPTION  │ This is a test key ring.                      
──────────────┼───────────────────────────────────────────────
 ACCESS SCOPE │ PUBLIC                                        
──────────────┼───────────────────────────────────────────────
 ALGORITHM    │ rsa_2048_oaep_sha256                          
──────────────┼───────────────────────────────────────────────
 EXPIRES AT   │ %-46s
──────────────┼───────────────────────────────────────────────
 KEYRING ID   │ %-46s
──────────────┼───────────────────────────────────────────────
 PROTECTION   │ software                                      
──────────────┼───────────────────────────────────────────────
 PUBLIC KEY   │ -----BEGIN PUBLIC KEY-----                    
              │ MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQ... 
              │ -----END PUBLIC KEY-----                      
──────────────┼───────────────────────────────────────────────
 PURPOSE      │ wrap_asymmetric_key                           

`,
				testWrappingKeyID,
				testTime,
				testTime,
				testKeyRingID),
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&params.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			var buf bytes.Buffer
			p.Cmd.SetOut(&buf)
			if err := outputResult(p, tt.outputFmt, tt.keyRing); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
			diff := cmp.Diff(buf.String(), tt.expected)
			if diff != "" {
				t.Fatalf("outputResult() output mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
