package list

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
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
	limitFlag    = "limit"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ServerId *string
	Limit    *int64
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all attached network interfaces of a server",
		Long:  "Lists all attached network interfaces of a server.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Lists all attached network interfaces of server with ID "xxx"`,
				"$ stackit server network-interface list --server-id xxx",
			),
			examples.NewExample(
				`Lists all attached network interfaces of server with ID "xxx" in JSON format`,
				"$ stackit server network-interface list --server-id xxx --output-format json",
			),
			examples.NewExample(
				`Lists up to 10 attached network interfaces of server with ID "xxx"`,
				"$ stackit server network-interface list --server-id xxx --limit 10",
			),
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("list attached network interfaces: %w", err)
			}

			if resp.Items == nil || len(*resp.Items) == 0 {
				serverLabel, err := iaasUtils.GetServerName(ctx, apiClient, model.ProjectId, *model.ServerId)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get server name: %v", err)
					serverLabel = *model.ServerId
				} else if serverLabel == "" {
					serverLabel = *model.ServerId
				}
				params.Printer.Info("No attached network interfaces found for server %q\n", serverLabel)
				return nil
			}

			// Truncate output
			items := *resp.Items
			if model.Limit != nil && len(items) > int(*model.Limit) {
				items = items[:*model.Limit]
			}

			return outputResult(params.Printer, model.OutputFormat, *model.ServerId, items)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), serverIdFlag, "Server ID")
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")

	err := flags.MarkFlagsRequired(cmd, serverIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	limit := flags.FlagToInt64Pointer(p, cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ServerId:        flags.FlagToStringPointer(p, cmd, serverIdFlag),
		Limit:           limit,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiListServerNicsRequest {
	return apiClient.ListServerNics(ctx, model.ProjectId, *model.ServerId)
}

func outputResult(p *print.Printer, outputFormat, serverId string, serverNics []iaas.NIC) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(serverNics, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal server network interfaces: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(serverNics, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal server network interfaces: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("NIC ID", "SERVER ID")

		for i := range serverNics {
			nic := serverNics[i]
			table.AddRow(utils.PtrString(nic.Id), serverId)
		}
		table.EnableAutoMergeOnColumns(2)

		p.Outputln(table.Render())
		return nil
	}
}
