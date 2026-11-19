package rb

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/object-storage/bucket"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

const bucketURIArg = "BUCKET_URI"

type inputModel struct {
	*globalflags.GlobalFlagModel
	bucket.ClientFlags
	URI string
}

func validateBucketURI(uri string) error {
	ref, err := bucket.ParseBucketURI(uri)
	if err != nil {
		return err
	}
	if !ref.IsBucket() {
		return fmt.Errorf("expected a bucket URI (s3://BUCKET or https://BUCKET.object.storage.REGION.onstackit.cloud), got %q", uri)
	}
	return nil
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("rb %s", bucketURIArg),
		Short: "Deletes a bucket",
		Long:  "Deletes a bucket in STACKIT Object Storage. The bucket must be empty.",
		Args:  args.SingleArg(bucketURIArg, validateBucketURI),
		Example: examples.Build(
			examples.NewExample(
				"Delete an empty bucket",
				"$ stackit alpha bucket rb s3://my-bucket"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			ref, err := bucket.ParseBucketURI(model.URI)
			if err != nil {
				return err
			}

			apiClient, err := bucket.NewClient(bucket.ClientConfigFromReference(model.ClientFlags, model.Region, ref))
			if err != nil {
				return err
			}

			prompt := fmt.Sprintf("Are you sure you want to delete bucket %q? All objects in the bucket must be deleted first", model.URI)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			err = spinner.Run(params.Printer, "Deleting bucket", func() error {
				return apiClient.DeleteBucket(ctx, ref.Bucket)
			})
			if err != nil {
				return fmt.Errorf("delete bucket: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, model.URI)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	bucket.ConfigureClientFlags(cmd.Flags())
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	model := inputModel{
		GlobalFlagModel: globalflags.Parse(p, cmd),
		ClientFlags:     bucket.ClientFlagsFromCommand(p, cmd),
		URI:             inputArgs[0],
	}
	p.DebugInputModel(model)
	return &model, nil
}

func outputResult(p *print.Printer, outputFormat, uri string) error {
	return p.OutputResult(outputFormat, struct {
		Bucket string `json:"bucket"`
	}{Bucket: uri}, func() error {
		p.Outputf("Deleted bucket %q\n", uri)
		return nil
	})
}
