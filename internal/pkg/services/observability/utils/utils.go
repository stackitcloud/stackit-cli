package utils

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	observability "github.com/stackitcloud/stackit-sdk-go/services/observability/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	service = "observability"
)

var (
	defaultStaticConfigs = []observability.CreateScrapeConfigPayloadStaticConfigsInner{
		{
			Targets: []string{
				"url-target",
			},
		},
	}
	DefaultCreateScrapeConfigPayload = observability.CreateScrapeConfigPayload{
		JobName:        "default-name",
		MetricsPath:    utils.Ptr("/metrics"),
		Scheme:         observability.CREATESCRAPECONFIGPAYLOADSCHEME_HTTPS,
		ScrapeInterval: "5m",
		ScrapeTimeout:  "2m",
		StaticConfigs:  defaultStaticConfigs,
	}
)

func ValidatePlanId(planId string, resp *observability.PlansResponse) error {
	if resp == nil {
		return fmt.Errorf("no Observability plans provided")
	}

	for i := range resp.Plans {
		plan := resp.Plans[i]
		if strings.EqualFold(plan.Id, planId) {
			return nil
		}
	}

	return &errors.ObservabilityInvalidPlanError{
		Service: service,
		Details: fmt.Sprintf("You provided plan ID %q, which is invalid.", planId),
	}
}

func LoadPlanId(planName string, resp *observability.PlansResponse) (string, error) {
	availablePlanNames := ""
	if resp == nil {
		return "", fmt.Errorf("no Observability plans provided")
	}

	for i := range resp.Plans {
		plan := resp.Plans[i]
		if plan.Name == nil {
			continue
		}
		if strings.EqualFold(*plan.Name, planName) {
			return plan.Id, nil
		}
		availablePlanNames = fmt.Sprintf("%s\n- %s", availablePlanNames, *plan.Name)
	}

	details := fmt.Sprintf("You provided plan name %q, which is invalid. Available plan names are: %s", planName, availablePlanNames)
	return "", &errors.ObservabilityInvalidPlanError{
		Service: service,
		Details: details,
	}
}

func MapToUpdateScrapeConfigPayload(resp *observability.GetScrapeConfigResponse) (*observability.UpdateScrapeConfigPayload, error) {
	if resp == nil {
		return nil, fmt.Errorf("no Observability scrape config provided")
	}

	data := resp.Data

	basicAuth := mapBasicAuth(data.BasicAuth)
	staticConfigs := mapStaticConfig(data.StaticConfigs)
	tlsConfig := mapTlsConfig(data.TlsConfig)
	metricsRelabelConfigs := mapMetricsRelabelConfig(data.MetricsRelabelConfigs)

	var params map[string]interface{}
	if data.Params != nil {
		params = mapParams(*data.Params)
	}

	var scheme observability.UpdateScrapeConfigPayloadScheme
	if data.Scheme != nil {
		scheme = observability.UpdateScrapeConfigPayloadScheme(*data.Scheme)
	}

	payload := observability.UpdateScrapeConfigPayload{
		BasicAuth:             basicAuth,
		BearerToken:           data.BearerToken,
		HonorLabels:           data.HonorLabels,
		HonorTimeStamps:       data.HonorTimeStamps,
		MetricsPath:           utils.PtrString(data.MetricsPath),
		MetricsRelabelConfigs: metricsRelabelConfigs,
		Params:                params,
		SampleLimit:           utils.ConvertInt32PToFloat32P(data.SampleLimit),
		Scheme:                scheme,
		ScrapeInterval:        data.ScrapeInterval,
		ScrapeTimeout:         data.ScrapeTimeout,
		StaticConfigs:         staticConfigs,
		TlsConfig:             tlsConfig,
	}

	if reflect.DeepEqual(payload, observability.UpdateScrapeConfigPayload{}) {
		return nil, fmt.Errorf("the provided Observability scrape config payload is empty")
	}

	return &payload, nil
}

func mapMetricsRelabelConfig(metricsRelabelConfigs []observability.MetricsRelabelConfig) []observability.UpdateScrapeConfigPayloadMetricsRelabelConfigsInner {
	if metricsRelabelConfigs == nil {
		return nil
	}
	var mappedConfigs []observability.UpdateScrapeConfigPayloadMetricsRelabelConfigsInner
	for _, config := range metricsRelabelConfigs {
		mappedConfig := observability.UpdateScrapeConfigPayloadMetricsRelabelConfigsInner{
			Action:       (*observability.UpdateScrapeConfigPayloadMetricsRelabelConfigsInnerAction)(config.Action),
			Modulus:      utils.ConvertInt32PToFloat32P(config.Modulus),
			Regex:        config.Regex,
			Replacement:  config.Replacement,
			Separator:    config.Separator,
			SourceLabels: config.SourceLabels,
			TargetLabel:  config.TargetLabel,
		}
		mappedConfigs = append(mappedConfigs, mappedConfig)
	}
	return mappedConfigs
}

func mapStaticConfig(staticConfigs []observability.StaticConfigs) []observability.UpdateScrapeConfigPayloadStaticConfigsInner {
	if staticConfigs == nil {
		return nil
	}
	var mappedConfigs []observability.UpdateScrapeConfigPayloadStaticConfigsInner
	for _, config := range staticConfigs {
		var labels map[string]interface{}
		if config.Labels != nil {
			labels = mapStaticConfigLabels(*config.Labels)
		}
		mappedConfig := observability.UpdateScrapeConfigPayloadStaticConfigsInner{
			Labels:  labels,
			Targets: config.Targets,
		}
		mappedConfigs = append(mappedConfigs, mappedConfig)
	}

	return mappedConfigs
}

func mapBasicAuth(basicAuth *observability.BasicAuth) *observability.UpdateScrapeConfigPayloadBasicAuth {
	if basicAuth == nil {
		return nil
	}

	var password, username *string
	if basicAuth.Password != "" {
		password = &basicAuth.Password
	}
	if basicAuth.Username != "" {
		username = &basicAuth.Username
	}

	return &observability.UpdateScrapeConfigPayloadBasicAuth{
		Password: password,
		Username: username,
	}
}

func mapTlsConfig(tlsConfig *observability.TLSConfig) *observability.UpdateScrapeConfigPayloadTlsConfig {
	if tlsConfig == nil {
		return nil
	}

	return &observability.UpdateScrapeConfigPayloadTlsConfig{
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

func GetInstanceName(ctx context.Context, apiClient observability.DefaultAPI, instanceId, projectId string) (string, error) {
	resp, err := apiClient.GetInstance(ctx, instanceId, projectId).Execute()
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

func GetPartialUpdateGrafanaConfigsPayload(ctx context.Context, apiClient observability.DefaultAPI, instanceId, projectId string, singleSignOn, publicReadAccess *bool) (*observability.UpdateGrafanaConfigsPayload, error) {
	currentConfigs, err := apiClient.GetGrafanaConfigs(ctx, instanceId, projectId).Execute()
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
