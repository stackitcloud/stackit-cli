package utils

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/argus"
)

var (
	testClient     = &argus.APIClient{}
	testProjectId  = uuid.NewString()
	testInstanceId = uuid.NewString()
	testPlanId     = uuid.NewString()
)

const (
	testInstanceName = "instance"
	testPlanName     = "Plan-Name-01"
)

var testPlansResponse = argus.PlansResponse{
	Plans: &[]argus.Plan{
		{
			Id:   utils.Ptr(testPlanId),
			Name: utils.Ptr(testPlanName),
		},
	},
}

func fixtureGetScrapeConfigResponse(mods ...func(*argus.GetScrapeConfigResponse)) *argus.GetScrapeConfigResponse {
	number := int64(1)
	resp := &argus.GetScrapeConfigResponse{
		Data: &argus.Job{
			BasicAuth: &argus.BasicAuth{
				Username: utils.Ptr("username"),
				Password: utils.Ptr("password"),
			},
			BearerToken:     utils.Ptr("bearerToken"),
			HonorLabels:     utils.Ptr(true),
			HonorTimeStamps: utils.Ptr(true),
			MetricsPath:     utils.Ptr("/metrics"),
			MetricsRelabelConfigs: &[]argus.MetricsRelabelConfig{
				{
					Action:       utils.Ptr("replace"),
					Modulus:      &number,
					Regex:        utils.Ptr("regex"),
					Replacement:  utils.Ptr("replacement"),
					Separator:    utils.Ptr("separator"),
					SourceLabels: &[]string{"sourceLabel"},
					TargetLabel:  utils.Ptr("targetLabel"),
				},
			},
			Params: &map[string][]string{
				"key":  {"value1", "value2"},
				"key2": {},
			},
			SampleLimit:    &number,
			Scheme:         utils.Ptr("scheme"),
			ScrapeInterval: utils.Ptr("interval"),
			ScrapeTimeout:  utils.Ptr("timeout"),
			StaticConfigs: &[]argus.StaticConfigs{
				{
					Labels: &map[string]string{
						"label":  "value",
						"label2": "value2",
					},
					Targets: &[]string{"target"},
				},
			},
			TlsConfig: &argus.TLSConfig{
				InsecureSkipVerify: utils.Ptr(true),
			},
		},
	}

	for _, mod := range mods {
		mod(resp)
	}

	return resp
}

func fixtureUpdateScrapeConfigPayload(mods ...func(*argus.UpdateScrapeConfigPayload)) *argus.UpdateScrapeConfigPayload {
	payload := &argus.UpdateScrapeConfigPayload{
		BasicAuth: &argus.CreateScrapeConfigPayloadBasicAuth{
			Username: utils.Ptr("username"),
			Password: utils.Ptr("password"),
		},
		BearerToken:     utils.Ptr("bearerToken"),
		HonorLabels:     utils.Ptr(true),
		HonorTimeStamps: utils.Ptr(true),
		MetricsPath:     utils.Ptr("/metrics"),
		MetricsRelabelConfigs: &[]argus.CreateScrapeConfigPayloadMetricsRelabelConfigsInner{
			{
				Action:       utils.Ptr("replace"),
				Modulus:      utils.Ptr(1.0),
				Regex:        utils.Ptr("regex"),
				Replacement:  utils.Ptr("replacement"),
				Separator:    utils.Ptr("separator"),
				SourceLabels: &[]string{"sourceLabel"},
				TargetLabel:  utils.Ptr("targetLabel"),
			},
		},
		Params: &map[string]interface{}{
			"key":  []string{"value1", "value2"},
			"key2": []string{},
		},
		SampleLimit:    utils.Ptr(1.0),
		Scheme:         utils.Ptr("scheme"),
		ScrapeInterval: utils.Ptr("interval"),
		ScrapeTimeout:  utils.Ptr("timeout"),
		StaticConfigs: &[]argus.UpdateScrapeConfigPayloadStaticConfigsInner{
			{
				Labels: &map[string]interface{}{
					"label":  "value",
					"label2": "value2",
				},
				Targets: &[]string{"target"},
			},
		},
		TlsConfig: &argus.CreateScrapeConfigPayloadHttpSdConfigsInnerOauth2TlsConfig{
			InsecureSkipVerify: utils.Ptr(true),
		},
	}

	for _, mod := range mods {
		mod(payload)
	}

	return payload
}

