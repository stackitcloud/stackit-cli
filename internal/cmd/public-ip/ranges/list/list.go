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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	limitFlag = "limit"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit *int64
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all STACKIT public-ip ranges",
		Long:  "Lists all STACKIT public-ip ranges.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Lists all STACKIT public-ip ranges`,
				"$ stackit public-ip ranges list",
			),
			examples.NewExample(
				`Lists all STACKIT public-ip ranges, piping to a tool like fzf for interactive selection`,
				"$ stackit public-ip ranges list -o pretty | fzf",
			),
			examples.NewExample(
				`Lists up to 10 STACKIT public-ip ranges`,
				"$ stackit public-ip ranges list --limit 10",
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
			req := apiClient.ListPublicIPRanges(ctx)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("list public IP ranges: %w", err)
			}
			publicIpRanges := *resp.Items

			// Truncate output
			if model.Limit != nil && len(publicIpRanges) > int(*model.Limit) {
				publicIpRanges = publicIpRanges[:*model.Limit]
			}

			return outputResult(params.Printer, model.OutputFormat, publicIpRanges)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

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

func outputResult(p *print.Printer, outputFormat string, publicIpRanges []iaas.PublicNetwork) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(publicIpRanges, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal public IP ranges: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(publicIpRanges, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal public IP ranges: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		if len(publicIpRanges) == 0 {
			p.Outputln("No public IP ranges found")
			return nil
		}

		for _, item := range publicIpRanges {
			if item.Cidr != nil && *item.Cidr != "" {
				p.Outputln(*item.Cidr)
			}
		}

		return nil
	}
}
