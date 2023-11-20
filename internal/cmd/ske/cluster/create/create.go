package create

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/client"
	skeUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
	"github.com/stackitcloud/stackit-sdk-go/services/ske/wait"
)

const (
	ProjectIdFlag = "project-id"
	NameFlag      = "name"
	PayloadFlag   = "payload"
)

type FlagModel struct {
	ProjectId string
	Name      string
	Payload   ske.CreateOrUpdateClusterPayload
}

var Cmd = &cobra.Command{
	Use:     "create",
	Short:   "Creates an SKE cluster",
	Long:    "Creates an SKE cluster",
	Example: `$ stackit ske cluster create --project-id xxx --payload @./payload.json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		model, err := ParseFlags(cmd, os.ReadFile)
		if err != nil {
			return err
		}

		// Configure API client
		apiClient, err := client.ConfigureClient(cmd)
		if err != nil {
			return fmt.Errorf("authentication failed, please run \"stackit auth login\" or \"stackit auth activate-service-account\"")
		}

		// Check if cluster exists
		exists, err := skeUtils.ClusterExists(ctx, apiClient, model.ProjectId, model.Name)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("cluster with name %s already exists", model.Name)
		}

		// Call API
		req, err := BuildRequest(ctx, model, apiClient)
		if err != nil {
			return fmt.Errorf("build SKE cluster creation request: %w", err)
		}
		resp, err := req.Execute()
		if err != nil {
			return fmt.Errorf("create SKE cluster: %w", err)
		}

		// Wait for async operation
		name := *resp.Name
		_, err = wait.CreateOrUpdateClusterWaitHandler(ctx, apiClient, model.ProjectId, name).WaitWithContext(ctx)
		if err != nil {
			return fmt.Errorf("wait for SKE cluster creation: %w", err)
		}

		fmt.Printf("Created cluster with name %s\n", name)
		return nil
	},
}

func init() {
	ConfigureFlags(Cmd)
}

func ConfigureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(NameFlag, "n", "", "Cluster name")
	cmd.Flags().String(PayloadFlag, "", "Request payload (JSON)")

	err := utils.MarkFlagsRequired(cmd, NameFlag)
	cobra.CheckErr(err)
	err = utils.MarkFlagsRequired(cmd, PayloadFlag)
	cobra.CheckErr(err)
}

type FileReaderFunc func(filename string) ([]byte, error)

func ParseFlags(cmd *cobra.Command, fileReaderFunc FileReaderFunc) (*FlagModel, error) {
	projectId := viper.GetString(config.ProjectIdKey)
	if projectId == "" {
		return nil, fmt.Errorf("project ID not set")
	}

	name := utils.FlagToStringValue(cmd, NameFlag)
	payloadString := utils.FlagToStringValue(cmd, PayloadFlag)
	payloadStringBytes := []byte(payloadString)

	var payload ske.CreateOrUpdateClusterPayload
	var err error
	trimmedPayloadString := strings.Trim(string(payloadString), `"'`)
	if strings.HasPrefix(trimmedPayloadString, "@") {
		trimmedPayloadString = strings.Trim(trimmedPayloadString[1:], `"'`)
		payloadStringBytes, err = fileReaderFunc(trimmedPayloadString)
		if err != nil {
			return nil, fmt.Errorf("read payload from file: %w", err)
		}
	}
	err = json.Unmarshal(payloadStringBytes, &payload)
	if err != nil {
		return nil, fmt.Errorf("encode payload: %w", err)
	}

	return &FlagModel{
		ProjectId: projectId,
		Name:      name,
		Payload:   payload,
	}, nil
}

func BuildRequest(ctx context.Context, model *FlagModel, apiClient *ske.APIClient) (ske.ApiCreateOrUpdateClusterRequest, error) {
	req := apiClient.CreateOrUpdateCluster(ctx, model.ProjectId, model.Name)

	req = req.CreateOrUpdateClusterPayload(model.Payload)
	return req, nil
}
