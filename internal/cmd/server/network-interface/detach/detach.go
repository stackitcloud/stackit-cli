package detach

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	serverIdFlag           = "server-id"
	networkInterfaceIdFlag = "network-interface-id"
	networkIdFlag          = "network-id"
	deleteFlag             = "delete"

	defaultDeleteFlag = false
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ServerId  *string
	NicId     *string
	NetworkId *string
	Delete    *bool
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detach",
		Short: "Detaches a network interface from a server",
		Long:  "Detaches a network interface from a server.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Detach a network interface with ID "xxx" from a server with ID "yyy"`,
				`$ stackit server network-interface detach --network-interface-id xxx --server-id yyy`,
			),
			examples.NewExample(
				`Detach and delete all network interfaces for network with ID "xxx" and detach them from a server with ID "yyy"`,
				`$ stackit server network-interface detach --network-id xxx --server-id yyy --delete`,
			),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer)
			if err != nil {
				return err
			}

			serverLabel, err := iaasUtils.GetServerName(ctx, apiClient, model.ProjectId, *model.ServerId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get server name: %v", err)
				serverLabel = *model.ServerId
			} else if serverLabel == "" {
				serverLabel = *model.ServerId
			}

			// if the delete flag is provided a network interface is detached and deleted
			if model.Delete != nil && *model.Delete {
				networkLabel, err := iaasUtils.GetNetworkName(ctx, apiClient, model.ProjectId, *model.NetworkId)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get network name: %v", err)
					networkLabel = *model.NetworkId
				}
				if !model.AssumeYes {
					prompt := fmt.Sprintf("Are you sure you want to detach and delete all network interfaces of network %q from server %q? (This cannot be undone)", networkLabel, serverLabel)
					err = params.Printer.PromptForConfirmation(prompt)
					if err != nil {
						return err
					}
				}
				// Call API
				req := buildRequestDetachAndDelete(ctx, model, apiClient)
				err = req.Execute()
				if err != nil {
					return fmt.Errorf("detach and delete network interfaces: %w", err)
				}
				params.Printer.Info("Detached and deleted all network interfaces of network %q from server %q\n", networkLabel, serverLabel)
				return nil
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to detach network interface %q from server %q?", *model.NicId, serverLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}
			// Call API
			req := buildRequestDetach(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("detach network interface: %w", err)
			}
			params.Printer.Info("Detached network interface %q from server %q\n", utils.PtrString(model.NicId), serverLabel)

			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), serverIdFlag, "Server ID")
	cmd.Flags().Var(flags.UUIDFlag(), networkInterfaceIdFlag, "Network Interface ID")
	cmd.Flags().Var(flags.UUIDFlag(), networkIdFlag, "Network ID")
	cmd.Flags().BoolP(deleteFlag, "b", defaultDeleteFlag, "If this is set all network interfaces will be deleted. (default false)")

	cmd.MarkFlagsRequiredTogether(deleteFlag, networkIdFlag)
	cmd.MarkFlagsMutuallyExclusive(deleteFlag, networkInterfaceIdFlag)
	cmd.MarkFlagsMutuallyExclusive(networkIdFlag, networkInterfaceIdFlag)

	err := flags.MarkFlagsRequired(cmd, serverIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	// if delete is not provided then network-interface-id is needed
	networkInterfaceId := flags.FlagToStringPointer(p, cmd, networkInterfaceIdFlag)
	deleteValue := flags.FlagToBoolPointer(p, cmd, deleteFlag)
	if deleteValue == nil && networkInterfaceId == nil {
		return nil, &cliErr.ServerNicDetachMissingNicIdError{Cmd: cmd}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ServerId:        flags.FlagToStringPointer(p, cmd, serverIdFlag),
		NetworkId:       flags.FlagToStringPointer(p, cmd, networkIdFlag),
		NicId:           networkInterfaceId,
		Delete:          deleteValue,
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

func buildRequestDetach(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiRemoveNicFromServerRequest {
	return apiClient.RemoveNicFromServer(ctx, model.ProjectId, *model.ServerId, *model.NicId)
}

func buildRequestDetachAndDelete(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiRemoveNetworkFromServerRequest {
	return apiClient.RemoveNetworkFromServer(ctx, model.ProjectId, *model.ServerId, *model.NetworkId)
}
