package describe

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	keyPairNameArg = "KEY_PAIR_NAME"

	publicKeyFlag = "public-key"

	maxLengthPublicKey = 50
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	KeyPairName string
	PublicKey   bool
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", keyPairNameArg),
		Short: "Describes a key pair",
		Long:  "Describes a key pair.",
		Args:  args.SingleArg(keyPairNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Get details about a key pair with name "KEY_PAIR_NAME"`,
				"$ stackit beta key-pair describe KEY_PAIR_NAME",
			),
			examples.NewExample(
				`Get only the SSH public key of a key pair with name "KEY_PAIR_NAME"`,
				"$ stackit beta key-pair describe KEY_PAIR_NAME --public-key",
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read key pair: %w", err)
			}

			return outputResult(p, model.OutputFormat, model.PublicKey, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Bool(publicKeyFlag, false, "Show only the public key")
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	keyPairName := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)

	model := inputModel{
		GlobalFlagModel: globalFlags,
		KeyPairName:     keyPairName,
		PublicKey:       flags.FlagToBoolValue(p, cmd, publicKeyFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiGetKeyPairRequest {
	return apiClient.GetKeyPair(ctx, model.KeyPairName)
}

func outputResult(p *print.Printer, outputFormat string, showOnlyPublicKey bool, keyPair *iaas.Keypair) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(keyPair, "", "  ")
		if showOnlyPublicKey {
			onlyPublicKey := map[string]string{
				"publicKey": *keyPair.PublicKey,
			}
			details, err = json.MarshalIndent(onlyPublicKey, "", "  ")
		}

		if err != nil {
			return fmt.Errorf("marshal key pair: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(keyPair, yaml.IndentSequence(true))
		if showOnlyPublicKey {
			onlyPublicKey := map[string]string{
				"publicKey": *keyPair.PublicKey,
			}
			details, err = yaml.MarshalWithOptions(onlyPublicKey, yaml.IndentSequence(true))
		}

		if err != nil {
			return fmt.Errorf("marshal key pair: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		if showOnlyPublicKey {
			p.Outputln(*keyPair.PublicKey)
			return nil
		}
		table := tables.NewTable()
		table.AddRow("KEY PAIR NAME", utils.PtrString(keyPair.Name))
		table.AddSeparator()

		if *keyPair.Labels != nil && len(*keyPair.Labels) > 0 {
			var labels []string
			for key, value := range *keyPair.Labels {
				labels = append(labels, fmt.Sprintf("%s: %s", key, value))
			}
			table.AddRow("LABELS", strings.Join(labels, "\n"))
			table.AddSeparator()
		}

		table.AddRow("FINGERPRINT", utils.PtrString(keyPair.Fingerprint))
		table.AddSeparator()

		truncatedPublicKey := (*keyPair.PublicKey)[:maxLengthPublicKey] + "..."
		table.AddRow("PUBLIC KEY", truncatedPublicKey)
		table.AddSeparator()

		table.AddRow("CREATED AT", utils.PtrString(keyPair.CreatedAt))
		table.AddSeparator()

		table.AddRow("UPDATED AT", utils.PtrString(keyPair.UpdatedAt))
		table.AddSeparator()

		p.Outputln(table.Render())
	}

	return nil
}
