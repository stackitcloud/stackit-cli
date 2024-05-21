package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/argus/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/argus"
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

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", instanceIdArg),
		Short: "Shows details of the Grafana configuration of an Argus instance",
		Long: fmt.Sprintf("%s\n%s\n%s",
			"Shows details of the Grafana configuration of an Argus instance.",
			`The Grafana dashboard URL and initial credentials (admin user and password) will be shown in the "pretty" output format. These credentials are only valid for first login. Please change the password after first login. After changing, the initial password is no longer valid.`,
			`The initial password is hidden by default, if you want to show it use the "--show-password" flag.`,
		),
		Args: args.SingleArg(instanceIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get details of the Grafana configuration of an Argus instance with ID "xxx"`,
				"$ stackit argus credentials describe xxx"),
			examples.NewExample(
				`Get details of the Grafana configuration of an Argus instance with ID "xxx" and show the initial admin password`,
				"$ stackit argus credentials describe xxx --show-password"),
			examples.NewExample(
				`Get details of the Grafana configuration of an Argus instance with ID "xxx" in JSON format`,
				"$ stackit argus credentials describe xxx --output-format json"),
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
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(grafanaConfigs, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal Grafana configs: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(grafanaConfigs, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal Grafana configs: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		initialAdminPassword := *instance.Instance.GrafanaAdminPassword
		if !inputModel.ShowPassword {
			initialAdminPassword = "<hidden>"
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
	}
}
