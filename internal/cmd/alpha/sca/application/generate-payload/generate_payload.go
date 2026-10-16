package generatepayload

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	sca "github.com/stackitcloud/stackit-sdk-go/services/sca/v1alphaapi"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/fileutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sca/client"
	scautils "github.com/stackitcloud/stackit-cli/internal/pkg/services/sca/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	applicationIDFlag = "application-id"
	environmentIDFlag = "environment-id"
	filePathFlag      = "file-path"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ApplicationID *string
	EnvironmentID *string
	FilePath      *string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-payload",
		Short: "Generates a payload to create/update SCA applications",
		Long: fmt.Sprintf("%s\n%s",
			"Generates a JSON payload with values to be used as --payload input for application creation or update.",
			"See https://docs.api.stackit.cloud/documentation/sca/version/v1alpha#tag/Applications/operation/Applications_CreateApplication for information regarding the payload structure.",
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Generate a payload with default values, and adapt it with custom values for the different configuration options`,
				`$ stackit alpha sca application generate-payload --file-path ./payload.json`,
				`<Modify payload in file, if needed>`,
				`$ stackit alpha sca application create-from-payload --name application-name --payload @./payload.json`,
			),
			examples.NewExample(
				`Generate a payload with values of an application, and adapt it with custom values for the different configuration options`,
				`$ stackit alpha sca application generate-payload --application-id xxx --file-path ./payload.json`,
				`<Modify payload in file>`,
				`$ stackit alpha sca application update-from-payload --payload @./payload.json`,
			),
			examples.NewExample(
				`Generate a payload with values of an application, and preview it in the terminal`,
				`$ stackit alpha sca application generate-payload --application-id xxx`,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			var payload *sca.CreateApplicationPayload
			if model.ApplicationID == nil {
				payload = scautils.GetDefaultPayload()
			} else {
				apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
				if err != nil {
					return err
				}
				req := apiClient.DefaultAPI.GetApplication(ctx, model.ProjectId, *model.EnvironmentID, *model.ApplicationID)
				resp, err := req.Execute()
				if err != nil {
					return fmt.Errorf("read application: %w", err)
				}

				payload = &sca.CreateApplicationPayload{
					Secrets:    resp.Secrets,
					Registry:   resp.Registry,
					Scaling:    resp.Scaling,
					Network:    resp.Network,
					Containers: resp.Containers,
				}
			}

			return outputResult(params.Printer, model.FilePath, payload)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), applicationIDFlag, "")
	cmd.Flags().Var(flags.UUIDFlag(), environmentIDFlag, "")
	cmd.Flags().StringP(filePathFlag, "f", "", "If set, writes the payload to the given file. If unset, writes the payload to the standard output")
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	applicationID := flags.FlagToStringPointer(p, cmd, applicationIDFlag)
	if applicationID != nil && globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	environmentID := flags.FlagToStringPointer(p, cmd, environmentIDFlag)
	if environmentID == nil && globalFlags.ProjectId != "" {
		environmentID = &globalFlags.ProjectId
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ApplicationID:   applicationID,
		EnvironmentID:   environmentID,
		FilePath:        flags.FlagToStringPointer(p, cmd, filePathFlag),
	}

	return &model, nil
}

func outputResult(p *print.Printer, filePath *string, payload *sca.CreateApplicationPayload) error {
	if payload == nil {
		return fmt.Errorf("application payload is empty")
	}
	payloadBytes, err := json.MarshalIndent(*payload, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal create application payload: %w", err)
	}

	if filePath != nil {
		err = fileutils.WriteToFile(utils.PtrString(filePath), string(payloadBytes))
		if err != nil {
			return fmt.Errorf("write create application payload to the file: %w", err)
		}
	} else {
		p.Outputln(string(payloadBytes))
	}

	return nil
}
