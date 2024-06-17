package list

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/zalando/go-keyring"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	sdkConfig "github.com/stackitcloud/stackit-sdk-go/core/config"
	"github.com/stackitcloud/stackit-sdk-go/services/resourcemanager"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &resourcemanager.APIClient{}
var testParentId = uuid.NewString()
var testProjectIdLike = uuid.NewString()
var testCreationTimeAfter = "2023-01-01T00:00:00Z"

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

func fixtureRequest(mods ...func(request *resourcemanager.ApiListProjectsRequest)) resourcemanager.ApiListProjectsRequest {
	request := testClient.ListProjects(testCtx)
	request = request.ContainerParentId(testParentId)

	testCreationTimeAfter, err := time.Parse(creationTimeAfterFormat, testCreationTimeAfter)
	if err != nil {
		return resourcemanager.ApiListProjectsRequest{}
	}
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

			if tt.projectIdLikevalues != nil {
				for _, value := range *tt.projectIdLikevalues {
					err := cmd.Flags().Set(projectIdLikeFlag, value)
					if err != nil {
						if !tt.isValid {
							return
						}
						t.Fatalf("setting flag --%s=%s: %v", projectIdLikeFlag, value, err)
					}
				}
			}

			err = cmd.ValidateRequiredFlags()
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flags: %v", err)
			}

			err = cmd.ValidateFlagGroups()
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating one of required flags: %v", err)
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
	keyring.MockInit()
	err := auth.SetAuthField(auth.USER_EMAIL, "test@test.com")
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
		projectIdLike   []string
		offset          int
		expectedRequest resourcemanager.ApiListProjectsRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			offset:          1,
			expectedRequest: fixtureRequest().Offset(1),
		},
		{
			description:     "base 2",
			model:           fixtureInputModel(),
			offset:          10,
			expectedRequest: fixtureRequest().Offset(10),
		},
		{
			description: "fetch email from auth user",
			model: &inputModel{
				PageSize: pageSizeDefault,
			},
			offset:          1,
			expectedRequest: testClient.ListProjects(testCtx).Offset(1).Limit(pageSizeDefault).Member(authUserEmail),
		},
		{
			description:     "projectIdLike set",
			model:           fixtureInputModel(),
			projectIdLike:   []string{testProjectIdLike},
			offset:          0,
			expectedRequest: fixtureRequest().Offset(0).ContainerIds([]string{testProjectIdLike}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			if tt.projectIdLike != nil {
				tt.model.ProjectIdLike = tt.projectIdLike
			}
			request, err := buildRequest(testCtx, tt.model, testClient, tt.offset)
			if err != nil {
				t.Fatalf("Failed to build request: %v", err)
			}

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

func TestFetchProjects(t *testing.T) {
	tests := []struct {
		description         string
		model               *inputModel
		totalItems          int
		apiCallFails        bool
		expectedNumAPICalls int
		expectedNumItems    int
	}{
		{
			description:         "no limit and pageSize>totalItems",
			model:               fixtureInputModel(),
			totalItems:          10,
			expectedNumAPICalls: 1,
			apiCallFails:        false,
			expectedNumItems:    10,
		},
		{
			description:         "no limit and pageSize<totalItems",
			model:               fixtureInputModel(),
			totalItems:          170,
			expectedNumAPICalls: 4,
			apiCallFails:        false,
			expectedNumItems:    170,
		},
		{
			description:         "no limit and pageSize<totalItems 2",
			model:               fixtureInputModel(),
			totalItems:          100,
			expectedNumAPICalls: 3, // Last call will return no items
			apiCallFails:        false,
			expectedNumItems:    100,
		},
		{
			description:         "no limit and pageSize=totalItems",
			model:               fixtureInputModel(),
			totalItems:          50,
			expectedNumAPICalls: 2, // Last call will return no items
			apiCallFails:        false,
			expectedNumItems:    50,
		},
		{
			description: "limit<pageSize",
			model: fixtureInputModel(func(model *inputModel) {
				model.Limit = utils.Ptr(int64(10))
			}),
			totalItems:          50,
			expectedNumAPICalls: 1,
			apiCallFails:        false,
			expectedNumItems:    10,
		},
		{
			description: "limit>totalItems and pageSize>totalItems",
			model: fixtureInputModel(func(model *inputModel) {
				model.Limit = utils.Ptr(int64(200))
				model.PageSize = 300
			}),
			totalItems:          50,
			expectedNumAPICalls: 1,
			apiCallFails:        false,
			expectedNumItems:    50,
		},
		{
			description: "limit>totalItems and pageSize<totalItems",
			model: fixtureInputModel(func(model *inputModel) {
				model.Limit = utils.Ptr(int64(200))
				model.PageSize = 30
			}),
			totalItems:          50,
			expectedNumAPICalls: 2,
			apiCallFails:        false,
			expectedNumItems:    50,
		},
		{
			description:  "request fails",
			model:        fixtureInputModel(),
			totalItems:   100,
			apiCallFails: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			numAPICalls := 0
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				numAPICalls++

				w.Header().Set("Content-Type", "application/json")
				if tt.apiCallFails {
					w.WriteHeader(http.StatusInternalServerError)
					_, err := w.Write([]byte("{\"message\": \"Something bad happened\""))
					if err != nil {
						t.Errorf("Failed to write bad response: %v", err)
					}
					return
				}

				query := r.URL.Query()
				offsetStr := query.Get("offset")
				if offsetStr == "" {
					t.Errorf("Expected offset param page to be set")
				}
				offset, err := strconv.Atoi(offsetStr)
				if err != nil {
					t.Errorf("Failed to parse query param offset: %v", err)
				}
				limitStr := query.Get("limit")
				if limitStr == "" {
					t.Errorf("Expected query param limit to be set")
				}
				limit, err := strconv.Atoi(limitStr)
				if err != nil {
					t.Errorf("Failed to parse query param limit: %v", err)
				}

				var numItemsToReturn int
				if offset >= tt.totalItems {
					numItemsToReturn = 0 // Total items reached
				} else if offset+limit < tt.totalItems {
					numItemsToReturn = limit // Full intermediate page
				} else {
					numItemsToReturn = tt.totalItems - offset // Last page
				}

				projects := make([]resourcemanager.Project, numItemsToReturn)
				mockedResp := resourcemanager.ListProjectsResponse{
					Items: &projects,
				}

				mockedRespBytes, err := json.Marshal(mockedResp)
				if err != nil {
					t.Fatalf("Failed to marshal mocked response: %v", err)
				}

				_, err = w.Write(mockedRespBytes)
				if err != nil {
					t.Errorf("Failed to write response: %v", err)
				}
			})
			mockedServer := httptest.NewServer(handler)
			defer mockedServer.Close()
			client, err := resourcemanager.NewAPIClient(
				sdkConfig.WithEndpoint(mockedServer.URL),
				sdkConfig.WithoutAuthentication(),
			)
			if err != nil {
				t.Fatalf("Failed to initialize client: %v", err)
			}

			projects, err := fetchProjects(testCtx, tt.model, client)
			if err != nil {
				if !tt.apiCallFails {
					t.Fatalf("did not fail on invalid input")
				}
				return
			}
			if err == nil && tt.apiCallFails {
				t.Fatalf("did not fail on invalid input")
			}
			if numAPICalls != tt.expectedNumAPICalls {
				t.Fatalf("Expected %d API calls, got %d", tt.expectedNumAPICalls, numAPICalls)
			}
			if len(projects) != tt.expectedNumItems {
				t.Fatalf("Expected %d projects, got %d", tt.totalItems, len(projects))
			}
		})
	}
}
