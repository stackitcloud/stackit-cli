package generatepayload

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	automation "github.com/stackitcloud/stackit-sdk-go/services/automation/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/fileutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/automation/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

const (
	automationIdFlag = "automation-id"
	filePathFlag     = "file-path"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	AutomationId string
	FilePath     *string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-payload",
		Short: "Generates a payload to update volume automation",
		Long: fmt.Sprintf("%s\n%s",
			"Generates a JSON payload with values to be used as --input for volume automation update.",
			"See https://docs.api.stackit.cloud/documentation/automation-service/version/v1#tag/Volume-Automations/operation/PartialUpdateVolumeAutomation for information regarding the input structure.",
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Generate a payload with values of a volume automation, and adapt it with custom values to different configuration options.`,
				"$ stackit beta volume automation generate-payload --automation-id XXX --file-path ./input-payload.json",
				"<Modify payload in file>",
				"$ stackit beta volume automation update xxx --input @./input-payload.json"),
			examples.NewExample(
				`Generate a payload with values of a volume automation, and preview it in the terminal.`,
				"$ stackit beta volume automation generate-payload --automation-id XXX"),
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
			resp, err := buildRequest(ctx, model, apiClient).Execute()
			if err != nil {
				return fmt.Errorf("get volume automation: %w", err)
			}

			return outputResult(params.Printer, model.FilePath, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), automationIdFlag, "Automation ID from which the input values should be used")
	cmd.Flags().StringP(filePathFlag, "f", "", "If set, writes the payload to the given file. If unset, writes the payload to the standard output")

	err := flags.MarkFlagsRequired(cmd, automationIdFlag)
	cobra.CheckErr(err)
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *automation.APIClient) automation.ApiGetVolumeAutomationRequest {
	return apiClient.DefaultAPI.GetVolumeAutomation(ctx, model.ProjectId, model.Region, model.AutomationId)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		AutomationId:    flags.FlagToStringValue(p, cmd, automationIdFlag),
		FilePath:        flags.FlagToStringPointer(p, cmd, filePathFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func outputResult(p *print.Printer, filePath *string, item *automation.VolumeAutomation) error {
	if item == nil {
		return fmt.Errorf("payload is nil")
	}

	if !item.Input.IsSet() || item.Input.Get() == nil {
		return fmt.Errorf("volume automation %q does not have an input configuration", item.Id)
	}

	inputBytes, err := json.MarshalIndent(item.Input, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	if filePath != nil {
		err = fileutils.WriteToFile(*filePath, string(inputBytes))
		if err != nil {
			return fmt.Errorf("write payload to the file: %w", err)
		}
	} else {
		p.Outputln(string(inputBytes))
	}

	return nil
}
