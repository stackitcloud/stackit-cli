package utils

import (
	"context"
	"net/http"

	"github.com/stackitcloud/stackit-sdk-go/core/oapierror"
	"github.com/stackitcloud/stackit-sdk-go/services/serviceenablement"
	"github.com/stackitcloud/stackit-sdk-go/services/serviceenablement/wait"
)

const (
	SKEServiceId = "cloud.stackit.ske"
)

type ServiceEnablementClient interface {
	GetServiceStatusExecute(ctx context.Context, projectId string, serviceId string) (*serviceenablement.ServiceStatus, error)
}

func ProjectEnabled(ctx context.Context, apiClient ServiceEnablementClient, projectId string) (bool, error) {
	project, err := apiClient.GetServiceStatusExecute(ctx, projectId, SKEServiceId)
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
	return *project.State == wait.ServiceStateEnabled, nil
}
