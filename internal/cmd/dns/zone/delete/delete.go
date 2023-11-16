package delete

import (
	"context"
	"fmt"

	"stackit/internal/pkg/args"
	"stackit/internal/pkg/confirm"
	"stackit/internal/pkg/errors"
	"stackit/internal/pkg/examples"
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
	zoneIdArg = "ZONE_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ZoneId string
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", zoneIdArg),
		Short: "Delete a DNS zone",
		Long:  "Delete a DNS zone",
		Args:  args.SingleArg(zoneIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete a DNS zone with ID "xxx"`,
				"$ stackit dns zone delete xxx"),
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
			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete zone %s? (This cannot be undone)", zoneLabel)
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
				return fmt.Errorf("delete DNS zone: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(cmd)
				s.Start("Deleting zone")
				_, err = wait.DeleteZoneWaitHandler(ctx, apiClient, model.ProjectId, model.ZoneId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for DNS zone deletion: %w", err)
				}
				s.Stop()
			}

			operationState := "Deleted"
			if model.Async {
				operationState = "Triggered deletion of"
			}
			cmd.Printf("%s zone %s\n", operationState, zoneLabel)
			return nil
		},
	}
	return cmd
}

func parseInput(cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	zoneId := inputArgs[0]
	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		ZoneId:          zoneId,
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *dns.APIClient) dns.ApiDeleteZoneRequest {
	req := apiClient.DeleteZone(ctx, model.ProjectId, model.ZoneId)
	return req
}
