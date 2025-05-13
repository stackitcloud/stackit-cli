package update

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/secrets-manager/client"
	secretsManagerUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/secrets-manager/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/secretsmanager"
)

const (
	instanceIdArg = "INSTANCE_ID"

	aclFlag = "acl"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId string

	Acls *[]string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", instanceIdArg),
		Short: "Updates a Secrets Manager instance",
		Long:  "Updates a Secrets Manager instance.",
		Args:  args.SingleArg(instanceIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update the range of IPs allowed to access a Secrets Manager instance with ID "xxx"`,
				"$ stackit secrets-manager instance update xxx --acl 1.2.3.0/24"),
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

			instanceLabel, err := secretsManagerUtils.GetInstanceName(ctx, apiClient, model.ProjectId, model.InstanceId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update instance %q?", instanceLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("update Secrets Manager instance: %w", err)
			}

			params.Printer.Info("Updated instance %q\n", instanceLabel)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.CIDRSliceFlag(), aclFlag, "List of IP networks in CIDR notation which are allowed to access this instance")
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	instanceId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	acls := flags.FlagToStringSlicePointer(p, cmd, aclFlag)

	if acls == nil {
		return nil, &cliErr.EmptyUpdateError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      instanceId,
		Acls:            acls,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *secretsmanager.APIClient) secretsmanager.ApiUpdateACLsRequest {
	req := apiClient.UpdateACLs(ctx, model.ProjectId, model.InstanceId)

	cidrs := []secretsmanager.UpdateACLPayload{}

	for _, acl := range *model.Acls {
		cidrs = append(cidrs, secretsmanager.UpdateACLPayload{Cidr: utils.Ptr(acl)})
	}

	req = req.UpdateACLsPayload(secretsmanager.UpdateACLsPayload{Cidrs: &cidrs})

	return req
}
