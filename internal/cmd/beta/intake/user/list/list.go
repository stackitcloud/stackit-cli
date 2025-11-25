package list

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/intake/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
)

const (
	intakeIdFlag = "intake-id"
	limitFlag    = "limit"
)

// inputModel struct holds all the input parameters for the command
type inputModel struct {
	*globalflags.GlobalFlagModel
	IntakeId string
	Limit    *int64
}

// NewCmd creates a new cobra command for listing Intake Users
func NewCmd(p *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all Intake Users for an Intake",
		Long:  "Lists all Intake Users for a specific Intake.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all Intake Users for an Intake with ID "xxx"`,
				`$ stackit beta intake user list --intake-id xxx`),
			examples.NewExample(
				`List up to 5 Intake Users for an Intake with ID "xxx"`,
				`$ stackit beta intake user list --intake-id xxx --limit 5`),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(p.Printer, cmd)
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
				return fmt.Errorf("list Intake Users: %w", err)
			}
			users := resp.GetIntakeUsers()

			// Truncate output
			if model.Limit != nil && len(users) > int(*model.Limit) {
				users = users[:*model.Limit]
			}

			projectLabel := model.ProjectId
			if len(users) == 0 {
				projectLabel, err = projectname.GetProjectName(ctx, p.Printer, p.CliVersion, cmd)
				if err != nil {
					p.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				}
			}

			return outputResult(p.Printer, model.OutputFormat, projectLabel, users)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(intakeIdFlag, "", "ID of the Intake")
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")

	err := flags.MarkFlagsRequired(cmd, intakeIdFlag)
	cobra.CheckErr(err)
}

// parseInput parses the command flags into a standardized model
func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
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

	limit := flags.FlagToInt64Pointer(p, cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &cliErr.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		IntakeId:        intakeId,
		Limit:           limit,
	}

	p.DebugInputModel(model)
	return &model, nil
}

// buildRequest creates the API request to list Intake Users
func buildRequest(ctx context.Context, model *inputModel, apiClient *intake.APIClient) intake.ApiListIntakeUsersRequest {
	req := apiClient.ListIntakeUsers(ctx, model.ProjectId, model.Region, model.IntakeId)
	return req
}

// outputResult formats the API response and prints it to the console
func outputResult(p *print.Printer, outputFormat, projectLabel string, users []intake.IntakeUserResponse) error {
	return p.OutputResult(outputFormat, users, func() error {
		if len(users) == 0 {
			p.Outputf("No intake users found for project %q\n", projectLabel)
			return nil
		}

		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "STATE", "TYPE", "USERNAME")
		for i := range users {
			user := users[i]
			table.AddRow(
				user.GetId(),
				user.GetDisplayName(),
				user.GetState(),
				user.GetType(),
				user.GetUser(),
			)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
