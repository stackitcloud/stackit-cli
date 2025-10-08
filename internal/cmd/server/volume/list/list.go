package list

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	serverIdFlag = "server-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ServerId *string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all server volumes",
		Long:  "Lists all server volumes.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all volumes for a server with ID "xxx"`,
				"$ stackit server volume list --server-id xxx"),
			examples.NewExample(
				`List all volumes for a server with ID "xxx" in JSON format`,
				"$ stackit server volumes list --server-id xxx --output-format json"),
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

			serverLabel, err := iaasUtils.GetServerName(ctx, apiClient, model.ProjectId, *model.ServerId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get server name: %v", err)
				serverLabel = *model.ServerId
			} else if serverLabel == "" {
				serverLabel = *model.ServerId
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("list server volumes: %w", err)
			}
			volumes := *resp.Items
			if len(volumes) == 0 {
				params.Printer.Info("No volumes found for server %s\n", serverLabel)
				return nil
			}

			// get volume names
			var volumeNames []string
			for i := range volumes {
				volumeLabel, err := iaasUtils.GetVolumeName(ctx, apiClient, model.ProjectId, *volumes[i].VolumeId)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get volume name: %v", err)
					volumeLabel = ""
				}
				volumeNames = append(volumeNames, volumeLabel)
			}

			return outputResult(params.Printer, model.OutputFormat, serverLabel, volumeNames, volumes)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().VarP(flags.UUIDFlag(), serverIdFlag, "s", "Server ID")

	err := flags.MarkFlagsRequired(cmd, serverIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ServerId:        flags.FlagToStringPointer(p, cmd, serverIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiListAttachedVolumesRequest {
	req := apiClient.ListAttachedVolumes(ctx, model.ProjectId, *model.ServerId)
	return req
}

func outputResult(p *print.Printer, outputFormat, serverLabel string, volumeNames []string, volumes []iaas.VolumeAttachment) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(volumes, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal server volume list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(volumes, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal server volume list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("SERVER ID", "SERVER NAME", "VOLUME ID", "VOLUME NAME")
		for i := range volumes {
			s := volumes[i]
			var volumeName string
			if len(volumeNames)-1 > i {
				volumeName = volumeNames[i]
			}
			table.AddRow(utils.PtrString(s.ServerId), serverLabel, utils.PtrString(s.VolumeId), volumeName)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	}
}
