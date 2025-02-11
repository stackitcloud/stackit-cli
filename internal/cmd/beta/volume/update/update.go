package update

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	volumeIdArg = "VOLUME_ID"

	nameFlag        = "name"
	descriptionFlag = "description"
	labelFlag       = "labels"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	VolumeId    string
	Name        *string
	Description *string
	Labels      *map[string]string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", volumeIdArg),
		Short: "Updates a volume",
		Long:  "Updates a volume.",
		Args:  args.SingleArg(volumeIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update volume with ID "xxx" with new name "volume-1-new"`,
				`$ stackit beta volume update xxx --name volume-1-new`,
			),
			examples.NewExample(
				`Update volume with ID "xxx" with new name "volume-1-new" and new description "volume-1-desc-new"`,
				`$ stackit beta volume update xxx --name volume-1-new --description volume-1-desc-new`,
			),
			examples.NewExample(
				`Update volume with ID "xxx" with new name "volume-1-new", new description "volume-1-desc-new" and label(s)`,
				`$ stackit beta volume update xxx --name volume-1-new --description volume-1-desc-new --labels key=value,foo=bar`,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			volumeLabel, err := iaasUtils.GetVolumeName(ctx, apiClient, model.ProjectId, model.VolumeId)
			if err != nil {
				p.Debug(print.ErrorLevel, "get volume name: %v", err)
				volumeLabel = model.VolumeId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update volume %q?", volumeLabel)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update volume: %w", err)
			}

			return outputResult(p, model, volumeLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(nameFlag, "n", "", "Volume name")
	cmd.Flags().String(descriptionFlag, "", "Volume description")
	cmd.Flags().StringToString(labelFlag, nil, "Labels are key-value string pairs which can be attached to a volume. E.g. '--labels key1=value1,key2=value2,...'")
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	volumeId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Name:            flags.FlagToStringPointer(p, cmd, nameFlag),
		VolumeId:        volumeId,
		Description:     flags.FlagToStringPointer(p, cmd, descriptionFlag),
		Labels:          flags.FlagToStringToStringPointer(p, cmd, labelFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiUpdateVolumeRequest {
	req := apiClient.UpdateVolume(ctx, model.ProjectId, model.VolumeId)

	var labelsMap *map[string]interface{}
	if model.Labels != nil && len(*model.Labels) > 0 {
		// convert map[string]string to map[string]interface{}
		labelsMap = utils.Ptr(map[string]interface{}{})
		for k, v := range *model.Labels {
			(*labelsMap)[k] = v
		}
	}

	payload := iaas.UpdateVolumePayload{
		Name:        model.Name,
		Description: model.Description,
		Labels:      labelsMap,
	}

	return req.UpdateVolumePayload(payload)
}

func outputResult(p *print.Printer, model *inputModel, volumeLabel string, volume *iaas.Volume) error {
	switch model.OutputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(volume, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal volume: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(volume, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal volume: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		p.Outputf("Updated volume %q.\n", volumeLabel)
		return nil
	}
}
