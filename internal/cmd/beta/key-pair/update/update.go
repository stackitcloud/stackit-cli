package update

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	keypairNameArg = "KEYPAIR_NAME"
	labelsFlag     = "labels"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Labels      *map[string]string
	KeypairName *string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update key pair",
		Long:  "Update key pair.",
		Args:  args.SingleArg(keypairNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Update the labels of a keypair KEYPAIR_NAME with "key=value,key1=value1"`,
				"$ stackit beta key-pair update --labels key=value,key1=value1",
			),
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

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update key pair %q?", *model.KeypairName)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return fmt.Errorf("update keypair: %w", err)
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update keypair: %w", err)
			}

			return outputResult(p, model, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringToString(labelsFlag, nil, "Labels are key-value string pairs which can be attached to a server. E.g. '--labels key1=value1,key2=value2,...'")

	err := cmd.MarkFlagRequired(labelsFlag)
	cobra.CheckErr(err)
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiUpdateKeyPairRequest {
	req := apiClient.UpdateKeyPair(ctx, *model.KeypairName)

	var labelsMap *map[string]interface{}
	if model.Labels != nil && len(*model.Labels) > 0 {
		// convert map[string]string to map[string]interface{}
		labelsMap = utils.Ptr(map[string]interface{}{})
		for k, v := range *model.Labels {
			(*labelsMap)[k] = v
		}
	}
	payload := iaas.UpdateKeyPairPayload{
		Labels: labelsMap,
	}
	return req.UpdateKeyPairPayload(payload)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	keypairName := inputArgs[0]
	globalFlags := globalflags.Parse(p, cmd)

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Labels:          flags.FlagToStringToStringPointer(p, cmd, labelsFlag),
		KeypairName:     utils.Ptr(keypairName),
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

func outputResult(p *print.Printer, model *inputModel, keypair *iaas.Keypair) error {
	switch model.OutputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(keypair, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal keypair: %w", err)
		}
		p.Outputln(string(details))
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(keypair, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal keypair: %w", err)
		}
		p.Outputln(string(details))
	default:
		p.Outputf("Updated labels of Keypair %q\n", *model.KeypairName)
	}
	return nil
}
