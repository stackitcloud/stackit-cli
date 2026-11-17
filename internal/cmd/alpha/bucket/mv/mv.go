package mv

import (
	"context"
	"fmt"
	"os"
	"strings"

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

type moveResult struct {
	Objects int `json:"objects"`
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("mv %s %s", srcArg, dstArg),
		Short: "Moves objects between local and Object Storage",
		Long: `Moves a single object or a directory tree, then deletes the source.

One of SRC and DESTINATION must be a bucket URI (s3://BUCKET[/KEY]) or a STACKIT
Object Storage URL (https://BUCKET.object.storage.REGION.onstackit.cloud[/KEY])
and the other a local path.
Bucket-to-bucket moves use a server-side copy.
Without --recursive, a single object is moved.
With --recursive, all objects under the prefix (or all files under the local directory) are moved.`,
		Args: cobra.ExactArgs(2),
		Example: examples.Build(
			examples.NewExample(
				"Move a file to a bucket",
				"$ stackit alpha bucket mv local/file.txt s3://my-bucket/file.txt"),
			examples.NewExample(
				"Move an object to a new key",
				"$ stackit alpha bucket mv s3://my-bucket/old/file.txt s3://my-bucket/new/file.txt"),
			examples.NewExample(
				"Move a prefix to another bucket recursively",
				"$ stackit alpha bucket mv s3://src-bucket/prefix/ s3://dst-bucket/prefix/ --recursive"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			srcIsRef := bucket.IsObjectRef(model.Source)
			dstIsRef := bucket.IsObjectRef(model.Destination)
			if !srcIsRef && !dstIsRef {
				return fmt.Errorf("one of %s and %s must be a bucket URI (s3://...) or a STACKIT Object Storage URL (https://...)", srcArg, dstArg)
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

			prompt := fmt.Sprintf("Are you sure you want to move %q to %q? The source will be deleted", model.Source, model.Destination)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			count, err := spinner.Run2(params.Printer, "Moving", func() (int, error) {
				return runMove(ctx, apiClient, model, srcIsRef)
			})
			if err != nil {
				return err
			}

			direction := fmt.Sprintf("Moved %d object(s) from %q to %q", count, model.Source, model.Destination)
			return outputResult(params.Printer, model.OutputFormat, direction, count)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	bucket.ConfigureClientFlags(cmd.Flags())
	cmd.Flags().Bool(recursiveFlag, false, "Move a directory tree or all objects under a prefix")
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

func runMove(ctx context.Context, apiClient *bucket.Client, model *inputModel, srcIsRef bool) (int, error) {
	if !srcIsRef {
		dstRef, err := bucket.ParseBucketURI(model.Destination)
		if err != nil {
			return 0, err
		}
		count, err := bucket.Upload(ctx, apiClient, model.Source, dstRef, model.Recursive)
		if err != nil {
			return 0, err
		}
		if err := deleteLocal(model.Source); err != nil {
			return count, fmt.Errorf("objects were uploaded to %q, but deleting the local source failed: %w", dstRef.String(), err)
		}
		return count, nil
	}

	srcRef, err := bucket.ParseBucketURI(model.Source)
	if err != nil {
		return 0, err
	}
	if bucket.IsObjectRef(model.Destination) {
		dstRef, err := bucket.ParseBucketURI(model.Destination)
		if err != nil {
			return 0, err
		}
		if isCrossEndpointMove(srcRef, dstRef) {
			return 0, fmt.Errorf("cross-region or cross-endpoint bucket-to-bucket moves are not supported: %q and %q resolve to different endpoints or regions", model.Source, model.Destination)
		}
		return moveBucketToBucket(ctx, apiClient, srcRef, model)
	}
	count, err := bucket.Download(ctx, apiClient, srcRef, model.Destination, model.Recursive)
	if err != nil {
		return 0, err
	}
	if err := deleteSource(ctx, apiClient, srcRef, model.Recursive); err != nil {
		return count, fmt.Errorf("objects were downloaded to %q, but deleting the source failed: %w", model.Destination, err)
	}
	return count, nil
}

// isCrossEndpointMove reports whether src and dst explicitly resolve to
// different Object Storage endpoints or regions. s3:// references carry no endpoint or
// region (they are derived from the --region flag), so only explicit
// mismatches are detected.
func isCrossEndpointMove(src, dst bucket.Reference) bool {
	return (src.Endpoint != "" && dst.Endpoint != "" && src.Endpoint != dst.Endpoint) ||
		(src.Region != "" && dst.Region != "" && src.Region != dst.Region)
}

func moveBucketToBucket(ctx context.Context, apiClient *bucket.Client, srcRef bucket.Reference, model *inputModel) (int, error) {
	dstRef, err := bucket.ParseBucketURI(model.Destination)
	if err != nil {
		return 0, err
	}
	srcs, err := bucket.ResolveObjectRefs(ctx, apiClient, srcRef, model.Recursive)
	if err != nil {
		return 0, err
	}

	moved := 0
	for _, src := range srcs {
		dst := destinationForMove(src, srcRef, dstRef, model.Recursive)
		if err := bucket.CopyObject(ctx, apiClient, src, dst); err != nil {
			return moved, err
		}
		moved++
	}
	_, err = bucket.DeleteObjects(ctx, apiClient, srcs)
	if err != nil {
		return moved, fmt.Errorf("%d object(s) were copied to %q, but deleting the source failed: %w", moved, dstRef.String(), err)
	}
	return moved, nil
}

func destinationForMove(src, srcRef, dstRef bucket.Reference, recursive bool) bucket.Reference {
	if !recursive {
		if dstRef.IsBucket() {
			return bucket.Reference{Bucket: dstRef.Bucket, Key: src.Key}
		}
		return dstRef
	}
	rel := strings.TrimPrefix(src.Key, srcRef.Key)
	rel = strings.TrimPrefix(rel, "/")
	base := strings.TrimSuffix(dstRef.Key, "/")
	if base == "" {
		return bucket.Reference{Bucket: dstRef.Bucket, Key: rel}
	}
	return bucket.Reference{Bucket: dstRef.Bucket, Key: base + "/" + rel}
}

func deleteSource(ctx context.Context, apiClient *bucket.Client, ref bucket.Reference, recursive bool) error {
	refs, err := bucket.ResolveObjectRefs(ctx, apiClient, ref, recursive)
	if err != nil {
		return err
	}
	_, err = bucket.DeleteObjects(ctx, apiClient, refs)
	return err
}

func deleteLocal(path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat local source %q: %w", path, err)
	}
	if stat.IsDir() {
		return os.RemoveAll(path)
	}
	return os.Remove(path)
}

func outputResult(p *print.Printer, outputFormat, direction string, objects int) error {
	return p.OutputResult(outputFormat, moveResult{Objects: objects}, func() error {
		p.Outputln(direction)
		return nil
	})
}
