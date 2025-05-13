package create

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/object-storage/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/objectstorage"
)

const (
	credentialsGroupNameFlag = "name"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	CredentialsGroupName string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a credentials group to hold Object Storage access credentials",
		Long:  "Creates a credentials group to hold Object Storage access credentials.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create credentials group to hold Object Storage access credentials`,
				"$ stackit object-storage credentials-group create --name example"),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a credentials group with name %q?", model.CredentialsGroupName)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create Object Storage credentials group: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(credentialsGroupNameFlag, "", "Name of the group holding credentials")

	err := flags.MarkFlagsRequired(cmd, credentialsGroupNameFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:      globalFlags,
		CredentialsGroupName: flags.FlagToStringValue(p, cmd, credentialsGroupNameFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *objectstorage.APIClient) objectstorage.ApiCreateCredentialsGroupRequest {
	req := apiClient.CreateCredentialsGroup(ctx, model.ProjectId, model.Region)
	req = req.CreateCredentialsGroupPayload(objectstorage.CreateCredentialsGroupPayload{
		DisplayName: utils.Ptr(model.CredentialsGroupName),
	})
	return req
}

func outputResult(p *print.Printer, outputFormat string, resp *objectstorage.CreateCredentialsGroupResponse) error {
	if resp == nil || resp.CredentialsGroup == nil {
		return fmt.Errorf("create createndials group response is empty")
	}

	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal Object Storage credentials group: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(resp, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal Object Storage credentials group: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		p.Outputf("Created credentials group %q. Credentials group ID: %s\n\n",
			utils.PtrString(resp.CredentialsGroup.DisplayName),
			utils.PtrString(resp.CredentialsGroup.CredentialsGroupId),
		)
		p.Outputf("URN: %s\n", utils.PtrString(resp.CredentialsGroup.Urn))
		return nil
	}
}
