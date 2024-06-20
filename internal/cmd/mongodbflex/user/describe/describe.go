package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/mongodbflex/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/mongodbflex"
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
		Short: "Shows details of a MongoDB Flex user",
		Long: fmt.Sprintf("%s\n%s\n%s",
			"Shows details of a MongoDB Flex user.",
			`The user password is hidden inside the "host" field and replaced with asterisks, as it is only visible upon creation. You can reset it by running:`,
			"  $ stackit mongodbflex user reset-password USER_ID --instance-id INSTANCE_ID",
		),
		Example: examples.Build(
			examples.NewExample(
				`Get details of a MongoDB Flex user with ID "xxx" of instance with ID "yyy"`,
				"$ stackit mongodbflex user describe xxx --instance-id yyy"),
			examples.NewExample(
				`Get details of a MongoDB Flex user with ID "xxx" of instance with ID "yyy" in JSON format`,
				"$ stackit mongodbflex user describe xxx --instance-id yyy --output-format json"),
		),
		Args: args.SingleArg(userIdArg, utils.ValidateUUID),
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
				return fmt.Errorf("get MongoDB Flex user: %w", err)
			}

			return outputResult(p, model.OutputFormat, *resp.Item)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *mongodbflex.APIClient) mongodbflex.ApiGetUserRequest {
	req := apiClient.GetUser(ctx, model.ProjectId, model.InstanceId, model.UserId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, user mongodbflex.InstanceResponseUser) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(user, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal MongoDB Flex user: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(user, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal MongoDB Flex user: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.AddRow("ID", *user.Id)
		table.AddSeparator()
		table.AddRow("USERNAME", *user.Username)
		table.AddSeparator()
		table.AddRow("ROLES", *user.Roles)
		table.AddSeparator()
		table.AddRow("DATABASE", *user.Database)
		table.AddSeparator()
		table.AddRow("HOST", *user.Host)
		table.AddSeparator()
		table.AddRow("PORT", *user.Port)

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
