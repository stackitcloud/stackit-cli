package create

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
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

	"github.com/spf13/cobra"
)

const (
	availabilityZoneFlag = "availability-zone"
	nameFlag             = "name"
	descriptionFlag      = "description"
	labelFlag            = "labels"
	performanceClassFlag = "performance-class"
	sizeFlag             = "size"
	sourceIdFlag         = "source-id"
	sourceTypeFlag       = "source-type"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	AvailabilityZone *string
	Name             *string
	Description      *string
	Labels           *map[string]string
	PerformanceClass *string
	Size             *int64
	SourceId         *string
	SourceType       *string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a volume",
		Long:  "Creates a volume.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a volume with availability zone "eu01-1" and size 64 GB`,
				`$ stackit volume create --availability-zone eu01-1 --size 64`,
			),
			examples.NewExample(
				`Create a volume with availability zone "eu01-1", size 64 GB and labels`,
				`$ stackit volume create --availability-zone eu01-1 --size 64 --labels key=value,foo=bar`,
			),
			examples.NewExample(
				`Create a volume with name "volume-1", from a source image with ID "xxx"`,
				`$ stackit volume create --availability-zone eu01-1 --name volume-1 --source-id xxx --source-type image`,
			),
			examples.NewExample(
				`Create a volume with availability zone "eu01-1", performance class "storage_premium_perf1" and size 64 GB`,
				`$ stackit volume create --availability-zone eu01-1 --performance-class storage_premium_perf1 --size 64`,
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

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a volume for project %q?", projectLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create volume : %w", err)
			}
			volumeId := *resp.Id

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Creating volume")
				_, err = wait.CreateVolumeWaitHandler(ctx, apiClient, model.ProjectId, model.Region, volumeId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for volume creation: %w", err)
				}
				s.Stop()
			}

			return outputResult(params.Printer, model, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(availabilityZoneFlag, "", "Availability zone")
	cmd.Flags().StringP(nameFlag, "n", "", "Volume name")
	cmd.Flags().String(descriptionFlag, "", "Volume description")
	cmd.Flags().StringToString(labelFlag, nil, "Labels are key-value string pairs which can be attached to a volume. E.g. '--labels key1=value1,key2=value2,...'")
	cmd.Flags().String(performanceClassFlag, "", "Performance class")
	cmd.Flags().Int64(sizeFlag, 0, "Volume size (GB). Either 'size' or the 'source-id' and 'source-type' flags must be given")
	cmd.Flags().String(sourceIdFlag, "", "ID of the source object of volume. Either 'size' or the 'source-id' and 'source-type' flags must be given")
	cmd.Flags().String(sourceTypeFlag, "", "Type of the source object of volume. Either 'size' or the 'source-id' and 'source-type' flags must be given")

	err := flags.MarkFlagsRequired(cmd, availabilityZoneFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:  globalFlags,
		AvailabilityZone: flags.FlagToStringPointer(p, cmd, availabilityZoneFlag),
		Name:             flags.FlagToStringPointer(p, cmd, nameFlag),
		Description:      flags.FlagToStringPointer(p, cmd, descriptionFlag),
		Labels:           flags.FlagToStringToStringPointer(p, cmd, labelFlag),
		PerformanceClass: flags.FlagToStringPointer(p, cmd, performanceClassFlag),
		Size:             flags.FlagToInt64Pointer(p, cmd, sizeFlag),
		SourceId:         flags.FlagToStringPointer(p, cmd, sourceIdFlag),
		SourceType:       flags.FlagToStringPointer(p, cmd, sourceTypeFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiCreateVolumeRequest {
	req := apiClient.CreateVolume(ctx, model.ProjectId, model.Region)
	source := &iaas.VolumeSource{
		Id:   model.SourceId,
		Type: model.SourceType,
	}

	payload := iaas.CreateVolumePayload{
		AvailabilityZone: model.AvailabilityZone,
		Name:             model.Name,
		Description:      model.Description,
		Labels:           utils.ConvertStringMapToInterfaceMap(model.Labels),
		PerformanceClass: model.PerformanceClass,
		Size:             model.Size,
	}

	if model.SourceId != nil && model.SourceType != nil {
		payload.Source = source
	}

	return req.CreateVolumePayload(payload)
}

func outputResult(p *print.Printer, model *inputModel, projectLabel string, volume *iaas.Volume) error {
	if volume == nil {
		return fmt.Errorf("volume response is empty")
	}
	return p.OutputResult(model.OutputFormat, volume, func() error {
		p.Outputf("Created volume for project %q.\nVolume ID: %s\n", projectLabel, utils.PtrString(volume.Id))
		return nil
	})
}
