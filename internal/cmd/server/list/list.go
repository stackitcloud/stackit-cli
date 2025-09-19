package list

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	sdkutils "github.com/stackitcloud/stackit-sdk-go/core/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	limitFlag         = "limit"
	labelSelectorFlag = "label-selector"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit         *int64
	LabelSelector *string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all servers of a project",
		Long:  "Lists all servers of a project.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Lists all servers`,
				"$ stackit server list",
			),
			examples.NewExample(
				`Lists all servers which contains the label xxx`,
				"$ stackit server list --label-selector xxx",
			),
			examples.NewExample(
				`Lists all servers in JSON format`,
				"$ stackit server list --output-format json",
			),
			examples.NewExample(
				`Lists up to 10 servers`,
				"$ stackit server list --limit 10",
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
				return fmt.Errorf("list servers: %w", err)
			}

			if resp.Items == nil || len(*resp.Items) == 0 {
				projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
					projectLabel = model.ProjectId
				}
				params.Printer.Info("No servers found for project %q\n", projectLabel)
				return nil
			}

			// Truncate output
			items := *resp.Items
			if model.Limit != nil && len(items) > int(*model.Limit) {
				items = items[:*model.Limit]
			}

			return outputResult(params.Printer, model.OutputFormat, items)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
	cmd.Flags().String(labelSelectorFlag, "", "Filter by label")
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
		Limit:           limit,
		LabelSelector:   flags.FlagToStringPointer(p, cmd, labelSelectorFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiListServersRequest {
	req := apiClient.ListServers(ctx, model.ProjectId)
	if model.LabelSelector != nil {
		req = req.LabelSelector(*model.LabelSelector)
	}
	req = req.Details(true)

	return req
}

func outputResult(p *print.Printer, outputFormat string, servers []iaas.Server) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(servers, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal server: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		// Convert each server to map with base64 encoded byte arrays
		convertedServers := make([]map[string]interface{}, len(servers))
		for i, server := range servers {
			converted, err := sdkutils.ConvertByteArraysToBase64(server)
			if err != nil {
				return fmt.Errorf("convert server for YAML: %w", err)
			}
			convertedServers[i] = converted
		}

		details, err := yaml.MarshalWithOptions(convertedServers, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal server: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("ID", "Name", "Status", "Machine Type", "Availability Zones", "Nic IPv4", "Public IPs")

		for i := range servers {
			server := servers[i]

			nicIPv4 := ""
			publicIPs := ""
			if server.Nics != nil && len(*server.Nics) > 0 {
				for i, nic := range *server.Nics {
					if nic.Ipv4 != nil || nic.PublicIp != nil {
						nicIPv4 += utils.PtrString(nic.Ipv4)
						publicIPs += utils.PtrString(nic.PublicIp)

						if i != len(*server.Nics)-1 {
							publicIPs += "\n"
							nicIPv4 += "\n"
						}
					}
				}
			}

			table.AddRow(
				utils.PtrString(server.Id),
				utils.PtrString(server.Name),
				utils.PtrString(server.Status),
				utils.PtrString(server.MachineType),
				utils.PtrString(server.AvailabilityZone),
				nicIPv4,
				publicIPs,
			)
		}

		p.Outputln(table.Render())
		return nil
	}
}
