package describe

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sqlserverflex/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/sqlserverflex"
)

const (
	userIdArg = "USER_ID"

	instanceIdFlag = "instance-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	InstanceId string
	UserId     string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", userIdArg),
		Short: "Shows details of a SQLServer Flex user",
		Long: fmt.Sprintf("%s\n%s\n%s",
			"Shows details of a SQLServer Flex user.",
			`The user password is only visible upon creation. You can reset it by running:`,
			"  $ stackit beta sqlserverflex user reset-password USER_ID --instance-id INSTANCE_ID",
		),
		Example: examples.Build(
			examples.NewExample(
				`Get details of a SQLServer Flex user with ID "xxx" of instance with ID "yyy"`,
				"$ stackit beta sqlserverflex user describe xxx --instance-id yyy"),
			examples.NewExample(
				`Get details of a SQLServer Flex user with ID "xxx" of instance with ID "yyy" in JSON format`,
				"$ stackit beta sqlserverflex user describe xxx --instance-id yyy --output-format json"),
		),
		Args: args.SingleArg(userIdArg, nil),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("get SQLServer Flex user: %w", err)
			}

			return outputResult(p, model.OutputFormat, resp.Item)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "ID of the instance")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	userId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
		UserId:          userId,
	}

	if p.IsVerbosityDebug() {
		modelStr, err := print.BuildDebugStrFromInputModel(model)
		if err != nil {
			p.Debug(print.ErrorLevel, "convert model to string for debugging: %v", err)
		} else {
			p.Debug(print.DebugLevel, "parsed input values: %s", modelStr)
		}
	}

	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *sqlserverflex.APIClient) sqlserverflex.ApiGetUserRequest {
	req := apiClient.GetUser(ctx, model.ProjectId, model.InstanceId, model.UserId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, user *sqlserverflex.UserResponseUser) error {
	if user == nil {
		return fmt.Errorf("user response is empty")
	}
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(user, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal SQLServer Flex user: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(user, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal SQLServer Flex user: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.AddRow("ID", utils.PtrString(user.Id))
		table.AddSeparator()
		table.AddRow("USERNAME", utils.PtrString(user.Username))
		if user.Roles != nil && len(*user.Roles) != 0 {
			table.AddSeparator()
			table.AddRow("ROLES", strings.Join(*user.Roles, "\n"))
		}
		if user.DefaultDatabase != nil && *user.DefaultDatabase != "" {
			table.AddSeparator()
			table.AddRow("DATABASE", *user.DefaultDatabase)
		}
		if user.Host != nil && *user.Host != "" {
			table.AddSeparator()
			table.AddRow("HOST", *user.Host)
		}
		if user.Port != nil {
			table.AddSeparator()
			table.AddRow("PORT", *user.Port)
		}

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
