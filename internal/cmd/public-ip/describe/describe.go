package describe

import (
	"context"
	"fmt"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

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

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", publicIpIdArg),
		Short: "Shows details of a Public IP",
		Long:  "Shows details of a Public IP.",
		Args:  args.SingleArg(publicIpIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Show details of a public IP with ID "xxx"`,
				"$ stackit public-ip describe xxx",
			),
			examples.NewExample(
				`Show details of a public IP with ID "xxx" in JSON format`,
				"$ stackit public-ip describe xxx --output-format json",
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
				return fmt.Errorf("read public IP: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, *resp)
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

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiGetPublicIPRequest {
	return apiClient.GetPublicIP(ctx, model.ProjectId, model.Region, model.PublicIpId)
}

func outputResult(p *print.Printer, outputFormat string, publicIp iaas.PublicIp) error {
	return p.OutputResult(outputFormat, publicIp, func() error {
		table := tables.NewTable()
		table.AddRow("ID", utils.PtrString(publicIp.Id))
		table.AddSeparator()
		table.AddRow("IP ADDRESS", utils.PtrString(publicIp.Ip))
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
	})
}
