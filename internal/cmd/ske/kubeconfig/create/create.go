package create

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/confirm"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/client"
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
	KubeconfigPath *string
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
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create kubeconfig for SKE cluster: %w", err)
			}

			// Create a config file in $HOME/.kube/config
			configPath := model.KubeconfigPath

			if configPath == nil {
				userHome, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("could not get user home directory: %w", err)
				}

				err = os.MkdirAll(fmt.Sprintf("%s/.kube", userHome), 0700)
				if err != nil {
					return fmt.Errorf("could not create kube directory: %w", err)
				}
				configPath = utils.Ptr(fmt.Sprintf("%s/.kube", userHome))
			}

			err = os.WriteFile(fmt.Sprintf("%s/config", *configPath), []byte(*resp.Kubeconfig), 0600)
			if err != nil {
				return fmt.Errorf("could not write kubeconfig file: %w", err)
			}

			fmt.Printf("Created kubeconfig file for cluster %s with expiration date %v.\n", model.ClusterName, *resp.ExpirationTimestamp)

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

	expirationTimeInput := flags.FlagToStringPointer(cmd, expirationFlag)
	var expirationTime *string
	if expirationTimeInput != nil {
		expirationTime = convertToSeconds(*expirationTimeInput)
		if expirationTime == nil {
			return nil, fmt.Errorf("invalid expiration time: %s", *expirationTimeInput)
		}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		ClusterName:     clusterName,
		KubeconfigPath:  flags.FlagToStringPointer(cmd, locationFlag),
		ExpirationTime:  expirationTime,
	}, nil
}

func convertToSeconds(timeStr string) *string {
	if len(timeStr) < 2 {
		return nil
	}

	unit := timeStr[len(timeStr)-1:]
	if _, err := strconv.Atoi(unit); err == nil {
		// If the last character is a digit, assume the whole string is a number of seconds
		return utils.Ptr(timeStr)
	}

	valueStr := timeStr[:len(timeStr)-1]
	value, err := strconv.ParseUint(valueStr, 10, 64)
	if err != nil {
		return nil
	}

	var multiplier uint64
	switch unit {
	// second
	case "s":
		multiplier = 1
	// minute
	case "m":
		multiplier = 60
	// hour
	case "h":
		multiplier = 60 * 60
	// day
	case "d":
		multiplier = 60 * 60 * 24
	// month, assume 30 days
	case "M":
		multiplier = 60 * 60 * 24 * 30
	default:
		return nil
	}

	result := uint64(value) * multiplier
	return utils.Ptr(strconv.FormatUint(result, 10))
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *ske.APIClient) ske.ApiCreateKubeconfigRequest {
	req := apiClient.CreateKubeconfig(ctx, model.ProjectId, model.ClusterName)

	payload := ske.CreateKubeconfigPayload{}

	if model.ExpirationTime != nil {
		payload.ExpirationSeconds = model.ExpirationTime
	}

	return req.CreateKubeconfigPayload(payload)
}
