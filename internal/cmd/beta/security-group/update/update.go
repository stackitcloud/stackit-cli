package update

import (
	"context"
	"fmt"

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

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Updates a security group",
		Long:  "Updates a named security group",
		Args:  args.SingleArg(groupNameArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(`Update the name of group "xxx"`, `$ stackit beta security-group update xxx --name my-new-name`),
			examples.NewExample(`Update the labels of group "xxx"`, `$ stackit beta security-group update xxx --labels label1=value1,label2=value2`),
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

			projectLabel, err := projectname.GetProjectName(ctx, p, cmd)
			if err != nil {
				p.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)

			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("update security group: %w", err)
			}
			p.Info("Updated security group \"%v\" for %q\n", model.Name, projectLabel)

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
	var labelsMap *map[string]any
	if model.Labels != nil && len(*model.Labels) > 0 {
		// convert map[string]string to map[string]interface{}
		labelsMap = utils.Ptr(map[string]interface{}{})
		for k, v := range *model.Labels {
			(*labelsMap)[k] = v
		}
	}
	payload.Labels = labelsMap
	payload.Name = model.Name
	request = request.UpdateSecurityGroupPayload(*payload)

	return request
}
