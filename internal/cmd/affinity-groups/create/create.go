package create

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
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
	nameFlag   = "name"
	policyFlag = "policy"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Name   string
	Policy string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates an affinity groups",
		Long:  `Creates an affinity groups.`,
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create an affinity group with name "AFFINITY_GROUP_NAME" and policy "soft-affinity"`,
				"$ stackit affinity-group create --name AFFINITY_GROUP_NAME --policy soft-affinity",
			),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer)
			if err != nil {
				return err
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create the affinity group %q?", model.Name)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			request := buildRequest(ctx, *model, apiClient)

			result, err := request.Execute()
			if err != nil {
				return fmt.Errorf("create affinity group: %w", err)
			}
			if resp := result; resp != nil {
				return outputResult(params.Printer, *model, *resp)
			}
			return fmt.Errorf("create affinity group: nil result")
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(nameFlag, "", "The name of the affinity group.")
	cmd.Flags().String(policyFlag, "", `The policy for the affinity group. Valid values for the policy are: "hard-affinity", "hard-anti-affinity", "soft-affinity", "soft-anti-affinity"`)

	if err := flags.MarkFlagsRequired(cmd, nameFlag, policyFlag); err != nil {
		cobra.CheckErr(err)
	}
}

func buildRequest(ctx context.Context, model inputModel, apiClient *iaas.APIClient) iaas.ApiCreateAffinityGroupRequest {
	req := apiClient.CreateAffinityGroup(ctx, model.ProjectId)
	req = req.CreateAffinityGroupPayload(
		iaas.CreateAffinityGroupPayload{
			Name:   utils.Ptr(model.Name),
			Policy: utils.Ptr(model.Policy),
		},
	)
	return req
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Name:            flags.FlagToStringValue(p, cmd, nameFlag),
		Policy:          flags.FlagToStringValue(p, cmd, policyFlag),
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

func outputResult(p *print.Printer, model inputModel, resp iaas.AffinityGroup) error {
	outputFormat := ""
	if model.GlobalFlagModel != nil {
		outputFormat = model.GlobalFlagModel.OutputFormat
	}
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal affinity group: %w", err)
		}
		p.Outputln(string(details))
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(resp, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal affinity group: %w", err)
		}
		p.Outputln(string(details))
	default:
		p.Outputf("Created affinity group %q with id %s\n", model.Name, utils.PtrString(resp.Id))
	}
	return nil
}
