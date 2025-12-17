package create

import (
	"context"
	"os"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &iaas.APIClient{}

var testPublicKey = "ssh-rsa <key>"
var testKeyPairName = "foobar_key"

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		publicKeyFlag: testPublicKey,
		labelFlag:     "foo=bar",
		nameFlag:      testKeyPairName,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			Verbosity: globalflags.VerbosityDefault,
		},
		Labels: utils.Ptr(map[string]string{
			"foo": "bar",
		}),
		PublicKey: utils.Ptr(testPublicKey),
		Name:      utils.Ptr(testKeyPairName),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *iaas.ApiCreateKeyPairRequest)) iaas.ApiCreateKeyPairRequest {
	request := testClient.CreateKeyPair(testCtx)
	request = request.CreateKeyPairPayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayload(mods ...func(payload *iaas.CreateKeyPairPayload)) iaas.CreateKeyPairPayload {
	payload := iaas.CreateKeyPairPayload{
		Labels: utils.Ptr(map[string]interface{}{
			"foo": "bar",
		}),
		PublicKey: utils.Ptr(testPublicKey),
		Name:      utils.Ptr(testKeyPairName),
	}
	for _, mod := range mods {
		mod(&payload)
	}
	return payload
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
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "required only",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, nameFlag)
				delete(flagValues, labelFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Name = nil
				model.Labels = nil
			}),
		},
		{
			description: "read public key from file",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[publicKeyFlag] = "@./template/id_ed25519.pub"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				file, err := os.ReadFile("./template/id_ed25519.pub")
				if err != nil {
					t.Fatal("could not create expected Model", err)
				}
				model.PublicKey = utils.Ptr(string(file))
			}),
		},
		{
			description: "no values",
			flagValues:  map[string]string{},
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
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest iaas.ApiCreateKeyPairRequest
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
				cmp.AllowUnexported(iaas.NullableString{}),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func Test_outputResult(t *testing.T) {
	type args struct {
		item         *iaas.Keypair
		outputFormat string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				item:         nil,
				outputFormat: "",
			},
			wantErr: true,
		},
		{
			name: "base",
			args: args{
				item:         &iaas.Keypair{},
				outputFormat: "",
			},
			wantErr: false,
		},
	}

	p := print.NewPrinter()
	p.Cmd = NewCmd(&types.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.outputFormat, tt.args.item); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
