package update

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

var projectIdFlag = globalflags.ProjectIdFlag

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &iaas.APIClient{}
var testProjectId = uuid.NewString()
var testNetworkId = uuid.NewString()
var testNicId = uuid.NewString()
var testSecurityGroup = uuid.NewString()

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testNicId,
	}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:        testProjectId,
		networkIdFlag:        testNetworkId,
		allowedAddressesFlag: "1.1.1.1,8.8.8.8,9.9.9.9",
		labelFlag:            "key=value",
		nameFlag:             "testNic",
		nicSecurityFlag:      "true",
		securityGroupsFlag:   testSecurityGroup,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	var allowedAddresses []iaas.AllowedAddressesInner = []iaas.AllowedAddressesInner{
		iaas.StringAsAllowedAddressesInner(utils.Ptr("1.1.1.1")),
		iaas.StringAsAllowedAddressesInner(utils.Ptr("8.8.8.8")),
		iaas.StringAsAllowedAddressesInner(utils.Ptr("9.9.9.9")),
	}
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Verbosity: globalflags.VerbosityDefault,
		},
		NetworkId:        utils.Ptr(testNetworkId),
		AllowedAddresses: utils.Ptr(allowedAddresses),
		Labels: utils.Ptr(map[string]string{
			"key": "value",
		}),
		Name:           utils.Ptr("testNic"),
		NicSecurity:    utils.Ptr(true),
		SecurityGroups: utils.Ptr([]string{testSecurityGroup}),
		NicId:          testNicId,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *iaas.ApiUpdateNicRequest)) iaas.ApiUpdateNicRequest {
	request := testClient.UpdateNic(testCtx, testProjectId, testNetworkId, testNicId)
	request = request.UpdateNicPayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayload(mods ...func(payload *iaas.UpdateNicPayload)) iaas.UpdateNicPayload {
	var allowedAddresses []iaas.AllowedAddressesInner = []iaas.AllowedAddressesInner{
		iaas.StringAsAllowedAddressesInner(utils.Ptr("1.1.1.1")),
		iaas.StringAsAllowedAddressesInner(utils.Ptr("8.8.8.8")),
		iaas.StringAsAllowedAddressesInner(utils.Ptr("9.9.9.9")),
	}
	payload := iaas.UpdateNicPayload{
		AllowedAddresses: utils.Ptr(allowedAddresses),
		Labels: utils.Ptr(map[string]interface{}{
			"key": "value",
		}),
		Name:           utils.Ptr("testNic"),
		NicSecurity:    utils.Ptr(true),
		SecurityGroups: utils.Ptr([]string{testSecurityGroup}),
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
			argValues:     fixtureArgValues(),
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},

		{
			description: "no values",
			argValues:   fixtureArgValues(),
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "network id missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, networkIdFlag)
			}),
			isValid: false,
		},
		{
			description: "network id invalid 1",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[networkIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "network id invalid 2",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[networkIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "allowed addresses missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, allowedAddressesFlag)
			}),
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.AllowedAddresses = nil
			}),
			isValid: true,
		},
		{
			description: "name to long",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[nameFlag] = "verylongstringwith66characterstotestthenameregexwithinthisunittest"
			}),
			isValid: false,
		},
		{
			description: "name invalid",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[nameFlag] = "test?"
			}),
			isValid: false,
		},
		{
			description: "name empty string invalid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[nameFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "security group uuid to short",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[securityGroupsFlag] = "d61a8564-c8dd-4ffb-bc15-143e7d0c85e"
			}),
			isValid: false,
		},
		{
			description: "security group uuid invalid",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[securityGroupsFlag] = "d61a8564-c8dd-4ffb-bc15-143e7d0c85e?"
			}),
			isValid: false,
		},
		{
			description: "nic argument missing",
			argValues:   []string{},
			isValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			p := print.NewPrinter()
			cmd := NewCmd(p)
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

			err = cmd.ValidateArgs(tt.argValues)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error parsing args: %v", err)
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
		expectedRequest iaas.ApiUpdateNicRequest
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

func Test_outputResult(t *testing.T) {
	type args struct {
		outputFormat string
		projectId    string
		nic          *iaas.NIC
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty",
			args:    args{},
			wantErr: true,
		},
		{
			name: "set empty nic",
			args: args{
				nic: &iaas.NIC{},
			},
			wantErr: false,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(p)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.outputFormat, tt.args.projectId, tt.args.nic); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
