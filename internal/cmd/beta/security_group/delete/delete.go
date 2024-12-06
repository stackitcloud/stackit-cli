package delete

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Id string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete a security group",
		Long:  "delete a security group by its internal id",
		Args:  cobra.ExactArgs(1),
		Example: examples.Build(
			examples.NewExample(`delete a named group`, `$ stackit beta security-group delete 43ad419a-c68b-4911-87cd-e05752ac1e31`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeDelete(cmd, p, args)
		},
	}

	return cmd
}

func executeDelete(cmd *cobra.Command, p *print.Printer, args []string) error {
	p.Info("executing delete command")
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

	if !model.AssumeYes {
		prompt := fmt.Sprintf("Are you sure you want to delete the security group %q?", model.Id)
		err = p.PromptForConfirmation(prompt)
		if err != nil {
			return err
		}
	}

	// Call API
	request := buildRequest(ctx, model, apiClient)

	operationState := "Enabled"
	if model.Async {
		operationState = "Triggered security group deletion"
	}
	p.Info("%s security group %q for %q\n", operationState, model.Id, model.ProjectId)

	if err := request.Execute(); err != nil {
		return fmt.Errorf("delete security group: %w", err)
	}

	return nil
}

func parseInput(p *print.Printer, cmd *cobra.Command, args []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	if len(args) != 1 {
		return nil,&errors.ArgValidationError{}
	}

	name := flags.FlagToStringValue(p, cmd, "name")
	if len(name) >= 64 {
		return nil, &errors.ArgValidationError{
			Arg:     "invalid name",
			Details: "name exceeds 63 characters in length",
		}
	}

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
	description := flags.FlagToStringValue(p, cmd, "description")
	if len(description) >= 128 {
		return nil, &errors.ArgValidationError{
			Arg:     "invalid description",
			Details: "description exceeds 127 characters in length",
		}
	}
	model := inputModel{
		GlobalFlagModel: globalFlags,
		Id:              args[0],
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiDeleteSecurityGroupRequest {
	request := apiClient.DeleteSecurityGroup(ctx, model.ProjectId, model.Id)
	return request
}
