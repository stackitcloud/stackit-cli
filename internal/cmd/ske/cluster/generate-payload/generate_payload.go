package generatepayload

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/client"
	skeUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
)

const (
	clusterNameFlag = "cluster-name"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ClusterName *string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-payload",
		Short: "Generates a payload to create/update SKE clusters",
		Long: fmt.Sprintf("%s\n%s",
			"Generates a JSON payload with values to be used as --payload input for cluster creation or update.",
			"See https://docs.api.stackit.cloud/documentation/ske/version/v1#tag/Cluster/operation/SkeService_CreateOrUpdateCluster for information regarding the payload structure.",
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Generate a payload with default values, and adapt it with custom values for the different configuration options`,
				`$ stackit ske cluster generate-payload > ./payload.json`,
				`<Modify payload in file, if needed>`,
				`$ stackit ske cluster create my-cluster --payload @./payload.json`),
			examples.NewExample(
				`Generate a payload with values of a cluster, and adapt it with custom values for the different configuration options`,
				`$ stackit ske cluster generate-payload --cluster-name my-cluster > ./payload.json`,
				`<Modify payload in file>`,
				`$ stackit ske cluster update my-cluster --payload @./payload.json`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(cmd)
			if err != nil {
				return err
			}

			var payload *ske.CreateOrUpdateClusterPayload
			if model.ClusterName == nil {
				payload, err = skeUtils.GetDefaultPayload(ctx, apiClient)
				if err != nil {
					return err
				}
			} else {
				req := buildRequest(ctx, model, apiClient)
				resp, err := req.Execute()
				if err != nil {
					return fmt.Errorf("read SKE cluster: %w", err)
				}
				payload = &ske.CreateOrUpdateClusterPayload{
					Extensions:  resp.Extensions,
					Hibernation: resp.Hibernation,
					Kubernetes:  resp.Kubernetes,
					Maintenance: resp.Maintenance,
					Nodepools:   resp.Nodepools,
					Status:      resp.Status,
				}
			}

			return outputResult(p, payload)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(clusterNameFlag, "n", "", "If set, generates the payload with the current state of the given cluster. If unset, generates the payload with default values")
}

func parseInput(cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(cmd)

	clusterName := flags.FlagToStringPointer(cmd, clusterNameFlag)
	// If clusterName is provided, projectId is needed as well
	if clusterName != nil && globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		ClusterName:     clusterName,
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *ske.APIClient) ske.ApiGetClusterRequest {
	req := apiClient.GetCluster(ctx, model.ProjectId, *model.ClusterName)
	return req
}

func outputResult(p *print.Printer, payload *ske.CreateOrUpdateClusterPayload) error {
	payloadBytes, err := json.MarshalIndent(*payload, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}
	p.Outputln(string(payloadBytes))

	return nil
}
