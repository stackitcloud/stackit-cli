package utils

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	kms "github.com/stackitcloud/stackit-sdk-go/services/kms/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

var (
	testProjectId     = uuid.NewString()
	testKeyRingId     = uuid.NewString()
	testKeyId         = uuid.NewString()
	testWrappingKeyId = uuid.NewString()
)

const (
	testRegion          = "eu01"
	testKeyName         = "my-test-key"
	testKeyRingName     = "my-key-ring"
	testWrappingKeyName = "my-wrapping-key"
)

type kmsClientMocked struct {
	getKeyFails         bool
	getKeyResp          *kms.Key
	getKeyRingFails     bool
	getKeyRingResp      *kms.KeyRing
	getWrappingKeyFails bool
	getWrappingKeyResp  *kms.WrappingKey
}

func (m *kmsClientMocked) newMock() kms.DefaultAPI {
	return kms.DefaultAPIServiceMock{
		GetKeyExecuteMock: utils.Ptr(func(_ kms.ApiGetKeyRequest) (*kms.Key, error) {
			if m.getKeyFails {
				return nil, fmt.Errorf("could not get key")
			}
			return m.getKeyResp, nil
		}),
		GetKeyRingExecuteMock: utils.Ptr(func(_ kms.ApiGetKeyRingRequest) (*kms.KeyRing, error) {
			if m.getKeyRingFails {
				return nil, fmt.Errorf("could not get key ring")
			}
			return m.getKeyRingResp, nil
		}),
		GetWrappingKeyExecuteMock: utils.Ptr(func(_ kms.ApiGetWrappingKeyRequest) (*kms.WrappingKey, error) {
			if m.getWrappingKeyFails {
				return nil, fmt.Errorf("could not get wrapping key")
			}
			return m.getWrappingKeyResp, nil
		}),
	}
}

func TestGetKeyName(t *testing.T) {
	keyName := testKeyName

	tests := []struct {
		description    string
		getKeyFails    bool
		getKeyResp     *kms.Key
		isValid        bool
		expectedOutput string
	}{
		{
			description: "base",
			getKeyResp: &kms.Key{
				DisplayName: keyName,
			},
			isValid:        true,
			expectedOutput: testKeyName,
		},
		{
			description: "get key fails",
			getKeyFails: true,
			isValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &kmsClientMocked{
				getKeyFails: tt.getKeyFails,
				getKeyResp:  tt.getKeyResp,
			}

			output, err := GetKeyName(context.Background(), client.newMock(), testProjectId, testRegion, testKeyRingId, testKeyId)

			if tt.isValid && err != nil {
				t.Errorf("failed on valid input: %v", err)
			}
			if !tt.isValid && err == nil {
				t.Errorf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}
			if output != tt.expectedOutput {
				t.Errorf("expected output to be %q, got %q", tt.expectedOutput, output)
			}
		})
	}
}

// TestGetKeyDeletionDate tests the GetKeyDeletionDate function.
func TestGetKeyDeletionDate(t *testing.T) {
	mockTime := time.Date(2025, 8, 20, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		description    string
		getKeyFails    bool
		getKeyResp     *kms.Key
		isValid        bool
		expectedOutput time.Time
	}{
		{
			description: "base",
			getKeyResp: &kms.Key{
				DeletionDate: &mockTime,
			},
			isValid:        true,
			expectedOutput: mockTime,
		},
		{
			description: "get key fails",
			getKeyFails: true,
			isValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &kmsClientMocked{
				getKeyFails: tt.getKeyFails,
				getKeyResp:  tt.getKeyResp,
			}

			output, err := GetKeyDeletionDate(context.Background(), client.newMock(), testProjectId, testRegion, testKeyRingId, testKeyId)

			if tt.isValid && err != nil {
				t.Errorf("failed on valid input: %v", err)
			}
			if !tt.isValid && err == nil {
				t.Errorf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}
			if !output.Equal(tt.expectedOutput) {
				t.Errorf("expected output to be %v, got %v", tt.expectedOutput, output)
			}
		})
	}
}

// TestGetKeyRingName tests the GetKeyRingName function.
func TestGetKeyRingName(t *testing.T) {
	keyRingName := testKeyRingName

	tests := []struct {
		description     string
		getKeyRingFails bool
		getKeyRingResp  *kms.KeyRing
		isValid         bool
		expectedOutput  string
	}{
		{
			description: "base",
			getKeyRingResp: &kms.KeyRing{
				DisplayName: keyRingName,
			},
			isValid:        true,
			expectedOutput: testKeyRingName,
		},
		{
			description:     "get key ring fails",
			getKeyRingFails: true,
			isValid:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &kmsClientMocked{
				getKeyRingFails: tt.getKeyRingFails,
				getKeyRingResp:  tt.getKeyRingResp,
			}

			output, err := GetKeyRingName(context.Background(), client.newMock(), testProjectId, testKeyRingId, testRegion)

			if tt.isValid && err != nil {
				t.Errorf("failed on valid input: %v", err)
			}
			if !tt.isValid && err == nil {
				t.Errorf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}
			if output != tt.expectedOutput {
				t.Errorf("expected output to be %q, got %q", tt.expectedOutput, output)
			}
		})
	}
}

func TestGetWrappingKeyName(t *testing.T) {
	wrappingKeyName := testWrappingKeyName
	tests := []struct {
		description         string
		getWrappingKeyFails bool
		getWrappingKeyResp  *kms.WrappingKey
		isValid             bool
		expectedOutput      string
	}{
		{
			description: "base",
			getWrappingKeyResp: &kms.WrappingKey{
				DisplayName: wrappingKeyName,
			},
			isValid:        true,
			expectedOutput: testWrappingKeyName,
		},
		{
			description:         "get wrapping key fails",
			getWrappingKeyFails: true,
			isValid:             false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &kmsClientMocked{
				getWrappingKeyFails: tt.getWrappingKeyFails,
				getWrappingKeyResp:  tt.getWrappingKeyResp,
			}

			output, err := GetWrappingKeyName(context.Background(), client.newMock(), testProjectId, testRegion, testKeyRingId, testWrappingKeyId)

			if tt.isValid && err != nil {
				t.Errorf("failed on valid input: %v", err)
			}
			if !tt.isValid && err == nil {
				t.Errorf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}
			if output != tt.expectedOutput {
				t.Errorf("expected output to be %q, got %q", tt.expectedOutput, output)
			}
		})
	}
}
