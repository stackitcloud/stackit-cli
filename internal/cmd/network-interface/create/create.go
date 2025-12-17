package create

import (
	"context"
	"fmt"
	"regexp"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	networkIdFlag        = "network-id"
	allowedAddressesFlag = "allowed-addresses"
	ipv4Flag             = "ipv4"
	ipv6Flag             = "ipv6"
	labelFlag            = "labels"
	nameFlag             = "name"
	securityGroupsFlag   = "security-groups"
	nicSecurityFlag      = "nic-security"

	nameRegex              = `^[A-Za-z0-9]+((-|_|\s|\.)[A-Za-z0-9]+)*$`
	maxNameLength          = 63
	securityGroupsRegex    = `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`
	securityGroupLength    = 36
	defaultNicSecurityFlag = true
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	NetworkId        string
	AllowedAddresses *[]iaas.AllowedAddressesInner
	Ipv4             *string
	Ipv6             *string
	Labels           *map[string]string
	Name             *string // <= 63 characters + regex  ^[A-Za-z0-9]+((-|_|\s|\.)[A-Za-z0-9]+)*$
	NicSecurity      *bool
	SecurityGroups   *[]string // = 36 characters + regex ^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a network interface",
		Long:  "Creates a network interface.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a network interface for network with ID "xxx"`,
				`$ stackit network-interface create --network-id xxx`,
			),
			examples.NewExample(
				`Create a network interface with allowed addresses, labels, a name, security groups and nic security enabled for network with ID "xxx"`,
				`$ stackit network-interface create --network-id xxx --allowed-addresses "1.1.1.1,8.8.8.8,9.9.9.9" --labels key=value,key2=value2 --name NAME --security-groups "UUID1,UUID2" --nic-security`,
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

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			} else if projectLabel == "" {
				projectLabel = model.ProjectId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a network interface for project %q?", projectLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create network interface: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, model.ProjectId, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), networkIdFlag, "Network ID")
	cmd.Flags().StringSlice(allowedAddressesFlag, nil, "List of allowed IPs")
	cmd.Flags().StringP(ipv4Flag, "i", "", "IPv4 address")
	cmd.Flags().StringP(ipv6Flag, "s", "", "IPv6 address")
	cmd.Flags().StringToString(labelFlag, nil, "Labels are key-value string pairs which can be attached to a network-interface. E.g. '--labels key1=value1,key2=value2,...'")
	cmd.Flags().StringP(nameFlag, "n", "", "Network interface name")
	cmd.Flags().BoolP(nicSecurityFlag, "b", defaultNicSecurityFlag, "If this is set to false, then no security groups will apply to this network interface.")
	cmd.Flags().StringSlice(securityGroupsFlag, nil, "List of security groups")

	err := flags.MarkFlagsRequired(cmd, networkIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	allowedAddresses := flags.FlagToStringSlicePointer(p, cmd, allowedAddressesFlag)
	var allowedAddressesInner []iaas.AllowedAddressesInner
	if allowedAddresses != nil && len(*allowedAddresses) > 0 {
		allowedAddressesInner = make([]iaas.AllowedAddressesInner, len(*allowedAddresses))
		for i, address := range *allowedAddresses {
			allowedAddressesInner[i].String = &address
		}
	}

	// check name length <= 63 and regex must apply
	name := flags.FlagToStringPointer(p, cmd, nameFlag)
	if name != nil {
		if len(*name) > maxNameLength {
			return nil, &errors.FlagValidationError{
				Flag:    nameFlag,
				Details: fmt.Sprintf("name %s is too long (maximum length is %d characters)", *name, maxNameLength),
			}
		}
		nameRegex := regexp.MustCompile(nameRegex)
		if !nameRegex.MatchString(*name) {
			return nil, &errors.FlagValidationError{
				Flag:    nameFlag,
				Details: fmt.Sprintf("name %s didn't match the required regex expression %s", *name, nameRegex),
			}
		}
	}

	// check security groups size and regex
	securityGroups := flags.FlagToStringSlicePointer(p, cmd, securityGroupsFlag)
	if securityGroups != nil && len(*securityGroups) > 0 {
		securityGroupsRegex := regexp.MustCompile(securityGroupsRegex)
		// iterate over them
		for _, value := range *securityGroups {
			if len(value) != securityGroupLength {
				return nil, &errors.FlagValidationError{
					Flag:    securityGroupsFlag,
					Details: fmt.Sprintf("security groups uuid %s does not match (must be %d characters long)", value, securityGroupLength),
				}
			}
			if !securityGroupsRegex.MatchString(value) {
				return nil, &errors.FlagValidationError{
					Flag:    securityGroupsFlag,
					Details: fmt.Sprintf("security groups uuid %s didn't match the required regex expression %s", value, securityGroupsRegex),
				}
			}
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		NetworkId:       flags.FlagToStringValue(p, cmd, networkIdFlag),
		Ipv4:            flags.FlagToStringPointer(p, cmd, ipv4Flag),
		Ipv6:            flags.FlagToStringPointer(p, cmd, ipv6Flag),
		Labels:          flags.FlagToStringToStringPointer(p, cmd, labelFlag),
		Name:            name,
		NicSecurity:     flags.FlagToBoolPointer(p, cmd, nicSecurityFlag),
		SecurityGroups:  securityGroups,
	}

	if allowedAddresses != nil {
		model.AllowedAddresses = utils.Ptr(allowedAddressesInner)
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiCreateNicRequest {
	req := apiClient.CreateNic(ctx, model.ProjectId, model.Region, model.NetworkId)

	payload := iaas.CreateNicPayload{
		AllowedAddresses: model.AllowedAddresses,
		Ipv4:             model.Ipv4,
		Ipv6:             model.Ipv6,
		Labels:           utils.ConvertStringMapToInterfaceMap(model.Labels),
		Name:             model.Name,
		NicSecurity:      model.NicSecurity,
		SecurityGroups:   model.SecurityGroups,
	}
	return req.CreateNicPayload(payload)
}

func outputResult(p *print.Printer, outputFormat, projectId string, nic *iaas.NIC) error {
	if nic == nil {
		return fmt.Errorf("nic is empty")
	}
	return p.OutputResult(outputFormat, nic, func() error {
		p.Outputf("Created network interface for project %q.\nNIC ID: %s\n", projectId, utils.PtrString(nic.Id))
		return nil
	})
}
