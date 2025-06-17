package update

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	snapshotIdArg = "SNAPSHOT_ID"
	nameFlag      = "name"
	labelsFlag    = "labels"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	SnapshotId string
	Name       *string
	Labels     map[string]string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", snapshotIdArg),
		Short: "Updates a snapshot",
		Long:  "Updates a snapshot by its ID.",
		Args:  args.SingleArg(snapshotIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update a snapshot name with ID "xxx"`,
				"$ stackit volume snapshot update xxx --name my-new-name"),
			examples.NewExample(
				`Update a snapshot labels with ID "xxx"`,
				"$ stackit volume snapshot update xxx --labels key1=value1,key2=value2"),
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

			// Get snapshot name for label
			snapshotLabel := model.SnapshotId
			snapshot, err := apiClient.GetSnapshot(ctx, model.ProjectId, model.SnapshotId).Execute()
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get snapshot name: %v", err)
			} else if snapshot != nil && snapshot.Name != nil {
				snapshotLabel = *snapshot.Name
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update snapshot %q?", snapshotLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("update snapshot: %w", err)
			}

			params.Printer.Outputf("Updated snapshot %q\n", snapshotLabel)
			return nil
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(nameFlag, "", "Name of the snapshot")
	cmd.Flags().StringToString(labelsFlag, nil, "Key-value string pairs as labels")
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	snapshotId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	name := flags.FlagToStringPointer(p, cmd, nameFlag)
	labels := flags.FlagToStringToStringPointer(p, cmd, labelsFlag)
	if labels == nil {
		labels = &map[string]string{}
	}

	if name == nil && len(*labels) == 0 {
		return nil, fmt.Errorf("either name or labels must be provided")
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		SnapshotId:      snapshotId,
		Name:            name,
		Labels:          *labels,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiUpdateSnapshotRequest {
	req := apiClient.UpdateSnapshot(ctx, model.ProjectId, model.SnapshotId)
	payload := iaas.NewUpdateSnapshotPayloadWithDefaults()
	payload.Name = model.Name
	payload.Labels = utils.ConvertStringMapToInterfaceMap(utils.Ptr(model.Labels))

	req = req.UpdateSnapshotPayload(*payload)
	return req
}
