package describe

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	AffinityGroupId string
}

const (
	affinityGroupId = "AFFINITY_GROUP_ID"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", affinityGroupId),
		Short: "Show details of an affinity group",
		Long:  `Show details of an affinity group.`,
		Args:  args.SingleArg(affinityGroupId, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get details about an affinity group with the ID "xxx"`,
				"$ stackit affinity-group describe xxx",
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
			request := buildRequest(ctx, *model, apiClient)
			result, err := request.Execute()
			if err != nil {
				return fmt.Errorf("get affinity group: %w", err)
			}

			if err := outputResult(params.Printer, *model, *result); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func buildRequest(ctx context.Context, model inputModel, apiClient *iaas.APIClient) iaas.ApiGetAffinityGroupRequest {
	return apiClient.GetAffinityGroup(ctx, model.ProjectId, model.Region, model.AffinityGroupId)
}

func parseInput(p *print.Printer, cmd *cobra.Command, cliArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		AffinityGroupId: cliArgs[0],
	}

	p.DebugInputModel(model)
	return &model, nil
}

func outputResult(p *print.Printer, model inputModel, resp iaas.AffinityGroup) error {
	var outputFormat string
	if model.GlobalFlagModel != nil {
		outputFormat = model.OutputFormat
	}

	return p.OutputResult(outputFormat, resp, func() error {
		table := tables.NewTable()

		if resp.HasId() {
			table.AddRow("ID", utils.PtrString(resp.Id))
			table.AddSeparator()
		}
		if resp.Name != nil {
			table.AddRow("NAME", utils.PtrString(resp.Name))
			table.AddSeparator()
		}
		if resp.Policy != nil {
			table.AddRow("POLICY", utils.PtrString(resp.Policy))
			table.AddSeparator()
		}
		if resp.HasMembers() {
			table.AddRow("Members", utils.JoinStringPtr(resp.Members, ", "))
			table.AddSeparator()
		}

		if err := table.Display(p); err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
