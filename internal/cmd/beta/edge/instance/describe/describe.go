package describe

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	edge "github.com/stackitcloud/stackit-sdk-go/services/edge/v1beta1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	instanceIdArg = "INSTANCE_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId string
}

// NewCmd https://aip.stackit.cloud/aip/general/0121/
// We have decided to eliminate the usage of display name flag
// To be the AIP compliant, and align with the standard CLI implementation, we will use the InstanceID arg
func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", instanceIdArg),
		Short: "Describes an Edge Cloud instance",
		Long:  "Describes a STACKIT Edge Cloud (STEC) instance.",
		Args:  args.SingleArg(instanceIdArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Describe an Edge Cloud instance with ID "xxx"`,
				`$ stackit beta edge-cloud instance describe xxx`),
			examples.NewExample(
				`Get details of an Edge Cloud instance with ID "xxx" in JSON format`,
				"$ stackit beta edge-cloud instance describe xxx --output-format json"),
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
				return fmt.Errorf("read Edge Cloud instance: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	instanceId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      instanceId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *edge.APIClient) edge.ApiGetInstanceRequest {
	return apiClient.DefaultAPI.GetInstance(ctx, model.ProjectId, model.Region, model.InstanceId)
}

func outputResult(p *print.Printer, outputFormat string, instance *edge.Instance) error {
	return p.OutputResult(outputFormat, instance, func() error {
		if instance == nil {
			return fmt.Errorf("instance response is empty")
		}

		table := tables.NewTable()
		table.AddRow("ID", instance.Id)
		table.AddSeparator()
		table.AddRow("NAME", instance.DisplayName)
		table.AddSeparator()
		if instance.HasDescription() {
			table.AddRow("DESCRIPTION", utils.PtrString(instance.Description))
			table.AddSeparator()
		}
		table.AddRow("CREATED", instance.Created)
		table.AddSeparator()
		table.AddRow("UI", instance.FrontendUrl)
		table.AddSeparator()
		table.AddRow("STATE", instance.Status)
		table.AddSeparator()
		table.AddRow("PLAN", instance.PlanId)
		table.AddSeparator()

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
