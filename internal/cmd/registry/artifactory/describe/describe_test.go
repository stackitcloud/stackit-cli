package describe

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	registry "github.com/stackitcloud/stackit-sdk-go/services/registry/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
)

const (
	projectIdFlag = globalflags.ProjectIdFlag
	testRegion    = "eu01"
)

type testCtxKey struct{}

var (
	testCtx           = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient        = &registry.APIClient{DefaultAPI: &registry.DefaultAPIService{}}
	testProjectId     = uuid.NewString()
	testArtifactoryId = uuid.NewString()
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:          testProjectId,
		globalflags.RegionFlag: testRegion,
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
			Region:    testRegion,
		},
		ArtifactoryId: testArtifactoryId,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *registry.ApiGetArtifactoryRequest)) registry.ApiGetArtifactoryRequest {
	request := testClient.DefaultAPI.GetArtifactory(testCtx, testProjectId, testRegion, testArtifactoryId)
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
			argValues:     []string{testArtifactoryId},
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "project id missing",
			argValues:   []string{testArtifactoryId},
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, projectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			argValues:   []string{testArtifactoryId},
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			argValues:   []string{testArtifactoryId},
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = "invalid-uuid"
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
		expectedRequest registry.ApiGetArtifactoryRequest
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
				cmp.AllowUnexported(tt.expectedRequest, registry.DefaultAPIService{}),
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
		artifactory *registry.Artifactory
	}
	desc := "docker hub mirror"
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "single item",
			args: args{
				artifactory: &registry.Artifactory{
					Id:                 testArtifactoryId,
					Name:               "my-artifactory",
					State:              registry.ARTIFACTORYSTATE_ACTIVE,
					Url:                "https://example.com",
					ConsumedQuotaBytes: 1024,
					CreateTime:         time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "with proxy",
			args: args{
				artifactory: &registry.Artifactory{
					Id:                 testArtifactoryId,
					Name:               "my-artifactory",
					State:              registry.ARTIFACTORYSTATE_ACTIVE,
					Url:                "https://example.com",
					ConsumedQuotaBytes: 1024,
					CreateTime:         time.Now(),
					Proxy: &registry.Proxy{
						Url:          "https://registry-1.docker.io",
						Description:  &desc,
						RegistryType: "docker-hub",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "nil artifactory",
			args: args{
				artifactory: nil,
			},
			wantErr: true,
		},
	}
	params := testparams.NewTestParams()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, "", tt.args.artifactory); (err != nil) != tt.wantErr {
				t.Errorf("TestOutputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
