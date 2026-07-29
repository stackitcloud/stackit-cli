package delete

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresflex/client"
	postgresflexUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresflex/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	postgresflex "github.com/stackitcloud/stackit-sdk-go/services/postgresflex/v3api"
	"github.com/stackitcloud/stackit-sdk-go/services/postgresflex/v3api/wait"
)

const (
	instanceIdArg = "INSTANCE_ID"

	forceDeleteFlag = "force" // Deprecated: Will be removed after 2026-01-31
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId  string
	ForceDelete bool
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", instanceIdArg),
		Short: "Deletes a PostgreSQL Flex instance",
		Long: fmt.Sprintf("%s\n%s\n%s",
			"Deletes a PostgreSQL Flex instance.",
			"By default, instances will be kept in a delayed deleted state for 7 days before being permanently deleted.",
			"Use the --force flag to force the immediate deletion of a delayed deleted instance.",
		),
		Args: args.SingleArg(instanceIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete a PostgreSQL Flex instance with ID "xxx"`,
				"$ stackit postgresflex instance delete xxx"),
			examples.NewExample(
				`Force the deletion of a delayed deleted PostgreSQL Flex instance with ID "xxx"`,
				"$ stackit postgresflex instance delete xxx --force"),
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

			instanceLabel, err := postgresflexUtils.GetInstanceName(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.InstanceId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceId
			}

			prompt := fmt.Sprintf("Are you sure you want to delete instance %q? (This cannot be undone)", instanceLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			delReq := buildDeleteRequest(ctx, model, apiClient)
			err = delReq.Execute()
			if err != nil {
				return fmt.Errorf("delete PostgreSQL Flex instance: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				err := spinner.Run(params.Printer, "Deleting instance", func() error {
					_, err = wait.DeleteInstanceWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.InstanceId).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for PostgreSQL Flex instance deletion: %w", err)
				}
			}

			operationState := "Deleted"
			if model.Async {
				operationState = "Triggered deletion of"
			}

			params.Printer.Outputf("%s instance %q\n", operationState, instanceLabel)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().BoolP(forceDeleteFlag, "f", false, "Force deletion of a delayed deleted instance")

	err := cmd.Flags().MarkDeprecated(forceDeleteFlag, "Force delete is the default option by now. This flag will be removed after 2027-01-31.")
	cobra.CheckErr(err)
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
		ForceDelete:     flags.FlagToBoolValue(p, cmd, forceDeleteFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildDeleteRequest(ctx context.Context, model *inputModel, apiClient *postgresflex.APIClient) postgresflex.ApiDeleteInstanceRequest {
	req := apiClient.DefaultAPI.DeleteInstance(ctx, model.ProjectId, model.Region, model.InstanceId)
	return req
}
