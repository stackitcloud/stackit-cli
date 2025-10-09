package update

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
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

	serverIdFlag            = "server-id"
	deleteOnTerminationFlag = "delete-on-termination"

	defaultDeleteOnTermination = false
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ServerId            *string
	VolumeId            string
	DeleteOnTermination *bool
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", volumeIdArg),
		Short: "Updates an attached volume of a server",
		Long:  "Updates an attached volume of a server.",
		Args:  args.SingleArg(volumeIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update a volume with ID "xxx" of a server with ID "yyy" and enables delete on termination`,
				`$ stackit server volume update xxx --server-id yyy --delete-on-termination`,
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

			volumeLabel, err := iaasUtils.GetVolumeName(ctx, apiClient, model.ProjectId, model.VolumeId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get volume name: %v", err)
				volumeLabel = model.VolumeId
			}

			serverLabel, err := iaasUtils.GetServerName(ctx, apiClient, model.ProjectId, *model.ServerId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get server name: %v", err)
				serverLabel = *model.ServerId
			} else if serverLabel == "" {
				serverLabel = *model.ServerId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update attached volume %q of server %q?", volumeLabel, serverLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update server volume: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, volumeLabel, serverLabel, *resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), serverIdFlag, "Server ID")
	cmd.Flags().BoolP(deleteOnTerminationFlag, "b", defaultDeleteOnTermination, "Delete the volume during the termination of the server. (default false)")

	err := flags.MarkFlagsRequired(cmd, serverIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	volumeId := inputArgs[0]
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:     globalFlags,
		ServerId:            flags.FlagToStringPointer(p, cmd, serverIdFlag),
		DeleteOnTermination: flags.FlagToBoolPointer(p, cmd, deleteOnTerminationFlag),
		VolumeId:            volumeId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiUpdateAttachedVolumeRequest {
	req := apiClient.UpdateAttachedVolume(ctx, model.ProjectId, *model.ServerId, model.VolumeId)
	payload := iaas.UpdateAttachedVolumePayload{
		DeleteOnTermination: model.DeleteOnTermination,
	}
	return req.UpdateAttachedVolumePayload(payload)
}

func outputResult(p *print.Printer, outputFormat, volumeLabel, serverLabel string, volume iaas.VolumeAttachment) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(volume, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal server volume: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(volume, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal server volume: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		p.Outputf("Updated attached volume %q of server %q\n", volumeLabel, serverLabel)
		return nil
	}
}
