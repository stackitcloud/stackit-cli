package create

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
	nameFlag             = "name"
	resourcePoolIdFlag   = "resource-pool-id"
	exportPolicyNameFlag = "export-policy-name"
	hardLimitFlag        = "hard-limit"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Name             string
	ExportPolicyName *string
	ResourcePoolId   string
	HardLimit        *int64
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a share",
		Long:  "Creates a share.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a share in a resource pool with ID "xxx", name "yyy" and no space hard limit`,
				"$ stackit beta sfs share create --resource-pool-id xxx --name yyy --hard-limit 0",
			),
			examples.NewExample(
				`Create a share in a resource pool with ID "xxx", name "yyy" and export policy with name "zzz"`,
				"$ stackit beta sfs share create --resource-pool-id xxx --name yyy --export-policy-name zzz --hard-limit 0",
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return fmt.Errorf("unable to parse input: %w", err)
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			resourcePoolLabel, err := sfsUtils.GetResourcePoolName(ctx, apiClient, model.ProjectId, model.Region, model.ResourcePoolId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get resource pool name: %v", err)
				resourcePoolLabel = model.ResourcePoolId
			} else if resourcePoolLabel == "" {
				resourcePoolLabel = model.ResourcePoolId
			}

			prompt := fmt.Sprintf("Are you sure you want to create a SFS share for resource pool %q?", resourcePoolLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create SFS share: %w", err)
			}
			var shareId string
			if resp != nil && resp.HasShare() && resp.Share.HasId() {
				shareId = *resp.Share.Id
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Creating share")
				_, err = wait.CreateShareWaitHandler(ctx, apiClient, model.ProjectId, model.Region, model.ResourcePoolId, shareId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("waiting for share creation: %w", err)
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
	cmd.Flags().String(nameFlag, "", "Share name")
	cmd.Flags().Var(flags.UUIDFlag(), resourcePoolIdFlag, "The resource pool the share is assigned to")
	cmd.Flags().String(exportPolicyNameFlag, "", "The export policy the share is assigned to")
	cmd.Flags().Int64(hardLimitFlag, 0, "The space hard limit for the share")

	err := flags.MarkFlagsRequired(cmd, nameFlag, resourcePoolIdFlag, hardLimitFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	hardLimit := flags.FlagToInt64Pointer(p, cmd, hardLimitFlag)
	if hardLimit != nil {
		if *hardLimit < 0 {
			return nil, &errors.FlagValidationError{
				Flag:    hardLimitFlag,
				Details: "must be a positive integer",
			}
		}
	}

	model := inputModel{
		GlobalFlagModel:  globalFlags,
		Name:             flags.FlagToStringValue(p, cmd, nameFlag),
		ResourcePoolId:   flags.FlagToStringValue(p, cmd, resourcePoolIdFlag),
		ExportPolicyName: flags.FlagToStringPointer(p, cmd, exportPolicyNameFlag),
		HardLimit:        hardLimit,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *sfs.APIClient) sfs.ApiCreateShareRequest {
	req := apiClient.CreateShare(ctx, model.ProjectId, model.Region, model.ResourcePoolId)
	req = req.CreateSharePayload(
		sfs.CreateSharePayload{
			Name:                    utils.Ptr(model.Name),
			ExportPolicyName:        sfs.NewNullableString(model.ExportPolicyName),
			SpaceHardLimitGigabytes: model.HardLimit,
		},
	)
	return req
}

func outputResult(p *print.Printer, outputFormat string, async bool, resourcePoolLabel string, item *sfs.CreateShareResponse) error {
	return p.OutputResult(outputFormat, item, func() error {
		if item == nil || item.Share == nil {
			p.Outputln("SFS share response is empty")
			return nil
		}
		operation := "Created"
		if async {
			operation = "Triggered creation of"
		}
		p.Outputf(
			"%s SFS Share %q in resource pool %q.\nShare ID: %s\n",
			operation,
			utils.PtrString(item.Share.Name),
			resourcePoolLabel,
			utils.PtrString(item.Share.Id),
		)
		return nil
	})
}
