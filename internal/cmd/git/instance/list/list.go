package list

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/git/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/git"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit *int64
}

const limitFlag = "limit"

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all instances of STACKIT Git.",
		Long:  "Lists all instances of STACKIT Git for the current project.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all STACKIT Git instances`,
				"$ stackit git instance list"),
			examples.NewExample(
				"Lists up to 10 STACKIT Git instances",
				"$ stackit git instance list --limit=10",
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("get STACKIT Git instances: %w", err)
			}
			instances := resp.GetInstances()

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			// Truncate output
			if model.Limit != nil && len(instances) > int(*model.Limit) {
				instances = (instances)[:*model.Limit]
			}

			return outputResult(params.Printer, model.OutputFormat, projectLabel, instances)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Limit the output to the first n elements")
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
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *git.APIClient) git.ApiListInstancesRequest {
	return apiClient.ListInstances(ctx, model.ProjectId)
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, instances []git.Instance) error {
	return p.OutputResult(outputFormat, instances, func() error {
		if len(instances) == 0 {
			p.Outputf("No instances found for project %q\n", projectLabel)
			return nil
		}

		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "URL", "VERSION", "STATE", "CREATED")
		for i := range instances {
			instance := (instances)[i]
			table.AddRow(
				utils.PtrString(instance.Id),
				utils.PtrString(instance.Name),
				utils.PtrString(instance.Url),
				utils.PtrString(instance.Version),
				utils.PtrString(instance.State),
				utils.PtrString(instance.Created),
			)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
