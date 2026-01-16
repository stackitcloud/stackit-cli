package create

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
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

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a snapshot from a volume",
		Long:  "Creates a snapshot from a volume.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a snapshot from a volume with ID "xxx"`,
				"$ stackit volume snapshot create --volume-id xxx"),
			examples.NewExample(
				`Create a snapshot from a volume with ID "xxx" and name "my-snapshot"`,
				"$ stackit volume snapshot create --volume-id xxx --name my-snapshot"),
			examples.NewExample(
				`Create a snapshot from a volume with ID "xxx" and labels`,
				"$ stackit volume snapshot create --volume-id xxx --labels key1=value1,key2=value2"),
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

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			// Get volume name for label
			volumeLabel, err := iaasUtils.GetVolumeName(ctx, apiClient, model.ProjectId, model.Region, model.VolumeID)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get volume name: %v", err)
				volumeLabel = model.VolumeID
			}

			prompt := fmt.Sprintf("Are you sure you want to create snapshot from volume %q? (This cannot be undone)", volumeLabel)
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

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Creating snapshot")
				resp, err = wait.CreateSnapshotWaitHandler(ctx, apiClient, model.ProjectId, model.Region, *resp.Id).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for snapshot creation: %w", err)
				}
				s.Stop()
			}

			operationState := "Created"
			if model.Async {
				operationState = "Triggered creation of"
			}
			params.Printer.Outputf("%s snapshot of %q in %q. Snapshot ID: %s\n", operationState, volumeLabel, projectLabel, utils.PtrString(resp.Id))
			return nil
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), volumeIdFlag, "ID of the volume from which a snapshot should be created")
	cmd.Flags().String(nameFlag, "", "Name of the snapshot")
	cmd.Flags().StringToString(labelsFlag, nil, "Key-value string pairs as labels")

	err := flags.MarkFlagsRequired(cmd, volumeIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	volumeID := flags.FlagToStringValue(p, cmd, volumeIdFlag)

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

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiCreateSnapshotRequest {
	req := apiClient.CreateSnapshot(ctx, model.ProjectId, model.Region)
	payload := iaas.NewCreateSnapshotPayloadWithDefaults()
	payload.VolumeId = &model.VolumeID
	payload.Name = model.Name
	payload.Labels = utils.ConvertStringMapToInterfaceMap(utils.Ptr(model.Labels))

	req = req.CreateSnapshotPayload(*payload)
	return req
}
