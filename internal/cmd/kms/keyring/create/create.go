package create

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/kms/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"

	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/kms"
	"github.com/stackitcloud/stackit-sdk-go/services/kms/wait"
)

const (
	keyRingNameFlag = "name"
	descriptionFlag = "description"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	KeyringName string
	Description string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a KMS key ring",
		Long:  "Creates a KMS key ring.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a KMS key ring with name "my-keyring"`,
				"$ stackit kms keyring create --name my-keyring"),
			examples.NewExample(
				`Create a KMS key ring with a description`,
				"$ stackit kms keyring create --name my-keyring --description my-description"),
			examples.NewExample(
				`Create a KMS key ring and print the result as YAML`,
				"$ stackit kms keyring create --name my-keyring -o yaml"),
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

			err = params.Printer.PromptForConfirmation("Are you sure you want to create a KMS key ring?")
			if err != nil {
				return err
			}

			// Call API
			req, _ := buildRequest(ctx, model, apiClient)

			keyRing, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create KMS key ring: %w", err)
			}

			// Prevent potential nil pointer dereference
			if keyRing == nil || keyRing.Id == nil {
				return fmt.Errorf("API call succeeded but returned an invalid response (missing key ring ID)")
			}

			keyRingId := *keyRing.Id

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Creating key ring")
				_, err = wait.CreateKeyRingWaitHandler(ctx, apiClient, model.ProjectId, model.Region, keyRingId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for KMS key ring creation: %w", err)
				}
				s.Stop()
			}

			return outputResult(params.Printer, model, keyRing)
		},
	}
	configureFlags(cmd)
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	keyringName := flags.FlagToStringValue(p, cmd, keyRingNameFlag)

	if keyringName == "" {
		return nil, &cliErr.DSAInputPlanError{
			Cmd: cmd,
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		KeyringName:     keyringName,
		Description:     flags.FlagToStringValue(p, cmd, descriptionFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

type kmsKeyringClient interface {
	CreateKeyRing(ctx context.Context, projectId string, regionId string) kms.ApiCreateKeyRingRequest
}

func buildRequest(ctx context.Context, model *inputModel, apiClient kmsKeyringClient) (kms.ApiCreateKeyRingRequest, error) {
	req := apiClient.CreateKeyRing(ctx, model.ProjectId, model.Region)

	req = req.CreateKeyRingPayload(kms.CreateKeyRingPayload{
		DisplayName: &model.KeyringName,

		// Description should be empty by default and only be overwritten with the descriptionFlag if it was passed.
		Description: &model.Description,
	})
	return req, nil
}

func outputResult(p *print.Printer, model *inputModel, resp *kms.KeyRing) error {
	if resp == nil {
		return fmt.Errorf("response is nil")
	}

	switch model.OutputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal KMS key ring: %w", err)
		}
		p.Outputln(string(details))

	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(resp, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal KMS key ring: %w", err)
		}
		p.Outputln(string(details))

	default:
		operationState := "Created"
		if model.Async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s key ring. KMS key ring ID: %s\n", operationState, utils.PtrString(resp.Id))
	}
	return nil
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(keyRingNameFlag, "", "Name of the KMS key ring")
	cmd.Flags().String(descriptionFlag, "", "Optional description of the key ring")

	err := flags.MarkFlagsRequired(cmd, keyRingNameFlag)
	cobra.CheckErr(err)
}
