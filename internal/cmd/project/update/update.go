package update

import (
	"context"
	"fmt"
	"regexp"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/resourcemanager/client"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/resourcemanager"
)

const (
	parentIdFlag = "parent-id"
	nameFlag     = "name"
	labelFlag    = "label"

	ownerRole       = "project.owner"
	labelKeyRegex   = `[A-ZÄÜÖa-zäüöß0-9_-]{1,64}`
	labelValueRegex = `^$|[A-ZÄÜÖa-zäüöß0-9_-]{1,64}`
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ParentId *string
	Name     *string
	Labels   *map[string]string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Updates a STACKIT project",
		Long:  "Updates a STACKIT project.",
		Args:  args.NoArgs,
		Example: examples.Build(

			examples.NewExample(
				`Update the name of the configured STACKIT project`,
				"$ stackit project update --name my-updated-project"),
			examples.NewExample(
				`Add labels to the configured STACKIT project`,
				"$ stackit project update --label key=value,foo=bar"),
			examples.NewExample(
				`Update the name of a STACKIT project by explicitly providing the project ID`,
				"$ stackit project update --name my-updated-project --project-id xxx"),
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

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update project %q?", projectLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("update project: %w", err)
			}

			params.Printer.Info("Updated project %q\n", projectLabel)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(parentIdFlag, "", "Parent resource identifier. Both container ID (user-friendly) and UUID are supported")
	cmd.Flags().String(nameFlag, "", "Project name")
	cmd.Flags().StringToString(labelFlag, nil, "Labels are key-value string pairs which can be attached to a project. A label can be provided with the format key=value and the flag can be used multiple times to provide a list of labels")
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	labels := flags.FlagToStringToStringPointer(p, cmd, labelFlag)
	parentId := flags.FlagToStringPointer(p, cmd, parentIdFlag)
	name := flags.FlagToStringPointer(p, cmd, nameFlag)

	if labels == nil && parentId == nil && name == nil {
		return nil, &errors.EmptyUpdateError{}
	}

	if labels != nil {
		labelKeyRegex := regexp.MustCompile(labelKeyRegex)
		labelValueRegex := regexp.MustCompile(labelValueRegex)
		for key, value := range *labels {
			if !labelKeyRegex.MatchString(key) {
				return nil, &errors.FlagValidationError{
					Flag:    labelFlag,
					Details: fmt.Sprintf("label key %s didn't match the required regex expression %s", key, labelKeyRegex),
				}
			}

			if !labelValueRegex.MatchString(value) {
				return nil, &errors.FlagValidationError{
					Flag:    labelFlag,
					Details: fmt.Sprintf("label value %s for key %s didn't match the required regex expression %s", value, key, labelValueRegex),
				}
			}
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ParentId:        parentId,
		Name:            name,
		Labels:          labels,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *resourcemanager.APIClient) resourcemanager.ApiPartialUpdateProjectRequest {
	req := apiClient.PartialUpdateProject(ctx, model.ProjectId)
	req = req.PartialUpdateProjectPayload(resourcemanager.PartialUpdateProjectPayload{
		ContainerParentId: model.ParentId,
		Name:              model.Name,
		Labels:            model.Labels,
	})

	return req
}
