package list

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/observability/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	observabilityUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/observability/utils"

	"github.com/spf13/cobra"
	observability "github.com/stackitcloud/stackit-sdk-go/services/observability/v1api"
)

const (
	limitFlag      = "limit"
	instanceIdFlag = "instance-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit      *int64
	InstanceId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all ACLs of an Observability instance",
		Long:  "Lists all IP networks in CIDR notation which are allowed to access an Observability instance.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all ACLs of Observability instance "xxx"`,
				"$ stackit observability acl list --instance-id xxx"),
			examples.NewExample(
				`List all ACLs of Observability instance "xxx" in JSON format`,
				"$ stackit observability acl list --instance-id xxx --output-format json"),
			examples.NewExample(
				`List up to 10 ACLs of Observability instance "xxx"`,
				"$ stackit observability acl list --instance-id xxx --limit 10"),
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
			req := buildRequest(ctx, model, apiClient.DefaultAPI)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("get Observability instance ACLs: %w", err)
			}

			acls := resp.GetAcl()

			instanceLabel, err := observabilityUtils.GetInstanceName(ctx, apiClient.DefaultAPI, model.InstanceId, model.ProjectId)
			if err != nil || instanceLabel == "" {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceId
			}

			// Truncate output
			if model.Limit != nil && len(acls) > int(*model.Limit) {
				acls = acls[:*model.Limit]
			}

			return outputResult(params.Printer, model.OutputFormat, instanceLabel, acls)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "Instance ID")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
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
		Limit:           limit,
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient observability.DefaultAPI) observability.ApiListACLRequest {
	req := apiClient.ListACL(ctx, model.InstanceId, model.ProjectId)
	return req
}

func outputResult(p *print.Printer, outputFormat, instanceLabel string, acls []string) error {
	return p.OutputResult(outputFormat, acls, func() error {
		if len(acls) == 0 {
			p.Outputf("No ACLs found for instance %q\n", instanceLabel)
			return nil
		}
		table := tables.NewTable()
		table.SetHeader("CIDR")
		for _, acl := range acls {
			table.AddRow(acl)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
