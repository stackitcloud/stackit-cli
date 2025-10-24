package create

import (
	"context"
	"fmt"

	iaasClient "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/runcommand/client"
	runcommandUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/runcommand/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/runcommand"
)

const (
	serverIdFlag            = "server-id"
	commandTemplateNameFlag = "template-name"
	paramsFlag              = "params"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	ServerId            string
	CommandTemplateName string
	Params              *map[string]string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a Server Command",
		Long:  "Creates a Server Command.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a server command for server with ID "xxx", template name "RunShellScript" and a script from a file (using the @{...} format)`,
				`$ stackit server command create --server-id xxx --template-name=RunShellScript --params script='@{/path/to/script.sh}'`),
			examples.NewExample(
				`Create a server command for server with ID "xxx", template name "RunShellScript" and a script provided on the command line`,
				`$ stackit server command create --server-id xxx --template-name=RunShellScript --params script='echo hello'`),
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

			serverLabel := model.ServerId
			// Get server name
			if iaasApiClient, err := iaasClient.ConfigureClient(params.Printer, params.CliVersion); err == nil {
				serverName, err := iaasUtils.GetServerName(ctx, iaasApiClient, model.ProjectId, model.Region, model.ServerId)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get server name: %v", err)
				} else if serverName != "" {
					serverLabel = serverName
				}
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a Command for server %s?", serverLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req, err := buildRequest(ctx, model, apiClient)
			if err != nil {
				return err
			}
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create Server Command: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, serverLabel, *resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().VarP(flags.UUIDFlag(), serverIdFlag, "s", "Server ID")
	cmd.Flags().StringP(commandTemplateNameFlag, "n", "", "Template name")
	cmd.Flags().StringToStringP(paramsFlag, "r", nil, "Params can be provided with the format key=value and the flag can be used multiple times to provide a list of labels")

	err := flags.MarkFlagsRequired(cmd, serverIdFlag, commandTemplateNameFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:     globalFlags,
		ServerId:            flags.FlagToStringValue(p, cmd, serverIdFlag),
		CommandTemplateName: flags.FlagToStringValue(p, cmd, commandTemplateNameFlag),
		Params:              flags.FlagToStringToStringPointer(p, cmd, paramsFlag),
	}
	parsedParams, err := runcommandUtils.ParseScriptParams(*model.Params)
	if err != nil {
		return nil, &cliErr.FlagValidationError{
			Flag:    paramsFlag,
			Details: err.Error(),
		}
	}
	model.Params = &parsedParams

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *runcommand.APIClient) (runcommand.ApiCreateCommandRequest, error) {
	req := apiClient.CreateCommand(ctx, model.ProjectId, model.ServerId, model.Region)
	req = req.CreateCommandPayload(runcommand.CreateCommandPayload{
		CommandTemplateName: &model.CommandTemplateName,
		Parameters:          model.Params,
	})
	return req, nil
}

func outputResult(p *print.Printer, outputFormat, serverLabel string, resp runcommand.NewCommandResponse) error {
	return p.OutputResult(outputFormat, resp, func() error {
		p.Outputf("Created server command for server %s. Command ID: %s\n", serverLabel, utils.PtrString(resp.Id))
		return nil
	})
}
