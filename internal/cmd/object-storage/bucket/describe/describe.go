package describe

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/object-storage/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

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

func NewCmd(params *params.CmdParams) *cobra.Command {
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
				`Get details of an Object Storage bucket with name "my-bucket" in JSON format`,
				"$ stackit object-storage bucket describe my-bucket --output-format json"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}
			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read Object Storage bucket: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp.Bucket)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	bucketName := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		BucketName:      bucketName,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *objectstorage.APIClient) objectstorage.ApiGetBucketRequest {
	req := apiClient.GetBucket(ctx, model.ProjectId, model.Region, model.BucketName)
	return req
}

func outputResult(p *print.Printer, outputFormat string, bucket *objectstorage.Bucket) error {
	if bucket == nil {
		return fmt.Errorf("bucket is empty")
	}

	return p.OutputResult(outputFormat, bucket, func() error {
		table := tables.NewTable()
		table.AddRow("Name", utils.PtrString(bucket.Name))
		table.AddSeparator()
		table.AddRow("Region", utils.PtrString(bucket.Region))
		table.AddSeparator()
		table.AddRow("URL (Path Style)", utils.PtrString(bucket.UrlPathStyle))
		table.AddSeparator()
		table.AddRow("URL (Virtual Hosted Style)", utils.PtrString(bucket.UrlVirtualHostedStyle))
		table.AddSeparator()
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
