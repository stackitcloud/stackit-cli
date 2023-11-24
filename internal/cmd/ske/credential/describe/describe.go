package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
)

const (
	clusterNameFlag = "cluster-name"
)

type flagModel struct {
	*globalflags.GlobalFlagModel
	ClusterName string
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "describe",
		Short:   "Get details of the credential associated to a SKE cluster",
		Long:    "Get details of the credential associated to a SKE cluster",
		Example: `$ stackit ske credential describe --project-id xxx --cluster-name xxx`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseFlags(cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(cmd)
			if err != nil {
				return fmt.Errorf("authentication failed, please run \"stackit auth login\" or \"stackit auth activate-service-account\"")
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("describe SKE credential: %w", err)
			}

			return outputResult(cmd, model.OutputFormat, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(clusterNameFlag, "", "Cluster name")

	err := utils.MarkFlagsRequired(cmd, clusterNameFlag)
	cobra.CheckErr(err)
}

func parseFlags(cmd *cobra.Command) (*flagModel, error) {
	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, fmt.Errorf("project ID not set")
	}

	clusterName := utils.FlagToStringValue(cmd, clusterNameFlag)
	if clusterName == "" {
		return nil, fmt.Errorf("cluster name can't be empty")
	}

	return &flagModel{
		GlobalFlagModel: globalFlags,
		ClusterName:     clusterName,
	}, nil
}

func buildRequest(ctx context.Context, model *flagModel, apiClient *ske.APIClient) ske.ApiGetCredentialsRequest {
	req := apiClient.GetCredentials(ctx, model.ProjectId, model.ClusterName)
	return req
}

func outputResult(cmd *cobra.Command, outputFormat string, credential *ske.CredentialsResponse) error {
	switch outputFormat {
	case globalflags.TableOutputFormat:
		table := tables.NewTable()
		table.SetHeader("SERVER", "TOKEN")
		table.AddRow(*credential.Server, *credential.Token)
		err := table.Render(cmd)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	default:
		details, err := json.MarshalIndent(credential, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal PostgreSQL credential: %w", err)
		}
		cmd.Println(string(details))

		return nil
	}
}
