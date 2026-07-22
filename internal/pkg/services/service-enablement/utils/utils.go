package utils

import (
	"context"
	"net/http"

	"github.com/stackitcloud/stackit-sdk-go/core/oapierror"
	serviceenablement "github.com/stackitcloud/stackit-sdk-go/services/serviceenablement/v2api"
)

const (
	SKEServiceId = "cloud.stackit.ske"
)

func ProjectEnabled(ctx context.Context, apiClient serviceenablement.DefaultAPI, projectId, region string) (bool, error) {
	project, err := apiClient.GetServiceStatusRegional(ctx, region, projectId, SKEServiceId).Execute()
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
