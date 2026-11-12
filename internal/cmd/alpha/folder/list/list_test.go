package list

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/zalando/go-keyring"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	resourcemanager "github.com/stackitcloud/stackit-sdk-go/services/resourcemanager/v0api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &resourcemanager.APIClient{DefaultAPI: &resourcemanager.DefaultAPIServiceMock{}}
var testParentId = uuid.NewString()
var testProjectIdLike = uuid.NewString()

const testCreationTimeAfter = "2023-01-01T00:00:00Z"

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		parentIdFlag:          testParentId,
		memberFlag:            "member",
		creationTimeAfterFlag: testCreationTimeAfter,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	testCreationTimeAfter, err := time.Parse(creationTimeAfterFormat, testCreationTimeAfter)
	if err != nil {
		return &inputModel{}
	}

	model := &inputModel{
		GlobalFlagModel:   &globalflags.GlobalFlagModel{Verbosity: globalflags.VerbosityDefault},
		ParentId:          utils.Ptr(testParentId),
		Member:            utils.Ptr("member"),
		CreationTimeAfter: utils.Ptr(testCreationTimeAfter),
		PageSize:          pageSizeDefault,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *resourcemanager.ApiListFoldersRequest)) resourcemanager.ApiListFoldersRequest {
	testCreationTimeAfter, err := time.Parse(creationTimeAfterFormat, testCreationTimeAfter)
	if err != nil {
		return resourcemanager.ApiListFoldersRequest{}
	}
	request := testClient.DefaultAPI.ListFolders(testCtx)
	request = request.CreationTimeStart(testCreationTimeAfter)
	request = request.Member("member")
	request = request.Limit(pageSizeDefault)
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description         string
		argValues           []string
		flagValues          map[string]string
		projectIdLikevalues *[]string
		isValid             bool
		expectedModel       *inputModel
	}{
		{
			description:   "base",
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "parentId empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[parentIdFlag] = ""
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.ParentId = utils.Ptr("")
			}),
		},
		{
			description: "member empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[memberFlag] = ""
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Member = utils.Ptr("")
			}),
		},
		{
			description:         "projectIdLike one value",
			flagValues:          fixtureFlagValues(),
			projectIdLikevalues: utils.Ptr([]string{testProjectIdLike}),
			isValid:             true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.ProjectIdLike = []string{testProjectIdLike}
			}),
		},
		{
			description:         "projectIdLike multiple values",
			flagValues:          fixtureFlagValues(),
			projectIdLikevalues: utils.Ptr([]string{testProjectIdLike, testProjectIdLike}),
			isValid:             true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.ProjectIdLike = []string{testProjectIdLike, testProjectIdLike}
			}),
		},
		{
			description:         "projectIdLike empty",
			flagValues:          fixtureFlagValues(),
			projectIdLikevalues: utils.Ptr([]string{}),
			isValid:             true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.ProjectIdLike = nil
			}),
		},
		{
			description: "no values",
			flagValues:  map[string]string{},
			isValid:     true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.ParentId = nil
				model.Member = nil
				model.CreationTimeAfter = nil
			}),
		},
		{
			description:         "projectIdLike invalid",
			flagValues:          fixtureFlagValues(),
			projectIdLikevalues: utils.Ptr([]string{""}),
			isValid:             false,
		},
		{
			description: "no values",
			flagValues:  map[string]string{},
			isValid:     true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.ParentId = nil
				model.Member = nil
				model.CreationTimeAfter = nil
			}),
		},
		{
			description: "limit invalid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[limitFlag] = "invalid"
			}),
			isValid: false,
		},
		{
			description: "limit invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[limitFlag] = "0"
			}),
			isValid: false,
		},
		{
			description: "creationTimeAfter empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[creationTimeAfterFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "creationTimeAfter invalid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[creationTimeAfterFlag] = "test"
			}),
			isValid: false,
		},
		{
			description: "creationTimeAfter invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[creationTimeAfterFlag] = "11:00 12/12/2023"
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInputWithAdditionalFlags(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, map[string][]string{
				projectIdLikeFlag: utils.GetSliceFromPointer(tt.projectIdLikevalues),
			}, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	keyring.MockInit()
	err := auth.SetAuthFlow(auth.AUTH_FLOW_USER_TOKEN)
	if err != nil {
		t.Fatalf("Failed to set auth flow: %v", err)
	}

	err = auth.SetAuthField(auth.USER_EMAIL, "test@test.com")
	if err != nil {
		t.Fatalf("Failed to set auth user email: %v", err)
	}

	authUserEmail, err := auth.GetAuthField(auth.USER_EMAIL)
	if err != nil {
		t.Fatalf("Failed to get auth user email: %v", err)
	}

	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest resourcemanager.ApiListFoldersRequest
	}{
		{
			description: "base",
			model: fixtureInputModel(func(model *inputModel) {
				model.PageSize = 1
			}),
			expectedRequest: fixtureRequest().Limit(1).ContainerParentId(testParentId),
		},
		{
			description: "base 2",
			model: fixtureInputModel(func(model *inputModel) {
				model.PageSize = 10
			}),
			expectedRequest: fixtureRequest().Limit(10).ContainerParentId(testParentId),
		},
		{
			description: "fetch email from auth user",
			model: fixtureInputModel(func(model *inputModel) {
				model.ParentId = nil
				model.Member = nil
			}),
			expectedRequest: fixtureRequest().Limit(pageSizeDefault).Member(authUserEmail),
		},
		{
			description: "projectIdLike set",
			model: fixtureInputModel(func(model *inputModel) {
				model.ProjectIdLike = []string{testProjectIdLike}
			}),
			expectedRequest: fixtureRequest().ContainerIds([]string{testProjectIdLike}).ContainerParentId(testParentId),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request, err := buildRequest(testCtx, tt.model, resourcemanager.DefaultAPIServiceMock{})
			if err != nil {
				t.Fatalf("Failed to build request: %v", err)
			}

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx, resourcemanager.DefaultAPIService{}),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

type mockSettings struct {
	listFolderFails bool
	totalItems      int
	served          int // items returned so far (simulated server-side offset)
	numApiCalls     int
}

// newApiMock simulates an offset/limit API that returns pages of
// pageSizeDefault items, and an empty page once all items are served.
func newApiMock(s *mockSettings) resourcemanager.DefaultAPI {
	return resourcemanager.DefaultAPIServiceMock{
		ListFoldersExecuteMock: utils.Ptr(func(_ resourcemanager.ApiListFoldersRequest) (*resourcemanager.ListFoldersResponse, error) {
			s.numApiCalls++
			if s.listFolderFails {
				return nil, fmt.Errorf("failed to list folders")
			}

			n := max(0, min(int(pageSizeDefault), s.totalItems-s.served))
			s.served += n

			return &resourcemanager.ListFoldersResponse{
				Items: make([]resourcemanager.ListFoldersResponseItemsInner, n),
			}, nil
		}),
	}
}

func TestFetchFolders(t *testing.T) {
	withLimit := func(limit int64) *inputModel {
		return fixtureInputModel(func(model *inputModel) {
			model.Limit = utils.Ptr(limit)
		})
	}

	tests := []struct {
		description         string
		model               *inputModel
		testSettings        mockSettings
		expectedNumAPICalls int
		expectedNumItems    int
	}{
		{
			description:         "no limit, totalItems<pageSize",
			model:               fixtureInputModel(),
			testSettings:        mockSettings{totalItems: 10},
			expectedNumAPICalls: 1,
			expectedNumItems:    10,
		},
		{
			description:         "no limit, totalItems>pageSize",
			model:               fixtureInputModel(),
			testSettings:        mockSettings{totalItems: 170},
			expectedNumAPICalls: 4, // 50 + 50 + 50 + 20
			expectedNumItems:    170,
		},
		{
			description:         "no limit, totalItems multiple of pageSize",
			model:               fixtureInputModel(),
			testSettings:        mockSettings{totalItems: 100},
			expectedNumAPICalls: 3, // last call returns no items
			expectedNumItems:    100,
		},
		{
			description:         "no limit, totalItems=pageSize",
			model:               fixtureInputModel(),
			testSettings:        mockSettings{totalItems: 50},
			expectedNumAPICalls: 2, // last call returns no items
			expectedNumItems:    50,
		},
		{
			description:         "limit<pageSize",
			model:               withLimit(10),
			testSettings:        mockSettings{totalItems: 50},
			expectedNumAPICalls: 1,
			expectedNumItems:    10,
		},
		{
			description:         "limit=totalItems=pageSize",
			model:               withLimit(50),
			testSettings:        mockSettings{totalItems: 50},
			expectedNumAPICalls: 1, // limit reached, no extra empty call
			expectedNumItems:    50,
		},
		{
			description:         "limit spans pages",
			model:               withLimit(70),
			testSettings:        mockSettings{totalItems: 170},
			expectedNumAPICalls: 2, // 50 + 20
			expectedNumItems:    70,
		},
		{
			description:         "limit>totalItems, totalItems<pageSize",
			model:               withLimit(200),
			testSettings:        mockSettings{totalItems: 30},
			expectedNumAPICalls: 1,
			expectedNumItems:    30,
		},
		{
			description:         "limit>totalItems, totalItems>pageSize",
			model:               withLimit(200),
			testSettings:        mockSettings{totalItems: 120},
			expectedNumAPICalls: 3, // 50 + 50 + 20
			expectedNumItems:    120,
		},
		{
			description:  "request fails",
			model:        fixtureInputModel(),
			testSettings: mockSettings{listFolderFails: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			settings := tt.testSettings

			folders, err := fetchFolders(testCtx, tt.model, newApiMock(&settings))

			if settings.listFolderFails {
				if err == nil {
					t.Fatalf("expected an error, got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if settings.numApiCalls != tt.expectedNumAPICalls {
				t.Fatalf("expected %d API calls, got %d", tt.expectedNumAPICalls, settings.numApiCalls)
			}
			if len(folders) != tt.expectedNumItems {
				t.Fatalf("expected %d folders, got %d", tt.expectedNumItems, len(folders))
			}
		})
	}
}

func Test_outputResult(t *testing.T) {
	type args struct {
		outputFormat string
		folders      []resourcemanager.ListFoldersResponseItemsInner
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty",
			args:    args{},
			wantErr: false,
		},
		{
			name: "base",
			args: args{
				"",
				[]resourcemanager.ListFoldersResponseItemsInner{{}},
			},
			wantErr: false,
		},
		{
			name: "complete",
			args: args{"", []resourcemanager.ListFoldersResponseItemsInner{
				{
					ContainerId:  "container-id1",
					CreationTime: time.Now(),
					UpdateTime:   time.Now(),
					Labels:       &map[string]string{"foo": "bar"},
					Name:         "some name",
					Parent: resourcemanager.Parent{
						Id: "parent-id",
					},
					FolderId: "folder-id1",
				},
				{
					ContainerId:  "container-id2",
					CreationTime: time.Now(),
					UpdateTime:   time.Now(),
					Labels:       &map[string]string{"foo": "bar"},
					Name:         "some name",
					Parent: resourcemanager.Parent{
						Id: "parent-id",
					},
					FolderId: "folder-id2",
				},
			}},
			wantErr: false},
	}

	params := testparams.NewTestParams()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.outputFormat, tt.args.folders); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
