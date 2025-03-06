package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
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

func NewCmd(p *print.Printer) *cobra.Command {
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
			request := buildRequest(ctx, *model, apiClient)
			result, err := request.Execute()
			if err != nil {
				return fmt.Errorf("get affinity group: %w", err)
			}

			if err := outputResult(p, *model, *result); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func buildRequest(ctx context.Context, model inputModel, apiClient *iaas.APIClient) iaas.ApiGetAffinityGroupRequest {
	return apiClient.GetAffinityGroup(ctx, model.ProjectId, model.AffinityGroupId)
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

func outputResult(p *print.Printer, model inputModel, resp iaas.AffinityGroup) error {
	var outputFormat string
	if model.GlobalFlagModel != nil {
		outputFormat = model.GlobalFlagModel.OutputFormat
	}
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal affinity group: %w", err)
		}
		p.Outputln(string(details))
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(resp, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal affinity group: %w", err)
		}
		p.Outputln(string(details))
	default:
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
	}
	return nil
}
