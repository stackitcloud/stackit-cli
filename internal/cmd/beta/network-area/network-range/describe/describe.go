package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"

	"github.com/spf13/cobra"
)

const (
	networkRangeIdArg = "NETWORK_RANGE_ID"

	organizationIdFlag = "organization-id"
	networkAreaIdFlag  = "network-area-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	OrganizationId *string
	NetworkAreaId  *string
	NetworkRangeId string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", networkRangeIdArg),
		Short: "Shows details of a network range in a STACKIT Network Area (SNA)",
		Long:  "Shows details of a network range in a STACKIT Network Area (SNA).",
		Args:  args.SingleArg(networkRangeIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Show details of a network range with ID "xxx" in a STACKIT Network Area with ID "yyy" in organization with ID "zzz"`,
				`$ stackit beta network-area network-range describe xxx --network-area-id yyy --organization-id zzz`,
			),
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("describe network range: %w", err)
			}

			return outputResult(p, model.OutputFormat, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "Organization ID")
	cmd.Flags().Var(flags.UUIDFlag(), networkAreaIdFlag, "STACKIT Network Area (SNA) ID")

	err := flags.MarkFlagsRequired(cmd, organizationIdFlag, networkAreaIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	networkRangeId := inputArgs[0]
	globalFlags := globalflags.Parse(p, cmd)

	model := inputModel{
		GlobalFlagModel: globalFlags,
		OrganizationId:  flags.FlagToStringPointer(p, cmd, organizationIdFlag),
		NetworkAreaId:   flags.FlagToStringPointer(p, cmd, networkAreaIdFlag),
		NetworkRangeId:  networkRangeId,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiGetNetworkAreaRangeRequest {
	req := apiClient.GetNetworkAreaRange(ctx, *model.OrganizationId, *model.NetworkAreaId, model.NetworkRangeId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, networkRange *iaas.NetworkRange) error {
	if networkRange == nil {
		return fmt.Errorf("network range is nil")
	}
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(networkRange, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal network range: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(networkRange, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal network range: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.AddRow("ID", utils.PtrString(networkRange.NetworkRangeId))
		table.AddSeparator()
		table.AddRow("Network range", utils.PtrString(networkRange.Prefix))

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	}
}
