package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/dns"
)

var (
	testProjectId   = uuid.NewString()
	testZoneId      = uuid.NewString()
	testRecordSetId = uuid.NewString()
)

const (
	testZoneName      = "zone"
	testRecordSetName = "record-set"
)

type dnsClientMocked struct {
	getZoneFails      bool
	getZoneResp       *dns.ZoneResponse
	getRecordSetFails bool
	getRecordSetResp  *dns.RecordSetResponse
}

func (m *dnsClientMocked) GetZoneExecute(_ context.Context, _, _ string) (*dns.ZoneResponse, error) {
	if m.getZoneFails {
		return nil, fmt.Errorf("could not get zone")
	}
	return m.getZoneResp, nil
}

func (m *dnsClientMocked) GetRecordSetExecute(_ context.Context, _, _, _ string) (*dns.RecordSetResponse, error) {
	if m.getRecordSetFails {
		return nil, fmt.Errorf("could not get record set")
	}
	return m.getRecordSetResp, nil
}

func TestGetZoneName(t *testing.T) {
	tests := []struct {
		description    string
		getZoneFails   bool
		getZoneResp    *dns.ZoneResponse
		isValid        bool
		expectedOutput string
	}{
		{
			description: "base",
			getZoneResp: &dns.ZoneResponse{
				Zone: &dns.Zone{
					Name: utils.Ptr(testZoneName),
				},
			},
			isValid:        true,
			expectedOutput: testZoneName,
		},
		{
			description:  "get zone fails",
			getZoneFails: true,
			isValid:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &dnsClientMocked{
				getZoneFails: tt.getZoneFails,
				getZoneResp:  tt.getZoneResp,
			}

			output, err := GetZoneName(context.Background(), client, testProjectId, testZoneId)

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

func TestGetRecordSetName(t *testing.T) {
	tests := []struct {
		description       string
		getRecordSetFails bool
		getRecordSetResp  *dns.RecordSetResponse
		isValid           bool
		expectedOutput    string
	}{
		{
			description: "base",
			getRecordSetResp: &dns.RecordSetResponse{
				Rrset: &dns.RecordSet{
					Name: utils.Ptr(testRecordSetName),
				},
			},
			isValid:        true,
			expectedOutput: testRecordSetName,
		},
		{
			description:       "get record set fails",
			getRecordSetFails: true,
			isValid:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &dnsClientMocked{
				getRecordSetFails: tt.getRecordSetFails,
				getRecordSetResp:  tt.getRecordSetResp,
			}

			output, err := GetRecordSetName(context.Background(), client, testProjectId, testZoneId, testRecordSetId)

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
