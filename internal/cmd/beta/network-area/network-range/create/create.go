package create

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
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"

	"github.com/spf13/cobra"
)

const (
	organizationIdFlag = "organization-id"
	networkAreaIdFlag  = "network-area-id"
	networkRangeFlag   = "network-range"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	OrganizationId *string
	NetworkAreaId  *string
	NetworkRange   *string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a network range in a STACKIT Network Area (SNA)",
		Long:  "Creates a network range in a STACKIT Network Area (SNA).",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a network range in a STACKIT Network Area with ID "xxx" in organization with ID "yyy"`,
				`$ stackit beta network-area network-range create --network-area-id xxx --organization-id yyy --network-range "1.1.1.0/24"`,
			),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			// Get network area label
			networkAreaLabel, err := iaasUtils.GetNetworkAreaName(ctx, apiClient, *model.OrganizationId, *model.NetworkAreaId)
			if err != nil {
				p.Debug(print.ErrorLevel, "get network area name: %v", err)
				networkAreaLabel = *model.NetworkAreaId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a network range for STACKIT Network Area (SNA) %q?", networkAreaLabel)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create network range: %w", err)
			}

			if resp.Items == nil || len(*resp.Items) == 0 {
				return fmt.Errorf("empty response from API")
			}

			networkRange, err := iaasUtils.GetNetworkRangeFromAPIResponse(*model.NetworkRange, resp.Items)
			if err != nil {
				return err
			}

			return outputResult(p, model, networkAreaLabel, networkRange)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "Organization ID")
	cmd.Flags().Var(flags.UUIDFlag(), networkAreaIdFlag, "STACKIT Network Area (SNA) ID")
	cmd.Flags().Var(flags.CIDRFlag(), networkRangeFlag, "Network range to create in CIDR notation")

	err := flags.MarkFlagsRequired(cmd, organizationIdFlag, networkAreaIdFlag, networkRangeFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	model := inputModel{
		GlobalFlagModel: globalFlags,
		OrganizationId:  flags.FlagToStringPointer(p, cmd, organizationIdFlag),
		NetworkAreaId:   flags.FlagToStringPointer(p, cmd, networkAreaIdFlag),
		NetworkRange:    flags.FlagToStringPointer(p, cmd, networkRangeFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiCreateNetworkAreaRangeRequest {
	req := apiClient.CreateNetworkAreaRange(ctx, *model.OrganizationId, *model.NetworkAreaId)
	payload := iaas.CreateNetworkAreaRangePayload{
		Ipv4: &[]iaas.NetworkRange{
			{
				Prefix: model.NetworkRange,
			},
		},
	}
	return req.CreateNetworkAreaRangePayload(payload)
}

func outputResult(p *print.Printer, model *inputModel, networkAreaLabel string, networkRange iaas.NetworkRange) error {
	switch model.OutputFormat {
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
		p.Outputf("Created network range for SNA %q.\nNetwork range ID: %s\n", networkAreaLabel, utils.PtrString(networkRange.NetworkRangeId))
		return nil
	}
}
