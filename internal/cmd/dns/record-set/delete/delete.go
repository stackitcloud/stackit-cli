package delete

import (
	"context"
	"fmt"

	"stackit/internal/pkg/args"
	"stackit/internal/pkg/confirm"
	"stackit/internal/pkg/errors"
	"stackit/internal/pkg/examples"
	"stackit/internal/pkg/flags"
	"stackit/internal/pkg/globalflags"
	"stackit/internal/pkg/services/dns/client"
	dnsUtils "stackit/internal/pkg/services/dns/utils"
	"stackit/internal/pkg/spinner"
	"stackit/internal/pkg/utils"

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

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", recordSetIdArg),
		Short: "Delete a DNS record set",
		Long:  "Delete a DNS record set",
		Args:  args.SingleArg(recordSetIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete DNS record set with ID "xxx" in zone with ID "yyy"`,
				"$ stackit dns record-set delete xxx --zone-id yyy"),
		),
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

			zoneLabel, err := dnsUtils.GetZoneName(ctx, apiClient, model.ProjectId, model.ZoneId)
			if err != nil {
				zoneLabel = model.ZoneId
			}

			recordSetLabel, err := dnsUtils.GetRecordSetName(ctx, apiClient, model.ProjectId, model.ZoneId, model.RecordSetId)
			if err != nil {
				recordSetLabel = model.RecordSetId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete record set %s of zone %s? (This cannot be undone)", recordSetLabel, zoneLabel)
				err = confirm.PromptForConfirmation(cmd, prompt)
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
				s := spinner.New(cmd)
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
			cmd.Printf("%s record set %s of zone %s\n", operationState, recordSetLabel, zoneLabel)
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

func parseInput(cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	recordSetId := inputArgs[0]

	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		ZoneId:          flags.FlagToStringValue(cmd, zoneIdFlag),
		RecordSetId:     recordSetId,
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *dns.APIClient) dns.ApiDeleteRecordSetRequest {
	req := apiClient.DeleteRecordSet(ctx, model.ProjectId, model.ZoneId, model.RecordSetId)
	return req
}
