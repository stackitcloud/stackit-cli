package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sfs/client"
	sfsUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/sfs/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/sfs"
	"github.com/stackitcloud/stackit-sdk-go/services/sfs/wait"
)

const (
	resourcePoolIdArg = "RESOURCE_POOL_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ResourcePoolId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Deletes a SFS resource pool",
		Long:  "Deletes a SFS resource pool.",
		Args:  args.SingleArg(resourcePoolIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete the SFS resource pool with ID "xxx"`,
				"$ stackit beta sfs resource-pool delete xxx"),
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

			resourcePoolName, err := sfsUtils.GetResourcePoolName(ctx, apiClient, model.ProjectId, model.Region, model.ResourcePoolId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get resource pool name: %v", err)
				resourcePoolName = model.ResourcePoolId
			}

			prompt := fmt.Sprintf("Are you sure you want to delete resource pool %q? (This cannot be undone)", resourcePoolName)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			resp, err := buildRequest(ctx, model, apiClient).Execute()
			if err != nil {
				return fmt.Errorf("delete SFS resource pool: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Delete resource pool")
				_, err = wait.DeleteResourcePoolWaitHandler(ctx, apiClient, model.ProjectId, model.Region, model.ResourcePoolId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for resource pool deletion: %w", err)
				}
				s.Stop()
			}

			return outputResult(params.Printer, model.OutputFormat, resourcePoolName, resp)
		},
	}
	return cmd
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *sfs.APIClient) sfs.ApiDeleteResourcePoolRequest {
	req := apiClient.DeleteResourcePool(ctx, model.ProjectId, model.Region, model.ResourcePoolId)
	return req
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	resourcePoolId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ResourcePoolId:  resourcePoolId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func outputResult(p *print.Printer, outputFormat, resourcePoolName string, response map[string]interface{}) error {
	return p.OutputResult(outputFormat, response, func() error {
		p.Outputf("Deleted resource pool %q\n", resourcePoolName)
		return nil
	})
}
