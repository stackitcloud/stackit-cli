package create

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
var testSecurityGroupId = uuid.NewString()
var testRemoteSecurityGroupId = uuid.NewString()

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:             testProjectId,
		securityGroupIdFlag:       testSecurityGroupId,
		directionFlag:             "ingress",
		descriptionFlag:           "example-description",
		etherTypeFlag:             "ether",
		icmpParameterCodeFlag:     "0",
		icmpParameterTypeFlag:     "8",
		ipRangeFlag:               "10.1.2.3",
		portRangeMaxFlag:          "24",
		portRangeMinFlag:          "22",
		remoteSecurityGroupIdFlag: testRemoteSecurityGroupId,
		protocolNumberFlag:        "1",
		protocolNameFlag:          "icmp",
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
		SecurityGroupId:       testSecurityGroupId,
		Direction:             utils.Ptr("ingress"),
		Description:           utils.Ptr("example-description"),
		EtherType:             utils.Ptr("ether"),
		IcmpParameterCode:     utils.Ptr(int64(0)),
		IcmpParameterType:     utils.Ptr(int64(8)),
		IpRange:               utils.Ptr("10.1.2.3"),
		PortRangeMax:          utils.Ptr(int64(24)),
		PortRangeMin:          utils.Ptr(int64(22)),
		RemoteSecurityGroupId: utils.Ptr(testRemoteSecurityGroupId),
		ProtocolNumber:        utils.Ptr(int64(1)),
		ProtocolName:          utils.Ptr("icmp"),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *iaas.ApiCreateSecurityGroupRuleRequest)) iaas.ApiCreateSecurityGroupRuleRequest {
	request := testClient.CreateSecurityGroupRule(testCtx, testProjectId, testSecurityGroupId)
	request = request.CreateSecurityGroupRulePayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixtureRequiredRequest(mods ...func(request *iaas.ApiCreateSecurityGroupRuleRequest)) iaas.ApiCreateSecurityGroupRuleRequest {
	request := testClient.CreateSecurityGroupRule(testCtx, testProjectId, testSecurityGroupId)
	request = request.CreateSecurityGroupRulePayload(iaas.CreateSecurityGroupRulePayload{
		Direction: utils.Ptr("ingress"),
	})
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayload(mods ...func(payload *iaas.CreateSecurityGroupRulePayload)) iaas.CreateSecurityGroupRulePayload {
	payload := iaas.CreateSecurityGroupRulePayload{
		Direction:   utils.Ptr("ingress"),
		Description: utils.Ptr("example-description"),
		Ethertype:   utils.Ptr("ether"),
		IcmpParameters: &iaas.ICMPParameters{
			Code: utils.Ptr(int64(0)),
			Type: utils.Ptr(int64(8)),
		},
		IpRange: utils.Ptr("10.1.2.3"),
		PortRange: &iaas.PortRange{
			Max: utils.Ptr(int64(24)),
			Min: utils.Ptr(int64(22)),
		},
		Protocol: &iaas.CreateProtocol{
			Int64:  utils.Ptr(int64(1)),
			String: utils.Ptr("icmp"),
		},
		RemoteSecurityGroupId: utils.Ptr(testRemoteSecurityGroupId),
	}
	for _, mod := range mods {
		mod(&payload)
	}
	return payload
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description   string
		flagValues    map[string]string
		isValid       bool
		expectedModel *inputModel
	}{
		{
			description: "base",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, portRangeMaxFlag)
				delete(flagValues, portRangeMinFlag)
				delete(flagValues, protocolNumberFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.PortRangeMax = nil
				model.PortRangeMin = nil
				model.ProtocolNumber = nil
			}),
		},
		{
			description: "required only",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, descriptionFlag)
				delete(flagValues, etherTypeFlag)
				delete(flagValues, icmpParameterCodeFlag)
				delete(flagValues, icmpParameterTypeFlag)
				delete(flagValues, ipRangeFlag)
				delete(flagValues, portRangeMaxFlag)
				delete(flagValues, portRangeMinFlag)
				delete(flagValues, remoteSecurityGroupIdFlag)
				delete(flagValues, protocolNumberFlag)
				delete(flagValues, protocolNameFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Description = nil
				model.EtherType = nil
				model.IcmpParameterCode = nil
				model.IcmpParameterType = nil
				model.IpRange = nil
				model.PortRangeMax = nil
				model.PortRangeMin = nil
				model.RemoteSecurityGroupId = nil
				model.ProtocolNumber = nil
				model.ProtocolName = nil
			}),
		},
		{
			description: "direction missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, directionFlag)
			}),
			isValid: false,
		},
		{
			description: "protocol is icmp and parameters are missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, icmpParameterCodeFlag)
				delete(flagValues, icmpParameterTypeFlag)
			}),
			isValid: false,
		},
		{
			description: "protocol is icmp and port range values are provided",
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "protocol is not icmp and port range values are provided",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[protocolNameFlag] = "not-icmp"
				delete(flagValues, icmpParameterCodeFlag)
				delete(flagValues, icmpParameterTypeFlag)
				delete(flagValues, protocolNumberFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.IcmpParameterCode = nil
				model.IcmpParameterType = nil
				model.ProtocolName = utils.Ptr("not-icmp")
				model.ProtocolNumber = nil
			}),
		},
		{
			description: "protocol is not icmp and icmp parameters are provided",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[protocolNameFlag] = "not-icmp"
			}),
			isValid: false,
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
			description: "security group id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, securityGroupIdFlag)
			}),
			isValid: false,
		},
		{
			description: "security group id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[securityGroupIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "security group id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[securityGroupIdFlag] = "invalid-uuid"
			}),
			isValid: false,
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

			err = cmd.ValidateFlagGroups()
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flag groups: %v", err)
			}
			err = cmd.ValidateRequiredFlags()
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flags: %v", err)
			}

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
			diff := cmp.Diff(model, tt.expectedModel)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestBuildRequest(t *testing.T) {
	var tests = []struct {
		description     string
		model           *inputModel
		expectedRequest iaas.ApiCreateSecurityGroupRuleRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "only direction and security group id in payload",
			model: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Verbosity: globalflags.VerbosityDefault,
				},
				Direction:       utils.Ptr("ingress"),
				SecurityGroupId: testSecurityGroupId,
			},
			expectedRequest: fixtureRequiredRequest(),
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
