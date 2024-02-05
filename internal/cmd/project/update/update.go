package update

import (
	"context"
	"fmt"
	"regexp"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/confirm"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
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

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a STACKIT project",
		Long:  "Update a STACKIT project.",
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
			model, err := parseInput(cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(cmd)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, cmd)
			if err != nil {
				projectLabel = model.ProjectId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update project %s?", projectLabel)
				err = confirm.PromptForConfirmation(cmd, prompt)
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

			cmd.Printf("Updated project %s\n", projectLabel)
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

func parseInput(cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	labels := flags.FlagToStringToStringPointer(cmd, labelFlag)
	parentId := flags.FlagToStringPointer(cmd, parentIdFlag)
	name := flags.FlagToStringPointer(cmd, nameFlag)

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

	return &inputModel{
		GlobalFlagModel: globalFlags,
		ParentId:        parentId,
		Name:            name,
		Labels:          labels,
	}, nil
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
