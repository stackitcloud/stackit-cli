package cp

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/object-storage/bucket"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

const (
	srcArg        = "SRC"
	dstArg        = "DESTINATION"
	recursiveFlag = "recursive"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	bucket.ClientFlags
	Source      string
	Destination string
	Recursive   bool
}

type copyResult struct {
	Objects int `json:"objects"`
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("cp %s %s", srcArg, dstArg),
		Short: "Copies objects between local and Object Storage",
		Long: `Copies a single object or a directory tree between local storage and Object Storage.

One of SRC and DESTINATION must be a bucket URI (s3://BUCKET[/KEY]) or a STACKIT
Object Storage URL (https://BUCKET.object.storage.REGION.onstackit.cloud[/KEY])
and the other a local path.
Without --recursive, a single object is transferred.
With --recursive, all objects under the prefix (download) or all files under the
local directory (upload) are transferred.`,
		Args: cobra.ExactArgs(2),
		Example: examples.Build(
			examples.NewExample(
				"Upload a file to a bucket",
				"$ stackit alpha bucket cp local/file.txt s3://my-bucket/file.txt"),
			examples.NewExample(
				"Download a single object",
				"$ stackit alpha bucket cp s3://my-bucket/file.txt local/file.txt"),
			examples.NewExample(
				"Upload a directory recursively",
				"$ stackit alpha bucket cp ./local/dir s3://my-bucket/prefix/ --recursive"),
			examples.NewExample(
				"Download a prefix recursively",
				"$ stackit alpha bucket cp s3://my-bucket/prefix/ ./local/dir --recursive"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			srcIsRef := bucket.IsObjectRef(model.Source)
			dstIsRef := bucket.IsObjectRef(model.Destination)
			if srcIsRef == dstIsRef {
				return fmt.Errorf("one of %s and %s must be a bucket URI (s3://...) or a STACKIT Object Storage URL (https://...) and the other a local path", srcArg, dstArg)
			}

			refURI := model.Destination
			if srcIsRef {
				refURI = model.Source
			}
			ref, err := bucket.ParseBucketURI(refURI)
			if err != nil {
				return err
			}

			apiClient, err := bucket.NewClient(bucket.ClientConfigFromReference(model.ClientFlags, model.Region, ref))
			if err != nil {
				return err
			}

			var count int
			var direction string
			if srcIsRef {
				count, err = spinner.Run2(params.Printer, "Downloading", func() (int, error) {
					return bucket.Download(ctx, apiClient, ref, model.Destination, model.Recursive)
				})
				if err != nil {
					return err
				}
				direction = fmt.Sprintf("Downloaded %d object(s) to %q", count, model.Destination)
			} else {
				count, err = spinner.Run2(params.Printer, "Uploading", func() (int, error) {
					return bucket.Upload(ctx, apiClient, model.Source, ref, model.Recursive)
				})
				if err != nil {
					return err
				}
				direction = fmt.Sprintf("Uploaded %d object(s) to %q", count, ref.String())
			}

			return outputResult(params.Printer, model.OutputFormat, direction, count)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	bucket.ConfigureClientFlags(cmd.Flags())
	cmd.Flags().Bool(recursiveFlag, false, "Transfer a directory tree or all objects under a prefix")
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	model := inputModel{
		GlobalFlagModel: globalflags.Parse(p, cmd),
		ClientFlags:     bucket.ClientFlagsFromCommand(p, cmd),
		Source:          inputArgs[0],
		Destination:     inputArgs[1],
		Recursive:       flags.FlagToBoolValue(p, cmd, recursiveFlag),
	}
	p.DebugInputModel(model)
	return &model, nil
}

func outputResult(p *print.Printer, outputFormat, direction string, objects int) error {
	return p.OutputResult(outputFormat, copyResult{Objects: objects}, func() error {
		p.Outputln(direction)
		return nil
	})
}
