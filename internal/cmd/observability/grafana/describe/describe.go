package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/observability/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/observability"
)

const (
	instanceIdArg    = "INSTANCE_ID"
	showPasswordFlag = "show-password"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId   string
	ShowPassword bool
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", instanceIdArg),
		Short: "Shows details of the Grafana configuration of an Observability instance",
		Long: fmt.Sprintf("%s\n%s\n%s",
			"Shows details of the Grafana configuration of an Observability instance.",
			`The Grafana dashboard URL and initial credentials (admin user and password) will be shown in the "pretty" output format. These credentials are only valid for first login. Please change the password after first login. After changing, the initial password is no longer valid.`,
			`The initial password is hidden by default, if you want to show it use the "--show-password" flag.`,
		),
		Args: args.SingleArg(instanceIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get details of the Grafana configuration of an Observability instance with ID "xxx"`,
				"$ stackit observability grafana describe xxx"),
			examples.NewExample(
				`Get details of the Grafana configuration of an Observability instance with ID "xxx" and show the initial admin password`,
				"$ stackit observability grafana describe xxx --show-password"),
			examples.NewExample(
				`Get details of the Grafana configuration of an Observability instance with ID "xxx" in JSON format`,
				"$ stackit observability grafana describe xxx --output-format json"),
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

			return outputResult(params.Printer, model.OutputFormat, model.ShowPassword, grafanaConfigsResp, instanceResp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().BoolP(showPasswordFlag, "s", false, "Show password in output")
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	instanceId := inputArgs[0]
	globalFlags := globalflags.Parse(p, cmd)

	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      instanceId,
		ShowPassword:    flags.FlagToBoolValue(p, cmd, showPasswordFlag),
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

func buildGetGrafanaConfigRequest(ctx context.Context, model *inputModel, apiClient *observability.APIClient) observability.ApiGetGrafanaConfigsRequest {
	req := apiClient.GetGrafanaConfigs(ctx, model.InstanceId, model.ProjectId)
	return req
}

func buildGetInstanceRequest(ctx context.Context, model *inputModel, apiClient *observability.APIClient) observability.ApiGetInstanceRequest {
	req := apiClient.GetInstance(ctx, model.InstanceId, model.ProjectId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, showPassword bool, grafanaConfigs *observability.GrafanaConfigs, instance *observability.GetInstanceResponse) error {
	if instance == nil || instance.Instance == nil {
		return fmt.Errorf("instance or instance content is nil")
	} else if grafanaConfigs == nil {
		return fmt.Errorf("grafanaConfigs is nil")
	}

	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(grafanaConfigs, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal Grafana configs: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(grafanaConfigs, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal Grafana configs: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		initialAdminPassword := utils.PtrString(instance.Instance.GrafanaAdminPassword)
		if !showPassword {
			initialAdminPassword = "<hidden>"
		}

		table := tables.NewTable()
		table.AddRow("GRAFANA DASHBOARD", utils.PtrString(instance.Instance.GrafanaUrl))
		table.AddSeparator()
		table.AddRow("PUBLIC READ ACCESS", utils.PtrString(grafanaConfigs.PublicReadAccess))
		table.AddSeparator()
		table.AddRow("SINGLE SIGN-ON", utils.PtrString(grafanaConfigs.UseStackitSso))
		table.AddSeparator()
		table.AddRow("INITIAL ADMIN USER (DEFAULT)", utils.PtrString(instance.Instance.GrafanaAdminUser))
		table.AddSeparator()
		table.AddRow("INITIAL ADMIN PASSWORD (DEFAULT)", initialAdminPassword)
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
