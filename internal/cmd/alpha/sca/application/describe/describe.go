package describe

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	sca "github.com/stackitcloud/stackit-sdk-go/services/sca/v1alphaapi"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sca/client"
	scautils "github.com/stackitcloud/stackit-cli/internal/pkg/services/sca/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	applicationIDArg  = "APPLICATION_ID"
	environmentIDFlag = "environment-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	EnvironmentID string
	ApplicationID string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Show details of a SCA application",
		Long:  "Show details of a STACKIT Container Applications (SCA) application.",
		Args:  args.SingleArg(applicationIDArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Get details of a SCA application with ID "xxx" from an environment with ID "yyy"`,
				"$ stackit alpha sca application describe xxx --environment-id yyy"),
			examples.NewExample(
				`Get details of all SCA application with ID "xxx" from an environment with ID "yyy" in JSON format`,
				"$ stackit alpha sca application describe xxx --environment-id yyy --output-format json"),
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

			globalFlags := globalflags.Parse(params.Printer, cmd)
			if globalFlags.ProjectId == "" {
				return &errors.ProjectIdError{}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("describe SCA application: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}

	configureFlags(cmd)
	return cmd
}

func outputResult(p *print.Printer, outputFormat string, application *sca.Application) error {
	return p.OutputResult(outputFormat, application, func() error {
		if application == nil {
			p.Outputf("No application found")
			return nil
		}

		table := tables.NewTable()
		table.SetTitle("Application")
		table.AddRow("ID", utils.PtrString(application.Id))
		table.AddSeparator()
		table.AddRow("NAME", application.DisplayName)
		table.AddSeparator()
		table.AddRow("STATUS", scautils.ApplicationStatusToStr(application.GetRuntimeStatus().CurrentStatus))
		table.AddSeparator()
		table.AddRow("STATE", scautils.ApplicationStateToStr(application.GetStopped()))
		table.AddSeparator()
		if application.RuntimeStatus != nil {
			table.AddRow("AVAILABLE INSTANCES", len(application.GetRuntimeStatus().Instances))
			table.AddSeparator()
		}
		table.AddRow("PUBLIC INGRESS", application.GetNetwork().PublicIngress)
		if application.Network.Port != nil {
			table.AddRow("PUBLIC PORT", *application.GetNetwork().Port)
		}
		if len(application.Network.IngressAcl) > 0 {
			table.AddRow("ACLs", strings.Join(application.Network.IngressAcl, "\n"))
		}
		if len(application.GetRuntimeStatus().Urls) > 0 {
			table.AddRow("PUBLIC URL", application.RuntimeStatus.Urls[0])
		}
		table.AddSeparator()

		tablesToDisplay := []tables.Table{table}

		switch application.Scaling.Type {
		case sca.SCALINGTYPE_SCALING_TYPE_MANUAL:
			table.AddRow("SCALING TYPE", "MANUAL")
			table.AddRow("INSTANCES", application.GetScaling().ManualScaling.Instances)
		case sca.SCALINGTYPE_SCALING_TYPE_AUTO:
			table.AddRow("SCALING TYPE", "AUTO")
			table.AddRow("SCALE TO ZERO", application.GetScaling().AutoScaling.GetAllowScaleToZero())
			table.AddRow("MIN INSTANCES", application.GetScaling().AutoScaling.MinInstances)
			table.AddRow("MAX INSTANCES", application.GetScaling().AutoScaling.MaxInstances)
			rules := application.GetScaling().AutoScaling.Rules
			if len(rules) > 0 {
				tablesToDisplay = append(tablesToDisplay, buildScalingRulesTable(rules))
			}
		}

		containersTable := tables.NewTable()
		containersTable.SetTitle("Application Containers")
		containersTable.SetHeader("NAME", "IMAGE", "CPU", "MEMORY", "COMMANDS", "ARGS")

		environmentVarTables := []tables.Table{}
		for _, c := range application.Containers {
			envVarTable := buildEnvVarsTable(&c)
			if envVarTable != nil {
				environmentVarTables = append(environmentVarTables, *envVarTable)
			}
			commands := "-"
			if len(c.Command) > 0 {
				commands = strings.Join(c.Command, " ")
			}
			cArgs := "-"
			if len(c.Args) > 0 {
				cArgs = strings.Join(c.Args, " ")
			}

			cpu := "-"
			if c.Cpu != nil {
				cpu = strconv.Itoa(int(*c.Cpu))
			}

			memory := "-"
			if c.Memory != nil {
				memory = strconv.Itoa(int(*c.Memory))
			}

			containersTable.AddRow(
				c.Name,
				c.Image,
				cpu,
				memory,
				commands,
				cArgs,
			)
		}
		tablesToDisplay = append(tablesToDisplay, containersTable)
		err := tables.DisplayTables(p, append(tablesToDisplay, environmentVarTables...))
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), environmentIDFlag, "Environment ID (if not set uses default environment)")
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	applicationID := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	environmentID := flags.FlagToStringValue(p, cmd, environmentIDFlag)
	if environmentID == "" {
		environmentID = globalFlags.ProjectId
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		EnvironmentID:   environmentID,
		ApplicationID:   applicationID,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *sca.APIClient) sca.ApiGetApplicationRequest {
	return apiClient.DefaultAPI.GetApplication(ctx, model.ProjectId, model.EnvironmentID, model.ApplicationID)
}

func buildScalingRulesTable(rules []sca.ScaleRule) tables.Table {
	table := tables.NewTable()
	table.SetTitle("Autoscaling rules")
	table.SetHeader("NAME", "TYPE", "PARAMETERS", "SECRETS MAPPING")
	for _, rule := range rules {
		ruleType := "-"
		parameters := "-"
		secretsMapping := "-"
		switch rule.Type {
		case sca.RULETYPE_RULE_TYPE_HTTP:
			ruleType = "HTTP"
			params := []string{}
			if rule.GetHttpRule().Concurrency != nil && *rule.GetHttpRule().Concurrency > 0 {
				params = append(params, fmt.Sprintf("Concurrency: %d", *rule.GetHttpRule().Concurrency))
			}
			if rule.GetHttpRule().Rps != nil && *rule.GetHttpRule().Rps > 0 {
				params = append(params, fmt.Sprintf("RPS: %d", *rule.GetHttpRule().Rps))
			}
			parameters = strings.Join(params, "\n")
			secretsMapping = "-"
		case sca.RULETYPE_RULE_TYPE_CUSTOM:
			ruleType = "CUSTOM"
			params := []string{}
			if rule.GetCustomRule().Type != "" {
				params = append(params, fmt.Sprintf("TYPE: %s", rule.GetCustomRule().Type))
			}
			for _, p := range rule.GetCustomRule().Parameters {
				params = append(params, fmt.Sprintf("%s: %s", p.Name, p.Value))
			}
			parameters = strings.Join(params, "\n")
			sm := make([]string, 0, len(rule.GetCustomRule().SecretsMapping))
			for _, s := range rule.GetCustomRule().SecretsMapping {
				sm = append(sm, fmt.Sprintf("%s: %s", s.Parameter, s.Secret))
			}
			if len(sm) > 0 {
				secretsMapping = strings.Join(sm, "\n")
			}
		}

		table.AddRow(
			rule.Name,
			ruleType,
			parameters,
			secretsMapping,
		)
		table.AddSeparator()
	}

	return table
}

func buildEnvVarsTable(container *sca.Container) *tables.Table {
	if len(container.EnvironmentVariables) == 0 {
		return nil
	}
	table := tables.NewTable()
	table.SetTitle(fmt.Sprintf("%s environment variables", container.Name))
	table.SetHeader("NAME", "ORIGIN", "VALUE / SECRET REF")
	for _, e := range container.EnvironmentVariables {
		origin := "MANUAL"
		if e.Origin != nil && *e.Origin == sca.ENVVARTYPE_ENV_FROM_SOURCE_TYPE_SECRET {
			origin = "SECRET"
		}
		table.AddRow(
			e.Key,
			origin,
			e.Value,
		)
		table.AddSeparator()
	}
	return &table
}
