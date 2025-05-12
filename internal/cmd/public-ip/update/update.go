package update

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	publicIpIdArg = "PUBLIC_IP_ID"

	labelFlag = "labels"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	PublicIpId string
	Labels     *map[string]string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", publicIpIdArg),
		Short: "Updates a Public IP",
		Long:  "Updates a Public IP.",
		Args:  args.SingleArg(publicIpIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update public IP with ID "xxx"`,
				`$ stackit public-ip update xxx`,
			),
			examples.NewExample(
				`Update public IP with ID "xxx" with new labels`,
				`$ stackit public-ip update xxx --labels key=value,foo=bar`,
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

			publicIpLabel, _, err := iaasUtils.GetPublicIP(ctx, apiClient, model.ProjectId, model.PublicIpId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get public IP: %v", err)
				publicIpLabel = model.PublicIpId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update public IP %q?", publicIpLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update public IP: %w", err)
			}

			return outputResult(params.Printer, model, publicIpLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringToString(labelFlag, nil, "Labels are key-value string pairs which can be attached to a public IP. E.g. '--labels key1=value1,key2=value2,...'")
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	publicIpId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	labels := flags.FlagToStringToStringPointer(p, cmd, labelFlag)

	if labels == nil {
		return nil, &errors.EmptyUpdateError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		PublicIpId:      publicIpId,
		Labels:          labels,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiUpdatePublicIPRequest {
	req := apiClient.UpdatePublicIP(ctx, model.ProjectId, model.PublicIpId)

	var labelsMap *map[string]interface{}
	if model.Labels != nil && len(*model.Labels) > 0 {
		// convert map[string]string to map[string]interface{}
		labelsMap = utils.Ptr(map[string]interface{}{})
		for k, v := range *model.Labels {
			(*labelsMap)[k] = v
		}
	}

	payload := iaas.UpdatePublicIPPayload{
		Labels: labelsMap,
	}

	return req.UpdatePublicIPPayload(payload)
}

func outputResult(p *print.Printer, model *inputModel, publicIpLabel string, publicIp *iaas.PublicIp) error {
	switch model.OutputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(publicIp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal public IP: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(publicIp, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal public IP: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		p.Outputf("Updated public IP %q.\n", publicIpLabel)
		return nil
	}
}
