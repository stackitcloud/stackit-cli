package list

import (
	"context"
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
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"

	"github.com/spf13/cobra"
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
		Short: "Lists all networks of a project",
		Long:  "Lists all network of a project.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Lists all networks`,
				"$ stackit network list",
			),
			examples.NewExample(
				`Lists all networks in JSON format`,
				"$ stackit network list --output-format json",
			),
			examples.NewExample(
				`Lists up to 10 networks`,
				"$ stackit network list --limit 10",
			),
			examples.NewExample(
				`Lists all networks which contains the label xxx`,
				"$ stackit network list --label-selector xxx",
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
				return fmt.Errorf("list networks: %w", err)
			}

			if resp.Items == nil || len(*resp.Items) == 0 {
				projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
					projectLabel = model.ProjectId
				} else if projectLabel == "" {
					projectLabel = model.ProjectId
				}
				params.Printer.Info("No networks found for project %q\n", projectLabel)
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

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiListNetworksRequest {
	req := apiClient.ListNetworks(ctx, model.ProjectId)
	if model.LabelSelector != nil {
		req = req.LabelSelector(*model.LabelSelector)
	}
	return req
}

func outputResult(p *print.Printer, outputFormat string, networks []iaas.Network) error {
	return p.OutputResult(outputFormat, networks, func() error {
		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "STATUS", "PUBLIC IP", "PREFIXES", "ROUTED")

		for _, network := range networks {
			publicIp := utils.PtrString(network.PublicIp)

			routed := false
			if network.Routed != nil {
				routed = *network.Routed
			}
			prefixes := utils.JoinStringPtr(network.Prefixes, ", ")

			table.AddRow(
				utils.PtrString(network.NetworkId),
				utils.PtrString(network.Name),
				utils.PtrString(network.State),
				publicIp,
				prefixes,
				routed,
			)
			table.AddSeparator()
		}

		p.Outputln(table.Render())
		return nil
	})
}
