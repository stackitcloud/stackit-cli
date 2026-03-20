package update

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/secretsmanager"
)

const (
	testACL1 = "1.2.3.4/24"
	testACL2 = "4.3.2.1/12"
)

var projectIdFlag = globalflags.ProjectIdFlag

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &secretsmanager.APIClient{}

var (
	testProjectId  = uuid.NewString()
	testInstanceId = uuid.NewString()
)

const (
	testKmsKeyId               = "key-id"
	testKmsKeyringId           = "keyring-id"
	testKmsKeyVersion          = int64(1)
	testKmsServiceAccountEmail = "my-service-account-1234567@sa.stackit.cloud"
)

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testInstanceId,
	}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag: testProjectId,
		aclFlag:       testACL1,
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
		},
		InstanceId: testInstanceId,
		Acls:       utils.Ptr([]string{testACL1}),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *secretsmanager.ApiUpdateACLsRequest)) secretsmanager.ApiUpdateACLsRequest {
	request := testClient.UpdateACLs(testCtx, testProjectId, testInstanceId)
	request = request.UpdateACLsPayload(secretsmanager.UpdateACLsPayload{
		Cidrs: utils.Ptr([]secretsmanager.UpdateACLPayload{
			{Cidr: utils.Ptr(testACL1)},
		})})

	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixtureUpdateInstanceRequest(mods ...func(request *secretsmanager.ApiUpdateInstanceRequest)) secretsmanager.ApiUpdateInstanceRequest {
	request := testClient.UpdateInstance(testCtx, testProjectId, testInstanceId)
	request = request.UpdateInstancePayload(secretsmanager.UpdateInstancePayload{
		KmsKey: &secretsmanager.KmsKeyPayload{
			KeyId:               utils.Ptr(testKmsKeyId),
			KeyRingId:           utils.Ptr(testKmsKeyringId),
			KeyVersion:          utils.Ptr(testKmsKeyVersion),
			ServiceAccountEmail: utils.Ptr(testKmsServiceAccountEmail),
		},
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
		aclValues     []string
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
			description: "no update flags",
			argValues:   fixtureArgValues(),
			flagValues: map[string]string{
				projectIdFlag: testProjectId,
			},
			isValid: false,
		},
		{
			description: "zero values",
			argValues:   fixtureArgValues(),
			flagValues: map[string]string{
				projectIdFlag: testProjectId,
				aclFlag:       "",
			},
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Acls = &[]string{}
			}),
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
			description: "instance id invalid 1",
			argValues:   []string{""},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "instance id invalid 2",
			argValues:   []string{"invalid-uuid"},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "kms key id without other required kms flags",
			argValues:   fixtureArgValues(),
			flagValues: map[string]string{
				projectIdFlag: testProjectId,
				kmsKeyIdFlag:  "key-id",
			},
			isValid: false,
		},
		{
			description: "acl flag conflicts with kms flags",
			argValues:   fixtureArgValues(),
			flagValues: map[string]string{
				projectIdFlag:              testProjectId,
				aclFlag:                    testACL1,
				kmsKeyIdFlag:               "key-id",
				kmsKeyringIdFlag:           "keyring-id",
				kmsKeyVersionFlag:          "1",
				kmsServiceAccountEmailFlag: "my-service-account-1234567@sa.stackit.cloud",
			},
			isValid: false,
		},
		{
			description: "repeated acl flags",
			argValues:   fixtureArgValues(),
			flagValues:  fixtureFlagValues(),
			aclValues:   []string{testACL1, testACL1},
			isValid:     true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Acls = utils.Ptr(
					append(*model.Acls, testACL1, testACL1))
			}),
		},
		{
			description: "repeated acl flag with list value",
			argValues:   fixtureArgValues(),
			flagValues:  fixtureFlagValues(),
			aclValues:   []string{"198.51.100.14/24,198.51.100.14/32"},
			isValid:     true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Acls = utils.Ptr(
					append(*model.Acls, "198.51.100.14/24", "198.51.100.14/32"),
				)
			}),
		},
		{
			description: "kms flags",
			argValues:   fixtureArgValues(),
			flagValues: map[string]string{
				projectIdFlag:              testProjectId,
				kmsKeyIdFlag:               testKmsKeyId,
				kmsKeyringIdFlag:           testKmsKeyringId,
				kmsKeyVersionFlag:          "1",
				kmsServiceAccountEmailFlag: testKmsServiceAccountEmail,
			},
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Acls = nil
				model.KmsKeyId = utils.Ptr(testKmsKeyId)
				model.KmsKeyringId = utils.Ptr(testKmsKeyringId)
				model.KmsKeyVersion = utils.Ptr(testKmsKeyVersion)
				model.KmsServiceAccountEmail = utils.Ptr(testKmsServiceAccountEmail)
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInputWithAdditionalFlags(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, map[string][]string{
				aclFlag: tt.aclValues,
			}, tt.isValid)
		})
	}
}
func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest secretsmanager.ApiUpdateACLsRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "multiple ACLs",
			model: fixtureInputModel(func(model *inputModel) {
				*model.Acls = append(*model.Acls, testACL2)
			}),
			expectedRequest: fixtureRequest().UpdateACLsPayload(secretsmanager.UpdateACLsPayload{
				Cidrs: utils.Ptr([]secretsmanager.UpdateACLPayload{
					{Cidr: utils.Ptr(testACL1)},
					{Cidr: utils.Ptr(testACL2)},
				})}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)
			aclRequest, ok := request.(secretsmanager.ApiUpdateACLsRequest)
			if !ok {
				t.Fatalf("expected ACL update request, got %T", request)
			}

			diff := cmp.Diff(aclRequest, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestBuildRequestKms(t *testing.T) {
	model := fixtureInputModel(func(model *inputModel) {
		model.Acls = nil
		model.KmsKeyId = utils.Ptr(testKmsKeyId)
		model.KmsKeyringId = utils.Ptr(testKmsKeyringId)
		model.KmsKeyVersion = utils.Ptr(testKmsKeyVersion)
		model.KmsServiceAccountEmail = utils.Ptr(testKmsServiceAccountEmail)
	})

	request := buildRequest(testCtx, model, testClient)
	updateRequest, ok := request.(secretsmanager.ApiUpdateInstanceRequest)
	if !ok {
		t.Fatalf("expected instance update request, got %T", request)
	}

	expectedRequest := fixtureUpdateInstanceRequest()
	diff := cmp.Diff(updateRequest, expectedRequest,
		cmp.AllowUnexported(expectedRequest),
		cmpopts.EquateComparable(testCtx),
	)
	if diff != "" {
		t.Fatalf("Data does not match: %s", diff)
	}
}
