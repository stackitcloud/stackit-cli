package create

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
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

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a key pair",
		Long:  "Creates a key pair.",
		Args:  cobra.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a new key pair with public-key "ssh-rsa xxx"`,
				"$ stackit key-pair create --public-key `ssh-rsa xxx`",
			),
			examples.NewExample(
				`Create a new key pair with public-key from file "/Users/username/.ssh/id_rsa.pub"`,
				"$ stackit key-pair create --public-key `@/Users/username/.ssh/id_rsa.pub`",
			),
			examples.NewExample(
				`Create a new key pair with name "KEY_PAIR_NAME" and public-key "ssh-rsa yyy"`,
				"$ stackit key-pair create --name KEY_PAIR_NAME --public-key `ssh-rsa yyy`",
			),
			examples.NewExample(
				`Create a new key pair with public-key "ssh-rsa xxx" and labels "key=value,key1=value1"`,
				"$ stackit key-pair create --public-key `ssh-rsa xxx` --labels key=value,key1=value1",
			),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd)
			if err != nil {
				return err
			}

			// Configure client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			if !model.AssumeYes {
				prompt := "Are your sure you want to create a key pair?"
				err = params.Printer.PromptForConfirmation(prompt)
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

			return outputResult(params.Printer, model.GlobalFlagModel.OutputFormat, resp)
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

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiCreateKeyPairRequest {
	req := apiClient.CreateKeyPair(ctx)

	payload := iaas.CreateKeyPairPayload{
		Name:      model.Name,
		Labels:    utils.ConvertStringMapToInterfaceMap(model.Labels),
		PublicKey: model.PublicKey,
	}

	return req.CreateKeyPairPayload(payload)
}

func outputResult(p *print.Printer, outputFormat string, item *iaas.Keypair) error {
	if item == nil {
		return fmt.Errorf("no key pair found")
	}

	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(item, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal key pair: %w", err)
		}
		p.Outputln(string(details))
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(item, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
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
