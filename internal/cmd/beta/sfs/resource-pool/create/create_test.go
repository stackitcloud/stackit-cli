package create

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-sdk-go/services/sfs"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &sfs.APIClient{}

var (
	testProjectId                          = uuid.NewString()
	testRegion                             = "eu02"
	testResourcePoolPerformanceClass       = "Standard"
	testResourcePoolSizeInGB         int64 = 50
	testResourcePoolAvailabilityZone       = "eu02-m"
	testResourcePoolName                   = "sfs-resource-pool-01"
	testResourcePoolIpAcl                  = []string{"10.88.135.144/28", "250.81.87.224/32"}
	testSnapshotsVisible                   = true
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		performanceClassFlag:      testResourcePoolPerformanceClass,
		sizeFlag:                  strconv.FormatInt(testResourcePoolSizeInGB, 10),
		ipAclFlag:                 strings.Join(testResourcePoolIpAcl, ","),
		availabilityZoneFlag:      testResourcePoolAvailabilityZone,
		nameFlag:                  testResourcePoolName,
		snapshotsVisibleFlag:      strconv.FormatBool(testSnapshotsVisible),
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
		PerformanceClass: testResourcePoolPerformanceClass,
		AvailabilityZone: testResourcePoolAvailabilityZone,
		Name:             testResourcePoolName,
		SizeInGB:         testResourcePoolSizeInGB,
		IpAcl:            testResourcePoolIpAcl,
		SnapshotsVisible: testSnapshotsVisible,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *sfs.ApiCreateResourcePoolRequest)) sfs.ApiCreateResourcePoolRequest {
	request := testClient.CreateResourcePool(testCtx, testProjectId, testRegion)
	request = request.CreateResourcePoolPayload(sfs.CreateResourcePoolPayload{
		Name:                &testResourcePoolName,
		PerformanceClass:    &testResourcePoolPerformanceClass,
		AvailabilityZone:    &testResourcePoolAvailabilityZone,
		IpAcl:               &testResourcePoolIpAcl,
		SizeGigabytes:       &testResourcePoolSizeInGB,
		SnapshotsAreVisible: &testSnapshotsVisible,
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
		ipAclValues   []string
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
			description: "ip acl missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, ipAclFlag)
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
			description: "performance class missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, performanceClassFlag)
			}),
			isValid: false,
		},
		{
			description: "size missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, sizeFlag)
			}),
			isValid: false,
		},
		{
			description: "availability zone missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, availabilityZoneFlag)
			}),
			isValid: false,
		},
		{
			description: "missing snapshot visible - fallback to false",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, snapshotsVisibleFlag)
			}),
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.SnapshotsVisible = false
			}),
			isValid: true,
		},
		{
			description: "project id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "repeated ip acl flags",
			flagValues:  fixtureFlagValues(),
			ipAclValues: []string{"198.51.100.14/24", "198.51.100.14/32"},
			isValid:     true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.IpAcl = append(model.IpAcl, "198.51.100.14/24", "198.51.100.14/32")
			}),
		},
		{
			description: "repeated ip acl flags with list value",
			flagValues:  fixtureFlagValues(),
			ipAclValues: []string{"198.51.100.14/24,198.51.100.14/32"},
			isValid:     true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.IpAcl = append(model.IpAcl, "198.51.100.14/24", "198.51.100.14/32")
			}),
		},
		{
			description: "invalid ip acl 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[ipAclFlag] = "foo-bar"
			}),
			isValid: false,
		},
		{
			description: "invalid ip acl 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[ipAclFlag] = "192.168.178.256/32"
			}),
			isValid: false,
		},
		{
			description: "invalid ip acl 3",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[ipAclFlag] = "192.168.178.255/32,"
			}),
			isValid: false,
		},
		{
			description: "invalid ip acl 4",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[ipAclFlag] = "192.168.178.255/32,"
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInputWithAdditionalFlags(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, map[string][]string{
				ipAclFlag: tt.ipAclValues,
			}, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest sfs.ApiCreateResourcePoolRequest
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
		projectLabel string
		resp         *sfs.CreateResourcePoolResponse
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
			name: "set empty response",
			args: args{
				resp: &sfs.CreateResourcePoolResponse{},
			},
			wantErr: false,
		},
		{
			name: "set response",
			args: args{
				resp: &sfs.CreateResourcePoolResponse{
					ResourcePool: &sfs.CreateResourcePoolResponseResourcePool{},
				},
			},
			wantErr: false,
		},
	}

	p := print.NewPrinter()
	p.Cmd = NewCmd(&types.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.outputFormat, tt.args.projectLabel, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
