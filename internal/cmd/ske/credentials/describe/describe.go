package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/client"
	skeUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
)

const (
	clusterNameArg = "CLUSTER_NAME"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ClusterName string
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", clusterNameArg),
		Short: "Shows details of the credentials associated to a SKE cluster",
		Long:  "Shows details of the credentials associated to a STACKIT Kubernetes Engine (SKE) cluster",
		Args:  args.SingleArg(clusterNameArg, nil),
		Deprecated: fmt.Sprintf("%s\n%s\n%s\n%s\n",
			"and will be removed in a future release.",
			"Please use the following command to obtain a kubeconfig file instead:",
			" $ stackit ske kubeconfig create CLUSTER_NAME",
			"For more information, visit: https://docs.stackit.cloud/stackit/en/how-to-rotate-ske-credentials-200016334.html",
		),
		Example: examples.Build(
			examples.NewExample(
				`Get details of the credentials associated to the SKE cluster with name "my-cluster"`,
				"$ stackit ske credentials describe my-cluster"),
			examples.NewExample(
				`Get details of the credentials associated to the SKE cluster with name "my-cluster" in a table format`,
				"$ stackit ske credentials describe my-cluster --output-format pretty"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(cmd)
			if err != nil {
				return err
			}

			// Check if SKE is enabled for this project
			enabled, err := skeUtils.ProjectEnabled(ctx, apiClient, model.ProjectId)
			if err != nil {
				return err
			}
			if !enabled {
				return fmt.Errorf("SKE isn't enabled for this project, please run 'stackit ske enable'")
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("get SKE credentials: %w", err)
			}

			return outputResult(cmd, model.OutputFormat, resp)
		},
	}
	return cmd
}

func parseInput(cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	clusterName := inputArgs[0]

	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		ClusterName:     clusterName,
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *ske.APIClient) ske.ApiGetCredentialsRequest {
	req := apiClient.GetCredentials(ctx, model.ProjectId, model.ClusterName)
	return req
}

func outputResult(cmd *cobra.Command, outputFormat string, credentials *ske.Credentials) error {
	switch outputFormat {
	case globalflags.PrettyOutputFormat:
		table := tables.NewTable()
		table.AddRow("SERVER", *credentials.Server)
		table.AddSeparator()
		table.AddRow("TOKEN", *credentials.Token)
		err := table.Display(cmd)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	default:
		details, err := json.MarshalIndent(credentials, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal SKE credentials: %w", err)
		}
		cmd.Println(string(details))

		return nil
	}
}
