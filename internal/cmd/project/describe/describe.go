package describe

import (
	"context"
	"encoding/json"
	"fmt"

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

func NewCmd(p *print.Printer) *cobra.Command {
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
			model, err := parseInput(cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(cmd)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read project details: %w", err)
			}

			return outputResult(cmd, model.OutputFormat, resp, p)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Bool(includeParentsFlag, false, "When true, the details of the parent resources will be included in the output")
}

func parseInput(cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	var projectId string
	if len(inputArgs) > 0 {
		projectId = inputArgs[0]
	}

	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" && projectId == "" {
		return nil, fmt.Errorf("Project ID needs to be provided either as an argument or as a flag")
	}

	if projectId == "" {
		projectId = globalFlags.ProjectId
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		ArgProjectId:    projectId,
		IncludeParents:  flags.FlagToBoolValue(cmd, includeParentsFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *resourcemanager.APIClient) resourcemanager.ApiGetProjectRequest {
	req := apiClient.GetProject(ctx, model.ArgProjectId)
	req.IncludeParents(model.IncludeParents)
	return req
}

func outputResult(cmd *cobra.Command, outputFormat string, project *resourcemanager.ProjectResponseWithParents, p *print.Printer) error {
	switch outputFormat {
	case globalflags.PrettyOutputFormat:
		table := tables.NewTable()
		table.AddRow("ID", *project.ProjectId)
		table.AddSeparator()
		table.AddRow("NAME", *project.Name)
		table.AddSeparator()
		table.AddRow("CREATION", *project.CreationTime)
		table.AddSeparator()
		table.AddRow("STATE", *project.LifecycleState)
		table.AddSeparator()
		table.AddRow("PARENT ID", *project.Parent.Id)
		err := table.Display(cmd)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	default:
		details, err := json.MarshalIndent(project, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal project details: %w", err)
		}
		p.Outputln(string(details))

		return nil
	}
}
