package create

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
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
	Payload     *ske.CreateOrUpdateClusterPayload
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("create %s", clusterNameArg),
		Short: "Creates an SKE cluster",
		Long: fmt.Sprintf("%s\n%s\n%s",
			"Creates a STACKIT Kubernetes Engine (SKE) cluster.",
			"The payload can be provided as a JSON string or a file path prefixed with \"@\".",
			"See https://docs.api.stackit.cloud/documentation/ske/version/v1#tag/Cluster/operation/SkeService_CreateOrUpdateCluster for information regarding the payload structure.",
		),
		Args: args.SingleArg(clusterNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Create an SKE cluster using default configuration`,
				"$ stackit ske cluster create my-cluster"),
			examples.NewExample(
				`Create an SKE cluster using an API payload sourced from the file "./payload.json"`,
				"$ stackit ske cluster create my-cluster --payload @./payload.json"),
			examples.NewExample(
				`Create an SKE cluster using an API payload provided as a JSON string`,
				`$ stackit ske cluster create my-cluster --payload "{...}"`),
			examples.NewExample(
				`Generate a payload with default values, and adapt it with custom values for the different configuration options`,
				`$ stackit ske cluster generate-payload > ./payload.json`,
				`<Modify payload in file, if needed>`,
				`$ stackit ske cluster create my-cluster --payload @./payload.json`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, cmd, p)
			if err != nil {
				projectLabel = model.ProjectId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a cluster for project %q?", projectLabel)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Check if SKE is enabled for this project
			enabled, err := skeUtils.ProjectEnabled(ctx, apiClient, model.ProjectId)
			if err != nil {
				return err
			}
			if !enabled {
				return fmt.Errorf("SKE isn't enabled for this project, please run 'stackit ske enable'")
			}

			// Check if cluster exists
			exists, err := skeUtils.ClusterExists(ctx, apiClient, model.ProjectId, model.ClusterName)
			if err != nil {
				return err
			}
			if exists {
				return fmt.Errorf("cluster with name %s already exists", model.ClusterName)
			}

			// Fill in default payload, if needed
			if model.Payload == nil {
				defaultPayload, err := skeUtils.GetDefaultPayload(ctx, apiClient)
				if err != nil {
					return fmt.Errorf("get default payload: %w", err)
				}
				model.Payload = defaultPayload
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create SKE cluster: %w", err)
			}
			name := *resp.Name

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(p)
				s.Start("Creating cluster")
				_, err = wait.CreateOrUpdateClusterWaitHandler(ctx, apiClient, model.ProjectId, name).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for SKE cluster creation: %w", err)
				}
				s.Stop()
			}

			operationState := "Created"
			if model.Async {
				operationState = "Triggered creation of"
			}
			p.Outputf("%s cluster for project %q. Cluster name: %s\n", operationState, projectLabel, name)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.ReadFromFileFlag(), payloadFlag, `Request payload (JSON). Can be a string or a file path, if prefixed with "@" (example: @./payload.json). If unset, will use a default payload (you can check it by running "stackit ske cluster generate-payload")`)
}

func parseInput(cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	clusterName := inputArgs[0]

	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	payloadValue := flags.FlagToStringPointer(cmd, payloadFlag)
	var payload *ske.CreateOrUpdateClusterPayload
	if payloadValue != nil {
		payload = &ske.CreateOrUpdateClusterPayload{}
		err := json.Unmarshal([]byte(*payloadValue), payload)
		if err != nil {
			return nil, fmt.Errorf("encode payload: %w", err)
		}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		ClusterName:     clusterName,
		Payload:         payload,
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *ske.APIClient) ske.ApiCreateOrUpdateClusterRequest {
	req := apiClient.CreateOrUpdateCluster(ctx, model.ProjectId, model.ClusterName)

	req = req.CreateOrUpdateClusterPayload(*model.Payload)
	return req
}
