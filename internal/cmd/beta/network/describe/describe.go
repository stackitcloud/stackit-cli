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
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"

	"github.com/spf13/cobra"
)

const (
	networkIdArg = "NETWORK_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	NetworkId string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Shows details of a network",
		Long:  "Shows details of a network.",
		Args:  args.SingleArg(networkIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Show details of a network with ID "xxx"`,
				"$ stackit beta network describe xxx",
			),
			examples.NewExample(
				`Show details of a network with ID "xxx" in JSON format`,
				"$ stackit beta network describe xxx --output-format json",
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
				return fmt.Errorf("read network: %w", err)
			}

			return outputResult(p, model.OutputFormat, resp)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	networkId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		NetworkId:       networkId,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiGetNetworkRequest {
	return apiClient.GetNetwork(ctx, model.ProjectId, model.NetworkId)
}

func outputResult(p *print.Printer, outputFormat string, network *iaas.Network) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(network, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal network: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(network, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal network: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		var nameservers []string
		if network.Nameservers != nil {
			for _, nameserver := range *network.Nameservers {
				nameservers = append(nameservers, nameserver)
			}
		}

		var prefixes []string
		if network.Prefixes != nil {
			for _, prefix := range *network.Prefixes {
				prefixes = append(prefixes, prefix)
			}
		}

		table := tables.NewTable()
		table.AddRow("ID", *network.NetworkId)
		table.AddSeparator()
		table.AddRow("NAME", *network.Name)
		table.AddSeparator()
		table.AddRow("STATE", *network.State)
		table.AddSeparator()
		table.AddRow("PUBLIC IP", *network.PublicIp)
		table.AddSeparator()
		if len(nameservers) > 0 {
			table.AddRow("NAME SERVERS", strings.Join(nameservers, ","))
		}
		table.AddSeparator()
		if len(prefixes) > 0 {
			table.AddRow("PREFIXES", strings.Join(prefixes, ","))
		}
		table.AddSeparator()

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	}
}
