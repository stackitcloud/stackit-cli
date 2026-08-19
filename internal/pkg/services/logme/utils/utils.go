package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"

	logme "github.com/stackitcloud/stackit-sdk-go/services/logme/v2api"
)

const (
	service = "logme"
)

// Deprecated: resolving a plan by --plan-name/--version will be removed after 2027-02-28. Use --plan-id instead.
func LoadPlanId(planName, version string, offerings *logme.ListOfferingsResponse) (*string, error) {
	availableVersions := ""
	availablePlanNames := ""
	isValidVersion := false
	for _, offer := range offerings.Offerings {
		if !strings.EqualFold(offer.Version, version) {
			availableVersions = fmt.Sprintf("%s\n- %s", availableVersions, offer.Version)
			continue
		}
		isValidVersion = true

		for _, plan := range offer.Plans {
			if strings.EqualFold(plan.Name, planName) {
				return &plan.Id, nil
			}
			availablePlanNames = fmt.Sprintf("%s\n- %s", availablePlanNames, plan.Name)
		}
	}

	if !isValidVersion {
		details := fmt.Sprintf("You provided version %q, which is invalid. Available versions are: %s", version, availableVersions)
		return nil, &errors.DSAInvalidPlanError{
			Service: service,
			Details: details,
		}
	}
	details := fmt.Sprintf("You provided plan_name %q for version %s, which is invalid. Available plan names for that version are: %s", planName, version, availablePlanNames)
	return nil, &errors.DSAInvalidPlanError{
		Service: service,
		Details: details,
	}
}

func GetInstanceName(ctx context.Context, apiClient logme.DefaultAPI, projectId, instanceId, region string) (string, error) {
	resp, err := apiClient.GetInstance(ctx, projectId, region, instanceId).Execute()
	if err != nil {
		return "", fmt.Errorf("get LogMe instance: %w", err)
	}
	return resp.Name, nil
}

func GetCredentialsUsername(ctx context.Context, apiClient logme.DefaultAPI, projectId, instanceId, credentialsId, region string) (string, error) {
	resp, err := apiClient.GetCredentials(ctx, projectId, instanceId, region, credentialsId).Execute()
	if err != nil {
		return "", fmt.Errorf("get LogMe credentials: %w", err)
	}
	return resp.Raw.Credentials.Username, nil
}
