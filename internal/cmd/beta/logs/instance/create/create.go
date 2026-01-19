package create

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-sdk-go/services/logs"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/logs/client"

	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/logs/wait"
)

const (
	displayNameFlag   = "display-name"
	retentionDaysFlag = "retention-days"
	aclFlag           = "acl"
	descriptionFlag   = "description"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	DisplayName   *string
	RetentionDays *int64
	ACL           *[]string
	Description   *string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a Logs instance",
		Long:  "Creates a Logs instance.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a Logs instance with name "my-instance" and retention time 10 days`,
				`$ stackit beta logs instance create --display-name "my-instance" --retention-days 10`),
			examples.NewExample(
				`Create a Logs instance with name "my-instance", retention time 10 days, and a description`,
				`$ stackit beta logs instance create --display-name "my-instance" --retention-days 10 --description "Description of the instance"`),
			examples.NewExample(
				`Create a Logs instance with name "my-instance", retention time 10 days, and restrict access to a specific range of IP addresses.`,
				`$ stackit beta logs instance create --display-name "my-instance" --retention-days 10 --acl 1.2.3.0/24`),
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

			prompt := fmt.Sprintf("Are you sure you want to create a Logs instance for project %q?", projectLabel)
			err = params.Printer.PromptForConfirmation(prompt)

			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			if err != nil {
				return err
			}
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create Logs instance: %w", err)
			}
			instanceId := *resp.Id

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Creating instance")
				_, err = wait.CreateLogsInstanceWaitHandler(ctx, apiClient, model.ProjectId, model.Region, instanceId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for LogMe instance creation: %w", err)
				}
				s.Stop()
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

	err := flags.MarkFlagsRequired(cmd, displayNameFlag, retentionDaysFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		DisplayName:     flags.FlagToStringPointer(p, cmd, displayNameFlag),
		RetentionDays:   flags.FlagToInt64Pointer(p, cmd, retentionDaysFlag),
		Description:     flags.FlagToStringPointer(p, cmd, descriptionFlag),
		ACL:             flags.FlagToStringSlicePointer(p, cmd, aclFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *logs.APIClient) logs.ApiCreateLogsInstanceRequest {
	req := apiClient.CreateLogsInstance(ctx, model.ProjectId, model.Region)

	req = req.CreateLogsInstancePayload(logs.CreateLogsInstancePayload{
		DisplayName:   model.DisplayName,
		Description:   model.Description,
		RetentionDays: model.RetentionDays,
		Acl:           model.ACL,
	})
	return req
}

func outputResult(p *print.Printer, model *inputModel, projectLabel string, resp *logs.LogsInstance) error {
	if resp == nil {
		return fmt.Errorf("create logs instance response is empty")
	}
	var outputFormat string
	var async bool

	if model.GlobalFlagModel != nil {
		outputFormat = model.OutputFormat
		async = model.Async
	}

	return p.OutputResult(outputFormat, resp, func() error {
		operationState := "Created"
		if async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s instance for project %q. Instance ID: %s\n", operationState, projectLabel, utils.PtrString(resp.Id))
		return nil
	})
}
