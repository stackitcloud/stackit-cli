package clone

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/confirm"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
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

	storageClassFlag      = "storage-class"
	storageSizeFlag       = "storage-size"
	recoveryTimestampFlag = "recovery-timestamp"
	recoveryDateFormat    = "2023-04-17T09:28:00+00:00"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	InstanceId   string
	StorageClass *string
	StorageSize  *int64
	RecoveryDate *string
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("clone %s", instanceIdArg),
		Short: "Clones a PostgreSQL Flex instance",
		Long:  "Clones a PostgreSQL Flex instance from a selected point in time.",
		Example: examples.Build(
			examples.NewExample(
				`Clone a PostgreSQL Flex instance with ID "xxx" . The recovery timestamp should be specified in UTC time following the format provided in the example.`,
				`$ stackit postgresflex instance clone xxx --recovery-timestamp 2023-04-17T09:28:00+00:00`),
			examples.NewExample(
				`Clone a PostgreSQL Flex instance with ID "xxx" from a selected recovery timestamp and specify storage class. If not specified, storage class from the existing instance will be used.`,
				`$ stackit postgresflex instance clone xxx --recovery-timestamp 2023-04-17T09:28:00+00:00 --storage-class premium-perf6-stackit`),
			examples.NewExample(
				`Clone a PostgreSQL Flex instance with ID "xxx" from a selected recovery timestamp and specify storage size. If not specified, storage size from the existing instance will be used.`,
				`$ stackit postgresflex instance clone xxx --recovery-timestamp 2023-04-17T09:28:00+00:00 --storage-size 10`),
		),
		Args: args.SingleArg(instanceIdArg, utils.ValidateUUID),
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
				prompt := fmt.Sprintf("Are you sure you want to clone instance %q?", instanceLabel)
				err = confirm.PromptForConfirmation(cmd, prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req, err := buildRequest(ctx, model, apiClient)
			if err != nil {
				return err
			}
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("clone PostgreSQL Flex instance: %w", err)
			}
			instanceId := *resp.InstanceId

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(cmd)
				s.Start("Cloning instance")
				_, err = wait.CreateInstanceWaitHandler(ctx, apiClient, model.ProjectId, instanceId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for PostgreSQL Flex instance cloning: %w", err)
				}
				s.Stop()
			}

			operationState := "Cloned"
			if model.Async {
				operationState = "Triggered cloning of"
			}

			cmd.Printf("%s instance from instance %q. New Instance ID: %s\n", operationState, instanceLabel, instanceId)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(recoveryTimestampFlag, "", "Recovery timestamp for the instance, in a date-time with the layout format, e.g. 2024-03-12T09:28:00+00:00")
	cmd.Flags().String(storageClassFlag, "", "Storage class")
	cmd.Flags().Int64(storageSizeFlag, 0, "Storage size (in GB)")

	err := flags.MarkFlagsRequired(cmd, recoveryTimestampFlag)
	cobra.CheckErr(err)
}

func parseInput(cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	instanceId := inputArgs[0]

	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      instanceId,
		StorageClass:    flags.FlagToStringPointer(cmd, storageClassFlag),
		StorageSize:     flags.FlagToInt64Pointer(cmd, storageSizeFlag),
		RecoveryDate:    flags.FlagToStringPointer(cmd, recoveryTimestampFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *postgresflex.APIClient) (postgresflex.ApiCloneInstanceRequest, error) {
	req := apiClient.CloneInstance(ctx, model.ProjectId, model.InstanceId)
	req = req.CloneInstancePayload(postgresflex.CloneInstancePayload{
		Class:     model.StorageClass,
		Size:      model.StorageSize,
		Timestamp: model.RecoveryDate,
	})
	return req, nil
}
