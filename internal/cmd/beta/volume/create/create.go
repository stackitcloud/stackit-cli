package create

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
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

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a volume",
		Long:  "Creates a volume.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a volume with availability zone "eu01-1" and size 64 GB`,
				`$ stackit beta volume create --availability-zone eu01-1 --size 64`,
			),
			examples.NewExample(
				`Create a volume with availability zone "eu01-1", size 64 GB and labels`,
				`$ stackit beta volume create --availability-zone eu01-1 --size 64 --labels key=value,foo=bar`,
			),
			examples.NewExample(
				`Create a volume with name "volume-1", from a source image with ID "xxx"`,
				`$ stackit beta volume create --availability-zone eu01-1 --name volume-1 --source-id xxx --source-type image`,
			),
			examples.NewExample(
				`Create a volume with availability zone "eu01-1", performance class "storage_premium_perf1" and size 64 GB`,
				`$ stackit beta volume create --availability-zone eu01-1 --performance-class storage_premium_perf1 --size 64`,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, p, cmd)
			if err != nil {
				p.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a volume for project %q?", projectLabel)
				err = p.PromptForConfirmation(prompt)
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
				s := spinner.New(p)
				s.Start("Creating volume")
				_, err = wait.CreateVolumeWaitHandler(ctx, apiClient, model.ProjectId, volumeId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for volume creation: %w", err)
				}
				s.Stop()
			}

			return outputResult(p, model, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(availabilityZoneFlag, "", "Availability zone")
	cmd.Flags().StringP(nameFlag, "n", "", "Volume name")
	cmd.Flags().String(descriptionFlag, "", "Volume description")
	cmd.Flags().StringToString(labelFlag, nil, "Labels are key-value string pairs which can be attached to a volume. A label can be provided with the format key=value. To provide a list of labels, key=value pairs must be seperated by commas(,) e.g. --labels key=value,foo=bar.")
	cmd.Flags().String(performanceClassFlag, "", "Performance class")
	cmd.Flags().Int64(sizeFlag, 0, "Volume size (GB). Either 'size' or the 'source-id' and 'source-type' flags must be given")
	cmd.Flags().String(sourceIdFlag, "", "ID of the source object of volume. Either 'size' or the 'source-id' and 'source-type' flags must be given")
	cmd.Flags().String(sourceTypeFlag, "", "Type of the source object of volume. Either 'size' or the 'source-id' and 'source-type' flags must be given")

	err := flags.MarkFlagsRequired(cmd, availabilityZoneFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiCreateVolumeRequest {
	req := apiClient.CreateVolume(ctx, model.ProjectId)
	source := &iaas.VolumeSource{
		Id:   model.SourceId,
		Type: model.SourceType,
	}

	var labelsMap *map[string]interface{}
	if model.Labels != nil && len(*model.Labels) > 0 {
		// convert map[string]string to map[string]interface{}
		labelsMap = utils.Ptr(map[string]interface{}{})
		for k, v := range *model.Labels {
			(*labelsMap)[k] = v
		}
	}

	payload := iaas.CreateVolumePayload{
		AvailabilityZone: model.AvailabilityZone,
		Name:             model.Name,
		Description:      model.Description,
		Labels:           labelsMap,
		PerformanceClass: model.PerformanceClass,
		Size:             model.Size,
	}

	if model.SourceId != nil && model.SourceType != nil {
		payload.Source = source
	}

	return req.CreateVolumePayload(payload)
}

func outputResult(p *print.Printer, model *inputModel, projectLabel string, volume *iaas.Volume) error {
	switch model.OutputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(volume, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal volume: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(volume, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal volume: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		p.Outputf("Created volume for project %q.\nVolume ID: %s\n", projectLabel, *volume.Id)
		return nil
	}
}
