package update

import (
	"context"
	"testing"
	"time"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	ske "github.com/stackitcloud/stackit-sdk-go/services/ske/v2api"
)

var projectIdFlag = globalflags.ProjectIdFlag

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &ske.APIClient{DefaultAPI: &ske.DefaultAPIService{}}
var testProjectId = uuid.NewString()
var testClusterName = "cluster"

const testRegion = "eu01"

var testPayload = ske.CreateOrUpdateClusterPayload{
	Kubernetes: ske.Kubernetes{
		Version: "1.25.15",
	},
	Nodepools: []ske.Nodepool{
		{
			Name: "np-name",
			Machine: ske.Machine{
				Image: ske.Image{
					Name:    "flatcar",
					Version: "3760.2.1",
				},
				Type: "b1.2",
			},
			Minimum:  int32(1),
			Maximum:  int32(2),
			MaxSurge: utils.Ptr(int32(1)),
			Volume: ske.Volume{
				Type: utils.Ptr("storage_premium_perf0"),
				Size: int32(40),
			},
			AvailabilityZones: []string{"eu01-3"},
			Cri:               &ske.CRI{Name: utils.Ptr("containerd")},
		},
	},
	Extensions: &ske.Extension{
		Acl: &ske.ACL{
			Enabled:      true,
			AllowedCidrs: []string{"0.0.0.0/0"},
		},
	},
	Maintenance: &ske.Maintenance{
		AutoUpdate: ske.MaintenanceAutoUpdate{
			KubernetesVersion:   utils.Ptr(true),
			MachineImageVersion: utils.Ptr(true),
		},
		TimeWindow: ske.TimeWindow{
			End:   time.Date(0, 1, 1, 5, 0, 0, 0, time.FixedZone("test-zone", 2*60*60)),
			Start: time.Date(0, 1, 1, 3, 0, 0, 0, time.FixedZone("test-zone", 2*60*60)),
		},
	},
}

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testClusterName,
	}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		payloadFlag: `{
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
					"version": "3760.2.1"
				  },
				  "type": "b1.2"
				},
				"minimum": 1,
				"maximum": 2,
				"maxSurge": 1,
				"volume": { "type": "storage_premium_perf0", "size": 40 },
				"cri": { "name": "containerd" },
				"availabilityZones": ["eu01-3"]
			  }
			],
			"extensions": { "acl": { "enabled": true, "allowedCidrs": ["0.0.0.0/0"] } },
			"maintenance": {
				"autoUpdate": {
				  "kubernetesVersion": true,
				  "machineImageVersion": true
				},
				"timeWindow": {
				  "end": "0000-01-01T05:00:00+02:00",
				  "start": "0000-01-01T03:00:00+02:00"
				}
			  }
		  }`,
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
			Region:    testRegion,
			Verbosity: globalflags.VerbosityDefault,
		},
		ClusterName: testClusterName,
		Payload:     testPayload,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *ske.ApiCreateOrUpdateClusterRequest)) ske.ApiCreateOrUpdateClusterRequest {
	request := testClient.DefaultAPI.CreateOrUpdateCluster(testCtx, testProjectId, testRegion, fixtureInputModel().ClusterName)
	request = request.CreateOrUpdateClusterPayload(testPayload)
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
			description: "no values",
			argValues:   []string{},
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "no arg values",
			argValues:   []string{},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "no flag values",
			argValues:   fixtureArgValues(),
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "project id missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, projectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "invalid json",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[payloadFlag] = "not json"
			}),
			isValid:       false,
			expectedModel: fixtureInputModel(),
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
		expectedRequest ske.ApiCreateOrUpdateClusterRequest
		isValid         bool
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
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx),
				cmpopts.EquateComparable(testClient.DefaultAPI),
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
		async        bool
		clusterName  string
		cluster      *ske.Cluster
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
			name: "empty cluster",
			args: args{
				cluster: &ske.Cluster{},
			},
			wantErr: false,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&types.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.outputFormat, tt.args.async, tt.args.clusterName, tt.args.cluster); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
