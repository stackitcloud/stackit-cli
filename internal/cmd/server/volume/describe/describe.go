package describe

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	volumeIdArg  = "VOLUME_ID"
	serverIdFlag = "server-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ServerId string
	VolumeId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", volumeIdArg),
		Short: "Describes a server volume attachment",
		Long:  "Describes a server volume attachment.",
		Args:  args.SingleArg(volumeIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get details of the attachment of volume with ID "xxx" to server with ID "yyy"`,
				`$ stackit server volume describe xxx --server-id yyy`,
			),
			examples.NewExample(
				`Get details of the attachment of volume with ID "xxx" to server with ID "yyy" in JSON format`,
				`$ stackit server volume describe xxx --server-id yyy --output-format json`,
			),
			examples.NewExample(
				`Get details of the attachment of volume with ID "xxx" to server with ID "yyy" in yaml format`,
				`$ stackit server volume describe xxx --server-id yyy --output-format yaml`,
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

			volumeLabel, err := iaasUtils.GetVolumeName(ctx, apiClient, model.ProjectId, model.Region, model.VolumeId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get volume name: %v", err)
				volumeLabel = model.VolumeId
			}

			serverLabel, err := iaasUtils.GetServerName(ctx, apiClient, model.ProjectId, model.Region, model.ServerId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get server name: %v", err)
				serverLabel = model.ServerId
			} else if serverLabel == "" {
				serverLabel = model.ServerId
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("describe server volume: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, serverLabel, volumeLabel, *resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), serverIdFlag, "Server ID")

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
		GlobalFlagModel: globalFlags,
		ServerId:        flags.FlagToStringValue(p, cmd, serverIdFlag),
		VolumeId:        volumeId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiGetAttachedVolumeRequest {
	req := apiClient.GetAttachedVolume(ctx, model.ProjectId, model.Region, model.ServerId, model.VolumeId)
	return req
}

func outputResult(p *print.Printer, outputFormat, serverLabel, volumeLabel string, volume iaas.VolumeAttachment) error {
	return p.OutputResult(outputFormat, volume, func() error {
		table := tables.NewTable()
		table.AddRow("SERVER ID", utils.PtrString(volume.ServerId))
		table.AddSeparator()
		table.AddRow("SERVER NAME", serverLabel)
		table.AddSeparator()
		table.AddRow("VOLUME ID", utils.PtrString(volume.VolumeId))
		table.AddSeparator()
		// check if name is set
		if volumeLabel != "" {
			table.AddRow("VOLUME NAME", volumeLabel)
			table.AddSeparator()
		}
		table.AddRow("DELETE ON TERMINATION", utils.PtrString(volume.DeleteOnTermination))
		table.AddSeparator()

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
