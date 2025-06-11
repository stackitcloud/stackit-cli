package update

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	nicIdArg = "NIC_ID"

	networkIdFlag        = "network-id"
	allowedAddressesFlag = "allowed-addresses"
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
	NicId            string
	NetworkId        *string
	AllowedAddresses *[]iaas.AllowedAddressesInner
	Labels           *map[string]string
	Name             *string // <= 63 characters + regex  ^[A-Za-z0-9]+((-|_|\s|\.)[A-Za-z0-9]+)*$
	NicSecurity      *bool
	SecurityGroups   *[]string // = 36 characters + regex ^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", nicIdArg),
		Short: "Updates a network interface",
		Long:  "Updates a network interface.",
		Args:  args.SingleArg(nicIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Updates a network interface with nic id "xxx" and network-id "yyy" to new allowed addresses "1.1.1.1,8.8.8.8,9.9.9.9" and new labels "key=value,key2=value2"`,
				`$ stackit network-interface update xxx --network-id yyy --allowed-addresses "1.1.1.1,8.8.8.8,9.9.9.9" --labels key=value,key2=value2`,
			),
			examples.NewExample(
				`Updates a network interface with nic id "xxx" and network-id "yyy" with new name "nic-name-new"`,
				`$ stackit network-interface update xxx --network-id yyy --name nic-name-new`,
			),
			examples.NewExample(
				`Updates a network interface with nic id "xxx" and network-id "yyy" to include the security group "zzz"`,
				`$ stackit network-interface update xxx --network-id yyy --security-groups zzz`,
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

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update the network interface %q?", model.NicId)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update network interface: %w", err)
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
	cmd.Flags().StringToString(labelFlag, nil, "Labels are key-value string pairs which can be attached to a network-interface. E.g. '--labels key1=value1,key2=value2,...'")
	cmd.Flags().StringP(nameFlag, "n", "", "Network interface name")
	cmd.Flags().BoolP(nicSecurityFlag, "b", defaultNicSecurityFlag, "If this is set to false, then no security groups will apply to this network interface.")
	cmd.Flags().StringSlice(securityGroupsFlag, nil, "List of security groups")

	err := flags.MarkFlagsRequired(cmd, networkIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	nicId := inputArgs[0]
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

	// check name length and regex must apply
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
		NicId:           nicId,
		NetworkId:       flags.FlagToStringPointer(p, cmd, networkIdFlag),
		Labels:          flags.FlagToStringToStringPointer(p, cmd, labelFlag),
		Name:            name,
		NicSecurity:     flags.FlagToBoolPointer(p, cmd, nicSecurityFlag),
		SecurityGroups:  securityGroups,
	}

	if allowedAddresses != nil {
		model.AllowedAddresses = utils.Ptr(allowedAddressesInner)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiUpdateNicRequest {
	req := apiClient.UpdateNic(ctx, model.ProjectId, *model.NetworkId, model.NicId)

	payload := iaas.UpdateNicPayload{
		AllowedAddresses: model.AllowedAddresses,
		Labels:           utils.ConvertStringMapToInterfaceMap(model.Labels),
		Name:             model.Name,
		NicSecurity:      model.NicSecurity,
		SecurityGroups:   model.SecurityGroups,
	}
	return req.UpdateNicPayload(payload)
}

func outputResult(p *print.Printer, outputFormat, projectId string, nic *iaas.NIC) error {
	if nic == nil {
		return fmt.Errorf("nic is empty")
	}
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(nic, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal network interface: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(nic, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal network interface: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		p.Outputf("Updated network interface for project %q.\n", projectId)
		return nil
	}
}
