package delete

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sqlserverflex/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/sqlserverflex"
)

const (
	databaseNameArg = "DATABASE_NAME"

	instanceIdFlag = "instance-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	DatabaseName string
	InstanceId   string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", databaseNameArg),
		Short: "Deletes a SQLServer Flex database",
		Long: fmt.Sprintf("%s\n%s",
			"Deletes a SQLServer Flex database.",
			`This operation cannot be triggered asynchronously (the "--async" flag will have no effect).`,
		),
		Args: args.SingleArg(databaseNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Delete a SQLServer Flex database with name "my-database" of instance with ID "xxx"`,
				"$ stackit beta sqlserverflex database delete my-database --instance-id xxx"),
		),
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

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete database %q? (This cannot be undone)", model.DatabaseName)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			s := spinner.New(p)
			s.Start("Deleting database")
			err = req.Execute()
			if err != nil {
				s.StopWithError()
				return fmt.Errorf("delete SQLServer Flex database: %w", err)
			}
			s.Stop()

			p.Info("Deleted database %q\n", model.DatabaseName)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "SQLServer Flex instance ID")
	err := flags.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	databaseName := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		DatabaseName:    databaseName,
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *sqlserverflex.APIClient) sqlserverflex.ApiDeleteDatabaseRequest {
	req := apiClient.DeleteDatabase(ctx, model.ProjectId, model.InstanceId, model.DatabaseName, model.Region)
	return req
}
