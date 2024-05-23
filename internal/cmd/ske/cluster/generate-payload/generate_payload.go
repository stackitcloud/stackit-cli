package generatepayload

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/fileutils"
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
	filePathFlag    = "file-path"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ClusterName *string
	FilePath    *string
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
				`$ stackit ske cluster generate-payload --file-path ./payload.json`,
				`<Modify payload in file, if needed>`,
				`$ stackit ske cluster create my-cluster --payload @./payload.json`),
			examples.NewExample(
				`Generate a payload with values of a cluster, and adapt it with custom values for the different configuration options`,
				`$ stackit ske cluster generate-payload --cluster-name my-cluster --file-path ./payload.json`,
				`<Modify payload in file>`,
				`$ stackit ske cluster update my-cluster --payload @./payload.json`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
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

			return outputResult(p, model, payload)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(clusterNameFlag, "n", "", "If set, generates the payload with the current state of the given cluster. If unset, generates the payload with default values")
	cmd.Flags().StringP(filePathFlag, "f", "", "If set, writes the payload to the given file. If unset, writes the payload to the standard output")
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	clusterName := flags.FlagToStringPointer(p, cmd, clusterNameFlag)
	// If clusterName is provided, projectId is needed as well
	if clusterName != nil && globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ClusterName:     clusterName,
		FilePath:        flags.FlagToStringPointer(p, cmd, filePathFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *ske.APIClient) ske.ApiGetClusterRequest {
	req := apiClient.GetCluster(ctx, model.ProjectId, *model.ClusterName)
	return req
}

func outputResult(p *print.Printer, model *inputModel, payload *ske.CreateOrUpdateClusterPayload) error {
	payloadBytes, err := json.MarshalIndent(*payload, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	if model.FilePath != nil {
		err = fileutils.FileOutput(*model.FilePath, string(payloadBytes))
		if err != nil {
			return fmt.Errorf("write payload to the file: %w", err)
		}
	} else {
		p.Outputln(string(payloadBytes))
	}

	return nil
}
