package list

import (
	"context"
	"fmt"

	ufw "github.com/stackitcloud/stackit-sdk-go/services/ufw/v1api"

	serviceEnablementClient "github.com/stackitcloud/stackit-cli/internal/pkg/services/service-enablement/client"
	serviceEnablementUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/service-enablement/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/ufw/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
)

const (
	limitFlag = "limit"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit *int64
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all UFW rules",
		Long:  "Lists all STACKIT Unified Firewall (UFW) rules.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all UFW rules`,
				"$ stackit ufw rules list"),
			examples.NewExample(
				`List all UFW rules in JSON format`,
				"$ stackit ufw rules list --output-format json"),
			examples.NewExample(
				`List up to 10 UFW rules`,
				"$ stackit ufw rules list --limit 10"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			serviceEnablementApiClient, err := serviceEnablementClient.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				enabled, enabledErr := serviceEnablementUtils.ProjectEnabled(ctx, serviceEnablementApiClient.DefaultAPI, model.ProjectId, model.Region)
				if enabledErr != nil {
					return fmt.Errorf("check if project is enabled failed: %w", enabledErr)
				}
				if !enabled {
					return &errors.ServiceDisabledError{
						Service: "ufw",
					}
				}
				return fmt.Errorf("get UFW rules: %w", err)
			}
			rules := resp.Rules

			if model.Limit != nil && len(rules) > int(*model.Limit) {
				rules = rules[:*model.Limit]
			}

			projectLabel := model.ProjectId
			if len(rules) == 0 {
				projectLabel, err = projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				}
			}

			if resp.Rules == nil {
				params.Printer.Info("(...) %s", projectLabel)
				return nil
			}

			return outputResult(params.Printer, model.OutputFormat, projectLabel, rules)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	if globalFlags.Region == "" {
		return nil, &errors.RegionError{}
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
		Limit:           flags.FlagToInt64Pointer(p, cmd, limitFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *ufw.APIClient) ufw.ApiListRulesRequest {
	req := apiClient.DefaultAPI.ListRules(ctx, model.ProjectId, model.Region)
	return req
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, resources []ufw.RuleResponse) error {
	return p.OutputResult(outputFormat, resources, func() error {
		if len(resources) == 0 {
			p.Outputf("No rules found for project %q\n", projectLabel)
			return nil
		}

		table := tables.NewTable()
		table.SetHeader("PRODUCT", "SOURCE", "DEPLOYMENT TARGET", "PROTOCOL",
			"DIRECTION", "PORT RANGE", "ETHER TYPE", "STATUS")
		for i := range resources {
			resource := resources[i]
			table.AddRow(resource.Product, resource.SourceIP, resource.InstanceName, resource.Protocol, resource.Direction,
				resource.PortRange, resource.EtherType, resource.Status)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
