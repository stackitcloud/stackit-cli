package update

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/cdn/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	sdkUtils "github.com/stackitcloud/stackit-sdk-go/core/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/cdn"
)

const (
	argDistributionID                = "DISTRIBUTION_ID"
	flagRegions                      = "regions"
	flagHTTP                         = "http"
	flagHTTPOriginURL                = "http-origin-url"
	flagHTTPGeofencing               = "http-geofencing"
	flagHTTPOriginRequestHeaders     = "http-origin-request-headers"
	flagBucket                       = "bucket"
	flagBucketURL                    = "bucket-url"
	flagBucketCredentialsAccessKeyID = "bucket-credentials-access-key-id"
	flagBucketRegion                 = "bucket-region"
	flagBlockedCountries             = "blocked-countries"
	flagBlockedIPs                   = "blocked-ips"
	flagDefaultCacheDuration         = "default-cache-duration"
	flagLoki                         = "loki"
	flagLokiUsername                 = "loki-username"
	flagLokiPushURL                  = "loki-push-url"
	flagMonthlyLimitBytes            = "monthly-limit-bytes"
	flagOptimizer                    = "optimizer"
)

type bucketInputModel struct {
	URL         string
	AccessKeyID string
	Password    string
	Region      string
}

type httpInputModel struct {
	Geofencing           *map[string][]string
	OriginRequestHeaders *map[string]string
	OriginURL            string
}

type lokiInputModel struct {
	Password string
	Username string
	PushURL  string
}

type inputModel struct {
	*globalflags.GlobalFlagModel
	DistributionID       string
	Regions              []cdn.Region
	Bucket               *bucketInputModel
	HTTP                 *httpInputModel
	BlockedCountries     []string
	BlockedIPs           []string
	DefaultCacheDuration string
	MonthlyLimitBytes    *int64
	Loki                 *lokiInputModel
	Optimizer            *bool
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update",
		Short:   "Update a CDN distribution",
		Long:    "Update a CDN distribution by its ID, allowing replacement of its regions.",
		Args:    args.SingleArg(argDistributionID, utils.ValidateUUID),
		Example: examples.Build(
		// TODO
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}
			if model.Bucket != nil {
				pw, err := params.Printer.PromptForPassword("enter your secret access key for the object storage bucket: ")
				if err != nil {
					return fmt.Errorf("reading secret access key: %w", err)
				}
				model.Bucket.Password = pw
			}
			if model.Loki != nil {
				pw, err := params.Printer.PromptForPassword("enter your password for the loki log sink: ")
				if err != nil {
					return fmt.Errorf("reading loki password: %w", err)
				}
				model.Loki.Password = pw
			}

			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}
			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update a CDN distribution for project %q?", projectLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			req := buildRequest(ctx, apiClient, model)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update CDN distribution: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.EnumSliceFlag(false, []string{}, sdkUtils.EnumSliceToStringSlice(cdn.AllowedRegionEnumValues)...), flagRegions, fmt.Sprintf("Regions in which content should be cached, multiple of: %q", cdn.AllowedRegionEnumValues))
	cmd.Flags().Bool(flagHTTP, false, "Use HTTP backend")
	cmd.Flags().String(flagHTTPOriginURL, "", "Origin URL for HTTP backend")
	cmd.Flags().StringSlice(flagHTTPOriginRequestHeaders, []string{}, "Origin request headers for HTTP backend in the format 'HeaderName: HeaderValue', repeatable. WARNING: do not store sensitive values in the headers!")
	cmd.Flags().StringArray(flagHTTPGeofencing, []string{}, "Geofencing rules for HTTP backend in the format 'https://example.com US,DE'. URL and countries have to be quoted. Repeatable.")
	cmd.Flags().Bool(flagBucket, false, "Use Object Storage backend")
	cmd.Flags().String(flagBucketURL, "", "Bucket URL for Object Storage backend")
	cmd.Flags().String(flagBucketCredentialsAccessKeyID, "", "Access Key ID for Object Storage backend")
	cmd.Flags().String(flagBucketRegion, "", "Region for Object Storage backend")
	cmd.Flags().StringSlice(flagBlockedCountries, []string{}, "Comma-separated list of ISO 3166-1 alpha-2 country codes to block (e.g., 'US,DE,FR')")
	cmd.Flags().StringSlice(flagBlockedIPs, []string{}, "Comma-separated list of IPv4 addresses to block (e.g., '10.0.0.8,127.0.0.1')")
	cmd.Flags().String(flagDefaultCacheDuration, "", "ISO8601 duration string for default cache duration (e.g., 'PT1H30M' for 1 hour and 30 minutes)")
	cmd.Flags().Bool(flagLoki, false, "Enable Loki log sink for the CDN distribution")
	cmd.Flags().String(flagLokiUsername, "", "Username for log sink")
	cmd.Flags().String(flagLokiPushURL, "", "Push URL for log sink")
	cmd.Flags().Int64(flagMonthlyLimitBytes, 0, "Monthly limit in bytes for the CDN distribution")
	cmd.Flags().Bool(flagOptimizer, false, "Enable optimizer for the CDN distribution (paid feature).")
	cmd.MarkFlagsMutuallyExclusive(flagHTTP, flagBucket)
}

