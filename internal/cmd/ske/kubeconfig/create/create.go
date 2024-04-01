package create

import (
	"context"
	"fmt"
	"os"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/confirm"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/client"
	skeUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
)

const (
	clusterNameArg = "CLUSTER_NAME"

	expirationFlag = "expiration"
	locationFlag   = "location"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ClusterName    string
	Location       *string
	ExpirationTime *string
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("create %s", clusterNameArg),
		Short: "Creates a kubeconfig for an SKE cluster",
		Long: fmt.Sprintf("%s\n%s",
			"Creates a kubeconfig for a STACKIT Kubernetes Engine (SKE) cluster.",
			"By default the kubeconfig is created in the .kube folder, in the user's home directory. The kubeconfig file will be overwritten if it already exists."),
		Args: args.SingleArg(clusterNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Create a kubeconfig for the SKE cluster with name "my-cluster"`,
				"$ stackit ske kubeconfig create my-cluster"),
			examples.NewExample(
				`Create a kubeconfig for the SKE cluster with name "my-cluster" and set the expiration time to 30 days`,
				"$ stackit ske kubeconfig create my-cluster --expiration 30d"),
			examples.NewExample(
				`Create a kubeconfig for the SKE cluster with name "my-cluster" and set the expiration time to 2 months`,
				"$ stackit ske kubeconfig create my-cluster --expiration 2M"),
			examples.NewExample(
				`Create a kubeconfig for the SKE cluster with name "my-cluster" in a custom location`,
				"$ stackit ske kubeconfig create my-cluster --location /path/to/config"),
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

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a kubeconfig for SKE cluster %q? This will OVERWRITE your current configuration, if it exists.", model.ClusterName)
				err = confirm.PromptForConfirmation(cmd, prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req, err := buildRequest(ctx, model, apiClient)
			if err != nil {
				return fmt.Errorf("build kubeconfig create request: %w", err)
			}
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create kubeconfig for SKE cluster: %w", err)
			}

			// Create the config file
			if resp.Kubeconfig == nil {
				return fmt.Errorf("no kubeconfig returned from the API")
			}

			configPath, err := writeConfigFile(model.Location, *resp.Kubeconfig)
			if err != nil {
				return fmt.Errorf("write kubeconfig file: %w", err)
			}

			fmt.Printf("Created kubeconfig file for cluster %s in %q, with expiration date %v\n", model.ClusterName, configPath, *resp.ExpirationTimestamp)

			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(expirationFlag, "e", "", "Expiration time for the kubeconfig in seconds(s), minutes(m), hours(h), days(d) or months(M). Example: 30d. By default, expiration time is 1h")
	cmd.Flags().String(locationFlag, "", "Folder location to store the kubeconfig file. By default, the kubeconfig is created in the .kube folder, in the user's home directory.")
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
		Location:        flags.FlagToStringPointer(cmd, locationFlag),
		ExpirationTime:  flags.FlagToStringPointer(cmd, expirationFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *ske.APIClient) (ske.ApiCreateKubeconfigRequest, error) {
	req := apiClient.CreateKubeconfig(ctx, model.ProjectId, model.ClusterName)

	payload := ske.CreateKubeconfigPayload{}

	if model.ExpirationTime != nil {
		expirationTime, err := skeUtils.ConvertToSeconds(*model.ExpirationTime)
		if err != nil {
			return req, fmt.Errorf("parsing expiration time: %w", err)
		}

		payload.ExpirationSeconds = expirationTime
	}

	return req.CreateKubeconfigPayload(payload), nil
}

func writeConfigFile(configPath *string, data string) (string, error) {
	if configPath == nil {
		userHome, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("get user home directory: %w", err)
		}

		err = os.MkdirAll(fmt.Sprintf("%s/.kube", userHome), 0o700)
		if err != nil {
			return "", fmt.Errorf("create kube directory: %w", err)
		}
		configPath = utils.Ptr(fmt.Sprintf("%s/.kube", userHome))
	}

	writeLocation := fmt.Sprintf("%s/config", *configPath)

	err := os.WriteFile(writeLocation, []byte(data), 0o600)
	if err != nil {
		return "", fmt.Errorf("write file: %w", err)
	}
	return writeLocation, nil
}
