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

type ArgusClient interface {
	GetInstanceExecute(ctx context.Context, instanceId, projectId string) (*argus.GetInstanceResponse, error)
	GetGrafanaConfigsExecute(ctx context.Context, instanceId, projectId string) (*argus.GrafanaConfigs, error)
	UpdateGrafanaConfigs(ctx context.Context, instanceId string, projectId string) argus.ApiUpdateGrafanaConfigsRequest
}

func ValidatePlanId(planId string, resp *argus.PlansResponse) error {
	if resp == nil {
		return fmt.Errorf("no Argus plans provided")
	}

	for i := range *resp.Plans {
		plan := (*resp.Plans)[i]
		if plan.Id != nil && strings.EqualFold(*plan.Id, planId) {
			return nil
		}
	}

	return &errors.ArgusInvalidPlanError{
		Service: service,
		Details: fmt.Sprintf("You provided plan ID %q, which is invalid.", planId),
	}
}

func LoadPlanId(planName string, resp *argus.PlansResponse) (*string, error) {
	availablePlanNames := ""
	if resp == nil {
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

	details := fmt.Sprintf("You provided plan name %q, which is invalid. Available plan names are: %s", planName, availablePlanNames)
	return nil, &errors.ArgusInvalidPlanError{
		Service: service,
		Details: details,
	}
}

func GetInstanceName(ctx context.Context, apiClient ArgusClient, instanceId, projectId string) (string, error) {
	resp, err := apiClient.GetInstanceExecute(ctx, instanceId, projectId)
	if err != nil {
		return "", fmt.Errorf("get Argus instance: %w", err)
	}
	return *resp.Name, nil
}

func ToPayloadGenericOAuth(respOAuth *argus.GrafanaOauth) *argus.UpdateGrafanaConfigsPayloadGenericOauth {
	if respOAuth == nil {
		return nil
	}
	return &argus.UpdateGrafanaConfigsPayloadGenericOauth{
		ApiUrl:              respOAuth.ApiUrl,
		AuthUrl:             respOAuth.AuthUrl,
		Enabled:             respOAuth.Enabled,
		Name:                respOAuth.Name,
		OauthClientId:       respOAuth.OauthClientId,
		OauthClientSecret:   respOAuth.OauthClientSecret,
		RoleAttributePath:   respOAuth.RoleAttributePath,
		RoleAttributeStrict: respOAuth.RoleAttributeStrict,
		Scopes:              respOAuth.Scopes,
		TokenUrl:            respOAuth.TokenUrl,
		UsePkce:             respOAuth.UsePkce,
	}
}

func GetPartialUpdateGrafanaConfigsPayload(ctx context.Context, apiClient ArgusClient, instanceId, projectId string, singleSignOn, publicReadAccess *bool) (*argus.UpdateGrafanaConfigsPayload, error) {
	currentConfigs, err := apiClient.GetGrafanaConfigsExecute(ctx, instanceId, projectId)
	if err != nil {
		return nil, fmt.Errorf("get current Grafana configs: %w", err)
	}
	if currentConfigs == nil {
		return nil, fmt.Errorf("no Grafana configs found for instance %q", instanceId)
	}

	payload := &argus.UpdateGrafanaConfigsPayload{
		GenericOauth:     ToPayloadGenericOAuth(currentConfigs.GenericOauth),
		PublicReadAccess: currentConfigs.PublicReadAccess,
		UseStackitSso:    currentConfigs.UseStackitSso,
	}

	if singleSignOn != nil {
		payload.UseStackitSso = singleSignOn
	}
	if publicReadAccess != nil {
		payload.PublicReadAccess = publicReadAccess
	}

	return payload, nil
}
