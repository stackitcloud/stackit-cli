package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"

	"github.com/stackitcloud/stackit-sdk-go/services/opensearch"
)

func ValidatePlanId(service, planId string, offerings *opensearch.ListOfferingsResponse) error {
	for _, offer := range *offerings.Offerings {
		for _, plan := range *offer.Plans {
			if plan.Id != nil && strings.EqualFold(*plan.Id, planId) {
				return nil
			}
		}
	}

	return &errors.DSAInvalidPlanError{
		Service: service,
		Details: fmt.Sprintf("You provided plan ID %q, which is invalid.", planId),
	}
}

func LoadPlanId(planName, version string, offerings *opensearch.ListOfferingsResponse) (*string, error) {
	service := "opensearch"

	availableVersions := ""
	availablePlanNames := ""
	isValidVersion := false
	for _, offer := range *offerings.Offerings {
		if !strings.EqualFold(*offer.Version, version) {
			availableVersions = fmt.Sprintf("%s\n- %s", availableVersions, *offer.Version)
			continue
		}
		isValidVersion = true

		for _, plan := range *offer.Plans {
			if plan.Name == nil {
				continue
			}
			if strings.EqualFold(*plan.Name, planName) && plan.Id != nil {
				return plan.Id, nil
			}
			availablePlanNames = fmt.Sprintf("%s\n- %s", availablePlanNames, *plan.Name)
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

type OpenSearchClient interface {
	GetInstanceExecute(ctx context.Context, projectId, instanceId string) (*opensearch.Instance, error)
	GetCredentialsExecute(ctx context.Context, projectId, instanceId, credentialsId string) (*opensearch.CredentialsResponse, error)
}

func GetInstanceName(ctx context.Context, apiClient OpenSearchClient, projectId, instanceId string) (string, error) {
	resp, err := apiClient.GetInstanceExecute(ctx, projectId, instanceId)
	if err != nil {
		return "", fmt.Errorf("get OpenSearch instance: %w", err)
	}
	return *resp.Name, nil
}

func GetCredentialsUsername(ctx context.Context, apiClient OpenSearchClient, projectId, instanceId, credentialsId string) (string, error) {
	resp, err := apiClient.GetCredentialsExecute(ctx, projectId, instanceId, credentialsId)
	if err != nil {
		return "", fmt.Errorf("get OpenSearch credentials: %w", err)
	}
	return *resp.Raw.Credentials.Username, nil
}
