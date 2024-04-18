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

var (
	defaultStaticConfigs = []argus.CreateScrapeConfigPayloadStaticConfigsInner{
		{
			Targets: utils.Ptr([]string{
				"url-target",
			}),
		},
	}
	DefaultCreateScrapeConfigPayload = argus.CreateScrapeConfigPayload{
		JobName:        utils.Ptr("default-name"),
		MetricsPath:    utils.Ptr("/metrics"),
		Scheme:         utils.Ptr("https"),
		ScrapeInterval: utils.Ptr("5m"),
		ScrapeTimeout:  utils.Ptr("2m"),
		StaticConfigs:  utils.Ptr(defaultStaticConfigs),
	}
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
	staticConfigs := mapStaticConfig(data.StaticConfigs)
	tlsConfig := mapTlsConfig(data.TlsConfig)
	metricsRelabelConfigs := mapMetricsRelabelConfig(data.MetricsRelabelConfigs)

	var params *map[string]interface{}
	if data.Params != nil {
		params = utils.Ptr(mapParams(*data.Params))
	}

	payload := argus.UpdateScrapeConfigPayload{
		BasicAuth:             basicAuth,
		BearerToken:           data.BearerToken,
		HonorLabels:           data.HonorLabels,
		HonorTimeStamps:       data.HonorTimeStamps,
		MetricsPath:           data.MetricsPath,
		MetricsRelabelConfigs: metricsRelabelConfigs,
		Params:                params,
		SampleLimit:           utils.ConvertInt64PToFloat64P(data.SampleLimit),
		Scheme:                data.Scheme,
		ScrapeInterval:        data.ScrapeInterval,
		ScrapeTimeout:         data.ScrapeTimeout,
		StaticConfigs:         staticConfigs,
		TlsConfig:             tlsConfig,
	}

	if payload == (argus.UpdateScrapeConfigPayload{}) {
		return nil, fmt.Errorf("the provided Argus Scrape Sonfig payload is empty")
	}

	return &payload, nil
}

func mapMetricsRelabelConfig(metricsRelabelConfigs *[]argus.MetricsRelabelConfig) *[]argus.CreateScrapeConfigPayloadMetricsRelabelConfigsInner {
	if metricsRelabelConfigs == nil {
		return nil
	}
	configs := make([]argus.CreateScrapeConfigPayloadMetricsRelabelConfigsInner, 0)
	for _, config := range *metricsRelabelConfigs {
		mappedConfig := argus.CreateScrapeConfigPayloadMetricsRelabelConfigsInner{
			Action:       config.Action,
			Modulus:      utils.ConvertInt64PToFloat64P(config.Modulus),
			Regex:        config.Regex,
			Replacement:  config.Replacement,
			Separator:    config.Separator,
			SourceLabels: config.SourceLabels,
			TargetLabel:  config.TargetLabel,
		}
		configs = append(configs, mappedConfig)
	}
	return &configs
}

func mapStaticConfig(staticConfigs *[]argus.StaticConfigs) *[]argus.UpdateScrapeConfigPayloadStaticConfigsInner {
	if staticConfigs == nil {
		return nil
	}
	configs := make([]argus.UpdateScrapeConfigPayloadStaticConfigsInner, 0)
	for _, config := range *staticConfigs {
		var labels *map[string]interface{}
		if config.Labels != nil {
			labels = utils.Ptr(mapStaticConfigLabels(*config.Labels))
		}
		mappedConfig := argus.UpdateScrapeConfigPayloadStaticConfigsInner{
			Labels:  labels,
			Targets: config.Targets,
		}
		configs = append(configs, mappedConfig)
	}

	return &configs
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

func mapParams(params map[string][]string) map[string]interface{} {
	paramsMap := make(map[string]interface{})
	for k, v := range params {
		paramsMap[k] = v
	}
	return paramsMap
}

func mapStaticConfigLabels(labels map[string]string) map[string]interface{} {
	labelsMap := make(map[string]interface{})
	for k, v := range labels {
		labelsMap[k] = v
	}
	return labelsMap
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
