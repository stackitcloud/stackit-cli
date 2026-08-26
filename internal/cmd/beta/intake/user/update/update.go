package update

import (
	"context"
	"fmt"
	"io/fs"
	"strings"

	"github.com/spf13/cobra"
	intake "github.com/stackitcloud/stackit-sdk-go/services/intake/v1betaapi"
	"github.com/stackitcloud/stackit-sdk-go/services/intake/v1betaapi/wait"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/intake/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	userIdArg = "USER_ID"

	intakeIdFlag    = "intake-id"
	displayNameFlag = "display-name"
	descriptionFlag = "description"
	passwordFlag    = "password"
	userTypeFlag    = "type"
	labelsFlag      = "labels"

	interactivePasswordPlaceholder = "__INTERACTIVE__"
)

type secretUpdateFlag struct {
	printer  *print.Printer
	fs       fs.FS
	value    string
	isPrompt bool
}

func (f *secretUpdateFlag) String() string {
	return f.value
}

func (f *secretUpdateFlag) Set(value string) error {
	if value == interactivePasswordPlaceholder {
		f.isPrompt = true
		return nil
	}
	if strings.HasPrefix(value, "@") {
		path := strings.Trim(value[1:], `"'`)
		bytes, err := fs.ReadFile(f.fs, path)
		if err != nil {
			return fmt.Errorf("reading secret %s: %w", passwordFlag, err)
		}
		val := strings.TrimRight(string(bytes), "\r\n")
		if val == "" {
			return fmt.Errorf("the provided secret file %q is empty", path)
		}
		f.value = val
		return nil
	}
	f.printer.Warn("Passing a secret value on the command line is insecure and deprecated. This usage will stop working October 2026.\n")
	f.value = value
	return nil
}

func (f *secretUpdateFlag) Type() string {
	return "string"
}

type inputModel struct {
	*globalflags.GlobalFlagModel
	IntakeId    string
	UserId      string
	DisplayName *string
	Description *string
	Password    *string
	UserType    *string
	Labels      *map[string]string
}

func NewCmd(p *types.CmdParams) *cobra.Command {
	password := &secretUpdateFlag{
		printer: p.Printer,
		fs:      p.Fs,
	}

	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", userIdArg),
		Short: "Updates an Intake User",
		Long:  "Updates an Intake User. Only the specified fields are updated.",
		Args:  args.SingleArg(userIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update the display name of an Intake User`,
				`$ stackit beta intake user update xxx --intake-id yyy --display-name "new-user-name"`),
			examples.NewExample(
				`Update the password interactively for an Intake User`,
				`$ stackit beta intake user update xxx --intake-id yyy --password`),
			examples.NewExample(
				`Update the password and description for an Intake User from a file`,
				`$ stackit beta intake user update xxx --intake-id yyy --password @./secret.txt --description "Updated description"`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p.Printer, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p.Printer, p.CliVersion)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update Intake User: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				err := spinner.Run(p.Printer, "Updating STACKIT Intake User", func() error {
					_, err = wait.UpdateIntakeUserWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.IntakeId, model.UserId).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for STACKIT Intake User update: %w", err)
				}
			}

			return outputResult(p.Printer, model, resp)
		},
	}
	configureFlags(cmd, password)
	return cmd
}

func configureFlags(cmd *cobra.Command, password *secretUpdateFlag) {
	cmd.Flags().Var(flags.UUIDFlag(), intakeIdFlag, "Intake ID")
	cmd.Flags().String(displayNameFlag, "", "Display name")
	cmd.Flags().String(descriptionFlag, "", "Description")
	cmd.Flags().Var(password, passwordFlag, "Password. Can be a string (deprecated) or a file path, if prefixed with '@' (example: @./secret.txt). If provided without a value, you will be prompted interactively. Must contain lower, upper, digits, and special characters (min 12 chars).")
	cmd.Flags().Lookup(passwordFlag).NoOptDefVal = interactivePasswordPlaceholder
	cmd.Flags().String(userTypeFlag, "", "Type of user. One of 'intake' or 'dead-letter'")
	cmd.Flags().StringToString(labelsFlag, nil, `Labels in key=value format, separated by commas. Example: --labels "key1=value1,key2=value2".`)

	err := flags.MarkFlagsRequired(cmd, intakeIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	userId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	password, err := parsePassword(p, cmd)
	if err != nil {
		return nil, err
	}

	model := &inputModel{
		GlobalFlagModel: globalFlags,
		IntakeId:        flags.FlagToStringValue(p, cmd, intakeIdFlag),
		UserId:          userId,
		DisplayName:     flags.FlagToStringPointer(p, cmd, displayNameFlag),
		Description:     flags.FlagToStringPointer(p, cmd, descriptionFlag),
		Password:        password,
		UserType:        flags.FlagToStringPointer(p, cmd, userTypeFlag),
		Labels:          flags.FlagToStringToStringPointer(p, cmd, labelsFlag),
	}

	if model.DisplayName == nil && model.Description == nil && model.Password == nil && model.UserType == nil && model.Labels == nil {
		return nil, &cliErr.EmptyUpdateError{}
	}

	p.DebugInputModel(model)
	return model, nil
}

func parsePassword(p *print.Printer, cmd *cobra.Command) (*string, error) {
	flag := cmd.Flag(passwordFlag)
	if flag == nil || !flag.Changed {
		return nil, nil
	}
	if secretFlag, ok := flag.Value.(*secretUpdateFlag); ok && secretFlag.isPrompt {
		input, err := p.PromptForPassword("enter new password: ")
		if err != nil {
			return nil, fmt.Errorf("prompt for password: %w", err)
		}
		input = strings.TrimRight(input, "\r\n")
		if input == "" {
			return nil, fmt.Errorf("password cannot be empty")
		}
		return &input, nil
	}
	val := strings.TrimRight(flag.Value.String(), "\r\n")
	if val == "" {
		return nil, fmt.Errorf("the provided password (or secret file) is empty")
	}
	return &val, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *intake.APIClient) intake.ApiUpdateIntakeUserRequest {
	req := apiClient.DefaultAPI.UpdateIntakeUser(ctx, model.ProjectId, model.Region, model.IntakeId, model.UserId)

	payload := intake.UpdateIntakeUserPayload{
		DisplayName: model.DisplayName,
		Description: model.Description,
		Password:    model.Password,
		Labels:      utils.PtrValue(model.Labels),
	}

	if model.UserType != nil {
		userType := intake.UserType(*model.UserType)
		payload.Type = &userType
	}

	req = req.UpdateIntakeUserPayload(payload)
	return req
}

func outputResult(p *print.Printer, model *inputModel, resp *intake.IntakeUserResponse) error {
	return p.OutputResult(model.OutputFormat, resp, func() error {
		operationState := "Updated"
		if model.Async {
			operationState = "Triggered update of"
		}
		p.Outputf("%s Intake User %s\n", operationState, model.UserId)
		return nil
	})
}
