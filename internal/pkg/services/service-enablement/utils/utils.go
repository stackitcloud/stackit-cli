package utils

import (
	"context"
	"net/http"

	"github.com/stackitcloud/stackit-sdk-go/core/oapierror"
	"github.com/stackitcloud/stackit-sdk-go/services/serviceenablement"
)

const (
	SKEServiceId = "cloud.stackit.ske"
)

type ServiceEnablementClient interface {
	GetServiceStatusRegionalExecute(ctx context.Context, region, projectId, serviceId string) (*serviceenablement.ServiceStatus, error)
}

func ProjectEnabled(ctx context.Context, apiClient ServiceEnablementClient, projectId, region string) (bool, error) {
	project, err := apiClient.GetServiceStatusRegionalExecute(ctx, region, projectId, SKEServiceId)
	if err != nil {
		oapiErr, ok := err.(*oapierror.GenericOpenAPIError) //nolint:errorlint //complaining that error.As should be used to catch wrapped errors, but this error should not be wrapped
		if !ok {
			return false, err
		}
		if oapiErr.StatusCode == http.StatusNotFound {
			return false, nil
		}
		return false, err
	}
	return *project.State == serviceenablement.SERVICESTATUSSTATE_ENABLED, nil
}
