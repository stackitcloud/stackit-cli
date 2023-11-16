package utils

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-sdk-go/services/dns"
)

type DNSClient interface {
	GetZoneExecute(ctx context.Context, projectId, zoneId string) (*dns.ZoneResponse, error)
	GetRecordSetExecute(ctx context.Context, projectId, zoneId, recordSetId string) (*dns.RecordSetResponse, error)
}

func GetZoneName(ctx context.Context, apiClient DNSClient, projectId, zoneId string) (string, error) {
	resp, err := apiClient.GetZoneExecute(ctx, projectId, zoneId)
	if err != nil {
		return "", fmt.Errorf("get DNS zone: %w", err)
	}
	return *resp.Zone.Name, nil
}

func GetRecordSetName(ctx context.Context, apiClient DNSClient, projectId, zoneId, recordSetId string) (string, error) {
	resp, err := apiClient.GetRecordSetExecute(ctx, projectId, zoneId, recordSetId)
	if err != nil {
		return "", fmt.Errorf("get DNS recordset: %w", err)
	}
	return *resp.Rrset.Name, nil
}
