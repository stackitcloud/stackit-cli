package reconcile

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
)

type testCtxKey struct{}

const (
	testRegion      = "eu01"
	testClusterName = "my-cluster"
)

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &ske.APIClient{}
var testProjectId = uuid.NewString()

func fixtureArgValues(mods ...func([]string)) []string {
	argValues := []string{
		testClusterName,
	}
	for _, m := range mods {
		m(argValues)
	}
	return argValues
}

func fixtureFlagValues(mods ...func(map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
	}
	for _, m := range mods {
		m(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(*inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Region:    testRegion,
			Verbosity: globalflags.VerbosityDefault,
		},
		ClusterName: testClusterName,
	}
	for _, m := range mods {
		m(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *ske.ApiTriggerReconcileRequest)) ske.ApiTriggerHibernateRequest {
	request := testClient.TriggerReconcile(testCtx, testProjectId, testRegion, testClusterName)
	for _, m := range mods {
		m(&request)
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
			description: "missing project id",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(fv map[string]string) {
				delete(fv, globalflags.ProjectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "invalid project id - empty string",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(fv map[string]string) {
				fv[globalflags.ProjectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "invalid uuid format",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(fv map[string]string) {
				fv[globalflags.ProjectIdFlag] = "not-a-uuid"
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			p := print.NewPrinter()
			cmd := &cobra.Command{}
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

			if len(tt.argValues) == 0 {
				_, err := parseInput(p, cmd, tt.argValues)
				if err == nil && !tt.isValid {
					t.Fatalf("expected error due to missing args")
				}
				return
			}

			model, err := parseInput(p, cmd, tt.argValues)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error parsing input: %v", err)
			}

			if !tt.isValid {
				t.Fatalf("did not fail on invalid input")
			}

			diff := cmp.Diff(model, tt.expectedModel)
			if diff != "" {
				t.Fatalf("input model mismatch:\n%s", diff)
			}
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest ske.ApiTriggerHibernateRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			got := buildRequest(testCtx, tt.model, testClient)
			want := tt.expectedRequest

			diff := cmp.Diff(got, want,
				cmpopts.EquateComparable(testCtx),
				cmp.AllowUnexported(want),
			)
			if diff != "" {
				t.Fatalf("request mismatch:\n%s", diff)
			}
		})
	}
}
