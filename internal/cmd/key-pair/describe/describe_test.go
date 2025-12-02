package describe

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "test")
var testClient = &iaas.APIClient{}
var testKeyPairName = "foobar"
var testPublicKeyFlag = "true"

func fixtureArgValues(mods ...func(argVales []string)) []string {
	argVales := []string{
		testKeyPairName,
	}
	for _, m := range mods {
		m(argVales)
	}
	return argVales
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{}
	for _, m := range mods {
		m(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			Verbosity: globalflags.VerbosityDefault,
		},
		KeyPairName: testKeyPairName,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *iaas.ApiGetKeyPairRequest)) iaas.ApiGetKeyPairRequest {
	request := testClient.GetKeyPair(testCtx, testKeyPairName)
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description   string
		argsValues    []string
		flagValues    map[string]string
		isValid       bool
		expectedModel *inputModel
	}{
		{
			description:   "base",
			argsValues:    fixtureArgValues(),
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "no values",
			argsValues:  []string{},
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "no arg values",
			argsValues:  []string{},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description:   "no flag values",
			argsValues:    fixtureArgValues(),
			flagValues:    map[string]string{},
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "set flag 'public-key' true",
			argsValues:  fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[publicKeyFlag] = testPublicKeyFlag
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.PublicKey = true
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argsValues, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description    string
		model          *inputModel
		expectedResult iaas.ApiGetKeyPairRequest
	}{
		{
			description:    "base",
			model:          fixtureInputModel(),
			expectedResult: fixtureRequest(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)

			diff := cmp.Diff(request, tt.expectedResult,
				cmp.AllowUnexported(tt.expectedResult),
				cmpopts.EquateComparable(testCtx),
			)
			if diff != "" {
				t.Fatalf("data does not match: %s", diff)
			}
		})
	}
}

func Test_outputResult(t *testing.T) {
	type args struct {
		outputFormat      string
		showOnlyPublicKey bool
		keyPair           iaas.Keypair
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "base",
			args: args{
				outputFormat:      "",
				showOnlyPublicKey: false,
				keyPair:           iaas.Keypair{},
			},
			wantErr: false,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&types.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.outputFormat, tt.args.showOnlyPublicKey, tt.args.keyPair); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
