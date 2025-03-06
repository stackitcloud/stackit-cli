package attach

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
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
	createFlag             = "create"
	networkIdFlag          = "network-id"

	defaultCreateFlag = false
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ServerId  *string
	NicId     *string
	NetworkId *string
	Create    *bool
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attach",
		Short: "Attaches a network interface to a server",
		Long:  "Attaches a network interface to a server.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Attach a network interface with ID "xxx" to a server with ID "yyy"`,
				`$ stackit server network-interface attach --network-interface-id xxx --server-id yyy`,
			),
			examples.NewExample(
				`Create a network interface for network with ID "xxx" and attach it to a server with ID "yyy"`,
				`$ stackit server network-interface attach --network-id xxx --server-id yyy --create`,
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

			serverLabel, err := iaasUtils.GetServerName(ctx, apiClient, model.ProjectId, *model.ServerId)
			if err != nil {
				p.Debug(print.ErrorLevel, "get server name: %v", err)
				serverLabel = *model.ServerId
			}

			// if the create flag is provided a network interface will be created and attached
			if model.Create != nil && *model.Create {
				networkLabel, err := iaasUtils.GetNetworkName(ctx, apiClient, model.ProjectId, *model.NetworkId)
				if err != nil {
					p.Debug(print.ErrorLevel, "get network name: %v", err)
					networkLabel = *model.NetworkId
				}
				if !model.AssumeYes {
					prompt := fmt.Sprintf("Are you sure you want to create a network interface for network %q and attach it to server %q?", networkLabel, serverLabel)
					err = p.PromptForConfirmation(prompt)
					if err != nil {
						return err
					}
				}
				// Call API
				req := buildRequestCreateAndAttach(ctx, model, apiClient)
				err = req.Execute()
				if err != nil {
					return fmt.Errorf("create and attach network interface: %w", err)
				}
				p.Info("Created a network interface for network %q and attached it to server %q\n", networkLabel, serverLabel)
				return nil
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to attach network interface %q to server %q?", *model.NicId, serverLabel)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}
			// Call API
			req := buildRequestAttach(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("attach network interface: %w", err)
			}
			p.Info("Attached network interface %q to server %q\n", utils.PtrString(model.NicId), serverLabel)

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
	cmd.Flags().BoolP(createFlag, "b", defaultCreateFlag, "If this is set a network interface will be created. (default false)")

	cmd.MarkFlagsRequiredTogether(createFlag, networkIdFlag)
	cmd.MarkFlagsMutuallyExclusive(createFlag, networkInterfaceIdFlag)
	cmd.MarkFlagsMutuallyExclusive(networkIdFlag, networkInterfaceIdFlag)

	err := flags.MarkFlagsRequired(cmd, serverIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	// if create is not provided then network-interface-id is needed
	networkInterfaceId := flags.FlagToStringPointer(p, cmd, networkInterfaceIdFlag)
	create := flags.FlagToBoolPointer(p, cmd, createFlag)
	if create == nil && networkInterfaceId == nil {
		return nil, &cliErr.ServerNicAttachMissingNicIdError{Cmd: cmd}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ServerId:        flags.FlagToStringPointer(p, cmd, serverIdFlag),
		NetworkId:       flags.FlagToStringPointer(p, cmd, networkIdFlag),
		NicId:           networkInterfaceId,
		Create:          create,
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

func buildRequestAttach(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiAddNicToServerRequest {
	return apiClient.AddNicToServer(ctx, model.ProjectId, *model.ServerId, *model.NicId)
}

func buildRequestCreateAndAttach(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiAddNetworkToServerRequest {
	return apiClient.AddNetworkToServer(ctx, model.ProjectId, *model.ServerId, *model.NetworkId)
}
