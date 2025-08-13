package list

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
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
		Short: "Lists all volumes of a project",
		Long:  "Lists all volumes of a project.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Lists all volumes`,
				"$ stackit volume list",
			),
			examples.NewExample(
				`Lists all volumes which contains the label xxx`,
				"$ stackit volume list --label-selector xxx",
			),
			examples.NewExample(
				`Lists all volumes in JSON format`,
				"$ stackit volume list --output-format json",
			),
			examples.NewExample(
				`Lists up to 10 volumes`,
				"$ stackit volume list --limit 10",
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
				return fmt.Errorf("list volumes: %w", err)
			}

			if resp.Items == nil || len(*resp.Items) == 0 {
				projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
					projectLabel = model.ProjectId
				}
				params.Printer.Info("No volumes found for project %q\n", projectLabel)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiListVolumesRequest {
	req := apiClient.ListVolumes(ctx, model.ProjectId)
	if model.LabelSelector != nil {
		req = req.LabelSelector(*model.LabelSelector)
	}

	return req
}

func outputResult(p *print.Printer, outputFormat string, volumes []iaas.Volume) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(volumes, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal volume: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(volumes, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal volume: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("ID", "Name", "Status", "Server", "Availability Zone", "Size (GB)")

		for i := range volumes {
			volume := volumes[i]
			table.AddRow(
				utils.PtrString(volume.Id),
				utils.PtrString(volume.Name),
				utils.PtrString(volume.Status),
				utils.PtrString(volume.ServerId),
				utils.PtrString(volume.AvailabilityZone),
				utils.PtrString(volume.Size),
			)
			table.AddSeparator()
		}

		p.Outputln(table.Render())
		return nil
	}
}
