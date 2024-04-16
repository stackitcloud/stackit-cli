package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/stackitcloud/stackit-sdk-go/services/argus"
)

const (
	service = "argus"
)

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

func MapToUpdateScrapeConfigPayload(resp *argus.GetScrapeConfigResponse) *argus.UpdateScrapeConfigPayload {
	data := resp.Data

	basicAuth := mapBasicAuth(data.BasicAuth)

	httpSdConfigs := make([]argus.CreateScrapeConfigPayloadHttpSdConfigsInner, 0)
	if data.HttpSdConfigs != nil {
		for _, config := range *data.HttpSdConfigs {
			httpSdConfigs = append(httpSdConfigs, mapHttpSdConfig(config))
		}
	}

	staticConfigs := make([]argus.UpdateScrapeConfigPayloadStaticConfigsInner, 0)
	if data.StaticConfigs != nil {
		for _, config := range *data.StaticConfigs {
			staticConfigs = append(staticConfigs, mapStaticConfig(config))
		}
	}

	tlsConfig := mapTlsConfig(data.TlsConfig)

	return &argus.UpdateScrapeConfigPayload{
		BasicAuth:       basicAuth,
		BearerToken:     data.BearerToken,
		HonorLabels:     data.HonorLabels,
		HonorTimeStamps: data.HonorTimeStamps,
		MetricsPath:     data.MetricsPath,
		// MetricsRelabelConfigs: metricsRelabelConfigs,
		// Params: 	   convertMapStringToInterface(data.Params),
		SampleLimit:    utils.Ptr(float64(*data.SampleLimit)),
		Scheme:         data.Scheme,
		ScrapeInterval: data.ScrapeInterval,
		ScrapeTimeout:  data.ScrapeTimeout,
		StaticConfigs:  &staticConfigs,
		TlsConfig:      tlsConfig,
	}
}

func mapOAuth2(oauth2 *argus.OAuth2) *argus.CreateScrapeConfigPayloadHttpSdConfigsInnerOauth2 {
	return nil
}

func mapHttpSdConfig(httpSdConfig argus.HTTPServiceSD) argus.CreateScrapeConfigPayloadHttpSdConfigsInner {
	oauth2 := mapOAuth2(httpSdConfig.Oauth2)
	tlsConfig := mapTlsConfig(httpSdConfig.TlsConfig)
	return argus.CreateScrapeConfigPayloadHttpSdConfigsInner{
		Oauth2:    oauth2,
		TlsConfig: tlsConfig,
	}

}

func mapStaticConfig(staticConfig argus.StaticConfigs) argus.UpdateScrapeConfigPayloadStaticConfigsInner {
	labels := convertMapStringToInterface(staticConfig.Labels)
	return argus.UpdateScrapeConfigPayloadStaticConfigsInner{
		Labels:  labels,
		Targets: staticConfig.Targets,
	}
}

func mapBasicAuth(basicAuth *argus.BasicAuth) *argus.CreateScrapeConfigPayloadBasicAuth {
	if basicAuth == nil {
		return nil
	}

	return &argus.CreateScrapeConfigPayloadBasicAuth{
		Password: basicAuth.Password,
		Username: basicAuth.Username,
	}
}

func mapTlsConfig(tlsConfig *argus.TLSConfig) *argus.CreateScrapeConfigPayloadHttpSdConfigsInnerOauth2TlsConfig {
	if tlsConfig == nil {
		return nil
	}

	return &argus.CreateScrapeConfigPayloadHttpSdConfigsInnerOauth2TlsConfig{
		InsecureSkipVerify: tlsConfig.InsecureSkipVerify,
	}
}

func convertMapStringToInterface(m *map[string]string) *map[string]interface{} {
	if m == nil {
		return nil
	}

	newMap := make(map[string]interface{})
	for k, v := range *m {
		newMap[k] = v
	}
	return &newMap
}

type ArgusClient interface {
	GetInstanceExecute(ctx context.Context, instanceId, projectId string) (*argus.GetInstanceResponse, error)
}

func GetInstanceName(ctx context.Context, apiClient ArgusClient, instanceId, projectId string) (string, error) {
	resp, err := apiClient.GetInstanceExecute(ctx, instanceId, projectId)
	if err != nil {
		return "", fmt.Errorf("get Argus instance: %w", err)
	}
	return *resp.Name, nil
}

func GetDefaultCreateScrapeConfigPayload(ctx context.Context, apiClient ArgusClient) (*argus.CreateScrapeConfigPayload, error) {
	return nil, nil
}
