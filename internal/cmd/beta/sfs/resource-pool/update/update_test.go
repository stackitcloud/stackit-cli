package update

import (
	"context"
	"slices"
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
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/sfs"
)

type testCtxKey struct{}

const (
	testRegion = "eu02"
)

var (
	testCtx    = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient = &sfs.APIClient{}

	testProjectId                          = uuid.NewString()
	testResourcePoolId                     = uuid.NewString()
	testResourcePoolIpAcl                  = []string{"10.88.135.144/28", "250.81.87.224/32"}
	testResourcePoolPerformanceClass       = "Standard"
	testResourcePoolSizeInGB         int64 = 50
	testSnapshotsVisible                   = true
)

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testResourcePoolId,
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
		performanceClassFlag:      testResourcePoolPerformanceClass,
		sizeFlag:                  strconv.FormatInt(testResourcePoolSizeInGB, 10),
		ipAclFlag:                 strings.Join(testResourcePoolIpAcl, ","),
		snapshotsVisibleFlag:      strconv.FormatBool(testSnapshotsVisible),
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	ipAclClone := slices.Clone(testResourcePoolIpAcl)

	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Region:    testRegion,
			Verbosity: globalflags.VerbosityDefault,
		},
		ResourcePoolId:   testResourcePoolId,
		SizeGigabytes:    &testResourcePoolSizeInGB,
		PerformanceClass: &testResourcePoolPerformanceClass,
		IpAcl:            &ipAclClone,
		SnapshotsVisible: &testSnapshotsVisible,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *sfs.ApiUpdateResourcePoolRequest)) sfs.ApiUpdateResourcePoolRequest {
	request := testClient.UpdateResourcePool(testCtx, testProjectId, testRegion, testResourcePoolId)
	request = request.UpdateResourcePoolPayload(sfs.UpdateResourcePoolPayload{
		IpAcl:               &testResourcePoolIpAcl,
		PerformanceClass:    &testResourcePoolPerformanceClass,
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
			description: "no values to update",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, sizeFlag)
				delete(flagValues, ipAclFlag)
				delete(flagValues, performanceClassFlag)
				delete(flagValues, snapshotsVisibleFlag)
			}),
			isValid: false,
		},
		{
			description: "update only size",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, ipAclFlag)
				delete(flagValues, performanceClassFlag)
				delete(flagValues, snapshotsVisibleFlag)
			}),
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.IpAcl = nil
				model.PerformanceClass = nil
				model.SnapshotsVisible = nil
			}),
			isValid: true,
		},
		{
			description: "update only snapshots visibility",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, ipAclFlag)
				delete(flagValues, performanceClassFlag)
				delete(flagValues, sizeFlag)
			}),
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.IpAcl = nil
				model.PerformanceClass = nil
				model.SizeGigabytes = nil
			}),
			isValid: true,
		},
		{
			description: "update only performance class",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, ipAclFlag)
				delete(flagValues, snapshotsVisibleFlag)
				delete(flagValues, sizeFlag)
			}),
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.IpAcl = nil
				model.SnapshotsVisible = nil
				model.SizeGigabytes = nil
			}),
			isValid: true,
		},
		{
			description: "update only ipAcl",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, performanceClassFlag)
				delete(flagValues, snapshotsVisibleFlag)
				delete(flagValues, sizeFlag)
			}),
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.PerformanceClass = nil
				model.SnapshotsVisible = nil
				model.SizeGigabytes = nil
			}),
			isValid: true,
		},
		{
			description: "project id missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
				flagValues[sizeFlag] = "50"
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = ""
				flagValues[sizeFlag] = "50"
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = "invalid-uuid"
				flagValues[sizeFlag] = "50"
			}),
			isValid: false,
		},
		{
			description: "resource pool id invalid 1",
			argValues:   []string{""},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "resource pool id invalid 2",
			argValues:   []string{"invalid-uuid"},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "repeated acl flags",
			argValues:   fixtureArgValues(),
			flagValues:  fixtureFlagValues(),
			ipAclValues: []string{"198.51.100.14/24", "198.51.100.14/32"},
			isValid:     true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				if model.IpAcl == nil {
					model.IpAcl = &[]string{}
				}
				*model.IpAcl = append(*model.IpAcl, "198.51.100.14/24", "198.51.100.14/32")
			}),
		},
		{
			description: "repeated ip acl flag with list value",
			argValues:   fixtureArgValues(),
			flagValues:  fixtureFlagValues(),
			ipAclValues: []string{"198.51.100.14/24,198.51.100.14/32"},
			isValid:     true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				if model.IpAcl == nil {
					model.IpAcl = &[]string{}
				}
				*model.IpAcl = append(*model.IpAcl, "198.51.100.14/24", "198.51.100.14/32")
			}),
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
		expectedRequest sfs.ApiUpdateResourcePoolRequest
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
		resp         *sfs.UpdateResourcePoolResponse
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
			name: "empty response",
			args: args{
				resp: &sfs.UpdateResourcePoolResponse{},
			},
			wantErr: false,
		},
		{
			name: "valid response with empty resource pool",
			args: args{
				resp: &sfs.UpdateResourcePoolResponse{
					ResourcePool: &sfs.UpdateResourcePoolResponseResourcePool{},
				},
			},
			wantErr: false,
		},
		{
			name: "valid response with name",
			args: args{
				resp: &sfs.UpdateResourcePoolResponse{
					ResourcePool: &sfs.UpdateResourcePoolResponseResourcePool{
						Name: utils.Ptr("example name"),
					},
				},
			},
			wantErr: false,
		},
	}

	p := print.NewPrinter()
	p.Cmd = NewCmd(&types.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.outputFormat, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
