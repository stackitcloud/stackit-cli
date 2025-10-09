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

type inputModel struct {
	*globalflags.GlobalFlagModel
	LabelSelector *string
	Limit         *int64
}

const (
	labelSelectorFlag = "label-selector"
	limitFlag         = "limit"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists images",
		Long:  "Lists images by their internal ID.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all images`,
				`$ stackit image list`,
			),
			examples.NewExample(
				`List images with label`,
				`$ stackit image list --label-selector ARM64,dev`,
			),
			examples.NewExample(
				`List the first 10 images`,
				`$ stackit image list --limit=10`,
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

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			} else if projectLabel == "" {
				projectLabel = model.ProjectId
			}

			// Call API
			request := buildRequest(ctx, model, apiClient)

			response, err := request.Execute()
			if err != nil {
				return fmt.Errorf("list images: %w", err)
			}

			if items := response.GetItems(); len(items) == 0 {
				params.Printer.Info("No images found for project %q", projectLabel)
			} else {
				if model.Limit != nil && len(items) > int(*model.Limit) {
					items = (items)[:*model.Limit]
				}
				if err := outputResult(params.Printer, model.OutputFormat, items); err != nil {
					return fmt.Errorf("output images: %w", err)
				}
			}

			return nil
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(labelSelectorFlag, "", "Filter by label")
	cmd.Flags().Int64(limitFlag, 0, "Limit the output to the first n elements")
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
		LabelSelector:   flags.FlagToStringPointer(p, cmd, labelSelectorFlag),
		Limit:           limit,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiListImagesRequest {
	request := apiClient.ListImages(ctx, model.ProjectId)
	if model.LabelSelector != nil {
		request = request.LabelSelector(*model.LabelSelector)
	}

	return request
}
func outputResult(p *print.Printer, outputFormat string, items []iaas.Image) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(items, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal image list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(items, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal image list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "OS", "ARCHITECTURE", "DISTRIBUTION", "VERSION", "LABELS")
		for i := range items {
			item := items[i]
			var (
				architecture string = "n/a"
				os           string = "n/a"
				distro       string = "n/a"
				version      string = "n/a"
			)
			if cfg := item.Config; cfg != nil {
				if v := cfg.Architecture; v != nil {
					architecture = *v
				}
				if v := cfg.OperatingSystem; v != nil {
					os = *v
				}
				if v := cfg.OperatingSystemDistro; v != nil && v.IsSet() {
					distro = *v.Get()
				}
				if v := cfg.OperatingSystemVersion; v != nil && v.IsSet() {
					version = *v.Get()
				}
			}
			table.AddRow(utils.PtrString(item.Id),
				utils.PtrString(item.Name),
				os,
				architecture,
				distro,
				version,
				utils.JoinStringKeysPtr(*item.Labels, ","))
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
