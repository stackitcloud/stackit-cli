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
		Short: "Shows details  of a SKE cluster",
		Long:  "Shows details  of a STACKIT Kubernetes Engine (SKE) cluster.",
		Args:  args.SingleArg(clusterNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Get details of an SKE cluster with name "my-cluster"`,
				"$ stackit ske cluster describe my-cluster"),
			examples.NewExample(
				`Get details of an SKE cluster with name "my-cluster" in a table format`,
				"$ stackit ske cluster describe my-cluster --output-format pretty"),
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read SKE cluster: %w", err)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *ske.APIClient) ske.ApiGetClusterRequest {
	req := apiClient.GetCluster(ctx, model.ProjectId, model.ClusterName)
	return req
}

func outputResult(cmd *cobra.Command, outputFormat string, cluster *ske.Cluster) error {
	switch outputFormat {
	case globalflags.PrettyOutputFormat:

		acl := []string{}
		if cluster.Extensions != nil && cluster.Extensions.Acl != nil {
			acl = *cluster.Extensions.Acl.AllowedCidrs
		}

		table := tables.NewTable()
		table.AddRow("NAME", *cluster.Name)
		table.AddSeparator()
		table.AddRow("STATE", *cluster.Status.Aggregated)
		table.AddSeparator()
		table.AddRow("VERSION", *cluster.Kubernetes.Version)
		table.AddSeparator()
		table.AddRow("ACL", acl)
		err := table.Display(cmd)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	default:
		details, err := json.MarshalIndent(cluster, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal SKE cluster: %w", err)
		}
		cmd.Println(string(details))

		return nil
	}
}
