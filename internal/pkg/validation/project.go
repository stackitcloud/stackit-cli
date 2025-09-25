package validation

import (
	"context"
	"fmt"
	"net/http"

	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/resourcemanager/client"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/core/oapierror"
)

// ValidateProject validates that the project ID is not empty, exists, and the user has access to it.
// It returns the project name for display purposes.
func ValidateProject(ctx context.Context, p *print.Printer, cliVersion string, cmd *cobra.Command, projectId string) (string, error) {
	// Check if project ID is empty
	if projectId == "" {
		return "", &errors.ProjectIdError{}
	}

	// Configure Resource Manager API client
	apiClient, err := client.ConfigureClient(p, cliVersion)
	if err != nil {
		return "", fmt.Errorf("configure resource manager client: %w", err)
	}

	// Try to get project details to validate existence and access
	req := apiClient.GetProject(ctx, projectId)
	resp, err := req.Execute()
	if err != nil {
		// Check for specific HTTP status codes
		if httpErr, ok := err.(*oapierror.GenericOpenAPIError); ok { //nolint:errorlint //complaining that error.As should be used to catch wrapped errors, but this error should not be wrapped
			switch httpErr.StatusCode {
			case http.StatusForbidden:
				// Try to get project name for better error message
				projectLabel := projectId
				if projectName, nameErr := projectname.GetProjectName(ctx, p, cliVersion, cmd); nameErr == nil {
					projectLabel = projectName
				}
				return "", &errors.ProjectNotFoundError{ProjectId: projectId, ProjectLabel: projectLabel}
			case http.StatusUnauthorized:
				return "", &errors.AuthError{}
			}
		}
		return "", fmt.Errorf("validate project: %w", err)
	}

	// Project exists and user has access, returning the project name
	return *resp.Name, nil
}
