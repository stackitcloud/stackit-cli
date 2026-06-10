package describe

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	git "github.com/stackitcloud/stackit-sdk-go/services/git/v1betaapi"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/git/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId string
}

const instanceIdArg = "INSTANCE_ID"

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", instanceIdArg),
		Short: "Describes STACKIT Git instance",
		Long:  "Describes a STACKIT Git instance by its internal ID.",
		Args:  args.SingleArg(instanceIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(`Describe instance "xxx"`, `$ stackit git describe xxx`),
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
			request := buildRequest(ctx, model, apiClient)

			instance, err := request.Execute()
			if err != nil {
				return fmt.Errorf("get instance: %w", err)
			}

			if err := outputResult(params.Printer, model.OutputFormat, instance); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, cliArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      cliArgs[0],
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *git.APIClient) git.ApiGetInstanceRequest {
	return apiClient.DefaultAPI.GetInstance(ctx, model.ProjectId, model.InstanceId)
}

func outputResult(p *print.Printer, outputFormat string, resp *git.Instance) error {
	if resp == nil {
		return fmt.Errorf("instance not found")
	}

	return p.OutputResult(outputFormat, resp, func() error {
		table := tables.NewTable()

		table.AddRow("ID", resp.Id)
		table.AddSeparator()
		table.AddRow("NAME", resp.Name)
		table.AddSeparator()
		table.AddRow("URL", resp.Url)
		table.AddSeparator()
		table.AddRow("VERSION", resp.Version)
		table.AddSeparator()
		table.AddRow("STATE", resp.State)
		table.AddSeparator()
		table.AddRow("CREATED", resp.Created)
		table.AddSeparator()

		if err := table.Display(p); err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
