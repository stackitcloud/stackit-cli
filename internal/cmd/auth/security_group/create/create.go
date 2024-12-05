package create

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
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
	Labels      map[string]any
	Description string
	Name        string
	Stateful    bool
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create security groups",
		Long:  "create security groups",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(`create a named group`, `$ stackit beta security-group create --name my-new-group`),
			examples.NewExample(`create a named group with labels`, `$ stackit beta security-group create --name my-new-group --labels label1=value1,label2=value2`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeCreate(cmd, p, args)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "the name of the security group. Must be <= 63 chars")
	cmd.Flags().String("description", "", "an optional description of the security group. Must be <= 127 chars")
	cmd.Flags().Bool("stateful", false, "create a stateful or a stateless security group")
	cmd.Flags().StringSlice("labels", nil, "a list of labels in the form <key>=<value>")

	if err := flags.MarkFlagsRequired(cmd, "name"); err != nil {
		cobra.CheckErr(err)
	}
}

func executeCreate(cmd *cobra.Command, p *print.Printer, _ []string) error {
	p.Info("executing create command")
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

	if !model.AssumeYes {
		prompt := fmt.Sprintf("Are you sure you want to create the security group %q?", model.Name)
		err = p.PromptForConfirmation(prompt)
		if err != nil {
			return err
		}
	}

	// Call API
	request := buildRequest(ctx, model, apiClient)
	_, err = request.Execute()
	if err != nil {
		return fmt.Errorf("create security group: %w", err)
	}

	operationState := "Enabled"
	if model.Async {
		operationState = "Triggered label creation"
	}
	p.Info("%s security group %q for %q\n", operationState, model.Name, model.ProjectId)

	group, err := request.Execute()
	if err != nil {
		return fmt.Errorf("create security group: %w", err)
	}
	if err:=outputResult(p, model, group);err != nil {
		return err
	}

	return nil
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
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
		Name:            name,

		Labels:      labels,
		Description: description,
		Stateful:    flags.FlagToBoolValue(p, cmd, "stateful"),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiCreateSecurityGroupRequest {
	request := apiClient.CreateSecurityGroup(ctx, model.ProjectId)
	payload := iaas.NewCreateSecurityGroupPayload(&model.Name)
	payload.Description = &model.Description
	if model.Labels != nil {
		// this check assure that we don't end up with a pointer to nil
		// which is a thing in go!
		payload.Labels = &model.Labels
	}
	payload.Name = &model.Name
	payload.Stateful = &model.Stateful
	request = request.CreateSecurityGroupPayload(*payload)

	return request

}

func outputResult(p *print.Printer, model *inputModel, resp *iaas.SecurityGroup) error {
	switch model.OutputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal security group: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(resp, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal security group: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		operationState := "Created"
		if model.Async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s security group %q\n", operationState, model.Name)
		return nil
	}
}
