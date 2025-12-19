package update

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

	resourcePoolIdFlag   = "resource-pool-id"
	exportPolicyNameFlag = "export-policy-name"
	hardLimitFlag        = "hard-limit"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ShareId          string
	ResourcePoolId   string
	ExportPolicyName *string
	HardLimit        *int64
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", shareIdArg),
		Short: "Updates a share",
		Long:  "Updates a share.",
		Args:  args.SingleArg(shareIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update share with ID "xxx" with new export-policy-name "yyy" in resource-pool "zzz"`,
				"$ stackit beta sfs share update xxx --export-policy-name yyy --resource-pool-id zzz",
			),
			examples.NewExample(
				`Update share with ID "xxx" with new space hard limit "50" in resource-pool "yyy"`,
				"$ stackit beta sfs share update xxx --hard-limit 50 --resource-pool-id yyy",
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

			resourcePoolLabel, err := sfsUtils.GetResourcePoolName(ctx, apiClient, model.ProjectId, model.Region, model.ResourcePoolId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get resource pool name: %v", err)
				resourcePoolLabel = model.ResourcePoolId
			} else if resourcePoolLabel == "" {
				resourcePoolLabel = model.ResourcePoolId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update SFS share %q for resource pool %q?", shareLabel, resourcePoolLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update SFS share: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Updating share")
				_, err = wait.UpdateShareWaitHandler(ctx, apiClient, model.ProjectId, model.Region, model.ResourcePoolId, model.ShareId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("waiting for share update: %w", err)
				}
				s.Stop()
			}

			return outputResult(params.Printer, model.OutputFormat, model.Async, resourcePoolLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), resourcePoolIdFlag, "The resource pool the share is assigned to")
	cmd.Flags().String(exportPolicyNameFlag, "", "The export policy the share is assigned to")
	cmd.Flags().Int64(hardLimitFlag, 0, "The space hard limit for the share")

	err := flags.MarkFlagsRequired(cmd, resourcePoolIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	shareId := inputArgs[0]
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	hardLimit := flags.FlagToInt64Pointer(p, cmd, hardLimitFlag)
	if hardLimit != nil && *hardLimit < 0 {
		return nil, &errors.FlagValidationError{
			Flag:    hardLimitFlag,
			Details: "must be a positive integer",
		}
	}

	model := inputModel{
		GlobalFlagModel:  globalFlags,
		ResourcePoolId:   flags.FlagToStringValue(p, cmd, resourcePoolIdFlag),
		ExportPolicyName: flags.FlagToStringPointer(p, cmd, exportPolicyNameFlag),
		HardLimit:        hardLimit,
		ShareId:          shareId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *sfs.APIClient) sfs.ApiUpdateShareRequest {
	req := apiClient.UpdateShare(ctx, model.ProjectId, model.Region, model.ResourcePoolId, model.ShareId)
	req = req.UpdateSharePayload(sfs.UpdateSharePayload{
		ExportPolicyName:        sfs.NewNullableString(model.ExportPolicyName),
		SpaceHardLimitGigabytes: model.HardLimit,
	})
	return req
}

func outputResult(p *print.Printer, outputFormat string, async bool, resourcePoolLabel string, item *sfs.UpdateShareResponse) error {
	return p.OutputResult(outputFormat, item, func() error {
		if item == nil || item.Share == nil {
			p.Outputln("SFS share response is empty")
			return nil
		}

		operation := "Updated"
		if async {
			operation = "Triggered update of"
		}
		p.Outputf(
			"%s SFS share %q in resource pool %q.\n",
			operation,
			utils.PtrString(item.Share.Name),
			resourcePoolLabel,
		)
		return nil
	})
}
