package describe

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	nicIdArg = "NIC_ID"

	networkIdFlag = "network-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	NetworkId string
	NicId     string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", nicIdArg),
		Short: "Describes a network interface",
		Long:  "Describes a network interface.",
		Args:  args.SingleArg(nicIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Describes network interface with nic id "xxx" and network ID "yyy"`,
				`$ stackit network-interface describe xxx --network-id yyy`,
			),
			examples.NewExample(
				`Describes network interface with nic id "xxx" and network ID "yyy" in JSON format`,
				`$ stackit network-interface describe xxx --network-id yyy --output-format json`,
			),
			examples.NewExample(
				`Describes network interface with nic id "xxx" and network ID "yyy" in yaml format`,
				`$ stackit network-interface describe xxx --network-id yyy --output-format yaml`,
			),
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
				return fmt.Errorf("describe network interface: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), networkIdFlag, "Network ID")

	err := flags.MarkFlagsRequired(cmd, networkIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	nicId := inputArgs[0]
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		NetworkId:       flags.FlagToStringValue(p, cmd, networkIdFlag),
		NicId:           nicId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiGetNicRequest {
	req := apiClient.GetNic(ctx, model.ProjectId, model.Region, model.NetworkId, model.NicId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, nic *iaas.NIC) error {
	if nic == nil {
		return fmt.Errorf("nic is empty")
	}
	return p.OutputResult(outputFormat, nic, func() error {
		table := tables.NewTable()
		table.AddRow("ID", utils.PtrString(nic.Id))
		table.AddSeparator()
		table.AddRow("NETWORK ID", utils.PtrString(nic.NetworkId))
		table.AddSeparator()
		if nic.Name != nil {
			table.AddRow("NAME", utils.PtrString(nic.Name))
			table.AddSeparator()
		}
		if nic.Ipv4 != nil {
			table.AddRow("IPV4", utils.PtrString(nic.Ipv4))
			table.AddSeparator()
		}
		if nic.Ipv6 != nil {
			table.AddRow("IPV6", utils.PtrString(nic.Ipv6))
			table.AddSeparator()
		}
		table.AddRow("MAC", utils.PtrString(nic.Mac))
		table.AddSeparator()
		table.AddRow("NIC SECURITY", utils.PtrString(nic.NicSecurity))
		if nic.AllowedAddresses != nil && len(*nic.AllowedAddresses) > 0 {
			var allowedAddresses []string
			for _, value := range *nic.AllowedAddresses {
				allowedAddresses = append(allowedAddresses, *value.String)
			}
			table.AddSeparator()
			table.AddRow("ALLOWED ADDRESSES", strings.Join(allowedAddresses, "\n"))
		}
		if nic.Labels != nil && len(*nic.Labels) > 0 {
			var labels []string
			for key, value := range *nic.Labels {
				labels = append(labels, fmt.Sprintf("%s: %s", key, value))
			}
			table.AddSeparator()
			table.AddRow("LABELS", strings.Join(labels, "\n"))
		}
		table.AddSeparator()
		table.AddRow("STATUS", utils.PtrString(nic.Status))
		table.AddSeparator()
		table.AddRow("TYPE", utils.PtrString(nic.Type))
		if nic.SecurityGroups != nil && len(*nic.SecurityGroups) > 0 {
			table.AddSeparator()
			table.AddRow("SECURITY GROUPS", strings.Join(*nic.SecurityGroups, "\n"))
		}

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
