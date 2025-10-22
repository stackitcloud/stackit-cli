package describe

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/resourcemanager/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/resourcemanager"
)

const (
	includeParentsFlag = "include-parents"

	projectIdArg = "PROJECT_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ArgProjectId   string
	IncludeParents bool
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Shows details of a STACKIT project",
		Long:  "Shows details of a STACKIT project.",
		Args:  args.SingleOptionalArg(projectIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get the details of the configured STACKIT project`,
				"$ stackit project describe"),
			examples.NewExample(
				`Get the details of a STACKIT project by explicitly providing the project ID`,
				"$ stackit project describe xxx"),
			examples.NewExample(
				`Get the details of the configured STACKIT project, including details of the parent resources`,
				"$ stackit project describe --include-parents"),
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
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read project details: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Bool(includeParentsFlag, false, "When true, the details of the parent resources will be included in the output")
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	var projectId string
	if len(inputArgs) > 0 {
		projectId = inputArgs[0]
	}

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" && projectId == "" {
		return nil, fmt.Errorf("Project ID needs to be provided either as an argument or as a flag")
	}

	if projectId == "" {
		projectId = globalFlags.ProjectId
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ArgProjectId:    projectId,
		IncludeParents:  flags.FlagToBoolValue(p, cmd, includeParentsFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *resourcemanager.APIClient) resourcemanager.ApiGetProjectRequest {
	req := apiClient.GetProject(ctx, model.ArgProjectId)
	req.IncludeParents(model.IncludeParents)
	return req
}

func outputResult(p *print.Printer, outputFormat string, project *resourcemanager.GetProjectResponse) error {
	if project == nil {
		return fmt.Errorf("response not set")
	}

	return p.OutputResult(outputFormat, project, func() error {
		table := tables.NewTable()
		table.AddRow("ID", utils.PtrString(project.ProjectId))
		table.AddSeparator()
		table.AddRow("NAME", utils.PtrString(project.Name))
		table.AddSeparator()
		table.AddRow("CREATION", utils.PtrString(project.CreationTime))
		table.AddSeparator()
		table.AddRow("STATE", utils.PtrString(project.LifecycleState))
		table.AddSeparator()
		if project.Parent != nil {
			table.AddRow("PARENT ID", utils.PtrString(project.Parent.Id))
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
