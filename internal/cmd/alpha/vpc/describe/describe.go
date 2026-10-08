package describe

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	iaas "github.com/stackitcloud/stackit-sdk-go/services/iaas/v2alpha1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaasalpha/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	vpcIdArg = "VPC_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	VpcId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", vpcIdArg),
		Short: "Shows details of a VPC",
		Long:  "Shows details of a VPC.",
		Args:  args.SingleArg(vpcIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Show details of a vpc with ID "xxx"`,
				"$ stackit alpha vpc describe xxx"),
			examples.NewExample(
				`Show details of a vpc with ID "xxx" in JSON format`,
				"$ stackit alpha vpc describe xxx --output-format json"),
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
				return fmt.Errorf("read vpc: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	vpcId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		VpcId:           vpcId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiGetVPCRequest {
	return apiClient.DefaultAPI.GetVPC(ctx, model.ProjectId, model.VpcId)
}

func outputResult(p *print.Printer, outputFormat string, resp *iaas.VPC) error {
	if resp == nil {
		return fmt.Errorf("nil response")
	}
	return p.OutputResult(outputFormat, resp, func() error {
		table := tables.NewTable()
		table.AddRow("ID", resp.Id)
		table.AddSeparator()
		table.AddRow("NAME", resp.Name)
		table.AddSeparator()
		table.AddRow("DESCRIPTION", resp.Description)
		table.AddSeparator()
		if len(resp.Labels) > 0 {
			var labels []string
			for key, value := range resp.Labels {
				labels = append(labels, fmt.Sprintf("%s: %s", key, value))
			}
			table.AddRow("LABELS", strings.Join(labels, "\n"))
			table.AddSeparator()
		}
		table.AddRow("CREATED AT", resp.CreatedAt.Format(time.DateTime))
		table.AddSeparator()
		table.AddRow("UPDATED AT", resp.UpdatedAt.Format(time.DateTime))
		table.AddSeparator()

		if err := table.Display(p); err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
