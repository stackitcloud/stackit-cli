package create

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/object-storage/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/object-storage/utils"
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
				prompt := fmt.Sprintf("Are you sure you want to create bucket %q? (This cannot be undone)", model.BucketName)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Check if the project is enabled before trying to create
			enabled, err := utils.ProjectEnabled(ctx, apiClient, model.ProjectId, model.Region)
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
				s := spinner.New(params.Printer)
				s.Start("Creating bucket")
				_, err = wait.CreateBucketWaitHandler(ctx, apiClient, model.ProjectId, model.Region, model.BucketName).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for Object Storage bucket creation: %w", err)
				}
				s.Stop()
			}

			return outputResult(params.Printer, model.OutputFormat, model.Async, model.BucketName, resp)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *objectstorage.APIClient) objectstorage.ApiCreateBucketRequest {
	req := apiClient.CreateBucket(ctx, model.ProjectId, model.Region, model.BucketName)
	return req
}

func outputResult(p *print.Printer, outputFormat string, async bool, bucketName string, resp *objectstorage.CreateBucketResponse) error {
	if resp == nil {
		return fmt.Errorf("create bucket response is empty")
	}

	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal Object Storage bucket: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(resp, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal Object Storage bucket: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		operationState := "Created"
		if async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s bucket %q\n", operationState, bucketName)
		return nil
	}
}
