package ls

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

	timeFormat = "2006-01-02 15:04:05"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	bucket.ClientFlags
	URI       string
	Recursive bool
}

// listResult is the JSON/YAML output of a bucket object listing
type listResult struct {
	Prefixes []string            `json:"prefixes,omitempty"`
	Objects  []bucket.ObjectInfo `json:"objects,omitempty"`
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("ls [%s]", objectURIArg),
		Short: "Lists buckets and objects",
		Long: `Lists all buckets, or the contents of a bucket.

Without an argument, lists all buckets.
With a bucket URI (s3://BUCKET), lists the "directories" and objects at the bucket root.
With an object URI (s3://BUCKET/PREFIX), lists the "directories" and objects under the prefix.
With --recursive, lists all objects under the prefix flat, including all nested keys.`,
		Args: args.SingleOptionalArg(objectURIArg, nil),
		Example: examples.Build(
			examples.NewExample(
				"List all buckets",
				"$ stackit alpha bucket ls"),
			examples.NewExample(
				"List the contents of a bucket",
				"$ stackit alpha bucket ls s3://my-bucket"),
			examples.NewExample(
				"Recursively list all objects under a prefix",
				"$ stackit alpha bucket ls s3://my-bucket/path/to/ --recursive"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			ref := bucket.Reference{}
			if model.URI != "" {
				ref, err = bucket.ParseBucketURI(model.URI)
				if err != nil {
					return err
				}
			}

			apiClient, err := bucket.NewClient(bucket.ClientConfigFromReference(model.ClientFlags, model.Region, ref))
			if err != nil {
				return err
			}

			if model.URI == "" {
				buckets, err := spinner.Run2(params.Printer, "Listing buckets", func() ([]bucket.BucketInfo, error) {
					return bucket.ListBuckets(ctx, apiClient)
				})
				if err != nil {
					return err
				}
				return outputResultBuckets(params.Printer, model.OutputFormat, buckets)
			}

			if model.Recursive {
				objects, err := spinner.Run2(params.Printer, "Listing objects", func() ([]bucket.ObjectInfo, error) {
					return bucket.ListObjects(ctx, apiClient, ref)
				})
				if err != nil {
					return err
				}
				return outputResultObjects(params.Printer, model.OutputFormat, nil, objects)
			}

			result, err := spinner.Run2(params.Printer, "Listing objects", func() (listResult, error) {
				prefixes, objects, err := bucket.ListPrefixes(ctx, apiClient, ref)
				return listResult{Prefixes: prefixes, Objects: objects}, err
			})
			if err != nil {
				return err
			}
			return outputResultObjects(params.Printer, model.OutputFormat, result.Prefixes, result.Objects)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	bucket.ConfigureClientFlags(cmd.Flags())
	cmd.Flags().Bool(recursiveFlag, false, "List all objects under the prefix recursively, including nested keys")
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	uri := ""
	if len(inputArgs) > 0 {
		uri = inputArgs[0]
	}
	model := inputModel{
		GlobalFlagModel: globalflags.Parse(p, cmd),
		ClientFlags:     bucket.ClientFlagsFromCommand(p, cmd),
		URI:             uri,
		Recursive:       flags.FlagToBoolValue(p, cmd, recursiveFlag),
	}
	p.DebugInputModel(model)
	return &model, nil
}

func outputResultBuckets(p *print.Printer, outputFormat string, buckets []bucket.BucketInfo) error {
	return p.OutputResult(outputFormat, buckets, func() error {
		for _, bucket := range buckets {
			creationDate := "n/a"
			if !bucket.CreationDate.IsZero() {
				creationDate = bucket.CreationDate.Format(timeFormat)
			}
			p.Outputf("%s\t%s\n", creationDate, bucket.Name)
		}
		return nil
	})
}

func outputResultObjects(p *print.Printer, outputFormat string, prefixes []string, objects []bucket.ObjectInfo) error {
	result := listResult{
		Prefixes: prefixes,
		Objects:  objects,
	}
	return p.OutputResult(outputFormat, result, func() error {
		for _, prefix := range prefixes {
			p.Outputln(prefix)
		}
		for _, object := range objects {
			modified := "n/a"
			if !object.LastModified.IsZero() {
				modified = object.LastModified.Format(timeFormat)
			}
			p.Outputf("%s\t%d\t%s\n", modified, object.Size, object.Key)
		}
		return nil
	})
}
