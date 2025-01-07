package update

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
	rmClient "github.com/stackitcloud/stackit-cli/internal/pkg/services/resourcemanager/client"
	rmUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/resourcemanager/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"

	"github.com/spf13/cobra"
)

const (
	areaIdArg = "AREA_ID"

	nameFlag                = "name"
	organizationIdFlag      = "organization-id"
	areaIdFlag              = "area-id"
	dnsNameServersFlag      = "dns-name-servers"
	defaultPrefixLengthFlag = "default-prefix-length"
	maxPrefixLengthFlag     = "max-prefix-length"
	minPrefixLengthFlag     = "min-prefix-length"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	AreaId              string
	Name                *string
	OrganizationId      *string
	DnsNameServers      *[]string
	DefaultPrefixLength *int64
	MaxPrefixLength     *int64
	MinPrefixLength     *int64
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", areaIdArg),
		Short: "Updates a STACKIT Network Area (SNA)",
		Long:  "Updates a STACKIT Network Area (SNA) in an organization.",
		Args:  args.SingleArg(areaIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update network area with ID "xxx" in organization with ID "yyy" with new name "network-area-1-new"`,
				"$ stackit beta network-area update xxx --organization-id yyy --name network-area-1-new",
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

			var orgLabel string
			rmApiClient, err := rmClient.ConfigureClient(p)
			if err == nil {
				orgLabel, err = rmUtils.GetOrganizationName(ctx, rmApiClient, *model.OrganizationId)
				if err != nil {
					p.Debug(print.ErrorLevel, "get organization name: %v", err)
					orgLabel = *model.OrganizationId
				}
			} else {
				p.Debug(print.ErrorLevel, "configure resource manager client: %v", err)
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update a network area for organization %q?", orgLabel)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update network area: %w", err)
			}

			return outputResult(p, model, orgLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(nameFlag, "n", "", "Network area name")
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "Organization ID")
	cmd.Flags().StringSlice(dnsNameServersFlag, nil, "List of DNS name server IPs")
	cmd.Flags().Int64(defaultPrefixLengthFlag, 0, "The default prefix length for networks in the network area")
	cmd.Flags().Int64(maxPrefixLengthFlag, 0, "The maximum prefix length for networks in the network area")
	cmd.Flags().Int64(minPrefixLengthFlag, 0, "The minimum prefix length for networks in the network area")

	err := flags.MarkFlagsRequired(cmd, organizationIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	areaId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)

	model := inputModel{
		GlobalFlagModel:     globalFlags,
		Name:                flags.FlagToStringPointer(p, cmd, nameFlag),
		OrganizationId:      flags.FlagToStringPointer(p, cmd, organizationIdFlag),
		AreaId:              areaId,
		DnsNameServers:      flags.FlagToStringSlicePointer(p, cmd, dnsNameServersFlag),
		DefaultPrefixLength: flags.FlagToInt64Pointer(p, cmd, defaultPrefixLengthFlag),
		MaxPrefixLength:     flags.FlagToInt64Pointer(p, cmd, maxPrefixLengthFlag),
		MinPrefixLength:     flags.FlagToInt64Pointer(p, cmd, minPrefixLengthFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiPartialUpdateNetworkAreaRequest {
	req := apiClient.PartialUpdateNetworkArea(ctx, *model.OrganizationId, model.AreaId)

	payload := iaas.PartialUpdateNetworkAreaPayload{
		Name: model.Name,
		AddressFamily: &iaas.UpdateAreaAddressFamily{
			Ipv4: &iaas.UpdateAreaIPv4{
				DefaultNameservers: model.DnsNameServers,
				DefaultPrefixLen:   model.DefaultPrefixLength,
				MaxPrefixLen:       model.MaxPrefixLength,
				MinPrefixLen:       model.MinPrefixLength,
			},
		},
	}

	return req.PartialUpdateNetworkAreaPayload(payload)
}

func outputResult(p *print.Printer, model *inputModel, projectLabel string, networkArea *iaas.NetworkArea) error {
	switch model.OutputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(networkArea, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal network area: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(networkArea, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal network area: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		p.Outputf("Updated STACKIT Network Area for project %q.\n", projectLabel)
		return nil
	}
}