func parseInput(p *print.Printer, cmd *cobra.Command, args []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}
	distributionID := args[0]

	regionStrings := flags.FlagToStringSliceValue(p, cmd, flagRegions)
	regions := make([]cdn.Region, 0, len(regionStrings))
	for _, regionStr := range regionStrings {
		regions = append(regions, cdn.Region(regionStr))
	}

	var http *httpInputModel
	if flags.FlagToBoolValue(p, cmd, flagHTTP) {
		originURL := flags.FlagToStringValue(p, cmd, flagHTTPOriginURL)

		var geofencing *map[string][]string
		geofencingInput := flags.FlagToStringArrayValue(p, cmd, flagHTTPGeofencing)
		if geofencingInput != nil {
			geofencing = parseGeofencing(p, geofencingInput)
		}

		var originRequestHeaders *map[string]string
		originRequestHeadersInput := flags.FlagToStringSliceValue(p, cmd, flagHTTPOriginRequestHeaders)
		if originRequestHeadersInput != nil {
			originRequestHeaders = parseOriginRequestHeaders(p, originRequestHeadersInput)
		}

		http = &httpInputModel{
			OriginURL:            originURL,
			Geofencing:           geofencing,
			OriginRequestHeaders: originRequestHeaders,
		}
	}

	var bucket *bucketInputModel
	if flags.FlagToBoolValue(p, cmd, flagBucket) {
		bucketURL := flags.FlagToStringValue(p, cmd, flagBucketURL)
		accessKeyID := flags.FlagToStringValue(p, cmd, flagBucketCredentialsAccessKeyID)
		region := flags.FlagToStringValue(p, cmd, flagBucketRegion)

		bucket = &bucketInputModel{
			URL:         bucketURL,
			AccessKeyID: accessKeyID,
			Password:    "",
			Region:      region,
		}
	}

	blockedCountries := flags.FlagToStringSliceValue(p, cmd, flagBlockedCountries)
	blockedIPs := flags.FlagToStringSliceValue(p, cmd, flagBlockedIPs)
	cacheDuration := flags.FlagToStringValue(p, cmd, flagDefaultCacheDuration)
	monthlyLimit := flags.FlagToInt64Pointer(p, cmd, flagMonthlyLimitBytes)

	var loki *lokiInputModel
	if flags.FlagToBoolValue(p, cmd, flagLoki) {
		loki = &lokiInputModel{
			Username: flags.FlagToStringValue(p, cmd, flagLokiUsername),
			PushURL:  flags.FlagToStringValue(p, cmd, flagLokiPushURL),
			Password: "",
		}
	}

	var optimizer *bool
	if cmd.Flags().Changed(flagOptimizer) {
		o := flags.FlagToBoolValue(p, cmd, flagOptimizer)
		optimizer = &o
	}

	model := inputModel{
		GlobalFlagModel:      globalFlags,
		DistributionID:       distributionID,
		Regions:              regions,
		HTTP:                 http,
		Bucket:               bucket,
		BlockedCountries:     blockedCountries,
		BlockedIPs:           blockedIPs,
		DefaultCacheDuration: cacheDuration,
		MonthlyLimitBytes:    monthlyLimit,
		Loki:                 loki,
		Optimizer:            optimizer,
	}

	p.DebugInputModel(model)
	return &model, nil
}

