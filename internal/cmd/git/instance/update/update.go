package update

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	git "github.com/stackitcloud/stackit-sdk-go/services/git/v1betaapi"
	"github.com/stackitcloud/stackit-sdk-go/services/git/v1betaapi/wait"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/git/client"
	gitUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/git/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	instanceIdArg = "INSTANCE_ID"

	aclFlag = "acl"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId string
	Acl        []string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", instanceIdArg),
		Short: "Updates STACKIT Git instance",
		Long:  "Updates a STACKIT Git instance.",
		Args:  args.SingleArg(instanceIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update the ACL of the instance with ID "xxx"`,
				`$ stackit git instance update xxx --acl 1.2.3.0/24`,
			),
			examples.NewExample(
				`Update the ACL of the instance with ID "xxx" with multiple ranges`,
				`$ stackit git instance update xxx --acl 1.2.3.0/24,4.5.6.0/24`,
			),
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

			instanceLabel, err := gitUtils.GetInstanceName(ctx, apiClient.DefaultAPI, model.ProjectId, model.InstanceId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceId
			} else if instanceLabel == "" {
				instanceLabel = model.InstanceId
			}

			prompt := fmt.Sprintf("Are you sure you want to update instance %q?", instanceLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			request := buildRequest(ctx, model, apiClient)

			resp, err := request.Execute()
			if err != nil {
				return fmt.Errorf("update STACKIT Git instance: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				err := spinner.Run(params.Printer, "Updating STACKIT Git instance", func() error {
					resp, err = wait.UpdateGitInstanceWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.InstanceId).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for STACKIT Git instance update: %w", err)
				}
			}

			return outputResult(params.Printer, model, instanceLabel, resp)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.CIDRSliceFlag(), aclFlag, "List of IP networks in CIDR notation which are allowed to access this instance. Replaces the existing ACL.")
}

func parseInput(p *print.Printer, cmd *cobra.Command, cliArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	acl := flags.FlagToStringSliceValue(p, cmd, aclFlag)

	if acl == nil {
		return nil, &errors.EmptyUpdateError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      cliArgs[0],
		Acl:             acl,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *git.APIClient) git.ApiPatchInstanceRequest {
	return apiClient.DefaultAPI.PatchInstance(ctx, model.ProjectId, model.InstanceId).PatchInstancePayload(git.PatchInstancePayload{
		Acl: model.Acl,
	})
}

func outputResult(p *print.Printer, model *inputModel, instanceLabel string, resp *git.Instance) error {
	if model == nil {
		return fmt.Errorf("input model is nil")
	}

	return p.OutputResult(model.OutputFormat, resp, func() error {
		if resp == nil {
			return fmt.Errorf("API resp is nil")
		}

		operationState := "Updated"
		if model.Async {
			operationState = "Triggered update of"
		}
		p.Outputf("%s instance %q\n", operationState, instanceLabel)
		return nil
	})
}
