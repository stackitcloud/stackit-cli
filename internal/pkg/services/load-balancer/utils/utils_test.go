package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/loadbalancer"
)

var (
	testProjectId      = uuid.NewString()
	testCredentialsRef = "credentials-test"
)

const (
	testCredentialsDisplayName = "name"
)

type loadBalancerClientMocked struct {
	getCredentialsFails bool
	getCredentialsResp  *loadbalancer.GetCredentialsResponse
}

func (m *loadBalancerClientMocked) GetCredentialsExecute(_ context.Context, _, _ string) (*loadbalancer.GetCredentialsResponse, error) {
	if m.getCredentialsFails {
		return nil, fmt.Errorf("could not get credentials")
	}
	return m.getCredentialsResp, nil
}

func TestGetCredentialsDisplayName(t *testing.T) {
	tests := []struct {
		description         string
		getCredentialsFails bool
		getCredentialsResp  *loadbalancer.GetCredentialsResponse
		isValid             bool
		expectedOutput      string
	}{
		{
			description: "base",
			getCredentialsResp: &loadbalancer.GetCredentialsResponse{
				Credential: &loadbalancer.CredentialsResponse{
					DisplayName: utils.Ptr(testCredentialsDisplayName),
				},
			},
			isValid:        true,
			expectedOutput: testCredentialsDisplayName,
		},
		{
			description:         "get credentials fails",
			getCredentialsFails: true,
			isValid:             false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &loadBalancerClientMocked{
				getCredentialsFails: tt.getCredentialsFails,
				getCredentialsResp:  tt.getCredentialsResp,
			}

			output, err := GetCredentialsDisplayName(context.Background(), client, testProjectId, testCredentialsRef)

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
