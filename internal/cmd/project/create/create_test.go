package create

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/resourcemanager"
	"github.com/zalando/go-keyring"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &resourcemanager.APIClient{}
var testParentId = uuid.NewString()
var testNetworkAreaId = uuid.NewString()
var testEmail = "email"

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		parentIdFlag:      testParentId,
		nameFlag:          "name",
		labelFlag:         "key=value",
		networkAreaIdFlag: testNetworkAreaId,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{Verbosity: globalflags.VerbosityDefault},
		ParentId:        utils.Ptr(testParentId),
		Name:            utils.Ptr(nameFlag),
		Labels: utils.Ptr(map[string]string{
			"key": "value",
		}),
		NetworkAreaId: utils.Ptr(testNetworkAreaId),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixturePayload(mods ...func(payload *resourcemanager.CreateProjectPayload)) resourcemanager.CreateProjectPayload {
	payload := resourcemanager.CreateProjectPayload{
		ContainerParentId: utils.Ptr(testParentId),
		Name:              utils.Ptr(nameFlag),
		Labels: utils.Ptr(map[string]string{
			"key":            "value",
			networkAreaLabel: testNetworkAreaId,
		}),
		Members: &[]resourcemanager.Member{
			{
				Role:    utils.Ptr(ownerRole),
				Subject: utils.Ptr(testEmail),
			},
		},
	}
	for _, mod := range mods {
		mod(&payload)
	}
	return payload
}

func fixtureRequest(mods ...func(request *resourcemanager.ApiCreateProjectRequest)) resourcemanager.ApiCreateProjectRequest {
	request := testClient.CreateProject(testCtx)
	request = request.CreateProjectPayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description   string
		flagValues    map[string]string
		labelValues   []string
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
			description: "no values",
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "parent id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, parentIdFlag)
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
			description: "multiple_labels",
			flagValues:  fixtureFlagValues(),
			labelValues: []string{"key=value", "foo=bar"},
			expectedModel: fixtureInputModel(
				func(model *inputModel) {
					model.Labels = &map[string]string{
						"key": "value",
						"foo": "bar",
					}
				}),
			isValid: true,
		},
		{
			description: "multiple_labels_2",
			flagValues:  fixtureFlagValues(),
			labelValues: []string{"key=value,foo=bar"},
			expectedModel: fixtureInputModel(
				func(model *inputModel) {
					model.Labels = &map[string]string{
						"key": "value",
						"foo": "bar",
					}
				}),
			isValid: true,
		},
		{
			description: "invalid_labels",
			flagValues:  fixtureFlagValues(),
			labelValues: []string{"key"},
			isValid:     false,
		},
		{
			description: "network_area_id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, networkAreaIdFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(
				func(model *inputModel) {
					model.NetworkAreaId = nil
				}),
		},
		{
			description: "network_area_id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[networkAreaIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "network_area_id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[networkAreaIdFlag] = "invalid-uuid"
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

			for _, value := range tt.labelValues {
				err := cmd.Flags().Set(labelFlag, value)
				if err != nil {
					if !tt.isValid {
						return
					}
					t.Fatalf("setting flag --%s=%s: %v", labelFlag, value, err)
				}
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
	tests := []struct {
		description     string
		model           *inputModel
		authFlow        auth.AuthFlow
		sa_email        *string
		user_email      *string
		expectedRequest resourcemanager.ApiCreateProjectRequest
		isValid         bool
	}{
		{
			description:     "base_sa_key",
			model:           fixtureInputModel(),
			authFlow:        auth.AUTH_FLOW_SERVICE_ACCOUNT_KEY,
			sa_email:        utils.Ptr(testEmail),
			expectedRequest: fixtureRequest(),
			isValid:         true,
		},
		{
			description:     "base_sa_token",
			model:           fixtureInputModel(),
			authFlow:        auth.AUTH_FLOW_SERVICE_ACCOUNT_TOKEN,
			sa_email:        utils.Ptr(testEmail),
			expectedRequest: fixtureRequest(),
			isValid:         true,
		},
		{
			description:     "base_user",
			model:           fixtureInputModel(),
			authFlow:        auth.AUTH_FLOW_USER_TOKEN,
			user_email:      utils.Ptr(testEmail),
			expectedRequest: fixtureRequest(),
			isValid:         true,
		},
		{
			description: "missing_network_area_id sa_key",
			model: fixtureInputModel(
				func(model *inputModel) {
					model.NetworkAreaId = nil
				}),
			authFlow: auth.AUTH_FLOW_SERVICE_ACCOUNT_KEY,
			sa_email: utils.Ptr(testEmail),
			expectedRequest: fixtureRequest().CreateProjectPayload(fixturePayload(
				func(payload *resourcemanager.CreateProjectPayload) {
					delete((*payload.Labels), networkAreaLabel)
				}),
			),
			isValid: true,
		},
		{
			description: "missing_network_area_id sa_token",
			model: fixtureInputModel(
				func(model *inputModel) {
					model.NetworkAreaId = nil
				}),
			authFlow: auth.AUTH_FLOW_SERVICE_ACCOUNT_TOKEN,
			sa_email: utils.Ptr(testEmail),
			expectedRequest: fixtureRequest().CreateProjectPayload(fixturePayload(
				func(payload *resourcemanager.CreateProjectPayload) {
					delete((*payload.Labels), networkAreaLabel)
				}),
			),
			isValid: true,
		},
		{
			description: "missing_network_area_id user",
			model: fixtureInputModel(
				func(model *inputModel) {
					model.NetworkAreaId = nil
				}),
			authFlow:   auth.AUTH_FLOW_USER_TOKEN,
			user_email: utils.Ptr(testEmail),
			expectedRequest: fixtureRequest().CreateProjectPayload(fixturePayload(
				func(payload *resourcemanager.CreateProjectPayload) {
					delete((*payload.Labels), networkAreaLabel)
				}),
			),
			isValid: true,
		},
		{
			description: "missing_auth_flow",
			model:       fixtureInputModel(),
			isValid:     false,
		},
		{
			description: "missing_email",
			model:       fixtureInputModel(),
			isValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			keyring.MockInit()
			err := auth.SetAuthFlow(tt.authFlow)
			if err != nil {
				t.Fatalf("Failed to set auth flow in storage: %v", err)
			}
			if tt.sa_email != nil {
				err := auth.SetAuthField(auth.SERVICE_ACCOUNT_EMAIL, *tt.sa_email)
				if err != nil {
					t.Fatalf("Failed to set service account email in storage: %v", err)
				}
			}
			if tt.user_email != nil {
				err := auth.SetAuthField(auth.USER_EMAIL, *tt.user_email)
				if err != nil {
					t.Fatalf("Failed to set user email in storage: %v", err)
				}
			}
			request, err := buildRequest(testCtx, tt.model, testClient)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error building request: %v", err)
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

func Test_outputResult(t *testing.T) {
	type args struct {
		model inputModel
		resp  *resourcemanager.Project
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty", args{model: inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{}}}, true},
		{"base", args{inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{}}, &resourcemanager.Project{}}, false},
	}

	p := print.NewPrinter()
	p.Cmd = NewCmd(p)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.model, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
