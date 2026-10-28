package update

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	registry "github.com/stackitcloud/stackit-sdk-go/services/registry/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	projectIdFlag = globalflags.ProjectIdFlag
	testRegion    = "eu01"
)

type testCtxKey struct{}

var (
	testCtx               = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient            = &registry.APIClient{DefaultAPI: &registry.DefaultAPIService{}}
	testProjectId         = uuid.NewString()
	testArtifactoryId     = uuid.NewString()
	testProxyUrl          = "https://registry-1.docker.io"
	testProxyDescription  = "docker hub proxy"
	testProxyAccessKey    = "my-key"
	testProxyAccessSecret = "my-secret"
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:          testProjectId,
		globalflags.RegionFlag: testRegion,
		proxyUrlFlag:           testProxyUrl,
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
		ProxyUrl:      utils.Ptr(testProxyUrl),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *registry.ApiPatchArtifactoryRequest)) registry.ApiPatchArtifactoryRequest {
	request := testClient.DefaultAPI.PatchArtifactory(testCtx, testProjectId, testRegion, testArtifactoryId).
		PatchArtifactoryPayload(registry.PatchArtifactoryPayload{
			Proxy: &registry.UpdateProxyRequest{
				Url: utils.Ptr(testProxyUrl),
			},
		})
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
			description: "with description",
			argValues:   []string{testArtifactoryId},
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, proxyUrlFlag)
				flagValues[proxyDescriptionFlag] = testProxyDescription
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.ProxyUrl = nil
				model.ProxyDescription = utils.Ptr(testProxyDescription)
			}),
		},
		{
			description: "with credentials",
			argValues:   []string{testArtifactoryId},
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, proxyUrlFlag)
				flagValues[proxyAccessKeyFlag] = testProxyAccessKey
				flagValues[proxyAccessSecretFlag] = testProxyAccessSecret
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.ProxyUrl = nil
				model.ProxyAccessKey = utils.Ptr(testProxyAccessKey)
				model.ProxyAccessSecret = utils.Ptr(testProxyAccessSecret)
			}),
		},
		{
			description: "proxy access key without secret",
			argValues:   []string{testArtifactoryId},
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[proxyAccessKeyFlag] = testProxyAccessKey
			}),
			isValid: false,
		},
		{
			description: "proxy access secret without key",
			argValues:   []string{testArtifactoryId},
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[proxyAccessSecretFlag] = testProxyAccessSecret
			}),
			isValid: false,
		},
		{
			description: "no flags set",
			argValues:   []string{testArtifactoryId},
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, proxyUrlFlag)
			}),
			isValid: false,
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
		expectedRequest registry.ApiPatchArtifactoryRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "with description",
			model: fixtureInputModel(func(model *inputModel) {
				model.ProxyUrl = nil
				model.ProxyDescription = utils.Ptr(testProxyDescription)
			}),
			expectedRequest: fixtureRequest(func(request *registry.ApiPatchArtifactoryRequest) {
				*request = testClient.DefaultAPI.PatchArtifactory(testCtx, testProjectId, testRegion, testArtifactoryId).
					PatchArtifactoryPayload(registry.PatchArtifactoryPayload{
						Proxy: &registry.UpdateProxyRequest{
							Description: utils.Ptr(testProxyDescription),
						},
					})
			}),
		},
		{
			description: "with credentials",
			model: fixtureInputModel(func(model *inputModel) {
				model.ProxyUrl = nil
				model.ProxyAccessKey = utils.Ptr(testProxyAccessKey)
				model.ProxyAccessSecret = utils.Ptr(testProxyAccessSecret)
			}),
			expectedRequest: fixtureRequest(func(request *registry.ApiPatchArtifactoryRequest) {
				*request = testClient.DefaultAPI.PatchArtifactory(testCtx, testProjectId, testRegion, testArtifactoryId).
					PatchArtifactoryPayload(registry.PatchArtifactoryPayload{
						Proxy: &registry.UpdateProxyRequest{
							Credentials: &registry.Credentials{
								AccessKey:    testProxyAccessKey,
								AccessSecret: testProxyAccessSecret,
							},
						},
					})
			}),
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
		model            *inputModel
		artifactoryLabel string
		resp             *registry.Artifactory
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "updated successfully sync",
			args: args{
				model:            fixtureInputModel(),
				artifactoryLabel: "my-artifactory",
				resp: &registry.Artifactory{
					Id: testArtifactoryId,
				},
			},
			wantErr: false,
		},
		{
			name: "updated successfully async",
			args: args{
				model: fixtureInputModel(func(model *inputModel) {
					model.Async = true
				}),
				artifactoryLabel: "my-artifactory",
				resp: &registry.Artifactory{
					Id: testArtifactoryId,
				},
			},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.model, tt.args.artifactoryLabel, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("TestOutputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
