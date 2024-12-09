package describe

import (
	"context"
	"encoding/json"
	"fmt"
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
	keypairNameArg = "KEYPAIR_NAME"

	publicKeyFlag = "public-key"

	maxLengthPublicKey = 50
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	KeypairName string
	PublicKey   bool
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Describe a keypair",
		Long:  "Describe a keypair.",
		Args:  args.SingleArg(keypairNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Get details about a keypair named "KEYPAIR_NAME"`,
				"$ stackit beta keypair describe KEYPAIR_NAME",
			),
			examples.NewExample(
				`Get only the SSH public key of a keypair with the name "KEYPAIR_NAME"`,
				"$ stackit beta keypair describe KEYPAIR_NAME --public-key",
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
				return fmt.Errorf("read keypair: %w", err)
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
	keypairName := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)

	model := inputModel{
		GlobalFlagModel: globalFlags,
		KeypairName:     keypairName,
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
	return apiClient.GetKeyPair(ctx, model.KeypairName)
}

func outputResult(p *print.Printer, outputFormat string, shpwOnlyPublicKey bool, keypair *iaas.Keypair) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(keypair, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal keypair: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(keypair, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal keypair: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		if shpwOnlyPublicKey {
			p.Outputln(*keypair.PublicKey)
			return nil
		}
		table := tables.NewTable()
		table.AddRow("KEYPAIR NAME", *keypair.Name)
		table.AddSeparator()

		if *keypair.Labels != nil && len(*keypair.Labels) > 0 {
			var labels []string
			for key, value := range *keypair.Labels {
				labels = append(labels, fmt.Sprintf("%s: %s", key, value))
			}
			table.AddRow("LABELS", strings.Join(labels, "\n"))
			table.AddSeparator()
		}

		table.AddRow("FINGERPRINT", *keypair.Fingerprint)
		table.AddSeparator()

		truncatedPublicKey := (*keypair.PublicKey)[:maxLengthPublicKey] + "..."
		table.AddRow("PUBLIC KEY", truncatedPublicKey)
		table.AddSeparator()

		table.AddRow("CREATED AT", *keypair.CreatedAt)
		table.AddSeparator()

		table.AddRow("UPDATED AT", *keypair.UpdatedAt)
		table.AddSeparator()

		p.Outputln(table.Render())
	}

	return nil
}
