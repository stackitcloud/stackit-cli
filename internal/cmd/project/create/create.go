package create

import (
	"context"
	"fmt"
	"regexp"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/confirm"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/resourcemanager/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

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
		Use:   "create",
		Short: "Create STACKIT projects",
		Long:  "Create STACKIT projects",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a STACKIT project`,
				"$ stackit project create --parent-id xxxx --name my-project"),
			examples.NewExample(
				`Create a STACKIT project with a set of labels`,
				"$ stackit project create --parent-id xxxx --name my-project --label key=value --label foo=bar"),
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

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a project under the parent with ID %s?", *model.ParentId)
				err = confirm.PromptForConfirmation(cmd, prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req, err := buildRequest(ctx, model, apiClient)
			if err != nil {
				return fmt.Errorf("build project creation request: %w", err)
			}
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create project: %w", err)
			}

			cmd.Printf("Created project under the parent with ID %s. Project ID: %s\n", *model.ParentId, *resp.ProjectId)
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

	err := flags.MarkFlagsRequired(cmd, parentIdFlag, nameFlag)
	cobra.CheckErr(err)
}

func parseInput(cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(cmd)

	labels := flags.FlagToStringToStringPointer(cmd, labelFlag)
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
		ParentId:        flags.FlagToStringPointer(cmd, parentIdFlag),
		Name:            flags.FlagToStringPointer(cmd, nameFlag),
		Labels:          labels,
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *resourcemanager.APIClient) (resourcemanager.ApiCreateProjectRequest, error) {
	req := apiClient.CreateProject(ctx)

	authFlow, err := auth.GetAuthFlow()
	if err != nil {
		return req, fmt.Errorf("get authentication flow: %w", err)
	}
	var email string
	switch authFlow {
	case auth.AUTH_FLOW_SERVICE_ACCOUNT_TOKEN:
		email, err = auth.GetAuthField(auth.SERVICE_ACCOUNT_EMAIL)
		if err != nil {
			return req, fmt.Errorf("get email of the service account that was used to authenticate: %w", err)
		}
	case auth.AUTH_FLOW_SERVICE_ACCOUNT_KEY:
		email, err = auth.GetAuthField(auth.SERVICE_ACCOUNT_EMAIL)
		if err != nil {
			return req, fmt.Errorf("get email of the service account that was used to authenticate: %w", err)
		}
	case auth.AUTH_FLOW_USER_TOKEN:
		email, err = auth.GetAuthField(auth.USER_EMAIL)
		if err != nil {
			return req, fmt.Errorf("get your user email from configuration: %w", err)
		}
	default:
		return req, fmt.Errorf("the configured authentication flow (%s) is not supported, please report this issue", authFlow)
	}

	if email == "" {
		return req, fmt.Errorf("the authenticated subject email cannot be empty, please report this issue")
	}

	req = req.CreateProjectPayload(resourcemanager.CreateProjectPayload{
		ContainerParentId: model.ParentId,
		Name:              model.Name,
		Labels:            model.Labels,
		Members: &[]resourcemanager.ProjectMember{
			{
				Role:    utils.Ptr(ownerRole),
				Subject: utils.Ptr(email),
			},
		},
	})

	return req, nil
}
