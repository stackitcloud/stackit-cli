package update

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	keyPairNameArg = "KEY_PAIR_NAME"
	labelsFlag     = "labels"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Labels      *map[string]string
	KeyPairName *string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", keyPairNameArg),
		Short: "Updates a key pair",
		Long:  "Updates a key pair.",
		Args:  args.SingleArg(keyPairNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Update the labels of a key pair with name "KEY_PAIR_NAME" with "key=value,key1=value1"`,
				"$ stackit key-pair update KEY_PAIR_NAME --labels key=value,key1=value1",
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model := parseInput(params.Printer, cmd, args)

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update key pair %q?", *model.KeyPairName)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return fmt.Errorf("update key pair: %w", err)
				}
			}

			// Call API
			req := buildRequest(ctx, &model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update key pair: %w", err)
			}
			if resp == nil {
				return fmt.Errorf("response is nil")
			}

			return outputResult(params.Printer, model, *resp)
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
	req := apiClient.UpdateKeyPair(ctx, *model.KeyPairName)

	payload := iaas.UpdateKeyPairPayload{
		Labels: utils.ConvertStringMapToInterfaceMap(model.Labels),
	}
	return req.UpdateKeyPairPayload(payload)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) inputModel {
	keyPairName := inputArgs[0]
	globalFlags := globalflags.Parse(p, cmd)

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Labels:          flags.FlagToStringToStringPointer(p, cmd, labelsFlag),
		KeyPairName:     utils.Ptr(keyPairName),
	}

	p.DebugInputModel(model)
	return model
}

func outputResult(p *print.Printer, model inputModel, keyPair iaas.Keypair) error {
	var outputFormat string
	if model.GlobalFlagModel != nil {
		outputFormat = model.GlobalFlagModel.OutputFormat
	}

	return p.OutputResult(outputFormat, keyPair, func() error {
		p.Outputf("Updated labels of key pair %q\n", utils.PtrString(model.KeyPairName))
		return nil
	})
}
