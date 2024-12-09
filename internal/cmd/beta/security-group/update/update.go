package update

import (
	"context"
	"fmt"
	"strings"

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
	Labels          *map[string]any
	Description     *string
	Name            *string
	SecurityGroupId string
}

const argNameGroupId = "argGroupId"

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "Update",
		Short: "Update a security group",
		Long:  "Update a named security group",
		Args:  args.SingleArg(argNameGroupId, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(`Update the name of a group`, `$ stackit beta security-group update 541d122f-0a5f-4bb0-94b9-b1ccbd7ba776 --name my-new-name`),
			examples.NewExample(`Update the labels of a group`, `$ stackit beta security-group update 541d122f-0a5f-4bb0-94b9-b1ccbd7ba776 --labels label1=value1,label2=value2`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeUpdate(cmd, p, args)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "the name of the security group. Must be <= 63 chars")
	cmd.Flags().String("description", "", "an optional description of the security group. Must be <= 127 chars")
	cmd.Flags().StringSlice("labels", nil, "a list of labels in the form <key>=<value>")
}

func executeUpdate(cmd *cobra.Command, p *print.Printer, args []string) error {
	p.Info("executing update command")
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
		return fmt.Errorf("^date security group: %w", err)
	}

	operationState := "Enabled"
	if model.Async {
		operationState = "Triggered enablement of"
	}
	p.Info("%s security group \"%v\" for %q\n", operationState, model.Name, projectLabel)
	return nil
}

func parseInput(p *print.Printer, cmd *cobra.Command, args []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	if err := cmd.ValidateArgs(args); err != nil {
		return nil, &errors.ArgValidationError{
			Arg:     argNameGroupId,
			Details: fmt.Sprintf("argument validation failed: %v", err),
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
	}
	if len(args) != 1 {
		return nil, &errors.ArgValidationError{
			Arg:     argNameGroupId,
			Details: "wrong number of arguments",
		}
	}
	model.SecurityGroupId = args[0]

	if cmd.Flags().Lookup("name").Changed {
		name := flags.FlagToStringValue(p, cmd, "name")
		if len(name) >= 64 {
			return nil, &errors.ArgValidationError{
				Arg:     "invalid name",
				Details: "name exceeds 63 characters in length",
			}
		}
		model.Name = &name
	}

	if cmd.Flags().Lookup("labels").Changed {
		labels := make(map[string]any)
		for _, label := range flags.FlagToStringSliceValue(p, cmd, "labels") {
			parts := strings.Split(label, "=")
			if len(parts) != 2 {
				return nil, &errors.ArgValidationError{
					Arg:     "labels",
					Details: "invalid label declaration. Must be in the form <key>=<value>",
				}
			}
			labels[parts[0]] = parts[1]
		}
		model.Labels = &labels
	}
	if cmd.Flags().Lookup("description").Changed {
		description := flags.FlagToStringValue(p, cmd, "description")
		if len(description) >= 128 {
			return nil, &errors.ArgValidationError{
				Arg:     "invalid description",
				Details: "description exceeds 127 characters in length",
			}
		}
		model.Description = &description
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
	payload.Labels = model.Labels
	payload.Name = model.Name
	request = request.UpdateSecurityGroupPayload(*payload)

	return request

}
