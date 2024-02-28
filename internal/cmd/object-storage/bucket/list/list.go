package list

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/object-storage/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/objectstorage"
)

const (
	limitFlag = "limit"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit *int64
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all Object Storage buckets",
		Long:  "Lists all Object Storage buckets.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all Object Storage buckets`,
				"$ stackit object-storage bucket list"),
			examples.NewExample(
				`List all Object Storage buckets in JSON format`,
				"$ stackit object-storage bucket list --output-format json"),
			examples.NewExample(
				`List up to 10 Object Storage buckets`,
				"$ stackit object-storage bucket list --limit 10"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(cmd)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("get Object Storage buckets: %w", err)
			}
			if resp.Buckets == nil || len(*resp.Buckets) == 0 {
				projectLabel, err := projectname.GetProjectName(ctx, cmd)
				if err != nil {
					projectLabel = model.ProjectId
				}
				cmd.Printf("No buckets found for project %s\n", projectLabel)
				return nil
			}
			buckets := *resp.Buckets

			// Truncate output
			if model.Limit != nil && len(buckets) > int(*model.Limit) {
				buckets = buckets[:*model.Limit]
			}

			return outputResult(cmd, model.OutputFormat, buckets)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
}

func parseInput(cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	limit := flags.FlagToInt64Pointer(cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		Limit:           flags.FlagToInt64Pointer(cmd, limitFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *objectstorage.APIClient) objectstorage.ApiListBucketsRequest {
	req := apiClient.ListBuckets(ctx, model.ProjectId)
	return req
}

func outputResult(cmd *cobra.Command, outputFormat string, buckets []objectstorage.Bucket) error {
	switch outputFormat {
	case globalflags.JSONOutputFormat:
		details, err := json.MarshalIndent(buckets, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal Object Storage bucket list: %w", err)
		}
		cmd.Println(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("NAME", "REGION", "URL (PATH STYLE)", "URL (VIRTUAL HOSTED STYLE)")
		for i := range buckets {
			bucket := buckets[i]
			table.AddRow(*bucket.Name, *bucket.Region, *bucket.UrlPathStyle, *bucket.UrlVirtualHostedStyle)
		}
		err := table.Display(cmd)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
