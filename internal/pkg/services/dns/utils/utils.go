package utils

import (
	"context"
	"fmt"
	"math"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
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

func GetRecordSetType(ctx context.Context, apiClient DNSClient, projectId, zoneId, recordSetId string) (*string, error) {
	resp, err := apiClient.GetRecordSetExecute(ctx, projectId, zoneId, recordSetId)
	if err != nil {
		return utils.Ptr(""), fmt.Errorf("get DNS recordset: %w", err)
	}
	return (*string)(resp.Rrset.Type), nil
}

func FormatTxtRecord(input string) (string, error) {
	length := float64(len(input))
	if length <= 255 {
		return input, nil
	}
	// Max length with quotes and white spaces is 4096. Without the quotes and white spaces the max length is 4049
	if length > 4049 {
		return "", fmt.Errorf("max input length is 4049. The length of the input is %v", length)
	}

	result := ""
	chunks := int(math.Ceil(length / 255))
	for i := range chunks {
		skip := 255 * i
		if i == chunks-1 {
			// Append the left record content
			result += fmt.Sprintf("%q", input[0+skip:])
		} else {
			// Add 255 characters of the record data quoted to the result
			result += fmt.Sprintf("%q ", input[0+skip:255+skip])
		}
	}

	return result, nil
}