// TODO both parseGeofencing and parseOriginRequestHeaders copied from create.go, move to another package and make public?
func parseGeofencing(p *print.Printer, geofencingInput []string) *map[string][]string {
	geofencing := make(map[string][]string)
	for _, in := range geofencingInput {
		firstSpace := strings.IndexRune(in, ' ')
		if firstSpace == -1 {
			p.Debug(print.ErrorLevel, "invalid geofencing entry (no space found): %q", in)
			continue
		}
		urlPart := in[:firstSpace]
		countriesPart := in[firstSpace+1:]
		geofencing[urlPart] = nil
		countries := strings.Split(countriesPart, ",")
		for _, country := range countries {
			country = strings.TrimSpace(country)
			geofencing[urlPart] = append(geofencing[urlPart], country)
		}
	}
	return &geofencing
}

func parseOriginRequestHeaders(p *print.Printer, originRequestHeadersInput []string) *map[string]string {
	originRequestHeaders := make(map[string]string)
	for _, in := range originRequestHeadersInput {
		parts := strings.Split(in, ":")
		if len(parts) != 2 {
			p.Debug(print.ErrorLevel, "invalid origin request header entry (no colon found): %q", in)
			continue
		}
		originRequestHeaders[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	return &originRequestHeaders
}

func buildRequest(ctx context.Context, apiClient *cdn.APIClient, model *inputModel) cdn.ApiPatchDistributionRequest {
	req := apiClient.PatchDistribution(ctx, model.ProjectId, model.DistributionID)
	payload := cdn.NewPatchDistributionPayload()
	cfg := &cdn.ConfigPatch{}
	payload.Config = cfg
	if len(model.Regions) > 0 {
		cfg.Regions = &model.Regions
	}
	if model.Bucket != nil {
		bucket := &cdn.BucketBackendPatch{
			Type: utils.Ptr("bucket"),
		}
		cfg.Backend = &cdn.ConfigPatchBackend{
			BucketBackendPatch: bucket,
		}
		if model.Bucket.URL != "" {
			bucket.BucketUrl = utils.Ptr(model.Bucket.URL)
		}
		if model.Bucket.AccessKeyID != "" {
			bucket.Credentials = cdn.NewBucketCredentials(
				model.Bucket.AccessKeyID,
				model.Bucket.Password,
			)
		}
		if model.Bucket.Region != "" {
			bucket.Region = utils.Ptr(model.Bucket.Region)
		}
	} else if model.HTTP != nil {
		http := &cdn.HttpBackendPatch{
			Type: utils.Ptr("http"),
		}
		cfg.Backend = &cdn.ConfigPatchBackend{
			HttpBackendPatch: http,
		}
		if model.HTTP.OriginRequestHeaders != nil {
			http.OriginRequestHeaders = model.HTTP.OriginRequestHeaders
		}
		if model.HTTP.Geofencing != nil {
			http.Geofencing = model.HTTP.Geofencing
		}
		if model.HTTP.OriginURL != "" {
			http.OriginUrl = utils.Ptr(model.HTTP.OriginURL)
		}
	}
	if len(model.BlockedCountries) > 0 {
		cfg.BlockedCountries = &model.BlockedCountries
	}
	if len(model.BlockedIPs) > 0 {
		cfg.BlockedIps = &model.BlockedIPs
	}
	if model.DefaultCacheDuration != "" {
		cfg.DefaultCacheDuration = cdn.NewNullableString(&model.DefaultCacheDuration)
	}
	if model.MonthlyLimitBytes != nil && *model.MonthlyLimitBytes > 0 {
		cfg.MonthlyLimitBytes = model.MonthlyLimitBytes
	}
	if model.Loki != nil {
		loki := &cdn.LokiLogSinkPatch{}
		cfg.LogSink = cdn.NewNullableConfigPatchLogSink(&cdn.ConfigPatchLogSink{
			LokiLogSinkPatch: loki,
		})
		if model.Loki.PushURL != "" {
			loki.PushUrl = utils.Ptr(model.Loki.PushURL)
		}
		if model.Loki.Username != "" {
			loki.Credentials = cdn.NewLokiLogSinkCredentials(
				model.Loki.Password,
				model.Loki.Username,
			)
		}
	}
	if model.Optimizer != nil {
		cfg.Optimizer = &cdn.OptimizerPatch{
			Enabled: model.Optimizer,
		}
	}
	req = req.PatchDistributionPayload(*payload)
	return req
}

func outputResult(p *print.Printer, outputFormat string, projectLabel string, resp *cdn.PatchDistributionResponse) error {
	if resp == nil {
		return fmt.Errorf("update distribution response is empty")
	}
	return p.OutputResult(outputFormat, resp, func() error {
		p.Outputf("Updated CDN distribution for %q. ID: %s\n", projectLabel, utils.PtrString(resp.Distribution.Id))
		return nil
	})
}
