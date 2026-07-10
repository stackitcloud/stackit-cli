package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"

	opensearch "github.com/stackitcloud/stackit-sdk-go/services/opensearch/v2api"
)

const (
	service = "opensearch"
)

func ValidatePlanId(planId string, offerings *opensearch.ListOfferingsResponse) error {
	for _, offer := range offerings.Offerings {
		for _, plan := range offer.Plans {
			if strings.EqualFold(plan.Id, planId) {
				return nil
			}
		}
	}

	return &errors.DSAInvalidPlanError{
		Service: service,
		Details: fmt.Sprintf("You provided plan ID %q, which is invalid.", planId),
	}
}

func LoadPlanId(planName, version string, offerings *opensearch.ListOfferingsResponse) (string, error) {
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
				return plan.Id, nil
			}
			availablePlanNames = fmt.Sprintf("%s\n- %s", availablePlanNames, plan.Name)
		}
	}

	if !isValidVersion {
		details := fmt.Sprintf("You provided version %q, which is invalid. Available versions are: %s", version, availableVersions)
		return "", &errors.DSAInvalidPlanError{
			Service: service,
			Details: details,
		}
	}
	details := fmt.Sprintf("You provided plan_name %q for version %s, which is invalid. Available plan names for that version are: %s", planName, version, availablePlanNames)
	return "", &errors.DSAInvalidPlanError{
		Service: service,
		Details: details,
	}
}

type OpenSearchClient interface {
	GetInstance(ctx context.Context, projectId, region, instanceId string) opensearch.ApiGetInstanceRequest
	GetCredentials(ctx context.Context, projectId, region, instanceId, credentialsId string) opensearch.ApiGetCredentialsRequest
}

func GetInstanceName(ctx context.Context, apiClient OpenSearchClient, projectId, region, instanceId string) (string, error) {
	resp, err := apiClient.GetInstance(ctx, projectId, region, instanceId).Execute()
	if err != nil {
		return "", fmt.Errorf("get OpenSearch instance: %w", err)
	}
	return resp.Name, nil
}

func GetCredentialsUsername(ctx context.Context, apiClient OpenSearchClient, projectId, region, instanceId, credentialsId string) (string, error) {
	resp, err := apiClient.GetCredentials(ctx, projectId, instanceId, region, credentialsId).Execute()
	if err != nil {
		return "", fmt.Errorf("get OpenSearch credentials: %w", err)
	}
	return resp.Raw.Credentials.Username, nil
}
