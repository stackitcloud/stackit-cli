package disable

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/argus/client"
	argusUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/argus/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/argus"
)

const (
	instanceIdFlag = "instance-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable",
		Short: "Disables single sign-on for Grafana on Argus instances",
		Long: fmt.Sprintf("%s\n%s",
			"Disables single sign-on for Grafana on Argus instances.",
			"When disabled for an instance, the generic OAuth2 authentication is used for that instance.",
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Disable single sign-on for Grafana on an Argus instance with ID "xxx"`,
				"$ stackit argus grafana single-sign-on disable --instance-id xxx"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			instanceLabel, err := argusUtils.GetInstanceName(ctx, apiClient, model.InstanceId, model.ProjectId)
			if err != nil || instanceLabel == "" {
				instanceLabel = model.InstanceId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to disable single sign-on for instance %q?", instanceLabel)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req, err := buildRequest(ctx, model, apiClient)
			if err != nil {
				return fmt.Errorf("build request: %w", err)
			}
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("disable single sign-on: %w", err)
			}

			p.Info("Disabled single sign-on for instance %q\n", instanceLabel)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "Instance ID")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)
}

func parseInput(cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringValue(cmd, instanceIdFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient argusUtils.ArgusClient) (argus.ApiUpdateGrafanaConfigsRequest, error) {
	req := apiClient.UpdateGrafanaConfigs(ctx, model.InstanceId, model.ProjectId)
	payload, err := argusUtils.GetPartialUpdateGrafanaConfigsPayload(ctx, apiClient, model.InstanceId, model.ProjectId, utils.Ptr(false), nil)
	if err != nil {
		return req, fmt.Errorf("build request payload: %w", err)
	}
	req = req.UpdateGrafanaConfigsPayload(*payload)
	return req, nil
}