type argusClientMocked struct {
	getInstanceFails       bool
	getInstanceResp        *argus.GetInstanceResponse
	getGrafanaConfigsFails bool
	getGrafanaConfigsResp  *argus.GrafanaConfigs
}

func (m *argusClientMocked) GetInstanceExecute(_ context.Context, _, _ string) (*argus.GetInstanceResponse, error) {
	if m.getInstanceFails {
		return nil, fmt.Errorf("could not get instance")
	}
	return m.getInstanceResp, nil
}

func (m *argusClientMocked) GetGrafanaConfigsExecute(_ context.Context, _, _ string) (*argus.GrafanaConfigs, error) {
	if m.getGrafanaConfigsFails {
		return nil, fmt.Errorf("could not get grafana configs")
	}
	return m.getGrafanaConfigsResp, nil
}

func (c *argusClientMocked) UpdateGrafanaConfigs(ctx context.Context, instanceId, projectId string) argus.ApiUpdateGrafanaConfigsRequest {
	return testClient.UpdateGrafanaConfigs(ctx, instanceId, projectId)
}

func fixtureGrafanaConfigs(mods ...func(gc *argus.GrafanaConfigs)) *argus.GrafanaConfigs {
	gc := argus.GrafanaConfigs{
		GenericOauth: &argus.GrafanaOauth{
			ApiUrl:              utils.Ptr("apiUrl"),
			AuthUrl:             utils.Ptr("authUrl"),
			Enabled:             utils.Ptr(true),
			Name:                utils.Ptr("name"),
			OauthClientId:       utils.Ptr("oauthClientId"),
			OauthClientSecret:   utils.Ptr("oauthClientSecret"),
			RoleAttributePath:   utils.Ptr("roleAttributePath"),
			RoleAttributeStrict: utils.Ptr(true),
			Scopes:              utils.Ptr("scopes"),
			TokenUrl:            utils.Ptr("tokenUrl"),
			UsePkce:             utils.Ptr(true),
		},
		PublicReadAccess: utils.Ptr(false),
		UseStackitSso:    utils.Ptr(false),
	}
	for _, mod := range mods {
		mod(&gc)
	}
	return &gc
}

func TestGetInstanceName(t *testing.T) {
	tests := []struct {
		description      string
		getInstanceFails bool
		getInstanceResp  *argus.GetInstanceResponse
		isValid          bool
		expectedOutput   string
	}{
		{
			description: "base",
			getInstanceResp: &argus.GetInstanceResponse{
				Name: utils.Ptr(testInstanceName),
			},
			isValid:        true,
			expectedOutput: testInstanceName,
		},
		{
			description:      "get instance fails",
			getInstanceFails: true,
			isValid:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &argusClientMocked{
				getInstanceFails: tt.getInstanceFails,
				getInstanceResp:  tt.getInstanceResp,
			}

			output, err := GetInstanceName(context.Background(), client, testInstanceId, testProjectId)

			if tt.isValid && err != nil {
				t.Errorf("failed on valid input")
			}
			if !tt.isValid && err == nil {
				t.Errorf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}
			if output != tt.expectedOutput {
				t.Errorf("expected output to be %s, got %s", tt.expectedOutput, output)
			}
		})
	}
}

