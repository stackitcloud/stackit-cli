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
	publicIpIdArg = "PUBLIC_IP_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	PublicIpId string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", publicIpIdArg),
		Short: "Shows details of a Public IP",
		Long:  "Shows details of a Public IP.",
		Args:  args.SingleArg(publicIpIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Show details of a public IP with ID "xxx"`,
				"$ stackit beta public-ip describe xxx",
			),
			examples.NewExample(
				`Show details of a public IP with ID "xxx" in JSON format`,
				"$ stackit beta public-ip describe xxx --output-format json",
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
				return fmt.Errorf("read public IP: %w", err)
			}

			return outputResult(p, model.OutputFormat, resp)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	publicIpId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		PublicIpId:      publicIpId,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiGetPublicIPRequest {
	return apiClient.GetPublicIP(ctx, model.ProjectId, model.PublicIpId)
}

func outputResult(p *print.Printer, outputFormat string, publicIp *iaas.PublicIp) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(publicIp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal public IP: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(publicIp, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal public IP: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.AddRow("ID", *publicIp.Id)
		table.AddSeparator()
		table.AddRow("IP ADDRESS", *publicIp.Ip)
		table.AddSeparator()

		if publicIp.NetworkInterface != nil {
			networkInterfaceId := *publicIp.GetNetworkInterface()
			table.AddRow("ASSOCIATED TO", networkInterfaceId)
			table.AddSeparator()
		}

		if publicIp.Labels != nil && len(*publicIp.Labels) > 0 {
			labels := []string{}
			for key, value := range *publicIp.Labels {
				labels = append(labels, fmt.Sprintf("%s: %s", key, value))
			}
			table.AddRow("LABELS", strings.Join(labels, "\n"))
		}

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	}
}
