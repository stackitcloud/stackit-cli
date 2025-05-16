package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/observability"
)

const (
	service = "observability"
)

type ObservabilityClient interface {
	GetInstanceExecute(ctx context.Context, instanceId, projectId string) (*observability.GetInstanceResponse, error)
	GetGrafanaConfigsExecute(ctx context.Context, instanceId, projectId string) (*observability.GrafanaConfigs, error)
	UpdateGrafanaConfigs(ctx context.Context, instanceId string, projectId string) observability.ApiUpdateGrafanaConfigsRequest
}

var (
	defaultStaticConfigs = []observability.CreateScrapeConfigPayloadStaticConfigsInner{
		{
			Targets: utils.Ptr([]string{
				"url-target",
			}),
		},
	}
	DefaultCreateScrapeConfigPayload = observability.CreateScrapeConfigPayload{
		JobName:        utils.Ptr("default-name"),
		MetricsPath:    utils.Ptr("/metrics"),
		Scheme:         observability.CREATESCRAPECONFIGPAYLOADSCHEME_HTTPS.Ptr(),
		ScrapeInterval: utils.Ptr("5m"),
		ScrapeTimeout:  utils.Ptr("2m"),
		StaticConfigs:  utils.Ptr(defaultStaticConfigs),
	}
)

func ValidatePlanId(planId string, resp *observability.PlansResponse) error {
	if resp == nil {
		return fmt.Errorf("no Observability plans provided")
	}

	for i := range *resp.Plans {
		plan := (*resp.Plans)[i]
		if plan.Id != nil && strings.EqualFold(*plan.Id, planId) {
			return nil
		}
	}

	return &errors.ObservabilityInvalidPlanError{
		Service: service,
		Details: fmt.Sprintf("You provided plan ID %q, which is invalid.", planId),
	}
}

func LoadPlanId(planName string, resp *observability.PlansResponse) (*string, error) {
	availablePlanNames := ""
	if resp == nil {
		return nil, fmt.Errorf("no Observability plans provided")
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
	return nil, &errors.ObservabilityInvalidPlanError{
		Service: service,
		Details: details,
	}
}

func MapToUpdateScrapeConfigPayload(resp *observability.GetScrapeConfigResponse) (*observability.UpdateScrapeConfigPayload, error) {
	if resp == nil || resp.Data == nil {
		return nil, fmt.Errorf("no Observability scrape config provided")
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

	payload := observability.UpdateScrapeConfigPayload{
		BasicAuth:             basicAuth,
		BearerToken:           data.BearerToken,
		HonorLabels:           data.HonorLabels,
		HonorTimeStamps:       data.HonorTimeStamps,
		MetricsPath:           data.MetricsPath,
		MetricsRelabelConfigs: metricsRelabelConfigs,
		Params:                params,
		SampleLimit:           utils.ConvertInt64PToFloat64P(data.SampleLimit),
		Scheme:                observability.UpdateScrapeConfigPayloadGetSchemeAttributeType(data.Scheme),
		ScrapeInterval:        data.ScrapeInterval,
		ScrapeTimeout:         data.ScrapeTimeout,
		StaticConfigs:         staticConfigs,
		TlsConfig:             tlsConfig,
	}

	if payload == (observability.UpdateScrapeConfigPayload{}) {
		return nil, fmt.Errorf("the provided Observability scrape config payload is empty")
	}

	return &payload, nil
}

func mapMetricsRelabelConfig(metricsRelabelConfigs *[]observability.MetricsRelabelConfig) *[]observability.CreateScrapeConfigPayloadMetricsRelabelConfigsInner {
	if metricsRelabelConfigs == nil {
		return nil
	}
	var mappedConfigs []observability.CreateScrapeConfigPayloadMetricsRelabelConfigsInner
	for _, config := range *metricsRelabelConfigs {
		mappedConfig := observability.CreateScrapeConfigPayloadMetricsRelabelConfigsInner{
			Action:       observability.CreateScrapeConfigPayloadMetricsRelabelConfigsInnerGetActionAttributeType(config.Action),
			Modulus:      utils.ConvertInt64PToFloat64P(config.Modulus),
			Regex:        config.Regex,
			Replacement:  config.Replacement,
			Separator:    config.Separator,
			SourceLabels: config.SourceLabels,
			TargetLabel:  config.TargetLabel,
		}
		mappedConfigs = append(mappedConfigs, mappedConfig)
	}
	return &mappedConfigs
}

func mapStaticConfig(staticConfigs *[]observability.StaticConfigs) *[]observability.UpdateScrapeConfigPayloadStaticConfigsInner {
	if staticConfigs == nil {
		return nil
	}
	var mappedConfigs []observability.UpdateScrapeConfigPayloadStaticConfigsInner
	for _, config := range *staticConfigs {
		var labels *map[string]interface{}
		if config.Labels != nil {
			labels = utils.Ptr(mapStaticConfigLabels(*config.Labels))
		}
		mappedConfig := observability.UpdateScrapeConfigPayloadStaticConfigsInner{
			Labels:  labels,
			Targets: config.Targets,
		}
		mappedConfigs = append(mappedConfigs, mappedConfig)
	}

	return &mappedConfigs
}

func mapBasicAuth(basicAuth *observability.BasicAuth) *observability.CreateScrapeConfigPayloadBasicAuth {
	if basicAuth == nil {
		return nil
	}

	return &observability.CreateScrapeConfigPayloadBasicAuth{
		Password: basicAuth.Password,
		Username: basicAuth.Username,
	}
}

func mapTlsConfig(tlsConfig *observability.TLSConfig) *observability.CreateScrapeConfigPayloadHttpSdConfigsInnerOauth2TlsConfig {
	if tlsConfig == nil {
		return nil
	}

	return &observability.CreateScrapeConfigPayloadHttpSdConfigsInnerOauth2TlsConfig{
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

func GetInstanceName(ctx context.Context, apiClient ObservabilityClient, instanceId, projectId string) (string, error) {
	resp, err := apiClient.GetInstanceExecute(ctx, instanceId, projectId)
	if err != nil {
		return "", fmt.Errorf("get Observability instance: %w", err)
	}
	return *resp.Name, nil
}

func ToPayloadGenericOAuth(respOAuth *observability.GrafanaOauth) *observability.UpdateGrafanaConfigsPayloadGenericOauth {
	if respOAuth == nil {
		return nil
	}
	return &observability.UpdateGrafanaConfigsPayloadGenericOauth{
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

func GetPartialUpdateGrafanaConfigsPayload(ctx context.Context, apiClient ObservabilityClient, instanceId, projectId string, singleSignOn, publicReadAccess *bool) (*observability.UpdateGrafanaConfigsPayload, error) {
	currentConfigs, err := apiClient.GetGrafanaConfigsExecute(ctx, instanceId, projectId)
	if err != nil {
		return nil, fmt.Errorf("get current Grafana configs: %w", err)
	}
	if currentConfigs == nil {
		return nil, fmt.Errorf("no Grafana configs found for instance %q", instanceId)
	}

	payload := &observability.UpdateGrafanaConfigsPayload{
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
