package create

import (
	"context"
	"encoding/json"
	"fmt"

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
	nameFlag      = "name"
	publicKeyFlag = "public-key"
	labelFlag     = "labels"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Name      *string
	PublicKey *string
	Labels    *map[string]string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a key pair",
		Long:  "Creates a key pair.",
		Args:  cobra.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a new key pair with public-key "ssh-rsa xxx"`,
				"$ stackit beta key-pair create --public-key `ssh-rsa xxx`",
			),
			examples.NewExample(
				`Create a new key pair with public-key from file "/Users/username/.ssh/id_rsa.pub"`,
				"$ stackit beta key-pair create --public-key `@/Users/username/.ssh/id_rsa.pub`",
			),
			examples.NewExample(
				`Create a new key pair with name "KEY_PAIR_NAME" and public-key "ssh-rsa yyy"`,
				"$ stackit beta key-pair create --name KEY_PAIR_NAME --public-key `ssh-rsa yyy`",
			),
			examples.NewExample(
				`Create a new key pair with public-key "ssh-rsa xxx" and labels "key=value,key1=value1"`,
				"$ stackit beta key-pair create --public-key `ssh-rsa xxx` --labels key=value,key1=value1",
			),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd)
			if err != nil {
				return err
			}

			// Configure client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			if !model.AssumeYes {
				prompt := "Are your sure you want to create a key pair?"
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create key pair: %w", err)
			}

			return outputResult(p, model, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(nameFlag, "", "Key pair name")
	cmd.Flags().Var(flags.ReadFromFileFlag(), publicKeyFlag, "Public key to be imported (format: ssh-rsa|ssh-ed25519)")
	cmd.Flags().StringToString(labelFlag, nil, "Labels are key-value string pairs which can be attached to a key pair. E.g. '--labels key1=value1,key2=value2,...'")

	err := cmd.MarkFlagRequired(publicKeyFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Labels:          flags.FlagToStringToStringPointer(p, cmd, labelFlag),
		Name:            flags.FlagToStringPointer(p, cmd, nameFlag),
		PublicKey:       flags.FlagToStringPointer(p, cmd, publicKeyFlag),
	}

	if p.IsVerbosityDebug() {
		modelStr, err := print.BuildDebugStrFromInputModel(model)
		if err != nil {
			p.Debug(print.ErrorLevel, "convert model to string fo debugging: %v", err)
		} else {
			p.Debug(print.DebugLevel, "parsed input values: %s", modelStr)
		}
	}

	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiCreateKeyPairRequest {
	req := apiClient.CreateKeyPair(ctx)

	var labelsMap *map[string]interface{}
	if model.Labels != nil && len(*model.Labels) > 0 {
		// convert map[string]string to map[string]interface{}
		labelsMap = utils.Ptr(map[string]interface{}{})
		for k, v := range *model.Labels {
			(*labelsMap)[k] = v
		}
	}

	payload := iaas.CreateKeyPairPayload{
		Name:      model.Name,
		Labels:    labelsMap,
		PublicKey: model.PublicKey,
	}

	return req.CreateKeyPairPayload(payload)
}

func outputResult(p *print.Printer, model *inputModel, item *iaas.Keypair) error {
	switch model.OutputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(item, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal key pair: %w", err)
		}
		p.Outputln(string(details))
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(item, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal key pair: %w", err)
		}
		p.Outputln(string(details))
	default:
		p.Outputf("Created key pair %q.\nkey pair Fingerprint: %q\n",
			utils.PtrString(item.Name),
			utils.PtrString(item.Fingerprint),
		)
	}
	return nil
}
