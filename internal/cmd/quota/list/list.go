package list

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists quotas",
		Long:  "Lists project quotas.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List available quotas`,
				`$ stackit quota list`,
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
				return fmt.Errorf("list quotas: %w", err)
			}

			if items := response.Quotas; items == nil {
				params.Printer.Info("No quotas found for project %q", projectLabel)
			} else {
				if err := outputResult(params.Printer, model.OutputFormat, items); err != nil {
					return fmt.Errorf("output quotas: %w", err)
				}
			}

			return nil
		},
	}

	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiListQuotasRequest {
	request := apiClient.ListQuotas(ctx, model.ProjectId)

	return request
}

func outputResult(p *print.Printer, outputFormat string, quotas *iaas.QuotaList) error {
	if quotas == nil {
		return fmt.Errorf("quotas is nil")
	}
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(quotas, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal quota list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(quotas, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal quota list: %w", err)
		}
		p.Outputln(string(details))

		return nil

	default:
		table := tables.NewTable()
		table.SetHeader("NAME", "LIMIT", "CURRENT USAGE", "PERCENT")
		if val := quotas.BackupGigabytes; val != nil {
			table.AddRow("Total size in GiB of backups [GiB]", conv(val.Limit), conv(val.Usage), percentage(val))
		}
		if val := quotas.Backups; val != nil {
			table.AddRow("Number of backups [Count]", conv(val.Limit), conv(val.Usage), percentage(val))
		}
		if val := quotas.Gigabytes; val != nil {
			table.AddRow("Total size in GiB of volumes and snapshots [GiB]", conv(val.Limit), conv(val.Usage), percentage(val))
		}
		if val := quotas.Networks; val != nil {
			table.AddRow("Number of networks [Count]", conv(val.Limit), conv(val.Usage), percentage(val))
		}
		if val := quotas.Nics; val != nil {
			table.AddRow("Number of network interfaces (nics) [Count]", conv(val.Limit), conv(val.Usage), percentage(val))
		}
		if val := quotas.PublicIps; val != nil {
			table.AddRow("Number of public IP addresses [Count]", conv(val.Limit), conv(val.Usage), percentage(val))
		}
		if val := quotas.Ram; val != nil {
			table.AddRow("Amount of server RAM in MiB [MiB]", conv(val.Limit), conv(val.Usage), percentage(val))
		}
		if val := quotas.SecurityGroupRules; val != nil {
			table.AddRow("Number of security group rules [Count]", conv(val.Limit), conv(val.Usage), percentage(val))
		}
		if val := quotas.SecurityGroups; val != nil {
			table.AddRow("Number of security groups [Count]", conv(val.Limit), conv(val.Usage), percentage(val))
		}
		if val := quotas.Snapshots; val != nil {
			table.AddRow("Number of snapshots [Count]", conv(val.Limit), conv(val.Usage), percentage(val))
		}
		if val := quotas.Vcpu; val != nil {
			table.AddRow("Number of server cores (vcpu) [Count]", conv(val.Limit), conv(val.Usage), percentage(val))
		}
		if val := quotas.Volumes; val != nil {
			table.AddRow("Number of volumes [Count]", conv(val.Limit), conv(val.Usage), percentage(val))
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}

func conv(n *int64) string {
	if n != nil {
		return strconv.FormatInt(*n, 10)
	}
	return "n/a"
}

func percentage(val interface {
	GetLimitOk() (int64, bool)
	GetUsageOk() (int64, bool)
}) string {
	a, aOk := val.GetLimitOk()
	b, bOk := val.GetUsageOk()
	if aOk && bOk {
		return fmt.Sprintf("%3.1f%%", 100.0/float64(a)*float64(b))
	}
	return "n/a"
}
