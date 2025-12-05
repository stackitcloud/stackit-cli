package describe

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
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
var testKeyID = uuid.NewString()
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
		KeyID:     testKeyID,
		KeyRingID: testKeyRingID,
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
			argValues:     []string{testKeyID},
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
			argValues:   []string{testKeyID},
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "invalid key id",
			argValues:   []string{"invalid-uuid"},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "missing key ring id",
			argValues:   []string{testKeyID},
			flagValues:  fixtureFlagValues(func(m map[string]string) { delete(m, flagKeyRingID) }),
			isValid:     false,
		},
		{
			description: "invalid key ring id",
			argValues:   []string{testKeyID},
			flagValues: fixtureFlagValues(func(m map[string]string) {
				m[flagKeyRingID] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "missing project id",
			argValues:   []string{testKeyID},
			flagValues:  fixtureFlagValues(func(m map[string]string) { delete(m, globalflags.ProjectIdFlag) }),
			isValid:     false,
		},
		{
			description: "invalid project id",
			argValues:   []string{testKeyID},
			flagValues:  fixtureFlagValues(func(m map[string]string) { m[globalflags.ProjectIdFlag] = "invalid-uuid" }),
			isValid:     false,
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
	want := testClient.GetKey(testCtx, testProjectId, testRegion, testKeyRingID, testKeyID)
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
		keyRing     *kms.Key
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
			keyRing: &kms.Key{
				AccessScope:  utils.Ptr(kms.ACCESSSCOPE_PUBLIC),
				Algorithm:    utils.Ptr(kms.ALGORITHM_AES_256_GCM),
				CreatedAt:    utils.Ptr(testTime),
				DeletionDate: nil,
				Description:  utils.Ptr("very secure and secret key"),
				DisplayName:  utils.Ptr("Test Key"),
				Id:           utils.Ptr(testKeyID),
				ImportOnly:   utils.Ptr(true),
				KeyRingId:    utils.Ptr(testKeyRingID),
				Protection:   utils.Ptr(kms.PROTECTION_SOFTWARE),
				Purpose:      utils.Ptr(kms.PURPOSE_SYMMETRIC_ENCRYPT_DECRYPT),
				State:        utils.Ptr(kms.KEYSTATE_ACTIVE),
			},
			expected: fmt.Sprintf(`
 ID            │ %-37s
───────────────┼──────────────────────────────────────
 DISPLAY NAME  │ Test Key                             
───────────────┼──────────────────────────────────────
 CREATED AT    │ %-37s
───────────────┼──────────────────────────────────────
 STATE         │ active                               
───────────────┼──────────────────────────────────────
 DESCRIPTION   │ very secure and secret key           
───────────────┼──────────────────────────────────────
 ACCESS SCOPE  │ PUBLIC                               
───────────────┼──────────────────────────────────────
 ALGORITHM     │ aes_256_gcm                          
───────────────┼──────────────────────────────────────
 DELETION DATE │                                      
───────────────┼──────────────────────────────────────
 IMPORT ONLY   │ true                                 
───────────────┼──────────────────────────────────────
 KEYRING ID    │ %-37s
───────────────┼──────────────────────────────────────
 PROTECTION    │ software                             
───────────────┼──────────────────────────────────────
 PURPOSE       │ symmetric_encrypt_decrypt            

`,
				testKeyID,
				testTime,
				testKeyRingID,
			),
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&types.CmdParams{Printer: p})
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
