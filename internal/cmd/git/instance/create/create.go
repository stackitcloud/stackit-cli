package create

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/git/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/git"
	"github.com/stackitcloud/stackit-sdk-go/services/git/wait"
)

const (
	nameFlag   = "name"
	flavorFlag = "flavor"
	aclFlag    = "acl"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Id     *string
	Name   string
	Flavor string
	Acl    []string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates STACKIT Git instance",
		Long:  "Create a STACKIT Git instance by name.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a instance with name 'my-new-instance'`,
				`$ stackit git instance create --name my-new-instance`,
			),
			examples.NewExample(
				`Create a instance with name 'my-new-instance' and flavor`,
				`$ stackit git instance create --name my-new-instance --flavor git-100`,
			),
			examples.NewExample(
				`Create a instance with name 'my-new-instance' and acl`,
				`$ stackit git instance create --name my-new-instance --acl 1.1.1.1/1`,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
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

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create the instance %q?", model.Name)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			request := buildRequest(ctx, model, apiClient)

			result, err := request.Execute()
			if err != nil {
				return fmt.Errorf("create stackit git instance: %w", err)
			}
			model.Id = result.Id

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Creating stackit git instance")
				_, err = wait.CreateGitInstanceWaitHandler(ctx, apiClient, model.ProjectId, *model.Id).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for stackit git Instance creation: %w", err)
				}
				s.Stop()
			}

			return outputResult(params.Printer, model, result)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(nameFlag, "", "The name of the instance.")
	cmd.Flags().String(flavorFlag, "", "Flavor of the instance.")
	cmd.Flags().StringSlice(aclFlag, []string{}, "Acl for the instance.")
	if err := flags.MarkFlagsRequired(cmd, nameFlag); err != nil {
		cobra.CheckErr(err)
	}
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}
	name := flags.FlagToStringValue(p, cmd, nameFlag)
	flavor := flags.FlagToStringValue(p, cmd, flavorFlag)
	acl := flags.FlagToStringSliceValue(p, cmd, aclFlag)

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Name:            name,
		Flavor:          flavor,
		Acl:             acl,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *git.APIClient) git.ApiCreateInstanceRequest {
	return apiClient.CreateInstance(ctx, model.ProjectId).CreateInstancePayload(createPayload(model))
}

func createPayload(model *inputModel) git.CreateInstancePayload {
	return git.CreateInstancePayload{
		Name:   &model.Name,
		Flavor: git.CreateInstancePayloadGetFlavorAttributeType(&model.Flavor),
		Acl:    &model.Acl,
	}
}

func outputResult(p *print.Printer, model *inputModel, resp *git.Instance) error {
	if model == nil {
		return fmt.Errorf("input model is nil")
	}
	var outputFormat string
	if model.GlobalFlagModel != nil {
		outputFormat = model.OutputFormat
	}

	return p.OutputResult(outputFormat, resp, func() error {
		p.Outputf("Created instance %q with id %s\n", model.Name, utils.PtrString(model.Id))
		return nil
	})
}
