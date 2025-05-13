package create

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sqlserverflex/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-sdk-go/services/sqlserverflex"

	"github.com/spf13/cobra"
)

const (
	databaseNameArg = "DATABASE_NAME"

	instanceIdFlag = "instance-id"
	ownerFlag      = "owner"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	DatabaseName string
	InstanceId   string
	Owner        string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("create %s", databaseNameArg),
		Short: "Creates a SQLServer Flex database",
		Long: fmt.Sprintf("%s\n%s",
			"Creates a SQLServer Flex database.",
			`This operation cannot be triggered asynchronously (the "--async" flag will have no effect).`,
		),
		Args: args.SingleArg(databaseNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Create a SQLServer Flex database with name "my-database" on instance with ID "xxx"`,
				"$ stackit beta sqlserverflex database create my-database --instance-id xxx --owner some-username"),
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

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create database %q? (This cannot be undone)", model.DatabaseName)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			s := spinner.New(params.Printer)
			s.Start("Creating database")
			resp, err := req.Execute()
			if err != nil {
				s.StopWithError()
				return fmt.Errorf("create SQLServer Flex database: %w", err)
			}
			s.Stop()

			return outputResult(params.Printer, model.OutputFormat, model.DatabaseName, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "SQLServer Flex instance ID")
	cmd.Flags().String(ownerFlag, "", "Username of the owner user")
	err := flags.MarkFlagsRequired(cmd, instanceIdFlag, ownerFlag)
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
		Owner:           flags.FlagToStringValue(p, cmd, ownerFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *sqlserverflex.APIClient) sqlserverflex.ApiCreateDatabaseRequest {
	req := apiClient.CreateDatabase(ctx, model.ProjectId, model.InstanceId, model.Region)
	payload := sqlserverflex.CreateDatabasePayload{
		Name: &model.DatabaseName,
		Options: &sqlserverflex.DatabaseDocumentationCreateDatabaseRequestOptions{
			Owner: &model.Owner,
		},
	}
	req = req.CreateDatabasePayload(payload)
	return req
}

func outputResult(p *print.Printer, outputFormat, databaseName string, resp *sqlserverflex.CreateDatabaseResponse) error {
	if resp == nil {
		return fmt.Errorf("sqlserverflex response is empty")
	}
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal SQLServer Flex database: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(resp, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal SQLServer Flex database: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		p.Outputf("Created database %q\n", databaseName)
		return nil
	}
}
