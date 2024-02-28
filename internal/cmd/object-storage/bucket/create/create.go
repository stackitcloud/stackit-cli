package create

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/confirm"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
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

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("create %s", bucketNameArg),
		Short: "Creates an Object Storage bucket",
		Long:  "Creates an Object Storage bucket.",
		Args:  args.SingleArg(bucketNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Create an Object Storage bucket with name "my-bucket"`,
				"$ stackit object-storage bucket create my-bucket"),
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

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create bucket %q? (This cannot be undone)", model.BucketName)
				err = confirm.PromptForConfirmation(cmd, prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("create Object Storage bucket: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(cmd)
				s.Start("Creating bucket")
				_, err = wait.CreateBucketWaitHandler(ctx, apiClient, model.ProjectId, model.BucketName).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for Object Storage bucket creation: %w", err)
				}
				s.Stop()
			}

			operationState := "Created"
			if model.Async {
				operationState = "Triggered creation of"
			}
			cmd.Printf("%s bucket %q\n", operationState, model.BucketName)
			return nil
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *objectstorage.APIClient) objectstorage.ApiCreateBucketRequest {
	req := apiClient.CreateBucket(ctx, model.ProjectId, model.BucketName)
	return req
}
