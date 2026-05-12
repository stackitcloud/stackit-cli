package update

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	sfs "github.com/stackitcloud/stackit-sdk-go/services/sfs/v1api"

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
)

const (
	snapshotNameArg = "SNAPSHOT_NAME"

	resourcePoolIdFlag  = "resource-pool-id"
	newSnapshotNameFlag = "name"
	commentFlag         = "comment"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ResourcePoolId  string
	SnapshotName    string
	NewSnapshotName string
	Comment         *string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", snapshotNameArg),
		Short: "Updates a new snapshot of a resource pool",
		Long:  "Updates a new snapshot of a resource pool.",
		Args:  args.SingleArg(snapshotNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Updates the name of a snapshot with name "snapshot-name" of a resource pool with ID "xxx"`,
				"$ stackit beta sfs snapshot update snapshot-name --resource-pool-id xxx --name new-snapshot-name",
			),
			examples.NewExample(
				`Updates the comment of a snapshot with name "snapshot-name" of a resource pool with ID "xxx"`,
				`$ stackit beta sfs snapshot update snapshot-name --resource-pool-id xxx --comment "snapshot-comment"`,
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

			resourcePoolLabel, err := sfsUtils.GetResourcePoolName(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.ResourcePoolId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get resource pool name: %v", err)
				resourcePoolLabel = model.ResourcePoolId
			} else if resourcePoolLabel == "" {
				resourcePoolLabel = model.ResourcePoolId
			}

			prompt := fmt.Sprintf("Are you sure you want to update the snapshot %q for resource pool %q?", model.SnapshotName, resourcePoolLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update snapshot: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, model.SnapshotName, resourcePoolLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(newSnapshotNameFlag, "", "Snapshot name")
	cmd.Flags().String(commentFlag, "", "A comment to add more information to the snapshot")
	cmd.Flags().Var(flags.UUIDFlag(), resourcePoolIdFlag, "The resource pool from which the snapshot should be updated")

	cmd.MarkFlagsOneRequired(newSnapshotNameFlag, commentFlag)
	err := flags.MarkFlagsRequired(cmd, resourcePoolIdFlag)
	cobra.CheckErr(err)
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *sfs.APIClient) sfs.ApiUpdateResourcePoolSnapshotRequest {
	req := apiClient.DefaultAPI.UpdateResourcePoolSnapshot(ctx, model.ProjectId, model.Region, model.ResourcePoolId, model.SnapshotName)

	payload := sfs.UpdateResourcePoolSnapshotPayload{
		Comment: *sfs.NewNullableString(model.Comment),
	}

	if model.NewSnapshotName != "" {
		payload.Name = *sfs.NewNullableString(utils.Ptr(model.NewSnapshotName))
	}
	return req.UpdateResourcePoolSnapshotPayload(payload)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	snapshotName := inputArgs[0]
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		SnapshotName:    snapshotName,
		NewSnapshotName: flags.FlagToStringValue(p, cmd, newSnapshotNameFlag),
		ResourcePoolId:  flags.FlagToStringValue(p, cmd, resourcePoolIdFlag),
		Comment:         flags.FlagToStringPointer(p, cmd, commentFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func outputResult(p *print.Printer, outputFormat, snapshotLabel, resourcePoolLabel string, resp *sfs.UpdateResourcePoolSnapshotResponse) error {
	return p.OutputResult(outputFormat, resp, func() error {
		if resp == nil || resp.ResourcePoolSnapshot == nil {
			p.Outputln("SFS snapshot response is empty")
			return nil
		}

		p.Outputf(
			"Updated snapshot %q for resource pool %q.\n",
			snapshotLabel,
			resourcePoolLabel,
		)

		if resp.ResourcePoolSnapshot.SnaplockExpiryTime.IsSet() && resp.ResourcePoolSnapshot.SnaplockExpiryTime.Get() != nil {
			p.Outputf("Snaplock expiry time: %s\n", utils.ConvertTimePToDateTimeString(resp.ResourcePoolSnapshot.SnaplockExpiryTime.Get()))
		}

		return nil
	})
}
