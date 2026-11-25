package mb

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

const (
	bucketURIArg  = "BUCKET_URI"
	defaultRegion = "eu01"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	bucket.ClientFlags
	URI string
}

type bucketResult struct {
	Bucket string `json:"bucket"`
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
		Use:   fmt.Sprintf("mb %s", bucketURIArg),
		Short: "Creates a bucket",
		Long:  "Creates a bucket in STACKIT Object Storage.",
		Args:  args.SingleArg(bucketURIArg, validateBucketURI),
		Example: examples.Build(
			examples.NewExample(
				"Create a bucket",
				"$ stackit alpha bucket mb s3://my-bucket"),
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
			region := model.Region
			if ref.Region != "" {
				region = ref.Region
			}

			apiClient, err := bucket.NewClient(bucket.ClientConfigFromReference(model.ClientFlags, model.Region, ref))
			if err != nil {
				return err
			}

			prompt := fmt.Sprintf("Are you sure you want to create bucket %q?", model.URI)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			err = spinner.Run(params.Printer, "Creating bucket", func() error {
				return apiClient.CreateBucket(ctx, ref.Bucket, locationConstraint(region))
			})
			if err != nil {
				return fmt.Errorf("create bucket: %w", err)
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

// locationConstraint returns the LocationConstraint for the region, which
// must be omitted for the default region (the constraint value is a region
// code, and the default region requires no constraint body)
func locationConstraint(region string) string {
	if region == "" || region == defaultRegion {
		return ""
	}
	return region
}

func outputResult(p *print.Printer, outputFormat, uri string) error {
	return p.OutputResult(outputFormat, bucketResult{Bucket: uri}, func() error {
		p.Outputf("Created bucket %q\n", uri)
		return nil
	})
}
