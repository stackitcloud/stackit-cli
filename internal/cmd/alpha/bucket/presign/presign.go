package presign

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/object-storage/bucket"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

const (
	objectURIArg = "BUCKET_URI"
	expiresFlag  = "expires"

	defaultExpires = time.Hour
	maxExpires     = 7 * 24 * time.Hour
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	bucket.ClientFlags
	URI     string
	Expires time.Duration
}

type presignResult struct {
	URL string `json:"url"`
}

func validateObjectURI(uri string) error {
	ref, err := bucket.ParseBucketURI(uri)
	if err != nil {
		return err
	}
	if ref.IsBucket() {
		return fmt.Errorf("expected an object URI (s3://BUCKET/KEY or https://BUCKET.object.storage.REGION.onstackit.cloud/KEY), got %q", uri)
	}
	return nil
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("presign %s", objectURIArg),
		Short: "Generates a presigned URL for an object",
		Long: `Generates a presigned URL that grants time-limited read access to an object
without exposing the Object Storage credentials.`,
		Args: args.SingleArg(objectURIArg, validateObjectURI),
		Example: examples.Build(
			examples.NewExample(
				"Generate a presigned URL valid for 1 hour (default)",
				"$ stackit alpha bucket presign s3://my-bucket/file.txt"),
			examples.NewExample(
				"Generate a presigned URL valid for 15 minutes",
				"$ stackit alpha bucket presign s3://my-bucket/file.txt --expires 15m"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}
			if model.Expires <= 0 || model.Expires > maxExpires {
				return fmt.Errorf("--%s must be between 0 and %s", expiresFlag, maxExpires)
			}

			ref, err := bucket.ParseBucketURI(model.URI)
			if err != nil {
				return err
			}

			apiClient, err := bucket.NewClient(bucket.ClientConfigFromReference(model.ClientFlags, model.Region, ref))
			if err != nil {
				return err
			}

			presignedURL, err := bucket.PresignGetObject(ctx, apiClient, ref, model.Expires)
			if err != nil {
				return err
			}

			return outputResult(params.Printer, model.OutputFormat, presignedURL)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	bucket.ConfigureClientFlags(cmd.Flags())
	cmd.Flags().Duration(expiresFlag, defaultExpires, "Expiration of the presigned URL (e.g. 15m, 1h, 72h). Maximum: 7 days")
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	model := inputModel{
		GlobalFlagModel: globalflags.Parse(p, cmd),
		ClientFlags:     bucket.ClientFlagsFromCommand(p, cmd),
		URI:             inputArgs[0],
	}
	expires, err := cmd.Flags().GetDuration(expiresFlag)
	if err != nil {
		return nil, fmt.Errorf("parse --%s flag: %w", expiresFlag, err)
	}
	model.Expires = expires
	p.DebugInputModel(model)
	return &model, nil
}

func outputResult(p *print.Printer, outputFormat, url string) error {
	return p.OutputResult(outputFormat, presignResult{URL: url}, func() error {
		p.Outputf("%s\n", url)
		return nil
	})
}
