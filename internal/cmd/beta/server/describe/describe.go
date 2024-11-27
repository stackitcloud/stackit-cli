package describe

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"

	"github.com/spf13/cobra"
)

const (
	serverIdArg = "SERVER_ID"
	detailsFlag = "details"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ServerId string
	Details  bool
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Shows details of a server",
		Long:  "Shows details of a server.",
		Args:  args.SingleArg(serverIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Show details of a server with ID "xxx"`,
				"$ stackit beta server describe xxx",
			),
			examples.NewExample(
				`Show detailed information of a server with ID "xxx"`,
				"$ stackit beta server describe xxx --details",
			),
			examples.NewExample(
				`Show details of a server with ID "xxx" in JSON format`,
				"$ stackit beta server describe xxx --output-format json",
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read server: %w", err)
			}

			return outputResult(p, model, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Bool(detailsFlag, false, "Show detailed information about server")
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
		Details:         flags.FlagToBoolValue(p, cmd, detailsFlag),
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

	if model.Details {
		req = req.Details(true)
	}

	return req
}

func outputResult(p *print.Printer, model *inputModel, server *iaas.Server) error {
	outputFormat := model.OutputFormat

	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(server, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal server: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(server, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal server: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.AddRow("ID", *server.Id)
		table.AddSeparator()
		table.AddRow("NAME", *server.Name)
		table.AddSeparator()
		table.AddRow("STATE", *server.Status)
		table.AddSeparator()
		table.AddRow("AVAILABILITY ZONE", *server.AvailabilityZone)
		table.AddSeparator()
		table.AddRow("BOOT VOLUME", *server.BootVolume.Id)
		table.AddSeparator()
		table.AddRow("POWER STATUS", *server.PowerStatus)
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

		if server.Volumes != nil && len(*server.Volumes) > 0 {
			volumes := []string{}
			volumes = append(volumes, *server.Volumes...)
			table.AddRow("VOLUMES", strings.Join(volumes, "\n"))
			table.AddSeparator()
		}

		if model.Details {
			if server.ServiceAccountMails != nil && len(*server.ServiceAccountMails) > 0 {
				emails := []string{}
				emails = append(emails, *server.ServiceAccountMails...)
				table.AddRow("SERVICE ACCOUNTS", strings.Join(emails, "\n"))
				table.AddSeparator()
			}

			if server.Nics != nil && len(*server.Nics) > 0 {
				nics := []string{}
				for _, nic := range *server.Nics {
					nics = append(nics, *nic.NicId)
				}
				table.AddRow("NICS", strings.Join(nics, "\n"))
				table.AddSeparator()
			}
		}

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	}
}
