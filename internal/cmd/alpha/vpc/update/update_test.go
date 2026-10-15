package update

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	iaas "github.com/stackitcloud/stackit-sdk-go/services/iaas/v2alpha1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
)

type testCtxKey struct{}

var (
	testCtx    = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient = &iaas.APIClient{DefaultAPI: &iaas.DefaultAPIService{}}

	testProjectId   = uuid.NewString()
	testVpcId       = uuid.NewString()
	testName        = "test-vpc"
	testDescription = "test-vpc-description"
	testLabels      = "environment=dev,endpoint=alpha"
)

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testVpcId,
	}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		nameFlag:                  testName,
		descriptionFlag:           testDescription,
		labelsFlag:                testLabels,
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
		VpcId:       testVpcId,
		Name:        new(testName),
		Description: new(testDescription),
		Labels: map[string]any{
			"environment": "dev",
			"endpoint":    "alpha",
		},
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixturePayload(mods ...func(payload *iaas.PartialUpdateVPCPayload)) (payload iaas.PartialUpdateVPCPayload) {
	payload = iaas.PartialUpdateVPCPayload{
		Name:        new(testName),
		Description: new(testDescription),
		Labels:      map[string]any{"environment": "dev", "endpoint": "alpha"},
	}
	for _, f := range mods {
		f(&payload)
	}
	return payload
}

func fixtureRequest(mods ...func(request *iaas.ApiPartialUpdateVPCRequest)) iaas.ApiPartialUpdateVPCRequest {
	request := testClient.DefaultAPI.PartialUpdateVPC(testCtx, testProjectId, testVpcId)

	request = request.PartialUpdateVPCPayload(fixturePayload())
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
			description: "no values - nothing to update",
			argValues:   fixtureArgValues(),
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "required fields only",
			argValues:   fixtureArgValues(),
			flagValues: map[string]string{
				globalflags.ProjectIdFlag: testProjectId,
				nameFlag:                  testName,
			},
			isValid: true,
			expectedModel: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Verbosity: globalflags.VerbosityDefault,
				},
				VpcId: testVpcId,
				Name:  new(testName),
			},
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
		expectedRequest iaas.ApiPartialUpdateVPCRequest
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
				cmp.AllowUnexported(tt.expectedRequest, iaas.DefaultAPIService{}),
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
		outputFormat string
		vpcLabel     string
		resp         *iaas.VPC
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
			name: "empty vpc response",
			args: args{
				vpcLabel: testVpcId,
				resp:     &iaas.VPC{},
			},
		},
	}
	params := testparams.NewTestParams()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.outputFormat, tt.args.vpcLabel, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
