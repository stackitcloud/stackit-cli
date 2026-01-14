package update

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/logs/client"
	logsUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/logs/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/logs"
)

const (
	argInstanceID = "INSTANCE_ID"

	displayNameFlag   = "display-name"
	retentionDaysFlag = "retention-days"
	aclFlag           = "acl"
	descriptionFlag   = "description"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceID    string
	DisplayName   *string
	RetentionDays *int64
	ACL           *[]string
	Description   *string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", argInstanceID),
		Short: "Updates a Logs instance",
		Long:  "Updates a Logs instance.",
		Args:  args.SingleArg(argInstanceID, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update the display name of the Logs instance with ID "xxx"`,
				"$ stackit beta logs instance update xxx --display-name new-name"),
			examples.NewExample(
				`Update the retention time of the Logs instance with ID "xxx"`,
				"$ stackit beta logs instance update xxx --retention-days 40"),
			examples.NewExample(
				`Update the ACL of the Logs instance with ID "xxx"`,
				"$ stackit beta logs instance update xxx --acl 1.2.3.0/24"),
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

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			} else if projectLabel == "" {
				projectLabel = model.ProjectId
			}

			instanceLabel, err := logsUtils.GetInstanceName(ctx, apiClient, model.ProjectId, model.Region, model.InstanceID)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceID
			}

			prompt := fmt.Sprintf("Are you sure you want to update instance %s?", instanceLabel)
			err = params.Printer.PromptForConfirmation(prompt)

			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update logs instance: %w", err)
			}

			return outputResult(params.Printer, model, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(displayNameFlag, "", "Display name")
	cmd.Flags().String(descriptionFlag, "", "Description")
	cmd.Flags().StringSlice(aclFlag, []string{}, "Access control list")
	cmd.Flags().Int64(retentionDaysFlag, 0, "The days for how long the logs should be stored before being cleaned up")
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	instanceId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	displayName := flags.FlagToStringPointer(p, cmd, displayNameFlag)
	retentionDays := flags.FlagToInt64Pointer(p, cmd, retentionDaysFlag)
	acl := flags.FlagToStringSlicePointer(p, cmd, aclFlag)
	description := flags.FlagToStringPointer(p, cmd, descriptionFlag)

	if displayName == nil && retentionDays == nil && acl == nil && description == nil {
		return nil, &errors.EmptyUpdateError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceID:      instanceId,
		DisplayName:     displayName,
		ACL:             acl,
		Description:     description,
		RetentionDays:   retentionDays,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *logs.APIClient) logs.ApiUpdateLogsInstanceRequest {
	req := apiClient.UpdateLogsInstance(ctx, model.ProjectId, model.Region, model.InstanceID)
	req = req.UpdateLogsInstancePayload(logs.UpdateLogsInstancePayload{
		DisplayName:   model.DisplayName,
		Acl:           model.ACL,
		RetentionDays: model.RetentionDays,
		Description:   model.Description,
	})
	return req
}

func outputResult(p *print.Printer, model *inputModel, projectLabel string, instance *logs.LogsInstance) error {
	if instance == nil {
		return fmt.Errorf("instance is nil")
	}
	return p.OutputResult(model.OutputFormat, instance, func() error {
		p.Outputf("Updated instance %q for project %q.\n", utils.PtrString(instance.DisplayName), projectLabel)
		return nil
	})
}
