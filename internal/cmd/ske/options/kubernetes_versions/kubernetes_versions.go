package kubernetes_versions

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
)

const (
	supportedFlag = "supported"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Supported bool
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubernetes-versions",
		Short: "Lists SKE provider options for kubernetes-versions",
		Long: fmt.Sprintf("%s\n%s",
			"Lists STACKIT Kubernetes Engine (SKE) provider options for kubernetes-versions.",
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List SKE options for kubernetes-versions`,
				"$ stackit ske options kubernetes-versions"),
			examples.NewExample(
				`List SKE options for supported kubernetes-versions`,
				"$ stackit ske options kubernetes-versions --supported"),
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

			// Call API
			req := buildRequest(ctx, apiClient, model)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("get SKE provider options: %w", err)
			}

			return outputResult(params.Printer, model, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Bool(supportedFlag, false, "List supported versions only")
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Supported:       flags.FlagToBoolValue(p, cmd, supportedFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, apiClient *ske.APIClient, model *inputModel) ske.ApiListProviderOptionsRequest {
	req := apiClient.ListProviderOptions(ctx, model.Region)
	if model.Supported {
		req = req.VersionState("SUPPORTED")
	}
	return req
}

func outputResult(p *print.Printer, model *inputModel, options *ske.ProviderOptions) error {
	if model == nil || model.GlobalFlagModel == nil {
		return fmt.Errorf("model is nil")
	} else if options == nil {
		return fmt.Errorf("options is nil")
	}

	return p.OutputResult(model.OutputFormat, options, func() error {
		versions := *options.KubernetesVersions

		table := tables.NewTable()
		table.SetHeader("VERSION", "STATE", "EXPIRATION DATE", "FEATURE GATES")
		for i := range versions {
			v := versions[i]
			featureGate, err := json.Marshal(*v.FeatureGates)
			if err != nil {
				return fmt.Errorf("marshal featureGates of Kubernetes version %q: %w", *v.Version, err)
			}
			expirationDate := ""
			if v.ExpirationDate != nil {
				expirationDate = v.ExpirationDate.Format(time.RFC3339)
			}
			table.AddRow(
				utils.PtrString(v.Version),
				utils.PtrString(v.State),
				expirationDate,
				string(featureGate))
		}

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("display output: %w", err)
		}
		return nil
	})
}
