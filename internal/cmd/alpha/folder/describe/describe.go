package describe

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/resourcemanager/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	resourcemanager "github.com/stackitcloud/stackit-sdk-go/services/resourcemanager/v0api"
)

const (
	includeParentsFlag = "include-parents"

	containerIdArg = "CONTAINER_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ArgContainerId string
	IncludeParents bool
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Shows details of a STACKIT folder",
		Long:  "Shows details of a STACKIT folder.",
		Args:  args.SingleOptionalArg(containerIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get the details of the configured STACKIT folder`,
				"$ stackit folder describe"),
			examples.NewExample(
				`Get the details of a STACKIT folder by explicitly providing the container ID`,
				"$ stackit folder describe xxx"),
			examples.NewExample(
				`Get the details of the configured STACKIT folder, including details of the parent resources`,
				"$ stackit folder describe --include-parents"),
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
				return fmt.Errorf("read folder details: %w", err)
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
	var containerId string
	if len(inputArgs) > 0 {
		containerId = inputArgs[0]
	}

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" && containerId == "" {
		return nil, fmt.Errorf("container ID needs to be provided either as an argument or as a flag")
	}

	if containerId == "" {
		containerId = globalFlags.ProjectId
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ArgContainerId:  containerId,
		IncludeParents:  flags.FlagToBoolValue(p, cmd, includeParentsFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *resourcemanager.APIClient) resourcemanager.ApiGetFolderDetailsRequest {
	req := apiClient.DefaultAPI.GetFolderDetails(ctx, model.ArgContainerId)
	req.IncludeParents(model.IncludeParents)
	return req
}

func outputResult(p *print.Printer, outputFormat string, folder *resourcemanager.GetFolderDetailsResponse) error {
	if folder == nil {
		return fmt.Errorf("response not set")
	}

	return p.OutputResult(outputFormat, folder, func() error {
		table := tables.NewTable()
		table.AddRow("ID", folder.FolderId)
		table.AddSeparator()
		table.AddRow("NAME", folder.Name)
		table.AddSeparator()
		table.AddRow("CREATION", folder.CreationTime)
		table.AddSeparator()
		table.AddRow("CONTAINER ID", folder.ContainerId)
		table.AddSeparator()
		table.AddRow("PARENT ID", folder.Parent.Id)

		var labels string
		if folder.Labels != nil {
			labels = utils.JoinStringMap(*folder.Labels, "=", ", ")
		}
		table.AddSeparator()
		table.AddRow("LABELS", labels)

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
