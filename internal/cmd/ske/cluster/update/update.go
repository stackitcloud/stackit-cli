package update

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/client"
	skeUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
	"github.com/stackitcloud/stackit-sdk-go/services/ske/wait"
)

const (
	clusterNameArg = "CLUSTER_NAME"

	payloadFlag = "payload"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ClusterName string
	Payload     ske.CreateOrUpdateClusterPayload
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", clusterNameArg),
		Short: "Updates an SKE cluster",
		Long: fmt.Sprintf("%s\n%s\n%s",
			"Updates a STACKIT Kubernetes Engine (SKE) cluster.",
			"The payload can be provided as a JSON string or a file path prefixed with \"@\".",
			"See https://docs.api.stackit.cloud/documentation/ske/version/v1#tag/Cluster/operation/SkeService_CreateOrUpdateCluster for information regarding the payload structure.",
		),
		Args: args.SingleArg(clusterNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Update an SKE cluster using an API payload sourced from the file "./payload.json"`,
				"$ stackit ske cluster update my-cluster --payload @./payload.json"),
			examples.NewExample(
				`Update an SKE cluster using an API payload provided as a JSON string`,
				`$ stackit ske cluster update my-cluster --payload "{...}"`),
			examples.NewExample(
				`Generate a payload with the current values of a cluster, and adapt it with custom values for the different configuration options`,
				`$ stackit ske cluster generate-payload --cluster-name my-cluster > ./payload.json`,
				`<Modify payload in file>`,
				`$ stackit ske cluster update my-cluster --payload @./payload.json`),
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

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update cluster %q?", model.ClusterName)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Check if cluster exists
			exists, err := skeUtils.ClusterExists(ctx, apiClient, model.ProjectId, model.Region, model.ClusterName)
			if err != nil {
				return err
			}
			if !exists {
				return fmt.Errorf("cluster with name %s does not exist", model.ClusterName)
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update SKE cluster: %w", err)
			}
			name := *resp.Name

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Updating cluster")
				_, err = wait.CreateOrUpdateClusterWaitHandler(ctx, apiClient, model.ProjectId, model.Region, name).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for SKE cluster update: %w", err)
				}
				s.Stop()
			}

			return outputResult(params.Printer, model.OutputFormat, model.Async, model.ClusterName, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.ReadFromFileFlag(), payloadFlag, `Request payload (JSON). Can be a string or a file path, if prefixed with "@". Example: @./payload.json`)

	err := flags.MarkFlagsRequired(cmd, payloadFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	clusterName := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	payloadString := flags.FlagToStringValue(p, cmd, payloadFlag)
	var payload ske.CreateOrUpdateClusterPayload
	err := json.Unmarshal([]byte(payloadString), &payload)
	if err != nil {
		return nil, fmt.Errorf("encode payload: %w", err)
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ClusterName:     clusterName,
		Payload:         payload,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *ske.APIClient) ske.ApiCreateOrUpdateClusterRequest {
	req := apiClient.CreateOrUpdateCluster(ctx, model.ProjectId, model.Region, model.ClusterName)

	req = req.CreateOrUpdateClusterPayload(model.Payload)
	return req
}

func outputResult(p *print.Printer, outputFormat string, async bool, clusterName string, cluster *ske.Cluster) error {
	if cluster == nil {
		return fmt.Errorf("cluster is nil")
	}

	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(cluster, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal SKE cluster: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(cluster, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal SKE cluster: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		operationState := "Updated"
		if async {
			operationState = "Triggered update of"
		}
		p.Info("%s cluster %q\n", operationState, clusterName)
		return nil
	}
}
