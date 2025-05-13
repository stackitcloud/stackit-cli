package delete

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/objectstorage"
	"github.com/stackitcloud/stackit-sdk-go/services/objectstorage/wait"
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
		Use:   fmt.Sprintf("delete %s", bucketNameArg),
		Short: "Deletes an Object Storage bucket",
		Long:  "Deletes an Object Storage bucket.",
		Args:  args.SingleArg(bucketNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Delete an Object Storage bucket with name "my-bucket"`,
				"$ stackit object-storage bucket delete my-bucket"),
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

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete bucket %q? (This cannot be undone)", model.BucketName)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete Object Storage bucket: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Deleting bucket")
				_, err = wait.DeleteBucketWaitHandler(ctx, apiClient, model.ProjectId, model.Region, model.BucketName).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for Object Storage bucket deletion: %w", err)
				}
				s.Stop()
			}

			operationState := "Deleted"
			if model.Async {
				operationState = "Triggered deletion of"
			}
			params.Printer.Info("%s bucket %q\n", operationState, model.BucketName)
			return nil
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

	if p.IsVerbosityDebug() {
		modelStr, err := print.BuildDebugStrFromInputModel(model)
		if err != nil {
			p.Debug(print.ErrorLevel, "convert model to string for debugging: %v", err)
		} else {
			p.Debug(print.DebugLevel, "parsed input values: %s", modelStr)
		}
	}

	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *objectstorage.APIClient) objectstorage.ApiDeleteBucketRequest {
	req := apiClient.DeleteBucket(ctx, model.ProjectId, model.Region, model.BucketName)
	return req
}
