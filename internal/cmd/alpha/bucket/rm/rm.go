package rm

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/object-storage/bucket"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

const (
	objectURIArg  = "BUCKET_URI"
	recursiveFlag = "recursive"
	forceFlag     = "force"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	bucket.ClientFlags
	URI       string
	Recursive bool
	Force     bool
}

func validateObjectURI(uri string) error {
	ref, err := bucket.ParseBucketURI(uri)
	if err != nil {
		return err
	}
	if ref.IsBucket() {
		return fmt.Errorf("expected an object or prefix URI (s3://BUCKET/KEY or a STACKIT Object Storage https URL): to delete an empty bucket, use 'stackit alpha bucket rb'")
	}
	return nil
}

type deleteResult struct {
	Objects int `json:"objects"`
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("rm %s", objectURIArg),
		Short: "Deletes objects",
		Long: `Deletes a single object, or all objects under a prefix with --recursive.

To delete an empty bucket, use 'stackit alpha bucket rb'.`,
		Args: args.SingleArg(objectURIArg, validateObjectURI),
		Example: examples.Build(
			examples.NewExample(
				"Delete a single object",
				"$ stackit alpha bucket rm s3://my-bucket/file.txt"),
			examples.NewExample(
				"Delete all objects under a prefix without prompting",
				"$ stackit alpha bucket rm s3://my-bucket/path/to/ --recursive --force"),
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

			refs, err := bucket.ResolveObjectRefs(ctx, apiClient, ref, model.Recursive)
			if err != nil {
				return err
			}

			if !model.Force {
				prompt := fmt.Sprintf("Are you sure you want to delete %d object(s) from %q?", len(refs), ref.String())
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			deleted, err := spinner.Run2(params.Printer, "Deleting objects", func() (int, error) {
				return bucket.DeleteObjects(ctx, apiClient, refs)
			})
			if err != nil {
				return err
			}

			return outputResult(params.Printer, model.OutputFormat, ref.String(), deleted)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	bucket.ConfigureClientFlags(cmd.Flags())
	cmd.Flags().Bool(recursiveFlag, false, "Delete all objects under the prefix")
	cmd.Flags().Bool(forceFlag, false, "Do not prompt for confirmation")
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	model := inputModel{
		GlobalFlagModel: globalflags.Parse(p, cmd),
		ClientFlags:     bucket.ClientFlagsFromCommand(p, cmd),
		URI:             inputArgs[0],
		Recursive:       flags.FlagToBoolValue(p, cmd, recursiveFlag),
		Force:           flags.FlagToBoolValue(p, cmd, forceFlag),
	}
	p.DebugInputModel(model)
	return &model, nil
}

func outputResult(p *print.Printer, outputFormat, uri string, objects int) error {
	return p.OutputResult(outputFormat, deleteResult{Objects: objects}, func() error {
		p.Outputf("Deleted %d object(s) from %q\n", objects, uri)
		return nil
	})
}
