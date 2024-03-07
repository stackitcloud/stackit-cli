package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/secrets-manager/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/secretsmanager"
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

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", userIdArg),
		Short: "Shows details of a Secrets Manager user",
		Long:  "Shows details of a Secrets Manager user.",
		Example: examples.Build(
			examples.NewExample(
				`Get details of a Secrets Manager user with ID "xxx" of instance with ID "yyy"`,
				"$ stackit secrets-manager user list xxx --instance-id yyy"),
			examples.NewExample(
				`Get details of a Secrets Manager user with ID "xxx" of instance with ID "yyy" in table format`,
				"$ stackit secrets-manager user list xxx --instance-id yyy --output-format pretty"),
		),
		Args: args.SingleArg(userIdArg, utils.ValidateUUID),
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("get Secrets Manager user: %w", err)
			}

			return outputResult(cmd, model.OutputFormat, *resp)
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

func parseInput(cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	userId := inputArgs[0]

	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringValue(cmd, instanceIdFlag),
		UserId:          userId,
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *secretsmanager.APIClient) secretsmanager.ApiGetUserRequest {
	req := apiClient.GetUser(ctx, model.ProjectId, model.InstanceId, model.UserId)
	return req
}

func outputResult(cmd *cobra.Command, outputFormat string, user secretsmanager.User) error {
	switch outputFormat {
	case globalflags.PrettyOutputFormat:
		table := tables.NewTable()
		table.AddRow("ID", *user.Id)
		table.AddSeparator()
		table.AddRow("USERNAME", *user.Username)
		table.AddSeparator()
		table.AddRow("DESCRIPTION", *user.Description)
		if user.Password != nil && *user.Password != "" {
			table.AddSeparator()
			table.AddRow("PASSWORD", *user.Password)
		}
		table.AddSeparator()
		table.AddRow("WRITE ACCESS", *user.Write)

		err := table.Display(cmd)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	default:
		details, err := json.MarshalIndent(user, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal Secrets Manager user: %w", err)
		}
		cmd.Println(string(details))

		return nil
	}
}
