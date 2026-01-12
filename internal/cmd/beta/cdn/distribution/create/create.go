package create

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/cdn/client"
	cdnUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/cdn/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	sdkUtils "github.com/stackitcloud/stackit-sdk-go/core/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/cdn"
)

const (
	flagRegion                       = "regions"
	flagHTTP                         = "http"
	flagHTTPOriginURL                = "http-origin-url"
	flagHTTPGeofencing               = "http-geofencing"
	flagHTTPOriginRequestHeaders     = "http-origin-request-headers"
	flagBucket                       = "bucket"
	flagBucketURL                    = "bucket-url"
	flagBucketCredentialsAccessKeyID = "bucket-credentials-access-key-id" //nolint:gosec // linter false positive
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

type httpInputModel struct {
	OriginURL            string
	Geofencing           *map[string][]string
	OriginRequestHeaders *map[string]string
}

type bucketInputModel struct {
	URL         string
	AccessKeyID string
	Password    string
	Region      string
}

type lokiInputModel struct {
	Username string
	Password string
	PushURL  string
}

type inputModel struct {
	*globalflags.GlobalFlagModel
	Regions              []cdn.Region
	HTTP                 *httpInputModel
	Bucket               *bucketInputModel
	BlockedCountries     []string
	BlockedIPs           []string
	DefaultCacheDuration string
	MonthlyLimitBytes    *int64
	Loki                 *lokiInputModel
	Optimizer            bool
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a CDN distribution",
		Long:  "Create a CDN distribution for a given originUrl in multiple regions.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a CDN distribution with an HTTP backend`,
				`$ stackit beta cdn create --http  --http-origin-url https://example.com \
--regions AF,EU`,
			),
			examples.NewExample(
				`Create a CDN distribution with an Object Storage backend`,
				`$ stackit beta cdn create --bucket --bucket-url https://bucket.example.com \
--bucket-credentials-access-key-id yyyy --bucket-region EU \
--regions AF,EU`,
			),
		),
		PreRun: func(cmd *cobra.Command, _ []string) {
			// either flagHTTP or flagBucket must be set, depending on which we mark other flags as required
			if flags.FlagToBoolValue(params.Printer, cmd, flagHTTP) {
				err := cmd.MarkFlagRequired(flagHTTPOriginURL)
				cobra.CheckErr(err)
			} else {
				err := flags.MarkFlagsRequired(cmd, flagBucketURL, flagBucketCredentialsAccessKeyID, flagBucketRegion)
				cobra.CheckErr(err)
			}
			// if user uses loki, mark related flags as required
			if flags.FlagToBoolValue(params.Printer, cmd, flagLoki) {
				err := flags.MarkFlagsRequired(cmd, flagLokiUsername, flagLokiPushURL)
				cobra.CheckErr(err)
			}
		},
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
				prompt := fmt.Sprintf("Are you sure you want to create a CDN distribution for project %q?", projectLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			req := buildRequest(ctx, model, apiClient)

			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create CDN distribution: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.EnumSliceFlag(false, []string{}, sdkUtils.EnumSliceToStringSlice(cdn.AllowedRegionEnumValues)...), flagRegion, fmt.Sprintf("Regions in which content should be cached, multiple of: %q", cdn.AllowedRegionEnumValues))
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
	cmd.MarkFlagsOneRequired(flagHTTP, flagBucket)
	err := flags.MarkFlagsRequired(cmd, flagRegion)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	regionStrings := flags.FlagToStringSliceValue(p, cmd, flagRegion)
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
			geofencing = cdnUtils.ParseGeofencing(p, geofencingInput)
		}

		var originRequestHeaders *map[string]string
		originRequestHeadersInput := flags.FlagToStringSliceValue(p, cmd, flagHTTPOriginRequestHeaders)
		if originRequestHeadersInput != nil {
			originRequestHeaders = cdnUtils.ParseOriginRequestHeaders(p, originRequestHeadersInput)
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

	optimizer := flags.FlagToBoolValue(p, cmd, flagOptimizer)

	model := inputModel{
		GlobalFlagModel:      globalFlags,
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

	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *cdn.APIClient) cdn.ApiCreateDistributionRequest {
	req := apiClient.CreateDistribution(ctx, model.ProjectId)
	var backend cdn.CreateDistributionPayloadGetBackendArgType
	if model.HTTP != nil {
		backend = cdn.CreateDistributionPayloadGetBackendArgType{
			HttpBackendCreate: &cdn.HttpBackendCreate{
				Geofencing:           model.HTTP.Geofencing,
				OriginRequestHeaders: model.HTTP.OriginRequestHeaders,
				OriginUrl:            &model.HTTP.OriginURL,
				Type:                 utils.Ptr("http"),
			},
		}
	} else {
		backend = cdn.CreateDistributionPayloadGetBackendArgType{
			BucketBackendCreate: &cdn.BucketBackendCreate{
				BucketUrl: &model.Bucket.URL,
				Credentials: cdn.NewBucketCredentials(
					model.Bucket.AccessKeyID,
					model.Bucket.Password,
				),
				Region: &model.Bucket.Region,
				Type:   utils.Ptr("bucket"),
			},
		}
	}

	payload := cdn.NewCreateDistributionPayload(
		backend,
		model.Regions,
	)
	if len(model.BlockedCountries) > 0 {
		payload.BlockedCountries = &model.BlockedCountries
	}
	if len(model.BlockedIPs) > 0 {
		payload.BlockedIps = &model.BlockedIPs
	}
	if model.DefaultCacheDuration != "" {
		payload.DefaultCacheDuration = utils.Ptr(model.DefaultCacheDuration)
	}
	if model.Loki != nil {
		payload.LogSink = &cdn.CreateDistributionPayloadGetLogSinkArgType{
			LokiLogSinkCreate: &cdn.LokiLogSinkCreate{
				Credentials: &cdn.LokiLogSinkCredentials{
					Password: &model.Loki.Password,
					Username: &model.Loki.Username,
				},
				PushUrl: &model.Loki.PushURL,
				Type:    utils.Ptr("loki"),
			},
		}
	}
	payload.MonthlyLimitBytes = model.MonthlyLimitBytes
	if model.Optimizer {
		payload.Optimizer = &cdn.CreateDistributionPayloadGetOptimizerArgType{
			Enabled: utils.Ptr(true),
		}
	}
	return req.CreateDistributionPayload(*payload)
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, resp *cdn.CreateDistributionResponse) error {
	if resp == nil {
		return fmt.Errorf("create distribution response is nil")
	}
	return p.OutputResult(outputFormat, resp, func() error {
		p.Outputf("Created CDN distribution for %q. ID: %s\n", projectLabel, utils.PtrString(resp.Distribution.Id))
		return nil
	})
}
