package describe

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/intake"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/intake/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	userIdArg    = "USER_ID"
	intakeIdFlag = "intake-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	IntakeId string
	UserId   string
}

func NewCmd(p *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", userIdArg),
		Short: "Shows details of an Intake User",
		Long:  "Shows details of an Intake User.",
		Args:  args.SingleArg(userIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get details of an Intake User with ID "xxx" from an Intake with ID "yyy"`,
				`$ stackit beta intake user describe xxx --intake-id yyy`),
			examples.NewExample(
				`Get details of an Intake User in JSON format`,
				`$ stackit beta intake user describe xxx --intake-id yyy --output-format json`),
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
				return fmt.Errorf("get Intake User: %w", err)
			}

			return outputResult(p.Printer, model.OutputFormat, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(intakeIdFlag, "", "ID of the Intake to which the user belongs")

	err := flags.MarkFlagsRequired(cmd, intakeIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	userId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	intakeId := flags.FlagToStringValue(p, cmd, intakeIdFlag)
	if intakeId == "" {
		return nil, &cliErr.FlagValidationError{
			Flag:    intakeIdFlag,
			Details: "can't be empty",
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		IntakeId:        intakeId,
		UserId:          userId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *intake.APIClient) intake.ApiGetIntakeUserRequest {
	req := apiClient.GetIntakeUser(ctx, model.ProjectId, model.Region, model.IntakeId, model.UserId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, user *intake.IntakeUserResponse) error {
	if user == nil {
		return fmt.Errorf("received nil response, could not display details")
	}

	return p.OutputResult(outputFormat, user, func() error {
		table := tables.NewTable()
		table.SetHeader("Attribute", "Value")

		table.AddRow("ID", user.GetId())
		table.AddRow("Name", user.GetDisplayName())
		table.AddRow("State", user.GetState())
		table.AddRow("Created", user.GetCreateTime())
		table.AddRow("Labels", user.GetLabels())
		table.AddRow("Type", user.GetType())
		table.AddRow("Username", user.GetUser())

		if description := user.GetDescription(); description != "" {
			table.AddRow("Description", description)
		}
		table.AddSeparator()

		clientConfig := user.GetClientConfig()
		table.AddRow("Java Client Config", clientConfig.GetJava())
		table.AddSeparator()
		table.AddRow("librdkafka Client Config", clientConfig.GetLibrdkafka())

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
