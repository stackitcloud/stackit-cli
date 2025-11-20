package create

import (
	"bytes"
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
	sdkUtils "github.com/stackitcloud/stackit-sdk-go/core/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/cdn"
	"k8s.io/utils/ptr"
)

var projectIdFlag = globalflags.ProjectIdFlag

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &cdn.APIClient{}
var testProjectId = uuid.NewString()

const testOriginURL = "https://example.com/somePath?foo=bar"
const testRegions = cdn.REGION_EU

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag: testProjectId,
		originURLFlag: testOriginURL,
		regionsFlag:   string(testRegions),
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func flagRegions(regions ...cdn.Region) func(flagValues map[string]string) {
	return func(flagValues map[string]string) {
		if len(regions) == 0 {
			delete(flagValues, regionsFlag)
			return
		}
		stringRegions := sdkUtils.EnumSliceToStringSlice(regions)
		flagValues[regionsFlag] = strings.Join(stringRegions, ",")
	}
}

func flagOriginURL(originURL string) func(flagValues map[string]string) {
	return func(flagValues map[string]string) {
		if originURL == "" {
			delete(flagValues, originURLFlag)
			return
		}
		flagValues[originURLFlag] = originURL
	}
}

func flagProjectID(id *string) func(flagValues map[string]string) {
	return func(flagValues map[string]string) {
		if id == nil {
			delete(flagValues, projectIdFlag)
			return
		}
		flagValues[projectIdFlag] = *id
	}
}

func fixtureModel(mods ...func(m *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Verbosity: globalflags.VerbosityDefault,
		},
		Regions:   []cdn.Region{testRegions},
		OriginURL: testOriginURL,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func modelRegions(regions ...cdn.Region) func(m *inputModel) {
	return func(m *inputModel) {
		m.Regions = regions
	}
}

func fixturePayload(mods ...func(p *cdn.CreateDistributionPayload)) cdn.CreateDistributionPayload {
	p := *cdn.NewCreateDistributionPayload(
		testOriginURL,
		[]cdn.Region{testRegions},
	)
	for _, mod := range mods {
		mod(&p)
	}
	return p
}

func payloadRegions(regions ...cdn.Region) func(p *cdn.CreateDistributionPayload) {
	return func(p *cdn.CreateDistributionPayload) {
		p.Regions = &regions
	}
}

func fixtureRequest(mods ...func(p *cdn.CreateDistributionPayload)) cdn.ApiCreateDistributionRequest {
	req := testClient.CreateDistribution(testCtx, testProjectId)
	req = req.CreateDistributionPayload(fixturePayload(mods...))
	return req
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description string
		argValues   []string
		flagValues  map[string]string
		isValid     bool
		expected    *inputModel
	}{
		{
			description: "base",
			flagValues:  fixtureFlagValues(),
			isValid:     true,
			expected:    fixtureModel(),
		},
		{
			description: "project id missing",
			flagValues:  fixtureFlagValues(flagProjectID(nil)),
			isValid:     false,
		},
		{
			description: "project id invalid 1",
			flagValues:  fixtureFlagValues(flagProjectID(utils.Ptr(""))),
			isValid:     false,
		},
		{
			description: "project id invalid 2",
			flagValues:  fixtureFlagValues(flagProjectID(utils.Ptr("invalid-uuid"))),
			isValid:     false,
		},
		{
			description: "origin url missing",
			flagValues:  fixtureFlagValues(flagOriginURL("")),
			isValid:     false,
		},
		{
			description: "origin url invalid",
			flagValues:  fixtureFlagValues(flagOriginURL("://invalid-url")),
			isValid:     false,
		},
		{
			description: "regions missing",
			flagValues:  fixtureFlagValues(flagRegions()),
			isValid:     false,
		},
		{
			description: "multiple regions",
			flagValues:  fixtureFlagValues(flagRegions(cdn.REGION_EU, cdn.REGION_AF)),
			isValid:     true,
			expected:    fixtureModel(modelRegions(cdn.REGION_EU, cdn.REGION_AF)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expected, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description string
		model       *inputModel
		expected    cdn.ApiCreateDistributionRequest
	}{
		{
			description: "base",
			model:       fixtureModel(),
			expected:    fixtureRequest(),
		},
		{
			description: "multiple regions",
			model:       fixtureModel(modelRegions(cdn.REGION_AF, cdn.REGION_EU)),
			expected:    fixtureRequest(payloadRegions(cdn.REGION_AF, cdn.REGION_EU)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)

			diff := cmp.Diff(request, tt.expected,
				cmp.AllowUnexported(tt.expected),
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
		outputFormat string
		response     *cdn.CreateDistributionResponse
		expected     string
		wantErr      bool
	}{
		{
			description:  "nil response",
			outputFormat: "table",
			response:     nil,
			wantErr:      true,
		},
		{
			description:  "table output",
			outputFormat: "table",
			response: &cdn.CreateDistributionResponse{
				Distribution: &cdn.Distribution{
					Id: ptr.To("dist-1234"),
				},
			},
			expected: fmt.Sprintf("Created CDN distribution for %q. Id: dist-1234\n", testProjectId),
		},
		{
			description:  "json output",
			outputFormat: "json",
			response: &cdn.CreateDistributionResponse{
				Distribution: &cdn.Distribution{
					Id: ptr.To("dist-1234"),
				},
			},
			expected: `{
  "distribution": {
    "config": null,
    "createdAt": null,
    "domains": null,
    "id": "dist-1234",
    "projectId": null,
    "status": null,
    "updatedAt": null
  }
}
`,
		},
	}

	p := print.NewPrinter()
	p.Cmd = NewCmd(&params.CmdParams{Printer: p})

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			p.Cmd.SetOut(buffer)
			if err := outputResult(p, tt.outputFormat, testProjectId, tt.response); (err != nil) != tt.wantErr {
				t.Fatalf("outputResult: %v", err)
			}
			if buffer.String() != tt.expected {
				t.Errorf("want:\n%s\ngot:\n%s", tt.expected, buffer.String())
			}
		})
	}
}
