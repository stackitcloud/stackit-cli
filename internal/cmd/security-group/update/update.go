package update

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Labels          *map[string]string
	Description     *string
	Name            *string
	SecurityGroupId string
}

const groupNameArg = "GROUP_ID"

const (
	nameArg        = "name"
	descriptionArg = "description"
	labelsArg      = "labels"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", groupNameArg),
		Short: "Updates a security group",
		Long:  "Updates a named security group",
		Args:  args.SingleArg(groupNameArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(`Update the name of group "xxx"`, `$ stackit security-group update xxx --name my-new-name`),
			examples.NewExample(`Update the labels of group "xxx"`, `$ stackit security-group update xxx --labels label1=value1,label2=value2`),
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
			}

			groupLabel, err := iaasUtils.GetSecurityGroupName(ctx, apiClient, model.ProjectId, model.SecurityGroupId)
			if err != nil {
				params.Printer.Warn("cannot retrieve groupname: %v", err)
				groupLabel = model.SecurityGroupId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update the security group %q?", groupLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)

			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update security group: %w", err)
			}
			params.Printer.Info("Updated security group \"%v\" for %q\n", utils.PtrString(resp.Name), projectLabel)

			return nil
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(nameArg, "", "The name of the security group.")
	cmd.Flags().String(descriptionArg, "", "An optional description of the security group.")
	cmd.Flags().StringToString(labelsArg, nil, "Labels are key-value string pairs which can be attached to a network-interface. E.g. '--labels key1=value1,key2=value2,...'")
}

func parseInput(p *print.Printer, cmd *cobra.Command, cliArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Labels:          flags.FlagToStringToStringPointer(p, cmd, labelsArg),
		Description:     flags.FlagToStringPointer(p, cmd, descriptionArg),
		Name:            flags.FlagToStringPointer(p, cmd, nameArg),
		SecurityGroupId: cliArgs[0],
	}

	if model.Labels == nil && model.Description == nil && model.Name == nil {
		return nil, fmt.Errorf("no flags have been passed")
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiUpdateSecurityGroupRequest {
	request := apiClient.UpdateSecurityGroup(ctx, model.ProjectId, model.SecurityGroupId)
	payload := iaas.NewUpdateSecurityGroupPayload()
	payload.Description = model.Description
	payload.Labels = utils.ConvertStringMapToInterfaceMap(model.Labels)
	payload.Name = model.Name
	request = request.UpdateSecurityGroupPayload(*payload)

	return request
}
