package orgname

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/resourcemanager/client"
)

// GetOrganizationName returns the name of an organization by its ID.
func GetOrganizationName(ctx context.Context, p *print.Printer, orgId string) (string, error) {
	apiClient, err := client.ConfigureClient(p)
	if err != nil {
		return "", fmt.Errorf("configure resource manager client: %w", err)
	}
	req := apiClient.GetOrganization(ctx, orgId)
	resp, err := req.Execute()
	if err != nil {
		return "", fmt.Errorf("read project details: %w", err)
	}

	return *resp.Name, nil
}
