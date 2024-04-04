package delete

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/confirm"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresflex/client"
	postgresflexUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresflex/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/postgresflex"
	"github.com/stackitcloud/stackit-sdk-go/services/postgresflex/wait"
)

const (
	instanceIdArg = "INSTANCE_ID"

	forceDeleteFlag = "force"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId  string
	ForceDelete bool
}

func NewCmd() *cobra.Command {
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
			model, err := parseInput(cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(cmd)
			if err != nil {
				return err
			}

			instanceLabel, err := postgresflexUtils.GetInstanceName(ctx, apiClient, model.ProjectId, model.InstanceId)
			if err != nil {
				instanceLabel = model.InstanceId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete instance %q? (This cannot be undone)", instanceLabel)
				err = confirm.PromptForConfirmation(cmd, prompt)
				if err != nil {
					return err
				}
			}

			toDelete, toForceDelete, err := getNextOperations(ctx, model, apiClient)
			if err != nil {
				return err
			}

			if toDelete {
				// Call API
				delReq := buildDeleteRequest(ctx, model, apiClient)
				err = delReq.Execute()
				if err != nil {
					return fmt.Errorf("delete PostgreSQL Flex instance: %w", err)
				}

				// Wait for async operation, if async mode not enabled
				if !model.Async {
					s := spinner.New(cmd)
					s.Start("Deleting instance")
					_, err = wait.DeleteInstanceWaitHandler(ctx, apiClient, model.ProjectId, model.InstanceId).WaitWithContext(ctx)
					if err != nil {
						return fmt.Errorf("wait for PostgreSQL Flex instance deletion: %w", err)
					}
					s.Stop()
				}
			}

			if toForceDelete {
				// Call API
				forceDelReq := buildForceDeleteRequest(ctx, model, apiClient)
				err = forceDelReq.Execute()
				if err != nil {
					return fmt.Errorf("force delete PostgreSQL Flex instance: %w", err)
				}

				// Wait for async operation, if async mode not enabled
				if !model.Async {
					s := spinner.New(cmd)
					s.Start("Forcing deletion of instance")
					_, err = wait.ForceDeleteInstanceWaitHandler(ctx, apiClient, model.ProjectId, model.InstanceId).WaitWithContext(ctx)
					if err != nil {
						return fmt.Errorf("wait for PostgreSQL Flex instance force deletion: %w", err)
					}
					s.Stop()
				}
			}

			operationState := "Deleted"
			if toForceDelete {
				operationState = "Forcefully deleted"
			}
			if model.Async {
				operationState = "Triggered deletion of"
				if toForceDelete {
					operationState = "Triggered forced deletion of"
				}
			}

			cmd.Printf("%s instance %q\n", operationState, instanceLabel)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().BoolP(forceDeleteFlag, "f", false, "Force deletion of a delayed deleted instance")
}

func parseInput(cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	instanceId := inputArgs[0]

	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      instanceId,
		ForceDelete:     flags.FlagToBoolValue(cmd, forceDeleteFlag),
	}, nil
}

func buildDeleteRequest(ctx context.Context, model *inputModel, apiClient *postgresflex.APIClient) postgresflex.ApiDeleteInstanceRequest {
	req := apiClient.DeleteInstance(ctx, model.ProjectId, model.InstanceId)
	return req
}

func buildForceDeleteRequest(ctx context.Context, model *inputModel, apiClient *postgresflex.APIClient) postgresflex.ApiForceDeleteInstanceRequest {
	req := apiClient.ForceDeleteInstance(ctx, model.ProjectId, model.InstanceId)
	return req
}

type PostgreSQLFlexClient interface {
	GetInstanceExecute(ctx context.Context, projectId, instanceId string) (*postgresflex.InstanceResponse, error)
	ListVersionsExecute(ctx context.Context, projectId string) (*postgresflex.ListVersionsResponse, error)
	GetUserExecute(ctx context.Context, projectId, instanceId, userId string) (*postgresflex.GetUserResponse, error)
}

func getNextOperations(ctx context.Context, model *inputModel, apiClient PostgreSQLFlexClient) (toDelete, toForceDelete bool, err error) {
	instanceStatus, err := postgresflexUtils.GetInstanceStatus(ctx, apiClient, model.ProjectId, model.InstanceId)
	if err != nil {
		return false, false, fmt.Errorf("get PostgreSQL Flex instance status: %w", err)
	}

	if instanceStatus == wait.InstanceStateDeleted {
		if !model.ForceDelete {
			return false, false, fmt.Errorf("instance is already deleted, use --force to force the deletion of a delayed deleted instance")
		}

		return false, model.ForceDelete, nil
	}

	return true, model.ForceDelete, nil
}
