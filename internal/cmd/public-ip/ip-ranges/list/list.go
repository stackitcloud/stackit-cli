package list

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all STACKIT cloud public IP ranges.",
		Long:  "Lists all STACKIT cloud public IP ranges.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Lists all STACKIT cloud public IP ranges`,
				"$ stackit public-ip ranges list",
			),
			examples.NewExample(
				`Lists all STACKIT cloud public IP ranges, piping to a tool like fzf for interactive selection`,
				"$ stackit public-ip ip-ranges list -o pretty | fzf",
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
				return fmt.Errorf("list public ip ranges s: %w", err)
			}

			if resp.Items == nil || len(*resp.Items) == 0 {
				params.Printer.Info("No public IP ranges found\n")
				return nil
			}

			return outputResult(params.Printer, model.OutputFormat, *resp)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{GlobalFlagModel: globalFlags}

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

func outputResult(p *print.Printer, outputFormat string, networkListResponse iaas.PublicNetworkListResponse) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(networkListResponse, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal public IP: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(networkListResponse, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal public IP: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		var publicIps []string
		for _, item := range *networkListResponse.Items {
			if item.Cidr != nil || *item.Cidr != "" {
				publicIps = append(publicIps, *item.Cidr)
			}
		}
		p.Outputln(strings.Join(publicIps, "\n"))

		return nil
	}
}
