package create

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
	testName              = "my-registry"
	testProxyUrl          = "https://registry-1.docker.io"
	testProxyRegistryType = "docker-hub"
	testProxyDescription  = "Docker Hub proxy cache"
	testProxyAccessKey    = "my-user"
	testProxyAccessSecret = "my-secret"
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:          testProjectId,
		globalflags.RegionFlag: testRegion,
		nameFlag:               testName,
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
		Name: testName,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *registry.ApiCreateArtifactoryRequest)) registry.ApiCreateArtifactoryRequest {
	request := testClient.DefaultAPI.CreateArtifactory(testCtx, testProjectId, testRegion).
		CreateArtifactoryPayload(registry.CreateArtifactoryPayload{
			Name: testName,
		})
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description   string
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
			description: "with proxy minimal",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[proxyUrlFlag] = testProxyUrl
				flagValues[proxyRegistryTypeFlag] = testProxyRegistryType
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.ProxyUrl = utils.Ptr(testProxyUrl)
				model.ProxyRegistryType = utils.Ptr(testProxyRegistryType)
			}),
		},
		{
			description: "with proxy full",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[proxyUrlFlag] = testProxyUrl
				flagValues[proxyRegistryTypeFlag] = testProxyRegistryType
				flagValues[proxyDescriptionFlag] = testProxyDescription
				flagValues[proxyAccessKeyFlag] = testProxyAccessKey
				flagValues[proxyAccessSecretFlag] = testProxyAccessSecret
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.ProxyUrl = utils.Ptr(testProxyUrl)
				model.ProxyRegistryType = utils.Ptr(testProxyRegistryType)
				model.ProxyDescription = utils.Ptr(testProxyDescription)
				model.ProxyAccessKey = utils.Ptr(testProxyAccessKey)
				model.ProxyAccessSecret = utils.Ptr(testProxyAccessSecret)
			}),
		},
		{
			description: "proxy url without registry type",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[proxyUrlFlag] = testProxyUrl
			}),
			isValid: false,
		},
		{
			description: "proxy registry type without url",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[proxyRegistryTypeFlag] = testProxyRegistryType
			}),
			isValid: false,
		},
		{
			description: "proxy access key without secret",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[proxyUrlFlag] = testProxyUrl
				flagValues[proxyRegistryTypeFlag] = testProxyRegistryType
				flagValues[proxyAccessKeyFlag] = testProxyAccessKey
			}),
			isValid: false,
		},
		{
			description: "proxy description without proxy url",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[proxyDescriptionFlag] = testProxyDescription
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
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, []string{}, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest registry.ApiCreateArtifactoryRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "with proxy full",
			model: fixtureInputModel(func(model *inputModel) {
				model.ProxyUrl = utils.Ptr(testProxyUrl)
				model.ProxyRegistryType = utils.Ptr(testProxyRegistryType)
				model.ProxyDescription = utils.Ptr(testProxyDescription)
				model.ProxyAccessKey = utils.Ptr(testProxyAccessKey)
				model.ProxyAccessSecret = utils.Ptr(testProxyAccessSecret)
			}),
			expectedRequest: testClient.DefaultAPI.CreateArtifactory(testCtx, testProjectId, testRegion).
				CreateArtifactoryPayload(registry.CreateArtifactoryPayload{
					Name: testName,
					Proxy: &registry.CreateProxyRequest{
						Url:          testProxyUrl,
						RegistryType: testProxyRegistryType,
						Description:  utils.Ptr(testProxyDescription),
						Credentials: &registry.Credentials{
							AccessKey:    testProxyAccessKey,
							AccessSecret: testProxyAccessSecret,
						},
					},
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
		model        *inputModel
		projectLabel string
		id           string
		resp         *registry.Artifactory
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty response",
			args: args{
				model: fixtureInputModel(),
				resp:  nil,
			},
			wantErr: true,
		},
		{
			name: "created successfully sync",
			args: args{
				model:        fixtureInputModel(),
				projectLabel: "example-project",
				id:           "xxx",
				resp: &registry.Artifactory{
					Id: "xxx",
				},
			},
			wantErr: false,
		},
		{
			name: "created successfully async",
			args: args{
				model: fixtureInputModel(func(model *inputModel) {
					model.Async = true
				}),
				projectLabel: "example-project",
				id:           "xxx",
				resp: &registry.Artifactory{
					Id: "xxx",
				},
			},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.model, tt.args.projectLabel, tt.args.id, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("TestOutputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
