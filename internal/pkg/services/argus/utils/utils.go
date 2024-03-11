package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"

	"github.com/stackitcloud/stackit-sdk-go/services/argus"
)

const (
	service = "argus"
)

func ValidatePlanId(planId string, resp *argus.PlansResponse) error {
	if resp == nil {
		// Should not happen, check is done before calling this function
		return fmt.Errorf("no Argus plans provided")
	}

	for i := range *resp.Plans {
		plan := (*resp.Plans)[i]
		if plan.Id != nil && strings.EqualFold(*plan.Id, planId) {
			return nil
		}
	}

	return &errors.DSAInvalidPlanError{
		Service: service,
		Details: fmt.Sprintf("You provided plan ID %q, which is invalid.", planId),
	}
}

func LoadPlanId(planName string, resp *argus.PlansResponse) (*string, error) {
	availablePlanNames := ""
	if resp == nil {
		// Should not happen, check is done before calling this function
		return nil, fmt.Errorf("no Argus plans provided")
	}

	for i := range *resp.Plans {
		plan := (*resp.Plans)[i]
		if plan.Name == nil {
			continue
		}
		if strings.EqualFold(*plan.Name, planName) && plan.Id != nil {
			return plan.Id, nil
		}
		availablePlanNames = fmt.Sprintf("%s\n- %s", availablePlanNames, *plan.Name)
	}

	details := fmt.Sprintf("You provided plan_name %q, which is invalid. Available plan names are: %s", planName, availablePlanNames)
	return nil, &errors.DSAInvalidPlanError{
		Service: service,
		Details: details,
	}
}

type ArgusClient interface {
	GetInstanceExecute(ctx context.Context, projectId, instanceId string) (*argus.Instance, error)
}

func GetInstanceName(ctx context.Context, apiClient ArgusClient, projectId, instanceId string) (string, error) {
	resp, err := apiClient.GetInstanceExecute(ctx, projectId, instanceId)
	if err != nil {
		return "", fmt.Errorf("get Argus instance: %w", err)
	}
	return *resp.Name, nil
}
