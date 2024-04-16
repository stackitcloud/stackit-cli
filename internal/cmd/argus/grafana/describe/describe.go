package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/argus/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/argus"
)

const (
	instanceIdFlag = "instance-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Shows details of the Grafana configuration of an Argus instance",
		Long:  "Shows details of the Grafana configuration of an Argus instance.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Get details of the Grafana configuration of an Argus instance with ID "xxx"`,
				"$ stackit argus credentials describe --instance-id xxx"),
			examples.NewExample(
				`Get details of the Grafana configuration of an Argus instance with ID "xxx" in a table format`,
				"$ stackit argus credentials describe --instance-id xxx --output-format pretty"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(cmd)
			if err != nil {
				return err
			}
			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			// Call API
			grafanaConfigsReq := buildGetGrafanaConfigRequest(ctx, model, apiClient)
			grafanaConfigsResp, err := grafanaConfigsReq.Execute()
			if err != nil {
				return fmt.Errorf("get Grafana configs: %w", err)
			}
			instanceReq := buildGetInstanceRequest(ctx, model, apiClient)
			instanceResp, err := instanceReq.Execute()
			if err != nil {
				return fmt.Errorf("get instance: %w", err)
			}

			return outputResult(p, model.OutputFormat, grafanaConfigsResp, instanceResp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "Instance ID")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)
}

func parseInput(cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringValue(cmd, instanceIdFlag),
	}, nil
}

func buildGetGrafanaConfigRequest(ctx context.Context, model *inputModel, apiClient *argus.APIClient) argus.ApiGetGrafanaConfigsRequest {
	req := apiClient.GetGrafanaConfigs(ctx, model.InstanceId, model.ProjectId)
	return req
}

func buildGetInstanceRequest(ctx context.Context, model *inputModel, apiClient *argus.APIClient) argus.ApiGetInstanceRequest {
	req := apiClient.GetInstance(ctx, model.InstanceId, model.ProjectId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, grafanaConfigs *argus.GrafanaConfigs, instance *argus.GetInstanceResponse) error {
	switch outputFormat {
	case globalflags.PrettyOutputFormat:

		table := tables.NewTable()
		table.AddRow("GRAFANA DASHBOARD", *instance.Instance.GrafanaUrl)
		table.AddSeparator()
		table.AddRow("PUBLIC READ ACCESS", *grafanaConfigs.PublicReadAccess)
		table.AddSeparator()
		table.AddRow("SINGLE SIGN-ON", *grafanaConfigs.UseStackitSso)
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	default:
		details, err := json.MarshalIndent(grafanaConfigs, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal Grafana configs: %w", err)
		}
		p.Outputln(string(details))

		return nil
	}
}
