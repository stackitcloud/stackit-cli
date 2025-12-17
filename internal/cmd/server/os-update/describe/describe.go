package describe

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/serverosupdate/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/serverupdate"
)

const (
	updateIdArg  = "UPDATE_ID"
	serverIdFlag = "server-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ServerId string
	UpdateId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", updateIdArg),
		Short: "Shows details of a Server os-update",
		Long:  "Shows details of a Server os-update.",
		Args:  args.SingleArg(updateIdArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Get details of a Server os-update with id "my-os-update-id"`,
				"$ stackit server os-update describe my-os-update-id"),
			examples.NewExample(
				`Get details of a Server os-update with id "my-os-update-id" in JSON format`,
				"$ stackit server os-update describe my-os-update-id --output-format json"),
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read server os-update: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, *resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().VarP(flags.UUIDFlag(), serverIdFlag, "s", "Server ID")

	err := flags.MarkFlagsRequired(cmd, serverIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	updateId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ServerId:        flags.FlagToStringValue(p, cmd, serverIdFlag),
		UpdateId:        updateId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *serverupdate.APIClient) serverupdate.ApiGetUpdateRequest {
	req := apiClient.GetUpdate(ctx, model.ProjectId, model.ServerId, model.UpdateId, model.Region)
	return req
}

func outputResult(p *print.Printer, outputFormat string, update serverupdate.Update) error {
	return p.OutputResult(outputFormat, update, func() error {
		table := tables.NewTable()
		table.AddRow("ID", utils.PtrString(update.Id))
		table.AddSeparator()
		table.AddRow("STATUS", utils.PtrString(update.Status))
		table.AddSeparator()
		installedUpdates := utils.PtrStringDefault(update.InstalledUpdates, "n/a")
		table.AddRow("INSTALLED UPDATES", installedUpdates)
		table.AddSeparator()
		failedUpdates := utils.PtrStringDefault(update.FailedUpdates, "n/a")
		table.AddRow("FAILED UPDATES", failedUpdates)

		table.AddRow("START DATE", utils.PtrString(update.StartDate))
		table.AddSeparator()
		table.AddRow("END DATE", utils.PtrString(update.EndDate))
		table.AddSeparator()

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
