package detach

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	serverIdFlag        = "server-id"
	securityGroupIdFlag = "security-group-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ServerId        string
	SecurityGroupId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detach",
		Short: "Detaches a security group from a server",
		Long:  "Detaches a security group from a server.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Detach a security group with ID "xxx" from a server with ID "yyy"`,
				`$ stackit server security-group detach --server-id yyy --security-group-id xxx`,
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

			serverLabel, err := iaasUtils.GetServerName(ctx, apiClient, model.ProjectId, model.Region, model.ServerId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get server name: %v", err)
				serverLabel = model.ServerId
			} else if serverLabel == "" {
				serverLabel = model.ServerId
			}

			securityGroupLabel, err := iaasUtils.GetSecurityGroupName(ctx, apiClient, model.ProjectId, model.Region, model.SecurityGroupId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get security group name: %v", err)
				securityGroupLabel = model.SecurityGroupId
			}

			prompt := fmt.Sprintf("Are you sure you want to detach security group %q from server %q?", securityGroupLabel, serverLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			if err := req.Execute(); err != nil {
				return fmt.Errorf("detach security group from server: %w", err)
			}

			params.Printer.Outputf("Detached security group %q from server %q\n", securityGroupLabel, serverLabel)

			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), serverIdFlag, "Server ID")
	cmd.Flags().Var(flags.UUIDFlag(), securityGroupIdFlag, "Security Group ID")

	err := flags.MarkFlagsRequired(cmd, serverIdFlag, securityGroupIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ServerId:        flags.FlagToStringValue(p, cmd, serverIdFlag),
		SecurityGroupId: flags.FlagToStringValue(p, cmd, securityGroupIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiRemoveSecurityGroupFromServerRequest {
	req := apiClient.RemoveSecurityGroupFromServer(ctx, model.ProjectId, model.Region, model.ServerId, model.SecurityGroupId)
	return req
}
