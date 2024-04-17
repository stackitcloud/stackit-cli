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

func MapToUpdateScrapeConfigPayload(resp *argus.GetScrapeConfigResponse) (*argus.UpdateScrapeConfigPayload, error) {
	if resp == nil || resp.Data == nil {
		return nil, fmt.Errorf("no Argus scrape config provided")
	}

	data := resp.Data

	basicAuth := mapBasicAuth(data.BasicAuth)

	var staticConfigs *[]argus.UpdateScrapeConfigPayloadStaticConfigsInner
	if data.StaticConfigs != nil {
		configs := make([]argus.UpdateScrapeConfigPayloadStaticConfigsInner, 0)
		for _, config := range *data.StaticConfigs {
			newConfig, err := mapStaticConfig(config)
			if err != nil {
				return nil, fmt.Errorf("map static config: %w", err)
			}
			configs = append(configs, newConfig)
		}
		staticConfigs = &configs
	}

	tlsConfig := mapTlsConfig(data.TlsConfig)

	var metricsRelabelConfigs *[]argus.CreateScrapeConfigPayloadMetricsRelabelConfigsInner
	if data.MetricsRelabelConfigs != nil {
		configs := make([]argus.CreateScrapeConfigPayloadMetricsRelabelConfigsInner, 0)
		for _, config := range *data.MetricsRelabelConfigs {
			configs = append(configs, mapMetricsRelabelConfig(config))
		}
		metricsRelabelConfigs = &configs
	}

	var params *map[string]interface{}
	var err error

	if data.Params != nil {
		params, err = convertMapAnyToInterface(*data.Params)
		if err != nil {
			return nil, fmt.Errorf("convert params: %w", err)
		}
	}

	return &argus.UpdateScrapeConfigPayload{
		BasicAuth:             basicAuth,
		BearerToken:           data.BearerToken,
		HonorLabels:           data.HonorLabels,
		HonorTimeStamps:       data.HonorTimeStamps,
		MetricsPath:           data.MetricsPath,
		MetricsRelabelConfigs: metricsRelabelConfigs,
		Params:                params,
		SampleLimit:           convertIntToFloat64(data.SampleLimit),
		Scheme:                data.Scheme,
		ScrapeInterval:        data.ScrapeInterval,
		ScrapeTimeout:         data.ScrapeTimeout,
		StaticConfigs:         staticConfigs,
		TlsConfig:             tlsConfig,
	}, nil
}

func convertIntToFloat64(i *int64) *float64 {
	if i == nil {
		return nil
	}
	f := float64(*i)
	return &f
}

func mapMetricsRelabelConfig(metricsRelabelConfig argus.MetricsRelabelConfig) argus.CreateScrapeConfigPayloadMetricsRelabelConfigsInner {
	return argus.CreateScrapeConfigPayloadMetricsRelabelConfigsInner{
		Action:       metricsRelabelConfig.Action,
		Modulus:      convertIntToFloat64(metricsRelabelConfig.Modulus),
		Regex:        metricsRelabelConfig.Regex,
		Replacement:  metricsRelabelConfig.Replacement,
		Separator:    metricsRelabelConfig.Separator,
		SourceLabels: metricsRelabelConfig.SourceLabels,
		TargetLabel:  metricsRelabelConfig.TargetLabel,
	}
}

func mapStaticConfig(staticConfig argus.StaticConfigs) (argus.UpdateScrapeConfigPayloadStaticConfigsInner, error) {
	var labels *map[string]interface{}
	var err error
	if staticConfig.Labels != nil {
		labels, err = convertMapAnyToInterface(*staticConfig.Labels)

		if err != nil {
			return argus.UpdateScrapeConfigPayloadStaticConfigsInner{}, fmt.Errorf("convert labels: %w", err)
		}
	}

	return argus.UpdateScrapeConfigPayloadStaticConfigsInner{
		Labels:  labels,
		Targets: staticConfig.Targets,
	}, nil
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

func convertMapAnyToInterface(m interface{}) (*map[string]interface{}, error) {
	newMap := make(map[string]interface{})

	switch m.(type) {
	case map[string]string:
		convertedMap := m.(map[string]string)
		for k, v := range convertedMap {
			newMap[k] = v
		}
	case map[string][]string:
		convertedMap := m.(map[string][]string)
		for k, v := range convertedMap {
			newMap[k] = v
		}
	default:
		return nil, fmt.Errorf("unsupported map type")
	}

	return &newMap, nil
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

func GetDefaultCreateScrapeConfigPayload() *argus.CreateScrapeConfigPayload {
	staticConfigs := []argus.CreateScrapeConfigPayloadStaticConfigsInner{
		{
			Targets: utils.Ptr([]string{
				"url-target",
			}),
		},
	}
	return &argus.CreateScrapeConfigPayload{
		JobName:        utils.Ptr("default-name"),
		MetricsPath:    utils.Ptr("/metrics"),
		Scheme:         utils.Ptr("https"),
		ScrapeInterval: utils.Ptr("5m"),
		ScrapeTimeout:  utils.Ptr("2m"),
		StaticConfigs:  &staticConfigs,
	}
}
