package delete

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/dns/client"
	dnsUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/dns/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/dns"
	"github.com/stackitcloud/stackit-sdk-go/services/dns/wait"
)

const (
	recordSetIdArg = "RECORD_SET_ID"

	zoneIdFlag = "zone-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ZoneId      string
	RecordSetId string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", recordSetIdArg),
		Short: "Deletes a DNS record set",
		Long:  "Deletes a DNS record set.",
		Args:  args.SingleArg(recordSetIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete DNS record set with ID "xxx" in zone with ID "yyy"`,
				"$ stackit dns record-set delete xxx --zone-id yyy"),
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

			zoneLabel, err := dnsUtils.GetZoneName(ctx, apiClient, model.ProjectId, model.ZoneId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get zone name: %v", err)
				zoneLabel = model.ZoneId
			}

			recordSetLabel, err := dnsUtils.GetRecordSetName(ctx, apiClient, model.ProjectId, model.ZoneId, model.RecordSetId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get record set name: %v", err)
				recordSetLabel = model.RecordSetId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete record set %s of zone %s? (This cannot be undone)", recordSetLabel, zoneLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			if err != nil {
				return err
			}
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete DNS record set: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Deleting record set")
				_, err = wait.DeleteRecordSetWaitHandler(ctx, apiClient, model.ProjectId, model.ZoneId, model.RecordSetId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for DNS record set deletion: %w", err)
				}
				s.Stop()
			}

			operationState := "Deleted"
			if model.Async {
				operationState = "Triggered deletion of"
			}
			params.Printer.Info("%s record set %s of zone %s\n", operationState, recordSetLabel, zoneLabel)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), zoneIdFlag, "Zone ID")

	err := flags.MarkFlagsRequired(cmd, zoneIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	recordSetId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ZoneId:          flags.FlagToStringValue(p, cmd, zoneIdFlag),
		RecordSetId:     recordSetId,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *dns.APIClient) dns.ApiDeleteRecordSetRequest {
	req := apiClient.DeleteRecordSet(ctx, model.ProjectId, model.ZoneId, model.RecordSetId)
	return req
}
