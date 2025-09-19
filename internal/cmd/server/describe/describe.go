package describe

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	serverIdArg = "SERVER_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ServerId string
}

// convertServerForYAML creates a map with UserData as base64 string for YAML output
func convertServerForYAML(server *iaas.Server) map[string]interface{} {
	if server == nil {
		return nil
	}

	// Marshal to JSON first to get the correct format for UserData
	jsonData, err := json.Marshal(server)
	if err != nil {
		return nil
	}

	var serverMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &serverMap); err != nil {
		return nil
	}

	return serverMap
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", serverIdArg),
		Short: "Shows details of a server",
		Long:  "Shows details of a server.",
		Args:  args.SingleArg(serverIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Show details of a server with ID "xxx"`,
				"$ stackit server describe xxx",
			),
			examples.NewExample(
				`Show details of a server with ID "xxx" in JSON format`,
				"$ stackit server describe xxx --output-format json",
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read server: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	serverId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ServerId:        serverId,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiGetServerRequest {
	req := apiClient.GetServer(ctx, model.ProjectId, model.ServerId)
	req = req.Details(true)

	return req
}

func outputResult(p *print.Printer, outputFormat string, server *iaas.Server) error {
	if server == nil {
		return fmt.Errorf("api response is empty")
	}
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(server, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal server: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(convertServerForYAML(server), yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal server: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		content := []tables.Table{}

		table := tables.NewTable()
		table.SetTitle("Server")

		table.AddRow("ID", utils.PtrString(server.Id))
		table.AddSeparator()
		table.AddRow("NAME", utils.PtrString(server.Name))
		table.AddSeparator()
		table.AddRow("STATE", utils.PtrString(server.Status))
		table.AddSeparator()
		table.AddRow("AVAILABILITY ZONE", utils.PtrString(server.AvailabilityZone))
		table.AddSeparator()
		if server.BootVolume != nil && server.BootVolume.Id != nil {
			table.AddRow("BOOT VOLUME", *server.BootVolume.Id)
			table.AddSeparator()
			table.AddRow("DELETE ON TERMINATION", utils.PtrString(server.BootVolume.DeleteOnTermination))
			table.AddSeparator()
		}
		table.AddRow("POWER STATUS", utils.PtrString(server.PowerStatus))
		table.AddSeparator()

		if server.AffinityGroup != nil {
			table.AddRow("AFFINITY GROUP", *server.AffinityGroup)
			table.AddSeparator()
		}

		if server.ImageId != nil {
			table.AddRow("IMAGE", *server.ImageId)
			table.AddSeparator()
		}

		if server.KeypairName != nil {
			table.AddRow("KEYPAIR", *server.KeypairName)
			table.AddSeparator()
		}

		if server.MachineType != nil {
			table.AddRow("MACHINE TYPE", *server.MachineType)
			table.AddSeparator()
		}

		if server.Labels != nil && len(*server.Labels) > 0 {
			labels := []string{}
			for key, value := range *server.Labels {
				labels = append(labels, fmt.Sprintf("%s: %s", key, value))
			}
			table.AddRow("LABELS", strings.Join(labels, "\n"))
			table.AddSeparator()
		}

		if server.ServiceAccountMails != nil && len(*server.ServiceAccountMails) > 0 {
			table.AddRow("SERVICE ACCOUNTS", strings.Join(*server.ServiceAccountMails, "\n"))
			table.AddSeparator()
		}

		if server.Volumes != nil && len(*server.Volumes) > 0 {
			volumes := []string{}
			volumes = append(volumes, *server.Volumes...)
			table.AddRow("VOLUMES", strings.Join(volumes, "\n"))
			table.AddSeparator()
		}

		content = append(content, table)

		if server.Nics != nil && len(*server.Nics) > 0 {
			nicsTable := tables.NewTable()
			nicsTable.SetTitle("Attached Network Interfaces")
			nicsTable.SetHeader("ID", "NETWORK ID", "NETWORK NAME", "IPv4", "PUBLIC IP")

			for _, nic := range *server.Nics {
				nicsTable.AddRow(
					utils.PtrString(nic.NicId),
					utils.PtrString(nic.NetworkId),
					utils.PtrString(nic.NetworkName),
					utils.PtrString(nic.Ipv4),
					utils.PtrString(nic.PublicIp),
				)
				nicsTable.AddSeparator()
			}

			content = append(content, nicsTable)
		}

		if server.MaintenanceWindow != nil {
			maintenanceWindow := tables.NewTable()
			maintenanceWindow.SetTitle("Maintenance Window")

			if server.MaintenanceWindow.Status != nil {
				maintenanceWindow.AddRow("STATUS", *server.MaintenanceWindow.Status)
				maintenanceWindow.AddSeparator()
			}
			if server.MaintenanceWindow.Details != nil {
				maintenanceWindow.AddRow("DETAILS", *server.MaintenanceWindow.Details)
				maintenanceWindow.AddSeparator()
			}
			if server.MaintenanceWindow.StartsAt != nil {
				maintenanceWindow.AddRow(
					"STARTS AT",
					utils.ConvertTimePToDateTimeString(server.MaintenanceWindow.StartsAt),
				)
				maintenanceWindow.AddSeparator()
			}
			if server.MaintenanceWindow.EndsAt != nil {
				maintenanceWindow.AddRow(
					"ENDS AT",
					utils.ConvertTimePToDateTimeString(server.MaintenanceWindow.EndsAt),
				)
				maintenanceWindow.AddSeparator()
			}

			content = append(content, maintenanceWindow)
		}

		err := tables.DisplayTables(p, content)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
