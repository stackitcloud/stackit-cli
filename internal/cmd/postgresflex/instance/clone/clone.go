package clone

import (
	"context"
	"fmt"
	"time"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresflex/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	postgresflex "github.com/stackitcloud/stackit-sdk-go/services/postgresflex/v3api"
	"github.com/stackitcloud/stackit-sdk-go/services/postgresflex/v3api/wait"
)

const (
	instanceIdArg = "INSTANCE_ID"

	storageClassFlag      = "storage-class"
	storageSizeFlag       = "storage-size"
	recoveryTimestampFlag = "recovery-timestamp"
	recoveryDateFormat    = "2006-01-02T15:04:05-07:00"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	InstanceId   string
	InstanceName *string
	StorageClass *string
	StorageSize  *int64
	RecoveryDate time.Time
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("clone %s", instanceIdArg),
		Short: "Clones a PostgreSQL Flex instance",
		Long: "Clones a PostgreSQL Flex instance from a selected point in time. " +
			"The new cloned instance will be an independent instance with the same settings as the original instance unless the flags are specified.",
		Example: examples.Build(
			examples.NewExample(
				`Clone a PostgreSQL Flex instance with ID "xxx" from a selected recovery timestamp.`,
				`$ stackit postgresflex instance clone xxx --recovery-timestamp 2023-04-17T09:28:00+00:00 --storage-size 10 --storage-class premium-perf6-stackit`),
		),
		Args: args.SingleArg(instanceIdArg, utils.ValidateUUID),
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

			instance, err := apiClient.DefaultAPI.GetInstance(ctx, model.ProjectId, model.Region, model.InstanceId).Execute()
			if err != nil {
				return fmt.Errorf("get PostgreSQL Flex instance: %w", err)
			}

			instanceLabel := instance.Name

			if model.StorageSize == nil {
				params.Printer.Warn("The --%s flag is not set. Using the storage size from the instance you're cloning. This behavior is deprecated, the --%s flag will be required after 2027-01-31.\n", storageSizeFlag, storageSizeFlag)

				if instance.Storage.Size == nil {
					return fmt.Errorf("could not read storage size for instance %s", model.InstanceId)
				}

				model.StorageSize = instance.Storage.Size
			}

			if model.StorageClass == nil {
				params.Printer.Warn("The --%s flag is not set. Using the storage class from the instance you're cloning. This behavior is deprecated, the --%s flag will be required after 2027-01-31.\n", storageClassFlag, storageClassFlag)

				if instance.Storage.Class == nil {
					return fmt.Errorf("could not read storage class for instance %s", model.InstanceId)
				}

				model.StorageClass = instance.Storage.Class
			}

			prompt := fmt.Sprintf("Are you sure you want to clone instance %q?", instanceLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req, err := buildRequest(ctx, model, apiClient.DefaultAPI)
			if err != nil {
				return err
			}
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("clone PostgreSQL Flex instance: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				err := spinner.Run(params.Printer, "Cloning instance", func() error {
					_, err = wait.CreateInstanceWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, resp.Id).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for PostgreSQL Flex instance cloning: %w", err)
				}
			}

			return outputResult(params.Printer, model.OutputFormat, model.Async, instanceLabel, resp.Id, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(recoveryTimestampFlag, "", "Recovery timestamp for the instance, in a date-time with the layout format YYYY-MM-DDTHH:mm:ss±HH:mm, e.g. 2006-01-02T15:04:05-07:00")
	cmd.Flags().String(storageClassFlag, "", "Storage class. If not specified, storage class from the existing instance will be used. This flag will be required after 2027-01-31.")
	cmd.Flags().Int64(storageSizeFlag, 0, "Storage size (in GB). If not specified, storage size from the existing instance will be used. This flag will be required after 2027-01-31.")

	// mark storage-size flag required here after 2027-01-31

	err := flags.MarkFlagsRequired(cmd, recoveryTimestampFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	instanceId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	recoveryTimestamp, err := flags.FlagToDateTimePointer(p, cmd, recoveryTimestampFlag, recoveryDateFormat)
	if err != nil {
		return nil, &cliErr.FlagValidationError{
			Flag:    recoveryTimestampFlag,
			Details: err.Error(),
		}
	} else if recoveryTimestamp == nil {
		return nil, &cliErr.FlagValidationError{
			Flag:    recoveryTimestampFlag,
			Details: fmt.Sprintf("the --%s flag is required", recoveryTimestampFlag),
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      instanceId,
		StorageClass:    flags.FlagToStringPointer(p, cmd, storageClassFlag),
		StorageSize:     flags.FlagToInt64Pointer(p, cmd, storageSizeFlag),
		RecoveryDate:    *recoveryTimestamp,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient postgresflex.DefaultAPI) (postgresflex.ApiCloneInstanceRequest, error) {
	if model.StorageSize == nil {
		return postgresflex.ApiCloneInstanceRequest{}, fmt.Errorf("storage size is nil")
	}
	if model.StorageClass == nil {
		return postgresflex.ApiCloneInstanceRequest{}, fmt.Errorf("storage class is nil")
	}

	payload := postgresflex.CloneInstancePayload{
		InstanceOverrides: postgresflex.CloneInstanceOverrides{
			Class: *model.StorageClass,
			Size:  *model.StorageSize,
			Name:  model.InstanceName,
		},
		PointInTime: model.RecoveryDate,
	}

	return apiClient.CloneInstance(ctx, model.ProjectId, model.Region, model.InstanceId).CloneInstancePayload(payload), nil
}

func outputResult(p *print.Printer, outputFormat string, async bool, instanceLabel, instanceId string, resp *postgresflex.CloneInstanceResponse) error {
	return p.OutputResult(outputFormat, resp, func() error {
		if resp == nil {
			return fmt.Errorf("response not set")
		}

		operationState := "Cloned"
		if async {
			operationState = "Triggered cloning of"
		}

		p.Outputf("%s instance from instance %q. New Instance ID: %s\n", operationState, instanceLabel, instanceId)
		return nil
	})
}