func TestLoadPlanId(t *testing.T) {
	tests := []struct {
		description    string
		planName       string
		plansResponse  *argus.PlansResponse
		isValid        bool
		expectedOutput string
	}{
		{
			description:    "base case",
			planName:       testPlanName,
			plansResponse:  utils.Ptr(testPlansResponse),
			expectedOutput: testPlanId,
			isValid:        true,
		},
		{
			description:    "different casing",
			planName:       strings.ToLower(testPlanName),
			plansResponse:  utils.Ptr(testPlansResponse),
			expectedOutput: testPlanId,
			isValid:        true,
		},
		{
			description:   "empty plan name",
			planName:      "",
			plansResponse: utils.Ptr(testPlansResponse),
			isValid:       false,
		},
		{
			description:   "unexisting plan name",
			planName:      "another plan name",
			plansResponse: utils.Ptr(testPlansResponse),
			isValid:       false,
		},
		{
			description: "unable to fetch plans",
			isValid:     false,
		},
		{
			description: "no available plans",
			planName:    testPlanName,
			plansResponse: &argus.PlansResponse{
				Plans: &[]argus.Plan{},
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			output, err := LoadPlanId(tt.planName, tt.plansResponse)

			if tt.isValid && err != nil {
				t.Errorf("failed on valid input")
			}
			if !tt.isValid && err == nil {
				t.Errorf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}
			if *output != tt.expectedOutput {
				t.Errorf("expected output to be %s, got %s", tt.expectedOutput, *output)
			}
		})
	}
}

func TestValidatePlanId(t *testing.T) {
	tests := []struct {
		description   string
		planId        string
		plansResponse *argus.PlansResponse
		isValid       bool
	}{
		{
			description:   "base case",
			planId:        testPlanId,
			plansResponse: utils.Ptr(testPlansResponse),
			isValid:       true,
		},
		{
			description:   "different casing",
			planId:        strings.ToLower(testPlanId),
			plansResponse: utils.Ptr(testPlansResponse),
			isValid:       true,
		},
		{
			description:   "empty plan id",
			planId:        "",
			plansResponse: utils.Ptr(testPlansResponse),
			isValid:       false,
		},
		{
			description:   "unexisting plan id",
			planId:        uuid.NewString(),
			plansResponse: utils.Ptr(testPlansResponse),
			isValid:       false,
		},
		{
			description: "unable to fetch plans",
			isValid:     false,
		},
		{
			description: "no available plans",
			planId:      testPlanId,
			plansResponse: &argus.PlansResponse{
				Plans: &[]argus.Plan{},
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := ValidatePlanId(tt.planId, tt.plansResponse)

			if tt.isValid && err != nil {
				t.Errorf("failed on valid input")
			}
			if !tt.isValid && err == nil {
				t.Errorf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}
		})
	}
}

