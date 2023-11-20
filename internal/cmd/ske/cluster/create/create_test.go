package create

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &ske.APIClient{}
var testProjectId = uuid.NewString()

var testPayload = ske.CreateOrUpdateClusterPayload{
	Kubernetes: &ske.Kubernetes{
		Version: utils.Ptr("1.25.15"),
	},
	Nodepools: &[]ske.Nodepool{
		{
			Name: utils.Ptr("np-name"),
			Machine: &ske.Machine{
				Image: &ske.Image{
					Name:    utils.Ptr("flatcar"),
					Version: utils.Ptr("3602.2.1"),
				},
				Type: utils.Ptr("b1.2"),
			},
			Minimum:  utils.Ptr(int64(1)),
			Maximum:  utils.Ptr(int64(2)),
			MaxSurge: utils.Ptr(int64(1)),
			Volume: &ske.Volume{
				Type: utils.Ptr("storage_premium_perf0"),
				Size: utils.Ptr(int64(40)),
			},
			AvailabilityZones: &[]string{"eu01-3"},
			Cri:               &ske.CRI{Name: utils.Ptr("cri")},
		},
	},
	Extensions: &ske.Extension{
		Acl: &ske.ACL{
			Enabled:      utils.Ptr(true),
			AllowedCidrs: &[]string{"0.0.0.0/0"},
		},
	},
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		ProjectIdFlag: testProjectId,
		NameFlag:      "example-name",
		PayloadFlag: `{
			"name": "cli-jp",
			"kubernetes": {
			  "version": "1.25.15"
			},
			"nodepools": [
			  {
				"name": "np-name",
				"machine": {
				  "image": {
					"name": "flatcar",
					"version": "3602.2.1"
				  },
				  "type": "b1.2"
				},
				"minimum": 1,
				"maximum": 2,
				"maxSurge": 1,
				"volume": { "type": "storage_premium_perf0", "size": 40 },
				"cri": { "name": "cri" },
				"availabilityZones": ["eu01-3"]
			  }
			],
			"extensions": { "acl": { "enabled": true, "allowedCidrs": ["0.0.0.0/0"] } }
		  }`,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureFlagModel(mods ...func(model *FlagModel)) *FlagModel {
	model := &FlagModel{
		ProjectId: testProjectId,
		Name:      "example-name",
		Payload:   testPayload,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *ske.ApiCreateOrUpdateClusterRequest)) ske.ApiCreateOrUpdateClusterRequest {
	request := testClient.CreateOrUpdateCluster(testCtx, testProjectId, (*fixtureFlagModel()).Name)
	request = request.CreateOrUpdateClusterPayload(testPayload)
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func TestParseFlags(t *testing.T) {
	tests := []struct {
		description        string
		flagValues         map[string]string
		payloadFileContent []byte
		readFileFails      bool
		isValid            bool
		expectedModel      *FlagModel
	}{
		{
			description:   "base",
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureFlagModel(),
		},
		{
			description: "base from file",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[PayloadFlag] = "@" + fixtureFlagValues()[PayloadFlag]
			}),
			payloadFileContent: []byte(fixtureFlagValues()[PayloadFlag]),
			isValid:            true,
			expectedModel:      fixtureFlagModel(),
		},
		{
			description: "no values",
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "project id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, ProjectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[ProjectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[ProjectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "read file fails",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[PayloadFlag] = "@" + fixtureFlagValues()[PayloadFlag]
			}),
			payloadFileContent: []byte(fixtureFlagValues()[PayloadFlag]),
			readFileFails:      true,
			isValid:            false,
			expectedModel:      fixtureFlagModel(),
		},
		{
			description: "invalid json",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[PayloadFlag] = "@" + fixtureFlagValues()[PayloadFlag]
			}),
			payloadFileContent: []byte(`not json`),
			isValid:            false,
			expectedModel:      fixtureFlagModel(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			cmd := &cobra.Command{}

			// Flag defined in root command
			err := testutils.ConfigureBindUUIDFlag(cmd, ProjectIdFlag, config.ProjectIdKey)
			if err != nil {
				t.Fatalf("configure global flag --%s: %v", ProjectIdFlag, err)
			}

			ConfigureFlags(cmd)

			for flag, value := range tt.flagValues {
				err := cmd.Flags().Set(flag, value)
				if err != nil {
					if !tt.isValid {
						return
					}
					t.Fatalf("setting flag --%s=%s: %v", flag, value, err)
				}
			}

			err = cmd.ValidateRequiredFlags()
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flags: %v", err)
			}

			mockFileReader := func(filename string) ([]byte, error) {
				if tt.readFileFails {
					return nil, fmt.Errorf("could not read file")
				}
				return tt.payloadFileContent, nil
			}

			model, err := ParseFlags(cmd, mockFileReader)
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
		model           *FlagModel
		expectedRequest ske.ApiCreateOrUpdateClusterRequest
		isValid         bool
	}{
		{
			description:     "base",
			model:           fixtureFlagModel(),
			expectedRequest: fixtureRequest(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request, err := BuildRequest(testCtx, tt.model, testClient)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error building request: %v", err)
			}

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.IgnoreFields(ske.ApiCreateOrUpdateClusterRequest{}, "apiService", "ctx", "projectId"),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}
