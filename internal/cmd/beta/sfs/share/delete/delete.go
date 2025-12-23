package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
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
	shareIdArg = "SHARE_ID"

	resourcePoolIdFlag = "resource-pool-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ResourcePoolId string
	ShareId        string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", shareIdArg),
		Short: "Deletes a share",
		Long:  "Deletes a share.",
		Args:  args.SingleArg(shareIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete a share with ID "xxx" from a resource pool with ID "yyy"`,
				"$ stackit beta sfs share delete xxx --resource-pool-id yyy",
			),
		),
		RunE: func(cmd *cobra.Command, inputArgs []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, inputArgs)
			if err != nil {
				return fmt.Errorf("unable to parse input: %w", err)
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			shareLabel, err := sfsUtils.GetShareName(ctx, apiClient, model.ProjectId, model.Region, model.ResourcePoolId, model.ShareId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get share name: %v", err)
				shareLabel = model.ShareId
			} else if shareLabel == "" {
				shareLabel = model.ShareId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete SFS share %q? (This cannot be undone)", shareLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete SFS share: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Deleting share")
				_, err = wait.DeleteShareWaitHandler(ctx, apiClient, model.ProjectId, model.Region, model.ResourcePoolId, model.ShareId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("waiting for share deletion: %w", err)
				}
				s.Stop()
			}

			operation := "Deleted"
			if model.Async {
				operation = "Triggered deletion of"
			}

			params.Printer.Outputf("%s share %q\n", operation, shareLabel)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), resourcePoolIdFlag, "The resource pool the share is assigned to")

	err := flags.MarkFlagsRequired(cmd, resourcePoolIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	shareId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ShareId:         shareId,
		ResourcePoolId:  flags.FlagToStringValue(p, cmd, resourcePoolIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *sfs.APIClient) sfs.ApiDeleteShareRequest {
	return apiClient.DeleteShare(ctx, model.ProjectId, model.Region, model.ResourcePoolId, model.ShareId)
}
