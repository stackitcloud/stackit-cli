package create

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/object-storage/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/object-storage/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"

	"github.com/spf13/cobra"
	objectstorage "github.com/stackitcloud/stackit-sdk-go/services/objectstorage/v2api"
	"github.com/stackitcloud/stackit-sdk-go/services/objectstorage/v2api/wait"
)

const (
	bucketNameArg         = "BUCKET_NAME"
	objectLockEnabledFlag = "object-lock-enabled"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	BucketName        string
	ObjectLockEnabled bool
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("create %s", bucketNameArg),
		Short: "Creates an Object Storage bucket",
		Long:  "Creates an Object Storage bucket.",
		Args:  args.SingleArg(bucketNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Create an Object Storage bucket with name "my-bucket"`,
				"$ stackit object-storage bucket create my-bucket"),
			examples.NewExample(
				`Create an Object Storage bucket with enabled object-lock`,
				`$ stackit object-storage bucket create my-bucket --object-lock-enabled`),
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

			prompt := fmt.Sprintf("Are you sure you want to create bucket %q? (This cannot be undone)", model.BucketName)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Check if the project is enabled before trying to create
			enabled, err := utils.ProjectEnabled(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region)
			if err != nil {
				return fmt.Errorf("check if Object Storage is enabled: %w", err)
			}
			if !enabled {
				return &errors.ServiceDisabledError{
					Service: "object-storage",
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create Object Storage bucket: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				err := spinner.Run(params.Printer, "Creating bucket", func() error {
					_, err = wait.CreateBucketWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.BucketName).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for Object Storage bucket creation: %w", err)
				}
			}

			return outputResult(params.Printer, model.OutputFormat, model.Async, model.BucketName, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Bool(objectLockEnabledFlag, false, "is the object-lock enabled for the bucket")
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	bucketName := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:   globalFlags,
		BucketName:        bucketName,
		ObjectLockEnabled: flags.FlagToBoolValue(p, cmd, objectLockEnabledFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *objectstorage.APIClient) objectstorage.ApiCreateBucketRequest {
	req := apiClient.DefaultAPI.CreateBucket(ctx, model.ProjectId, model.Region, model.BucketName).ObjectLockEnabled(model.ObjectLockEnabled)
	return req
}

func outputResult(p *print.Printer, outputFormat string, async bool, bucketName string, resp *objectstorage.CreateBucketResponse) error {
	if resp == nil {
		return fmt.Errorf("create bucket response is empty")
	}

	return p.OutputResult(outputFormat, resp, func() error {
		operationState := "Created"
		if async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s bucket %q\n", operationState, bucketName)
		return nil
	})
}
