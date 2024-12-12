package create

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
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
	nameFlag        = "name"
	descriptionFlag = "description"
	statefulFlag    = "stateful"
	labelsFlag      = "labels"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Labels      *map[string]string
	Description *string
	Name        *string
	Stateful    *bool
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates security groups",
		Long:  "Creates security groups.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(`Create a named group`, `$ stackit beta security-group create --name my-new-group`),
			examples.NewExample(`Create a named group with labels`, `$ stackit beta security-group create --name my-new-group --labels label1=value1,label2=value2`),
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

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create the security group %q?", *model.Name)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			request := buildRequest(ctx, model, apiClient)

			group, err := request.Execute()
			if err != nil {
				return fmt.Errorf("create security group: %w", err)
			}

			if err := outputResult(p, model, group); err != nil {
				return err
			}

			return nil
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(nameFlag, "", "The name of the security group.")
	cmd.Flags().String(descriptionFlag, "", "An optional description of the security group.")
	cmd.Flags().Bool(statefulFlag, false, "Create a stateful or a stateless security group")
	cmd.Flags().StringToString(labelsFlag, nil, "Labels are key-value string pairs which can be attached to a network-interface. E.g. '--labels key1=value1,key2=value2,...'")

	if err := flags.MarkFlagsRequired(cmd, nameFlag); err != nil {
		cobra.CheckErr(err)
	}
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}
	name := flags.FlagToStringValue(p, cmd, nameFlag)

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Name:            &name,

		Labels:      flags.FlagToStringToStringPointer(p, cmd, labelsFlag),
		Description: flags.FlagToStringPointer(p, cmd, descriptionFlag),
		Stateful:    flags.FlagToBoolPointer(p, cmd, statefulFlag),
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

	var labelsMap *map[string]any
	if model.Labels != nil && len(*model.Labels) > 0 {
		// convert map[string]string to map[string]interface{}
		labelsMap = utils.Ptr(map[string]interface{}{})
		for k, v := range *model.Labels {
			(*labelsMap)[k] = v
		}
	}
	payload := iaas.CreateSecurityGroupPayload{
		Description: model.Description,
		Labels:      labelsMap,
		Name:        model.Name,
		Stateful:    model.Stateful,
	}

	return request.CreateSecurityGroupPayload(payload)
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
		p.Outputf("Created security group %q\n", *model.Name)
		return nil
	}
}
