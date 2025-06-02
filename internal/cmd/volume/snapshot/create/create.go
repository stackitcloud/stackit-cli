package create

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas/wait"
)

const (
	volumeIdFlag = "volume-id"
	nameFlag     = "name"
	labelsFlag   = "labels"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	VolumeID string
	Name     *string
	Labels   map[string]string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a snapshot from a volume",
		Long:  "Creates a snapshot from a volume.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a snapshot from a volume`,
				"$ stackit volume snapshot create --volume-id xxx --project-id xxx"),
			examples.NewExample(
				`Create a snapshot with a name`,
				"$ stackit volume snapshot create --volume-id xxx --name my-snapshot --project-id xxx"),
			examples.NewExample(
				`Create a snapshot with labels`,
				"$ stackit volume snapshot create --volume-id xxx --labels key1=value1,key2=value2 --project-id xxx"),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			// Get volume name for label
			volumeLabel := model.VolumeID
			volume, err := apiClient.GetVolume(ctx, model.ProjectId, model.VolumeID).Execute()
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get volume name: %v", err)
			} else if volume != nil && volume.Name != nil {
				volumeLabel = *volume.Name
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create snapshot from volume %q? (This cannot be undone)", volumeLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create snapshot: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Creating snapshot")
				resp, err = wait.CreateSnapshotWaitHandler(ctx, apiClient, model.ProjectId, *resp.Id).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for snapshot creation: %w", err)
				}
				s.Stop()
			}

			if model.Async {
				params.Printer.Info("Triggered snapshot of %q in %q. Snapshot ID: %s\n", volumeLabel, projectLabel, *resp.Id)
			} else {
				params.Printer.Info("Created snapshot of %q in %q. Snapshot ID: %s\n", volumeLabel, projectLabel, *resp.Id)
			}
			return nil
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(volumeIdFlag, "", "ID of the volume from which a snapshot should be created")
	cmd.Flags().String(nameFlag, "", "Name of the snapshot")
	cmd.Flags().StringToString(labelsFlag, nil, "Key-value string pairs as labels")

	err := flags.MarkFlagsRequired(cmd, volumeIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	volumeID := flags.FlagToStringValue(p, cmd, volumeIdFlag)
	if volumeID == "" {
		return nil, fmt.Errorf("volume-id is required")
	}
	if err := utils.ValidateUUID(volumeID); err != nil {
		return nil, fmt.Errorf("volume-id must be a valid UUID")
	}

	name := flags.FlagToStringPointer(p, cmd, nameFlag)
	labels := flags.FlagToStringToStringPointer(p, cmd, labelsFlag)
	if labels == nil {
		labels = &map[string]string{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		VolumeID:        volumeID,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiCreateSnapshotRequest {
	req := apiClient.CreateSnapshot(ctx, model.ProjectId)
	payload := iaas.NewCreateSnapshotPayloadWithDefaults()
	payload.VolumeId = &model.VolumeID
	payload.Name = model.Name

	// Convert labels to map[string]interface{}
	if len(model.Labels) > 0 {
		labelsMap := map[string]interface{}{}
		for k, v := range model.Labels {
			labelsMap[k] = v
		}
		payload.Labels = &labelsMap
	}

	req = req.CreateSnapshotPayload(*payload)
	return req
}
