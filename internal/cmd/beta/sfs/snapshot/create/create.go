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
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/sfs"
)

const (
	resourcePoolIdFlag = "resource-pool-id"
	nameFlag           = "name"
	commentFlag        = "comment"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ResourcePoolId string
	Name           string
	Comment        *string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a new snapshot of a resource pool",
		Long:  "Creates a new snapshot of a resource pool.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a new snapshot with name "snapshot-name" of a resource pool with ID "xxx"`,
				"$ stackit beta sfs snapshot create --name snapshot-name --resource-pool-id xxx",
			),
			examples.NewExample(
				`Create a new snapshot with name "snapshot-name" and comment "snapshot-comment" of a resource pool with ID "xxx"`,
				`$ stackit beta sfs snapshot create --name snapshot-name --resource-pool-id xxx --comment "snapshot-comment"`,
			),
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

			resourcePoolLabel, err := sfsUtils.GetResourcePoolName(ctx, apiClient, model.ProjectId, model.Region, model.ResourcePoolId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get resource pool name: %v", err)
				resourcePoolLabel = model.ResourcePoolId
			} else if resourcePoolLabel == "" {
				resourcePoolLabel = model.ResourcePoolId
			}

			prompt := fmt.Sprintf("Are you sure you want to create a snapshot for resource pool %q?", resourcePoolLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create snapshot: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, model.Name, resourcePoolLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(nameFlag, "", "Snapshot name")
	cmd.Flags().String(commentFlag, "", "A comment to add more information to the snapshot")
	cmd.Flags().Var(flags.UUIDFlag(), resourcePoolIdFlag, "The resource pool from which the snapshot should be created")

	err := flags.MarkFlagsRequired(cmd, resourcePoolIdFlag, nameFlag)
	cobra.CheckErr(err)
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *sfs.APIClient) sfs.ApiCreateResourcePoolSnapshotRequest {
	req := apiClient.CreateResourcePoolSnapshot(ctx, model.ProjectId, model.Region, model.ResourcePoolId)
	req = req.CreateResourcePoolSnapshotPayload(sfs.CreateResourcePoolSnapshotPayload{
		Name:    utils.Ptr(model.Name),
		Comment: sfs.NewNullableString(model.Comment),
	})
	return req
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Name:            flags.FlagToStringValue(p, cmd, nameFlag),
		ResourcePoolId:  flags.FlagToStringValue(p, cmd, resourcePoolIdFlag),
		Comment:         flags.FlagToStringPointer(p, cmd, commentFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func outputResult(p *print.Printer, outputFormat, snapshotLabel, resourcePoolLabel string, resp *sfs.CreateResourcePoolSnapshotResponse) error {
	return p.OutputResult(outputFormat, resp, func() error {
		if resp == nil || resp.ResourcePoolSnapshot == nil {
			p.Outputln("SFS snapshot response is empty")
			return nil
		}

		p.Outputf(
			"Created snapshot %q for resource pool %q.\n",
			snapshotLabel,
			resourcePoolLabel,
		)
		return nil
	})
}
