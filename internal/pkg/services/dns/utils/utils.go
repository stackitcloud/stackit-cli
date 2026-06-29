package utils

import (
	"context"
	"fmt"
	"math"

	dns "github.com/stackitcloud/stackit-sdk-go/services/dns/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func GetZoneName(ctx context.Context, apiClient dns.DefaultAPI, projectId, zoneId string) (string, error) {
	resp, err := apiClient.GetZone(ctx, projectId, zoneId).Execute()
	if err != nil {
		return "", fmt.Errorf("get DNS zone: %w", err)
	}
	return resp.Zone.Name, nil
}

func GetRecordSetName(ctx context.Context, apiClient dns.DefaultAPI, projectId, zoneId, recordSetId string) (string, error) {
	resp, err := apiClient.GetRecordSet(ctx, projectId, zoneId, recordSetId).Execute()
	if err != nil {
		return "", fmt.Errorf("get DNS recordset: %w", err)
	}
	return resp.Rrset.Name, nil
}

func GetRecordSetType(ctx context.Context, apiClient dns.DefaultAPI, projectId, zoneId, recordSetId string) (*string, error) {
	resp, err := apiClient.GetRecordSet(ctx, projectId, zoneId, recordSetId).Execute()
	if err != nil {
		return utils.Ptr(""), fmt.Errorf("get DNS recordset: %w", err)
	}
	return utils.Ptr(string(resp.Rrset.Type)), nil
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
