package create

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	associatedResourceIdFlag = "associated-resource-id"
	labelFlag                = "labels"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	AssociatedResourceId *string
	Labels               *map[string]string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a Public IP",
		Long:  "Creates a Public IP.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a public IP`,
				`$ stackit beta public-ip create`,
			),
			examples.NewExample(
				`Create a public IP with associated resource ID "xxx"`,
				`$ stackit beta public-ip create --associated-resource-id xxx`,
			),
			examples.NewExample(
				`Create a public IP with associated resource ID "xxx" and labels`,
				`$ stackit beta public-ip create --associated-resource-id xxx --labels key=value,foo=bar`,
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

			projectLabel, err := projectname.GetProjectName(ctx, p, cmd)
			if err != nil {
				p.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a public IP for project %q?", projectLabel)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create public IP: %w", err)
			}

			return outputResult(p, model, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), associatedResourceIdFlag, "Associates the public IP with a network interface or virtual IP (ID)")
	cmd.Flags().StringToString(labelFlag, nil, "Labels are key-value string pairs which can be attached to a public IP. E.g. '--labels key1=value1,key2=value2,...'")
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:      globalFlags,
		AssociatedResourceId: flags.FlagToStringPointer(p, cmd, associatedResourceIdFlag),
		Labels:               flags.FlagToStringToStringPointer(p, cmd, labelFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiCreatePublicIPRequest {
	req := apiClient.CreatePublicIP(ctx, model.ProjectId)

	var labelsMap *map[string]interface{}
	if model.Labels != nil && len(*model.Labels) > 0 {
		// convert map[string]string to map[string]interface{}
		labelsMap = utils.Ptr(map[string]interface{}{})
		for k, v := range *model.Labels {
			(*labelsMap)[k] = v
		}
	}

	payload := iaas.CreatePublicIPPayload{
		NetworkInterface: iaas.NewNullableString(model.AssociatedResourceId),
		Labels:           labelsMap,
	}

	return req.CreatePublicIPPayload(payload)
}

func outputResult(p *print.Printer, model *inputModel, projectLabel string, publicIp *iaas.PublicIp) error {
	switch model.OutputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(publicIp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal public IP: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(publicIp, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal public IP: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		p.Outputf("Created public IP for project %q.\nPublic IP ID: %s\n", projectLabel, *publicIp.Id)
		return nil
	}
}
