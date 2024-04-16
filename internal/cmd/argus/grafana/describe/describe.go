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
	instanceIdFlag   = "instance-id"
	showPasswordFlag = "show-password"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId   string
	ShowPassword bool
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Shows details of the Grafana configuration of an Argus instance",
		Long: fmt.Sprintf("%s\n%s\n%s",
			"Shows details of the Grafana configuration of an Argus instance.",
			`The initial admin user and password will be shown in the "pretty" output format. These credentials are only valid for first login. Please change the password after first login. After changing, the initial password is no longer valid.`,
			`The initial password is hidden by default, if you want to see it use the "--show-password" flag.`,
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Get details of the Grafana configuration of an Argus instance with ID "xxx"`,
				"$ stackit argus credentials describe --instance-id xxx"),
			examples.NewExample(
				`Get details of the Grafana configuration of an Argus instance with ID "xxx" in a table format`,
				"$ stackit argus credentials describe --instance-id xxx --output-format pretty"),
			examples.NewExample(
				`Get details of the Grafana configuration of an Argus instance with ID "xxx" and show the initial admin password`,
				"$ stackit argus credentials describe --instance-id xxx --output-format pretty --show-password"),
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

			return outputResult(p, model, grafanaConfigsResp, instanceResp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "Instance ID")
	cmd.Flags().Bool(showPasswordFlag, false, `Show the initial admin password in the "pretty" output format`)

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
		ShowPassword:    flags.FlagToBoolValue(cmd, showPasswordFlag),
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

func outputResult(p *print.Printer, inputModel *inputModel, grafanaConfigs *argus.GrafanaConfigs, instance *argus.GetInstanceResponse) error {
	switch inputModel.OutputFormat {
	case globalflags.PrettyOutputFormat:
		initialAdminPassword := "<hidden>"
		if inputModel.ShowPassword {
			initialAdminPassword = *instance.Instance.GrafanaAdminPassword
		}

		table := tables.NewTable()
		table.AddRow("GRAFANA DASHBOARD", *instance.Instance.GrafanaUrl)
		table.AddSeparator()
		table.AddRow("PUBLIC READ ACCESS", *grafanaConfigs.PublicReadAccess)
		table.AddSeparator()
		table.AddRow("SINGLE SIGN-ON", *grafanaConfigs.UseStackitSso)
		table.AddSeparator()
		table.AddRow("INITIAL ADMIN USER (DEFAULT)", *instance.Instance.GrafanaAdminUser)
		table.AddSeparator()
		table.AddRow("INITIAL ADMIN PASSWORD (DEFAULT)", initialAdminPassword)
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
