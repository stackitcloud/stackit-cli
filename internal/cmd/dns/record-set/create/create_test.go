package create

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/dns"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &dns.APIClient{}
var testProjectId = uuid.NewString()
var testZoneId = uuid.NewString()

var recordTxtOver255Char = []string{
	"foobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoo",
	"foobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoo",
	"foobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobar",
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		zoneIdFlag:                testZoneId,
		commentFlag:               "comment",
		nameFlag:                  "example.com",
		recordFlag:                "1.1.1.1",
		ttlFlag:                   "3600",
		typeFlag:                  "SOA", // Non-default value
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
		ZoneId:  testZoneId,
		Name:    utils.Ptr("example.com"),
		Comment: utils.Ptr("comment"),
		Records: []string{"1.1.1.1"},
		TTL:     utils.Ptr(int64(3600)),
		Type:    "SOA",
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *dns.ApiCreateRecordSetRequest)) dns.ApiCreateRecordSetRequest {
	request := testClient.CreateRecordSet(testCtx, testProjectId, testZoneId)
	request = request.CreateRecordSetPayload(dns.CreateRecordSetPayload{
		Name:    utils.Ptr("example.com"),
		Comment: utils.Ptr("comment"),
		Records: &[]dns.RecordPayload{
			{Content: utils.Ptr("1.1.1.1")},
		},
		Ttl:  utils.Ptr(int64(3600)),
		Type: dns.CREATERECORDSETPAYLOADTYPE_SOA.Ptr(),
	})
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func TestParseInput(t *testing.T) {
	var tests = []struct {
		description      string
		argValues        []string
		flagValues       map[string]string
		recordFlagValues []string
		isValid          bool
		expectedModel    *inputModel
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
			description: "required fields only",
			flagValues: map[string]string{
				globalflags.ProjectIdFlag: testProjectId,
				zoneIdFlag:                testZoneId,
				nameFlag:                  "example.com",
				recordFlag:                "1.1.1.1",
			},
			isValid: true,
			expectedModel: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Verbosity: globalflags.VerbosityDefault,
				},
				ZoneId:  testZoneId,
				Name:    utils.Ptr("example.com"),
				Records: []string{"1.1.1.1"},
				Type:    defaultType,
			},
		},
		{
			description: "zero values",
			flagValues: map[string]string{
				globalflags.ProjectIdFlag: testProjectId,
				zoneIdFlag:                testZoneId,
				commentFlag:               "",
				nameFlag:                  "",
				recordFlag:                "1.1.1.1",
				ttlFlag:                   "0",
			},
			isValid: true,
			expectedModel: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Verbosity: globalflags.VerbosityDefault,
				},
				ZoneId:  testZoneId,
				Name:    utils.Ptr(""),
				Comment: utils.Ptr(""),
				Records: []string{"1.1.1.1"},
				TTL:     utils.Ptr(int64(0)),
				Type:    defaultType,
			},
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
			description: "zone id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, zoneIdFlag)
			}),
			isValid: false,
		},
		{
			description: "zone id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[zoneIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "zone id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[zoneIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "name missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, nameFlag)
			}),
			isValid: false,
		},
		{
			description: "records missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, recordFlag)
			}),
			isValid: false,
		},
		{
			description: "type missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, typeFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Type = defaultType
			}),
		},
		{
			description: "type invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[typeFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "type invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[typeFlag] = "a"
			}),
			isValid: false,
		},
		{
			description:      "repeated primary flags",
			flagValues:       fixtureFlagValues(),
			recordFlagValues: []string{"1.2.3.4", "5.6.7.8"},
			isValid:          true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Records = append(model.Records, "1.2.3.4", "5.6.7.8")
			}),
		},
		{
			description:      "repeated primary flags with list value",
			flagValues:       fixtureFlagValues(),
			recordFlagValues: []string{"1.2.3.4,5.6.7.8"},
			isValid:          true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Records = append(model.Records, "1.2.3.4", "5.6.7.8")
			}),
		},
		{
			description: "TXT record with > 255 characters",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[typeFlag] = string(txtType)
				flagValues[recordFlag] = strings.Join(recordTxtOver255Char, "")
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				var content string
				for idx, val := range recordTxtOver255Char {
					content += fmt.Sprintf("%q", val)
					if idx != len(recordTxtOver255Char)-1 {
						content += " "
					}
				}

				model.Records = []string{content}
				model.Type = txtType
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInputWithAdditionalFlags(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, map[string][]string{
				recordFlag: tt.recordFlagValues,
			}, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest dns.ApiCreateRecordSetRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "required fields only",
			model: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Verbosity: globalflags.VerbosityDefault,
				},
				ZoneId:  testZoneId,
				Name:    utils.Ptr("example.com"),
				Records: []string{"1.1.1.1"},
				Type:    defaultType,
			},
			expectedRequest: testClient.CreateRecordSet(testCtx, testProjectId, testZoneId).
				CreateRecordSetPayload(dns.CreateRecordSetPayload{
					Name: utils.Ptr("example.com"),
					Records: &[]dns.RecordPayload{
						{Content: utils.Ptr("1.1.1.1")},
					},
					Type: utils.Ptr(defaultType),
				}),
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

func TestOutputResult(t *testing.T) {
	type args struct {
		model     *inputModel
		zoneLabel string
		resp      *dns.RecordSetResponse
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
			name: "only record set as argument",
			args: args{
				model: fixtureInputModel(),
				resp:  &dns.RecordSetResponse{Rrset: &dns.RecordSet{}},
			},
			wantErr: false,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&params.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.model, tt.args.zoneLabel, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