func TestMapToUpdateScrapeConfigPayload(t *testing.T) {
	tests := []struct {
		description     string
		resp            *argus.GetScrapeConfigResponse
		expectedPayload *argus.UpdateScrapeConfigPayload
		isValid         bool
	}{
		{
			description:     "base case",
			resp:            fixtureGetScrapeConfigResponse(),
			expectedPayload: fixtureUpdateScrapeConfigPayload(),
			isValid:         true,
		},
		{
			description: "nil response",
			resp:        nil,
			isValid:     false,
		},
		{
			description: "nil data",
			resp: &argus.GetScrapeConfigResponse{
				Data: nil,
			},
			isValid: false,
		},
		{
			description: "empty data",
			resp: &argus.GetScrapeConfigResponse{
				Data: &argus.Job{},
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			payload, err := MapToUpdateScrapeConfigPayload(tt.resp)

			if tt.isValid && err != nil {
				t.Errorf("failed on valid input")
			}
			if !tt.isValid && err == nil {
				t.Errorf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}

			diff := cmp.Diff(*payload, *tt.expectedPayload,
				cmp.AllowUnexported(*tt.expectedPayload),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestMapMetricsRelabelConfig(t *testing.T) {
	tests := []struct {
		description string
		config      *[]argus.MetricsRelabelConfig
		expected    *[]argus.CreateScrapeConfigPayloadMetricsRelabelConfigsInner
	}{
		{
			description: "base case",
			config: &[]argus.MetricsRelabelConfig{
				{
					Action:       utils.Ptr("replace"),
					Modulus:      utils.Int64Ptr(1),
					Regex:        utils.Ptr("regex"),
					Replacement:  utils.Ptr("replacement"),
					Separator:    utils.Ptr("separator"),
					SourceLabels: utils.Ptr([]string{"sourceLabel", "sourceLabel2"}),
					TargetLabel:  utils.Ptr("targetLabel"),
				},
			},
			expected: &[]argus.CreateScrapeConfigPayloadMetricsRelabelConfigsInner{
				{
					Action:       utils.Ptr("replace"),
					Modulus:      utils.Float64Ptr(1.0),
					Regex:        utils.Ptr("regex"),
					Replacement:  utils.Ptr("replacement"),
					Separator:    utils.Ptr("separator"),
					SourceLabels: utils.Ptr([]string{"sourceLabel", "sourceLabel2"}),
					TargetLabel:  utils.Ptr("targetLabel"),
				},
			},
		},
		{
			description: "empty data",
			config:      &[]argus.MetricsRelabelConfig{},
			expected:    nil,
		},
		{
			description: "nil",
			config:      nil,
			expected:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			output := mapMetricsRelabelConfig(tt.config)

			if tt.expected == nil && output == nil || *output == nil {
				return
			}

			diff := cmp.Diff(*output, *tt.expected)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestMapStaticConfig(t *testing.T) {
	tests := []struct {
		description string
		config      *[]argus.StaticConfigs
		expected    *[]argus.UpdateScrapeConfigPayloadStaticConfigsInner
	}{
		{
			description: "base case",
			config: &[]argus.StaticConfigs{
				{
					Labels: &map[string]string{
						"label":  "value",
						"label2": "value2",
					},
					Targets: &[]string{"target", "target2"},
				},
			},
			expected: &[]argus.UpdateScrapeConfigPayloadStaticConfigsInner{
				{
					Labels: utils.Ptr(map[string]interface{}{
						"label":  "value",
						"label2": "value2",
					}),
					Targets: utils.Ptr([]string{"target", "target2"}),
				},
			},
		},
		{
			description: "empty data",
			config:      &[]argus.StaticConfigs{},
			expected:    nil,
		},
		{
			description: "nil",
			config:      nil,
			expected:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			output := mapStaticConfig(tt.config)

			if tt.expected == nil && (output == nil || *output == nil) {
				return
			}

			diff := cmp.Diff(*output, *tt.expected)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestMapBasicAuth(t *testing.T) {
	tests := []struct {
		description string
		auth        *argus.BasicAuth
		expected    *argus.CreateScrapeConfigPayloadBasicAuth
	}{
		{
			description: "base case",
			auth: &argus.BasicAuth{
				Username: utils.Ptr("username"),
				Password: utils.Ptr("password"),
			},
			expected: &argus.CreateScrapeConfigPayloadBasicAuth{
				Username: utils.Ptr("username"),
				Password: utils.Ptr("password"),
			},
		},
		{
			description: "empty data",
			auth:        &argus.BasicAuth{},
			expected:    &argus.CreateScrapeConfigPayloadBasicAuth{},
		},
		{
			description: "nil",
			auth:        nil,
			expected:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			output := mapBasicAuth(tt.auth)

			if tt.expected == nil && output == nil && tt.auth == nil {
				return
			}

			diff := cmp.Diff(*output, *tt.expected)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestMapTlsConfig(t *testing.T) {
	tests := []struct {
		description string
		config      *argus.TLSConfig
		expected    *argus.CreateScrapeConfigPayloadHttpSdConfigsInnerOauth2TlsConfig
	}{
		{
			description: "base case",
			config: &argus.TLSConfig{
				InsecureSkipVerify: utils.Ptr(true),
			},
			expected: &argus.CreateScrapeConfigPayloadHttpSdConfigsInnerOauth2TlsConfig{
				InsecureSkipVerify: utils.Ptr(true),
			},
		},
		{
			description: "empty data",
			config:      &argus.TLSConfig{},
			expected:    &argus.CreateScrapeConfigPayloadHttpSdConfigsInnerOauth2TlsConfig{},
		},
		{
			description: "nil",
			config:      nil,
			expected:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			output := mapTlsConfig(tt.config)

			if tt.expected == nil && output == nil && tt.config == nil {
				return
			}

			diff := cmp.Diff(*output, *tt.expected)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestMapParams(t *testing.T) {
	tests := []struct {
		description string
		params      map[string][]string
		expected    map[string]interface{}
	}{
		{
			description: "base case",
			params: map[string][]string{
				"key":  {"value1", "value2"},
				"key2": {},
			},
			expected: map[string]interface{}{
				"key":  []string{"value1", "value2"},
				"key2": []string{},
			},
		},
		{
			description: "empty data",
			params:      map[string][]string{},
			expected:    map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			output := mapParams(tt.params)

			if tt.expected == nil && output == nil && tt.params == nil {
				return
			}

			diff := cmp.Diff(output, tt.expected)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestMapStaticConfigLabels(t *testing.T) {
	tests := []struct {
		description string
		labels      map[string]string
		expected    map[string]interface{}
	}{
		{
			description: "base case",
			labels: map[string]string{
				"label":  "value",
				"label2": "value2",
			},
			expected: map[string]interface{}{
				"label":  "value",
				"label2": "value2",
			},
		},
		{
			description: "empty data",
			labels:      map[string]string{},
			expected:    map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			output := mapStaticConfigLabels(tt.labels)

			diff := cmp.Diff(output, tt.expected)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestToPayloadGenericOAuth(t *testing.T) {
	tests := []struct {
		description string
		response    *argus.GrafanaOauth
		expected    *argus.UpdateGrafanaConfigsPayloadGenericOauth
	}{
		{
			description: "base",
			response: &argus.GrafanaOauth{
				ApiUrl:              utils.Ptr("apiUrl"),
				AuthUrl:             utils.Ptr("authUrl"),
				Enabled:             utils.Ptr(true),
				Name:                utils.Ptr("name"),
				OauthClientId:       utils.Ptr("oauthClientId"),
				OauthClientSecret:   utils.Ptr("oauthClientSecret"),
				RoleAttributePath:   utils.Ptr("roleAttributePath"),
				RoleAttributeStrict: utils.Ptr(true),
				Scopes:              utils.Ptr("scopes"),
				TokenUrl:            utils.Ptr("tokenUrl"),
				UsePkce:             utils.Ptr(true),
			},
			expected: &argus.UpdateGrafanaConfigsPayloadGenericOauth{
				ApiUrl:              utils.Ptr("apiUrl"),
				AuthUrl:             utils.Ptr("authUrl"),
				Enabled:             utils.Ptr(true),
				Name:                utils.Ptr("name"),
				OauthClientId:       utils.Ptr("oauthClientId"),
				OauthClientSecret:   utils.Ptr("oauthClientSecret"),
				RoleAttributePath:   utils.Ptr("roleAttributePath"),
				RoleAttributeStrict: utils.Ptr(true),
				Scopes:              utils.Ptr("scopes"),
				TokenUrl:            utils.Ptr("tokenUrl"),
				UsePkce:             utils.Ptr(true),
			},
		},
		{
			description: "nil response oauth",
			response:    nil,
			expected:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			output := ToPayloadGenericOAuth(tt.response)

			diff := cmp.Diff(output, tt.expected)
			if diff != "" {
				t.Errorf("expected output to be %+v, got %+v", tt.expected, output)
			}
		})
	}
}

func TestGetPartialUpdateGrafanaConfigsPayload(t *testing.T) {
	tests := []struct {
		description            string
		singleSignOn           *bool
		publicReadAccess       *bool
		getGrafanaConfigsFails bool
		getGrafanaConfigsResp  *argus.GrafanaConfigs
		isValid                bool
		expectedPayload        *argus.UpdateGrafanaConfigsPayload
	}{
		{
			description:           "enable both",
			singleSignOn:          utils.Ptr(true),
			publicReadAccess:      utils.Ptr(true),
			getGrafanaConfigsResp: fixtureGrafanaConfigs(),
			isValid:               true,
			expectedPayload: &argus.UpdateGrafanaConfigsPayload{
				GenericOauth:     ToPayloadGenericOAuth(fixtureGrafanaConfigs().GenericOauth),
				UseStackitSso:    utils.Ptr(true),
				PublicReadAccess: utils.Ptr(true),
			},
		},
		{
			description:           "disable both",
			singleSignOn:          utils.Ptr(false),
			publicReadAccess:      utils.Ptr(false),
			getGrafanaConfigsResp: fixtureGrafanaConfigs(),
			isValid:               true,
			expectedPayload: &argus.UpdateGrafanaConfigsPayload{
				GenericOauth:     ToPayloadGenericOAuth(fixtureGrafanaConfigs().GenericOauth),
				UseStackitSso:    utils.Ptr(false),
				PublicReadAccess: utils.Ptr(false),
			},
		},
		{
			description:           "enable single sign on",
			singleSignOn:          utils.Ptr(true),
			publicReadAccess:      nil,
			getGrafanaConfigsResp: fixtureGrafanaConfigs(),
			isValid:               true,
			expectedPayload: &argus.UpdateGrafanaConfigsPayload{
				GenericOauth:     ToPayloadGenericOAuth(fixtureGrafanaConfigs().GenericOauth),
				UseStackitSso:    utils.Ptr(true),
				PublicReadAccess: fixtureGrafanaConfigs().PublicReadAccess,
			},
		},
		{
			description:           "enable public read access",
			singleSignOn:          nil,
			publicReadAccess:      utils.Ptr(true),
			getGrafanaConfigsResp: fixtureGrafanaConfigs(),
			isValid:               true,
			expectedPayload: &argus.UpdateGrafanaConfigsPayload{
				GenericOauth:     ToPayloadGenericOAuth(fixtureGrafanaConfigs().GenericOauth),
				UseStackitSso:    fixtureGrafanaConfigs().UseStackitSso,
				PublicReadAccess: utils.Ptr(true),
			},
		},
		{
			description:      "disable single sign on",
			singleSignOn:     utils.Ptr(false),
			publicReadAccess: nil,
			getGrafanaConfigsResp: fixtureGrafanaConfigs(func(gc *argus.GrafanaConfigs) {
				gc.UseStackitSso = utils.Ptr(true)
			}),
			isValid: true,
			expectedPayload: &argus.UpdateGrafanaConfigsPayload{
				GenericOauth:     ToPayloadGenericOAuth(fixtureGrafanaConfigs().GenericOauth),
				UseStackitSso:    utils.Ptr(false),
				PublicReadAccess: fixtureGrafanaConfigs().PublicReadAccess,
			},
		},
		{
			description:      "disable public read access",
			singleSignOn:     nil,
			publicReadAccess: utils.Ptr(false),
			getGrafanaConfigsResp: fixtureGrafanaConfigs(func(gc *argus.GrafanaConfigs) {
				gc.PublicReadAccess = utils.Ptr(true)
			}),
			isValid: true,
			expectedPayload: &argus.UpdateGrafanaConfigsPayload{
				GenericOauth:     ToPayloadGenericOAuth(fixtureGrafanaConfigs().GenericOauth),
				UseStackitSso:    fixtureGrafanaConfigs().UseStackitSso,
				PublicReadAccess: utils.Ptr(false),
			},
		},
		{
			description:           "nil generic oauth",
			singleSignOn:          utils.Ptr(true),
			publicReadAccess:      utils.Ptr(true),
			getGrafanaConfigsResp: &argus.GrafanaConfigs{},
			isValid:               true,
			expectedPayload: &argus.UpdateGrafanaConfigsPayload{
				GenericOauth:     nil,
				UseStackitSso:    utils.Ptr(true),
				PublicReadAccess: utils.Ptr(true),
			},
		},
		{
			description:            "get grafana configs fails",
			singleSignOn:           utils.Ptr(true),
			publicReadAccess:       utils.Ptr(true),
			getGrafanaConfigsFails: true,
			isValid:                false,
		},
		{
			description:           "no grafana configs",
			singleSignOn:          utils.Ptr(true),
			publicReadAccess:      utils.Ptr(true),
			getGrafanaConfigsResp: nil,
			isValid:               false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &argusClientMocked{
				getGrafanaConfigsFails: tt.getGrafanaConfigsFails,
				getGrafanaConfigsResp:  tt.getGrafanaConfigsResp,
			}

			payload, err := GetPartialUpdateGrafanaConfigsPayload(context.Background(), client, testInstanceId, testProjectId, tt.singleSignOn, tt.publicReadAccess)

			if tt.isValid && err != nil {
				t.Errorf("failed on valid input")
			}
			if !tt.isValid && err == nil {
				t.Errorf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}

			diff := cmp.Diff(payload, tt.expectedPayload)
			if diff != "" {
				t.Errorf("expected output payload to be %+v, got %+v", tt.expectedPayload, payload)
			}
		})
	}
}
