package utils

import (
	"context"
	"fmt"

	scf "github.com/stackitcloud/stackit-sdk-go/services/scf/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func GetOrganizationName(ctx context.Context, apiClient scf.DefaultAPI, projectId, region, organizationId string) (string, error) {
	resp, err := apiClient.GetOrganization(ctx, projectId, region, organizationId).Execute()
	if err != nil {
		return "", fmt.Errorf("get STACKIT Cloud Foundry organization: %w", err)
	}
	if resp == nil {
		return "", fmt.Errorf("get STACKIT Cloud Foundry organization: empty response")
	}
	return resp.Name, nil
}

func GetSpaceName(ctx context.Context, apiClient scf.DefaultAPI, projectId, region, organizationId, spaceId string) (string, error) {
	resp, err := apiClient.GetSpace(ctx, projectId, region, organizationId, spaceId).Execute()
	if err != nil {
		return "", fmt.Errorf("get STACKIT Cloud Foundry space: %w", err)
	}
	if resp == nil {
		return "", fmt.Errorf("get STACKIT Cloud Foundry space: empty response")
	}
	return resp.Name, nil
}

// GetOrganizationLabel returns the name of the organization, falling back to its ID if the name can't be determined
func GetOrganizationLabel(ctx context.Context, p *print.Printer, apiClient scf.DefaultAPI, projectId, region, organizationId string) string {
	name, err := GetOrganizationName(ctx, apiClient, projectId, region, organizationId)
	if err != nil {
		p.Debug(print.ErrorLevel, "get organization name: %v", err)
		return organizationId
	}
	if name == "" {
		return organizationId
	}
	return name
}

// GetSpaceLabel returns the name of the space, falling back to its ID if the name can't be determined
func GetSpaceLabel(ctx context.Context, p *print.Printer, apiClient scf.DefaultAPI, projectId, region, organizationId, spaceId string) string {
	name, err := GetSpaceName(ctx, apiClient, projectId, region, organizationId, spaceId)
	if err != nil {
		p.Debug(print.ErrorLevel, "get space name: %v", err)
		return spaceId
	}
	if name == "" {
		return spaceId
	}
	return name
}

func UserLabel(userGuid, userName *string) string {
	if userGuid != nil {
		return *userGuid
	}
	return utils.PtrString(userName)
}
