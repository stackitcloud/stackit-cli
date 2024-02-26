package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/object-storage/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/objectstorage"
)

const (
	bucketNameArg = "BUCKET_NAME"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	BucketName string
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", bucketNameArg),
		Short: "Shows details of an Object Storage bucket",
		Long:  "Shows details of an Object Storage bucket.",
		Args:  args.SingleArg(bucketNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Get details of an Object Storage bucket with name "my-bucket"`,
				"$ stackit object-storage bucket describe my-bucket"),
			examples.NewExample(
				`Get details of an Object Storage bucket with name "my-bucket" in a table format`,
				"$ stackit object-storage bucket describe my-bucket --output-format pretty"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(cmd, args)
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
				return fmt.Errorf("read Object Storage bucket: %w", err)
			}

			return outputResult(cmd, model.OutputFormat, resp.Bucket)
		},
	}
	return cmd
}

func parseInput(cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	bucketName := inputArgs[0]

	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		BucketName:      bucketName,
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *objectstorage.APIClient) objectstorage.ApiGetBucketRequest {
	req := apiClient.GetBucket(ctx, model.ProjectId, model.BucketName)
	return req
}

func outputResult(cmd *cobra.Command, outputFormat string, bucket *objectstorage.Bucket) error {
	switch outputFormat {
	case globalflags.PrettyOutputFormat:
		table := tables.NewTable()
		table.AddRow("Name", *bucket.Name)
		table.AddSeparator()
		table.AddRow("Region", *bucket.Region)
		table.AddSeparator()
		table.AddRow("URL (Path Style)", *bucket.UrlPathStyle)
		table.AddSeparator()
		table.AddRow("URL (Virtual Hosted Style)", *bucket.UrlVirtualHostedStyle)
		table.AddSeparator()
		err := table.Display(cmd)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	default:
		details, err := json.MarshalIndent(bucket, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal Object Storage bucket: %w", err)
		}
		cmd.Println(string(details))

		return nil
	}
}
